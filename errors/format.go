// Copyright 2014, Hǎiliàng Wáng. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package errors

import (
	"bytes"
	"fmt"
	"path"
	"runtime"
	"strings"
)

var (
	// Sprint is used by checker to format actual and expected variables to
	// strings. It has a default implementation and can be replaced with an
	// external function.
	Sprint = func(v interface{}) string {
		return fmt.Sprintf("%v", v)
	}

	// IndentString is the default value for a level of indent.
	IndentString = "    "
)

// Pos represents a position in the source file.
type Pos struct {
	File string
	Line int
}

// GetPos get the current position of execution.
func GetPos(skip int) *Pos {
	_, file, line, _ := runtime.Caller(skip + 1)
	return &Pos{file, line}
}

// BasePath returns the base path of the source file.
func (pos *Pos) BasePath() string {
	return path.Base(pos.File)
}

// Decorate prefixes the string with the file and line of the call site
// and inserts the final newline if needed and indentation tabs for formatting.
func (pos *Pos) Decorate(s, indent string) string {
	var buf bytes.Buffer
	// Every line is indented at least one tab.
	buf.WriteString(indent)
	fmt.Fprintf(&buf, "%s:%d:", pos.BasePath(), pos.Line)

	if strings.Contains(s, "\n") {
		buf.WriteByte('\n')
		buf.WriteString(Indent(s, IndentString))
	} else {
		buf.WriteByte(' ')
		buf.WriteString(s)
	}
	return buf.String()
}

func toLines(s string) []string {
	lines := strings.Split(s, "\n")
	if l := len(lines); l > 1 && lines[l-1] == "" {
		lines = lines[:l-1]
	}
	return lines
}

// Indent splits s to lines and indent each line with argument indent.
func Indent(s, indent string) string {
	lines := toLines(s)
	for i, line := range lines {
		if len(line) > 0 {
			lines[i] = indent + line
		}
	}
	return strings.Join(lines, "\n")
}
