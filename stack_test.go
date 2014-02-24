// Copyright 2014, Hǎiliàng Wáng. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gspec

import (
	"testing"

	exp "github.com/hailiang/gspec/expectation"
)

/*
Story: Internal Tests
	Test internal types/functions
*/

func TestIDStack(t *testing.T) {
	expect := exp.Alias(exp.TFail(t))
	p := idStack{}
	p.push(funcID{p: 1})
	p.push(funcID{p: 2})
	expect(p.path).Equal(path{{p: 1}, {p: 2}})
	i := p.pop()
	expect(p.path).Equal(path{{p: 1}})
	expect(i).Equal(funcID{p: 2})
	i = p.pop()
	expect(p.path).Equal(path{})
	expect(func() { p.pop() }).Panic()
}
