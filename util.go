package gspec

import "strings"

// Unindent is a utility function that unindents Go's raw string literal by the
// indent guessed from the first nonblank line, so the raw string literal can be
// indented as normal code and looks better. A prefix and a suffix newline '\n'
// will be added if there is none.
func Unindent(s string) string {
	lines := strings.Split(s, "\n")
	indent := ""
	done := false
	for _, line := range lines {
		if strings.TrimSpace(line) != "" {
			for _, r := range line {
				if r == ' ' || r == '\t' {
					indent += string(r)
				} else {
					done = true
					break
				}
			}
		}
		if done {
			break
		}
	}
	for i := range lines {
		if strings.HasPrefix(lines[i], indent) {
			lines[i] = lines[i][len(indent):]
		} else {
			lines[i] = strings.TrimLeft(lines[i], "\t")
		}
	}
	if lines[0] == "" {
		lines = lines[1:]
	}
	s = strings.Join(lines, "\n")
	return s
}
