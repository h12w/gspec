// Copyright 2014, Hǎiliàng Wáng. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gspec

type groupStack struct {
	a TestGroups
}

func (s *groupStack) push(g *TestGroup) {
	s.a = append(s.a, g)
}
func (s *groupStack) pop() (g *TestGroup) {
	if len(s.a) == 0 {
		panic("call pop when groupStack is empty.")
	}
	s.a, g = s.a[:len(s.a)-1], s.a[len(s.a)-1]
	return
}

type idStack struct {
	path
}

func (p *idStack) push(i funcID) {
	p.path = append(p.path, i)
}
func (p *idStack) pop() (i funcID) {
	if len(p.path) == 0 {
		panic("call pop when idStack is empty.")
	}
	p.path, i = p.path[:len(p.path)-1], p.path[len(p.path)-1]
	return
}
