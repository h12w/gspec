// Copyright 2014, Hǎiliàng Wáng. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package extension

import (
	"bytes"
	"fmt"
	"io"
	"strings"
)

// T is an interface that allows a testing.T to be passed to GSpec.
type T interface {
	Fail()
	Parallel()
}

// Reporter is a interface to accept events from tests running.
type Reporter interface {
	Start()
	End(groups TestGroups)
	Progress(g *TestGroup, s *Stats)
}

// A TestGroup contains a test group's related data.
type TestGroup struct {
	Description string
	Error       error
	Children    TestGroups
}

// For loops through each leaf node of a TestGroup.
// path contains the path from root to leaf.
func (g *TestGroup) For(visit func(path TestGroups)) {
	g.each(&groupStack{}, visit)
}

func (g *TestGroup) each(s *groupStack, visit func(path TestGroups)) {
	s.push(g)
	defer s.pop()
	for _, child := range g.Children {
		child.each(s, visit)
	}
	if len(g.Children) == 0 {
		visit(s.a)
	}
}

// Stats contains statistics of tests running.
type Stats struct {
	Total int
	Ended int
}

// TestGroups is the type of slice of *TestGroup
type TestGroups []*TestGroup

// Write writes TestGroups from root to leaf.
func (gs TestGroups) Write(w io.Writer) {
	for i, g := range gs {
		indent := strings.Repeat("  ", i)
		fmt.Fprintln(w, indent+g.Description)
		if g.Error != nil {
			fmt.Fprintln(w, Indent(g.Error.Error(), indent+"  "))
		}
	}
}

// Indent splits s to lines and indent each line with argument indent.
func Indent(s, indent string) string {
	var buf bytes.Buffer
	lines := toLines(s)
	for _, line := range lines {
		buf.WriteString(indent)
		buf.WriteString(line)
		buf.WriteByte('\n')
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
