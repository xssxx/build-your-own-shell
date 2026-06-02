package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strings"
)

var builtin []string = []string{"echo", "exit", "type", "pwd", "cd"}

func parseArgs(input string) []string {
	var args []string
	var curr strings.Builder
	inSingle := false
	inDouble := false

	for i := 0; i < len(input); i++ {
		ch := input[i]

		switch {
		case ch == '\'' && !inDouble:
			inSingle = !inSingle
		case ch == '"' && !inSingle:
			inDouble = !inDouble
		case (ch == ' ' || ch == '\t') && !inSingle && !inDouble:
			if curr.Len() > 0 {
				args = append(args, curr.String())
				curr.Reset()
			}
		default:
			curr.WriteByte(ch)
		}

	}

	if curr.Len() > 0 {
		args = append(args, curr.String())
	}

	return args
}

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

func execBuiltin(cmd string, args []string) {
	switch cmd {
	case "pwd":
		dir, err := os.Getwd()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			return
		}
		fmt.Println(dir)
	case "cd":
		path := strings.TrimSpace(args[0])

		// handle home directory (~) manually
		if path == "~" || strings.HasPrefix(path, "~/") {
			home, err := os.UserHomeDir()
			if err != nil {
				fmt.Println("cd: " + path + ": No such file or directory")
			}

			if path == "~" {
				path = home
			} else {
				path = filepath.Join(home, path[2:])
			}
		}

		err := os.Chdir(path)
		if err != nil {
			fmt.Println("cd: " + path + ": No such file or directory")
			return
		}
	case "echo":
		for i := range args {
			args[i] = strings.TrimSpace(args[i])
		}
		fmt.Println(strings.Join(args, " "))
	default:
		fmt.Println(cmd + ": command not found")
	}
}
