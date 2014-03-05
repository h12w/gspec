// Copyright 2014, Hǎiliàng Wáng. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package core

import (
	"fmt"
	"testing"

	exp "github.com/hailiang/gspec/expectation"
)

func TestFuncUniqueID(t *testing.T) {
	f1 := func() {}
	f2 := func() {}
	if getFuncAddress(f1) != getFuncAddress(f1) {
		t.Fatalf("Does not return the same id for the same function.")
	}
	if getFuncAddress(f1) == getFuncAddress(f2) {
		t.Fatalf("Return the same id for different functions.")
	}
}

func TestFuncIDString(t *testing.T) {
	expect := exp.Alias(exp.TFail(t))
	expect((&funcID{0x12AB, 0}).String()).Equal("12AB")
	expect((&funcID{0x12AB, 3}).String()).Equal("12AB-3")

	id, err := parseFuncID("EF98")
	expect(err).Equal(nil)
	expect(id).Equal(funcID{0xEF98, 0})
	id, err = parseFuncID("98EF-5")
	expect(err).Equal(nil)
	expect(id).Equal(funcID{0x98EF, 5})

	id, err = parseFuncID("XYZ")
	expect(err).NotEqual(nil)
	expect(id).Equal(funcID{})
}

func TestP(t *testing.T) {
	if err := p(""); err != nil {
		t.Fatalf("fmt.Println return err %v", err)
	}
}

func p(v ...interface{}) error {
	_, err := fmt.Println(v...)
	return err
}
