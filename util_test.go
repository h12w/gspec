// Copyright 2014, Hǎiliàng Wáng. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gspec

import (
	"testing"
)

func TestFuncUniqueID(t *testing.T) {
	f1 := func() {}
	f2 := func() {}
	if getFuncID(f1) != getFuncID(f1) {
		t.Fatalf("Does not return the same id for the same function.")
	}
	if getFuncID(f1) == getFuncID(f2) {
		t.Fatalf("Return the same id for different functions.")
	}
}

func TestP(t *testing.T) {
	if err := p(""); err != nil {
		t.Fatalf("fmt.Println return err %v", err)
	}
}
