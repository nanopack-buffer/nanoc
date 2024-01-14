package generator

import "strings"

// Lines joins all the given strings with a newline character. Empty strings are ignored.
func Lines(lines ...string) string {
	var ls []string
	for _, l := range lines {
		if l != "" {
			ls = append(ls, l)
		}
	}
	return strings.Join(ls, "\n")
}
