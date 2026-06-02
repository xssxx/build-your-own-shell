package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

var _ = fmt.Print

func main() {
	readInput()
}

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
		} else {
			fmt.Println(command + ": command not found")
		}
	}
}
