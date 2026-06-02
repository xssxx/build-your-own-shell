package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strings"
)

// manually parse arguments to handle quote and backslash
func parseArgs(input string) []string {
	var args []string
	var curr strings.Builder
	inSingle := false
	inDouble := false

	for i := 0; i < len(input); i++ {
		ch := input[i]

		switch {
		case inSingle:
			if ch == '\'' {
				inSingle = false
			} else {
				curr.WriteByte(ch)
			}

		case inDouble:
			if ch == '\\' && i+1 < len(input) {
				next := input[i+1]
				if next == '"' || next == '\\' {
					curr.WriteByte(next)
					i++
				} else {
					curr.WriteByte(ch)
				}
			} else if ch == '"' {
				inDouble = false
			} else {
				curr.WriteByte(ch)
			}

		default:
			switch {
			case ch == '\\' && i+1 < len(input):
				curr.WriteByte(input[i+1])
				i++

			case ch == '\'':
				inSingle = true

			case ch == '"':
				inDouble = true

			case ch == ' ' || ch == '\t':
				if curr.Len() > 0 {
					args = append(args, curr.String())
					curr.Reset()
				}

			default:
				curr.WriteByte(ch)
			}
		}
	}

	if curr.Len() > 0 {
		args = append(args, curr.String())
	}

	return args
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
