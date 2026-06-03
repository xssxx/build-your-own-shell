package main

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"slices"
	"strings"

	"github.com/chzyer/readline"
)

var builtin []string = []string{"echo", "exit", "type", "pwd", "cd"}

type TabCompleter struct{}

func (t *TabCompleter) Do(line []rune, pos int) (newLine [][]rune, length int) {
	var candidates []string

	// get path of this program where it's running
	pathEnv := os.Getenv("PATH")
	for dir := range strings.SplitSeq(pathEnv, string(os.PathListSeparator)) {
		items, err := os.ReadDir(dir)
		if err != nil {
			continue
		}
		for _, item := range items {
			info, err := item.Info()
			if err != nil {
				continue
			}
			// check file permission bits if it is executable bit
			if !info.IsDir() && info.Mode().Perm()&0111 != 0 {
				candidates = append(candidates, info.Name())
			}
		}
	}

	for _, k := range builtin {
		candidates = append(candidates, k)
	}

	current := string(line[:pos])
	var suggestions [][]rune

	seen := map[string]bool{}
	var matches []string
	for _, cmd := range candidates {
		if strings.HasPrefix(cmd, current) && !seen[cmd] {
			seen[cmd] = true
			matches = append(matches, cmd)
		}
		// handle case if executable file start with `./` and it's not builtin command
		pathCmd := "./" + cmd
		if strings.HasPrefix(pathCmd, current) && !slices.Contains(builtin, cmd) && !seen[pathCmd] {
			seen[pathCmd] = true
			matches = append(matches, pathCmd)
		}
	}

	// print bell character if no matchs
	if len(matches) == 0 {
		fmt.Print("\a") // bell character: `\x07`, `\a`
		return nil, 0
	}

	for _, cmd := range matches {
		suffix := cmd[len(current):]

		if len(matches) == 1 {
			suffix += " "
		}

		suggestions = append(suggestions, []rune(suffix))
	}

	return suggestions, len(current)
}

func main() {
	rl, err := readline.NewEx(&readline.Config{
		Prompt:       "$ ",
		HistoryFile:  "/tmp/my_shell_history",
		HistoryLimit: 100,
		AutoComplete: &TabCompleter{},
	})
	rl.Config.AutoComplete = &TabCompleter{}
	if err != nil {
		panic(err)
	}
	defer rl.Close()

	for {
		cmd, err := rl.Readline()
		readline.AddHistory(cmd)
		readline.SetAutoComplete(rl.Config.AutoComplete)

		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			continue
		}

		cmd = strings.TrimSpace(cmd)
		fields := parseArgs(cmd)

		if len(fields[0]) == 0 {
			os.Exit(0)
		}

		cmdName := fields[0]
		args := fields[1:]

		var stdout io.Writer = os.Stdout
		var stderr io.Writer = os.Stderr
		var fileToClose *os.File

		// handle redirection
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

		// handle commands
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
