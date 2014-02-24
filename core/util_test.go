// Copyright 2014, Hǎiliàng Wáng. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package core

import (
	"fmt"
	"testing"
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

func TestP(t *testing.T) {
	if err := p(""); err != nil {
		t.Fatalf("fmt.Println return err %v", err)
	}
}

func p(v ...interface{}) error {
	_, err := fmt.Println(v...)
	return err
}
