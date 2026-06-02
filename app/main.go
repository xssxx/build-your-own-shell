package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"slices"
	"strings"
)

var builtin []string = []string{"echo", "exit", "type", "pwd", "cd"}

func main() {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("$ ")
		cmd, err := reader.ReadString('\n')
		cmd = strings.TrimSpace(cmd)
		fields := parseArgs(cmd)

		if len(fields[0]) == 0 {
			os.Exit(0)
		}

		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		commandName := fields[0]
		arguments := fields[1:]

		if strings.HasPrefix(commandName, "exit") {
			os.Exit(0)
		} else if strings.HasPrefix(commandName, "type") {
			handleType(fields[1])
		} else if slices.Contains(builtin, commandName) {
			execBuiltin(commandName, arguments)
		} else if _, err := exec.LookPath(commandName); err == nil {
			cmd := exec.Command(commandName, arguments...)
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			cmd.Run()
		} else {
			fmt.Println(cmd + ": command not found")
		}
	}
}
