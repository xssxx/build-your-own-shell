package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"slices"
	"strings"
)

var builtin []string = []string{"echo", "exit", "type", "pwd", "cd"}

func main() {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Fprint(os.Stdout, "$ ")
		cmd, err := reader.ReadString('\n')
		cmd = strings.TrimSpace(cmd)
		fields := parseArgs(cmd)

		if len(fields[0]) == 0 {
			os.Exit(0)
		}

		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			continue
		}

		cmdName := fields[0]
		args := fields[1:]

		var stdout io.Writer = os.Stdout
		var fileToClose *os.File

		if len(args) >= 2 && (args[len(args)-2] == ">" || args[len(args)-2] == "1>") {
			fileName := args[len(args)-1]
			file, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				continue
			}
			fileToClose = file
			stdout = file
			args = args[:len(args)-2]
		}

		if strings.HasPrefix(cmdName, "exit") {
			if fileToClose != nil {
				fileToClose.Close()
			}
			os.Exit(0)
		} else if strings.HasPrefix(cmdName, "type") {
			if len(fields) > 1 {
				handleType(fields[1])
			}
		} else if slices.Contains(builtin, cmdName) {
			execBuiltin(stdout, cmdName, args)
		} else if _, err := exec.LookPath(cmdName); err == nil {
			cmd := exec.Command(cmdName, args...)
			cmd.Stdout = stdout
			cmd.Stderr = os.Stderr
			cmd.Run()
		} else {
			fmt.Println(cmd + ": command not found")
		}

		if fileToClose != nil {
			fileToClose.Close()
			fileToClose = nil
		}

		stdout = os.Stdout
	}
}
