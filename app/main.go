package main

import (
	"bufio"
	"fmt"
	"os"
	"slices"
	"strings"
)

var _ = fmt.Print

func main() {
	readInput()
}

var builtin []string = []string{"echo", "exit", "type"}

func readInput() {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("$ ")
		command, err := reader.ReadString('\n')

		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		command = strings.TrimSpace(command)
		if command == "exit" {
			break
		} else if strings.HasPrefix(command, "echo ") {
			fmt.Println(command[5:])
		} else if strings.HasPrefix(command, "type ") {
			cmd := strings.TrimSpace(command[5:])
			if slices.Contains(builtin, cmd) {
				fmt.Println(cmd + " is a shell builtin")
			} else {
				fmt.Println(cmd + ": not found")
			}
		} else {
			fmt.Println(command + ": command not found")
		}
	}
}
