// Copyright 2014, Hǎiliàng Wáng. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package core

import (
	"testing"

	exp "github.com/hailiang/gspec/expectation"
)

/*
Story: Internal Tests
	Test internal types/functions
*/

func TestFuncIDString(t *testing.T) {
	expect := exp.Alias(exp.TFail(t))
	expect(serial(0).String()).Equal("0")
	expect(serial(3).String()).Equal("3")

	id, err := parseSerial("2")
	expect(err).Equal(nil)
	expect(id).Equal(serial(2))

	id, err = parseSerial("XYZ")
	expect(err).NotEqual(nil)
	expect(id).Equal(serial(0))
}

func TestPathSerialization(t *testing.T) {
	expect := exp.Alias(exp.TFailNow(t))

	var p path
	p.Set("0/1/2")
	expect(len(p)).Equal(3)
	expect(p[0]).Equal(serial(0))
	expect(p[1]).Equal(serial(1))
	expect(p[2]).Equal(serial(2))

	err := p.Set("UVW")
	expect(err).NotEqual(nil)

	p = path{0, 1, 2}
	expect(p.String()).Equal("0/1/2")
}

func TestIDStack(t *testing.T) {
	expect := exp.Alias(exp.TFail(t))
	p := serialStack{}
	p.push(serial(1))
	p.push(serial(2))
	expect(p.path).Equal(path{1, 2})
	i := p.pop()
	expect(p.path).Equal(path{1})
	expect(i).Equal(serial(2))
	i = p.pop()
	expect(p.path).Equal(path{})
	expect(func() { p.pop() }).Panic()
}
