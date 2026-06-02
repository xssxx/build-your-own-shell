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
		var stderr io.Writer = os.Stderr
		var fileToClose *os.File

		// redirection handle
		if len(args) >= 2 {
			op := args[len(args)-2]

			switch op {
			case ">", "1>":
				args = handleRedirection(args, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, &stdout, &fileToClose)
			case "2>":
				args = handleRedirection(args, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, &stderr, &fileToClose)
			case ">>", "1>>":
				args = handleRedirection(args, os.O_APPEND|os.O_CREATE|os.O_WRONLY, &stdout, &fileToClose)
			case "2>>":
				args = handleRedirection(args, os.O_APPEND|os.O_CREATE|os.O_WRONLY, &stderr, &fileToClose)
			}
		}

		// commands handle
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
			cmd.Stderr = stderr
			cmd.Run()
		} else {
			fmt.Println(cmd + ": command not found")
		}

		if fileToClose != nil {
			fileToClose.Close()
			fileToClose = nil
		}

		stdout = os.Stdout
		stderr = os.Stderr
	}
}
