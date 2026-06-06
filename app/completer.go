package main

import (
	"fmt"
	"os"
	"slices"
	"strings"
)

type TabCompleter struct {
	lastPrefix string
}

func (t *TabCompleter) Do(line []rune, pos int) ([][]rune, int) {
	current := string(line[:pos])
	matches := t.collectMatches(current)

	switch {
	case len(matches) == 0:
		t.lastPrefix = ""
		fmt.Print("\a")
		return nil, 0

	case len(matches) == 1:
		t.lastPrefix = ""
		return [][]rune{[]rune(matches[0][len(current):] + " ")}, len(current)

	default:
		if lcp := longestCommonPrefix(matches); len(lcp) > len(current) {
			t.lastPrefix = ""
			return [][]rune{[]rune(lcp[len(current):])}, len(current)
		}
		return t.handleAmbiguous(current, matches)
	}
}

func (t *TabCompleter) handleAmbiguous(current string, matches []string) ([][]rune, int) {
	if t.lastPrefix == current {
		t.lastPrefix = ""
		slices.Sort(matches)
		fmt.Print("\r\n" + strings.Join(matches, "  ") + "\r\n")
		return [][]rune{{}}, 0
	}
	t.lastPrefix = current
	fmt.Print("\a")
	return nil, 0
}

func (t *TabCompleter) collectMatches(current string) []string {
	candidates := append(executablesFromPATH(), builtin...)

	seen := map[string]bool{}
	var matches []string
	isBuiltin := slices.Contains(builtin, current)

	for _, cmd := range candidates {
		if strings.HasPrefix(cmd, current) && !seen[cmd] {
			seen[cmd] = true
			matches = append(matches, cmd)
		}
		if pathCmd := "./" + cmd; !isBuiltin && strings.HasPrefix(pathCmd, current) && !seen[pathCmd] {
			seen[pathCmd] = true
			matches = append(matches, pathCmd)
		}
	}
	return matches
}

func executablesFromPATH() []string {
	var result []string
	for dir := range strings.SplitSeq(os.Getenv("PATH"), string(os.PathListSeparator)) {
		items, err := os.ReadDir(dir)
		if err != nil {
			continue
		}
		for _, item := range items {
			info, err := item.Info()
			if err != nil || item.IsDir() || info.Mode().Perm()&0111 == 0 {
				continue
			}
			result = append(result, item.Name())
		}
	}
	return result
}

func longestCommonPrefix(strs []string) string {
	lcp := strs[0]
	for _, s := range strs[1:] {
		for !strings.HasPrefix(s, lcp) {
			lcp = lcp[:len(lcp)-1]
		}
	}
	return lcp
}
