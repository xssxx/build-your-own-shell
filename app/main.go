package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"slices"
	"strings"
)

func main() {
	builtin := []string{"echo", "exit", "type"}

	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("$ ")
		input, err := reader.ReadString('\n')

		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		input = strings.TrimSpace(input)
		inputParts := strings.Split(input, " ")
		cmd := inputParts[0]
		args := inputParts[1:]

		switch cmd {
		case "exit":
			os.Exit(0)
		case "echo":
			fmt.Println(strings.Join(args, " "))
		case "type":
			if slices.Contains(builtin, args[0]) {
				fmt.Println(args[0] + " is a shell builtin")
			} else if path, _ := exec.LookPath(args[0]); path != "" {
				fmt.Println(args[0] + " is " + path)
			} else {
				fmt.Println(args[0] + ": not found")
			}
		default:
			fmt.Println(cmd + ": command not found")
		}
	}
}
