package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"slices"
	"strings"
)

var builtin []string = []string{"echo", "exit", "type", "pwd"}

func main() {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("$ ")
		cmd, err := reader.ReadString('\n')
		cmd = strings.TrimSpace(cmd)
		fields := strings.Fields(cmd)

		if len(fields[0]) == 0 {
			os.Exit(0)
		}

		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		if strings.HasPrefix(fields[0], "exit") {
			os.Exit(0)
		} else if strings.HasPrefix(fields[0], "echo") {
			fmt.Println(cmd[5:])
		} else if strings.HasPrefix(fields[0], "type") {
			handleType(strings.TrimSpace(fields[1]))
		} else if slices.Contains(builtin, strings.TrimSpace(fields[0])) {
			execBuiltin(strings.TrimSpace(fields[0]))
		} else if _, err := exec.LookPath(fields[0]); err == nil {
			cmd := exec.Command(fields[0], fields[1:]...)
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			cmd.Run()
		} else {
			fmt.Println(cmd + ": command not found")
		}
	}
}

func handleType(cmd string) {
	if slices.Contains(builtin, cmd) {
		fmt.Printf("%s is a shell builtin\n", cmd)
		return
	} else if path, err := exec.LookPath(cmd); err == nil {
		fmt.Printf("%s is %s\n", cmd, path)
		return
	}

	fmt.Printf("%s: not found\n", cmd)
}

func execBuiltin(cmd string) {
	switch cmd {
	case "pwd":
		dir, err := os.Getwd()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
		fmt.Println(dir)
	}
}
