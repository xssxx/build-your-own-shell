package main

import (
	"bufio"
	"fmt"
	"os"
)

var _ = fmt.Print

func main() {
	readInput()
}

func readInput() {
	for {
		fmt.Print("$ ")
		command, err := bufio.NewReader(os.Stdin).ReadString('\n')

		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		fmt.Println(command[:len(command)-1] + ": command not found")
	}
}
