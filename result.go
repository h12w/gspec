// Copyright 2014, Hǎiliàng Wáng. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gspec

import (
	"fmt"
	"io"
	"strings"
)

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
		fmt.Fprintln(w, strings.Repeat("    ", i), g.Description)
		if g.Error != nil {
			fmt.Fprintln(w, g.Error.Error())
		}
	}
}
