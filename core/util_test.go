// Copyright 2014, Hǎiliàng Wáng. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package core

import (
	"fmt"
	"testing"

	exp "github.com/hailiang/gspec/expectation"
)

func TestFuncIDString(t *testing.T) {
	expect := exp.Alias(exp.TFail(t))
	expect(funcID(0).String()).Equal("0")
	expect(funcID(3).String()).Equal("3")

	id, err := parseFuncID("2")
	expect(err).Equal(nil)
	expect(id).Equal(funcID(2))

	id, err = parseFuncID("XYZ")
	expect(err).NotEqual(nil)
	expect(id).Equal(funcID(0))
}

func p(v ...interface{}) error {
	_, err := fmt.Println(v...)
	return err
}
