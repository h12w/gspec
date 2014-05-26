// Copyright 2014, Hǎiliàng Wáng. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package extension

import (
	"sort"
)

// For loops through each leaf node of a TestGroup.
// path contains the path from root to leaf.
func (g *TestGroup) For(visit func(path TestGroups) bool) bool {
	return g.each(&groupStack{}, visit)
}

func (g *TestGroup) each(s *groupStack, visit func(path TestGroups) bool) bool {
	s.push(g)
	defer s.pop()
	for _, child := range g.Children {
		if !child.each(s, visit) {
			return false
		}
	}
	if len(g.Children) == 0 {
		if !visit(s.a) {
			return false
		}
	}
	return true
}

// Sort sorts the elements by ID.
func (s TestGroups) Sort() {
	sort.Sort(ByID{s})
	for _, c := range s {
		c.Children.Sort()
	}
}

// Len implements Len method of sort.Interface.
func (s TestGroups) Len() int { return len(s) }

// Swap implements Swap method of sort.Interface.
func (s TestGroups) Swap(i, j int) { s[i], s[j] = s[j], s[i] }

// Less implements Less method of sort.Interface.
func (s ByID) Less(i, j int) bool { return s.TestGroups[i].ID < s.TestGroups[j].ID }
