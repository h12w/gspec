// Copyright 2014, Hǎiliàng Wáng. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package core

import (
	"testing"

	exp "h12.me/gspec/expectation"
)

/*
Story: Internal Tests
	Test internal types/functions
*/

func TestFuncIDString(t *testing.T) {
	expect := exp.Alias(exp.TFail(t.FailNow))
	expect(Serial(0).String()).Equal("0")
	expect(Serial(3).String()).Equal("3")

	id, err := parseSerial("2")
	expect(err).Equal(nil)
	expect(id).Equal(Serial(2))

	id, err = parseSerial("XYZ")
	expect(err).NotEqual(nil)
	expect(id).Equal(Serial(0))
}

func TestPathSerialization(t *testing.T) {
	expect := exp.Alias(exp.TFail(t.FailNow))

	var p Path
	p.Set("0/1/2")
	expect(len(p)).Equal(3)
	expect(p[0]).Equal(Serial(0))
	expect(p[1]).Equal(Serial(1))
	expect(p[2]).Equal(Serial(2))

	err := p.Set("UVW")
	expect(err).NotEqual(nil)

	p = Path{0, 1, 2}
	expect(p.String()).Equal("0/1/2")
}

func TestIDStack(t *testing.T) {
	expect := exp.Alias(exp.TFail(t.FailNow))
	p := serialStack{}
	p.push(Serial(1))
	p.push(Serial(2))
	expect(p.Path).Equal(Path{1, 2})
	i := p.pop()
	expect(p.Path).Equal(Path{1})
	expect(i).Equal(Serial(2))
	i = p.pop()
	expect(i).Equal(Serial(1))
	expect(p.Path).Equal(Path{})
	expect(func() { p.pop() }).Panic()
}
