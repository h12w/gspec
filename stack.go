// Copyright 2014, Hǎiliàng Wáng. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gspec

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
