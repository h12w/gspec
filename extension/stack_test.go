// Copyright 2014, Hǎiliàng Wáng. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package extension

import (
	"testing"

	exp "github.com/hailiang/gspec/expectation"
)

/*
Story: Internal Tests
	Test internal types/functions
*/

func TestGroupStack(t *testing.T) {
	expect := exp.Alias(exp.TFail(t))
	s := groupStack{}
	s.push(&TestGroup{Description: "a"})
	s.push(&TestGroup{Description: "b"})
	expect(s.a).Equal(
		TestGroups{
			{Description: "a"},
			{Description: "b"},
		})
	g := s.pop()
	expect(s.a).Equal(
		TestGroups{
			{Description: "a"},
		})
	expect(g).Equal(&TestGroup{Description: "b"})
	g = s.pop()
	expect(s.a).Equal(TestGroups{})
	expect(func() { s.pop() }).Panic()
}
