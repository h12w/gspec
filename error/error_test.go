// Copyright 2014, Hǎiliàng Wáng. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package error

import (
	"strconv"
	"testing"
)

func TestCompareError(t *testing.T) {
	// inline
	ce := Compare(0, 1, "to equal", "int", 0)
	exp := "error_test.go:14: expect int 0 to equal 1."
	if msg := ce.Error(); msg != exp {
		t.Errorf(`Expect error message "%s" but got "%s"`, exp, msg)
	}

	// multiple lines
	ce = Compare("\na\nb\n", "\na\nc\n", "to equal", "msg", 0)
	exp = `error_test.go:21:
    expect msg
        a
        b
    to equal
        a
        c
    .
`
	if msg := ce.Error(); msg != exp {
		t.Errorf("Expect error message\n%s\nbut got\n%s", strconv.Quote(exp),
			strconv.Quote(msg))
	}
}

func TestExpectError(t *testing.T) {
	ce := Expect("a", 0)
	exp := "error_test.go:38: expect a."
	if msg := ce.Error(); msg != exp {
		t.Errorf(`Expect error message "%s" but got "%s"`, exp, msg)
	}
}
