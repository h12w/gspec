// Copyright 2014, Hǎiliàng Wáng. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package errors

import (
	"testing"
)

func TestGetPos(t *testing.T) {
	var check = func(int, int, int) *Pos {
		return GetPos(1)
	}
	pos := check( // this is the line that should be restored.
		0,
		1,
		2,
	)
	exp := &Pos{"format_test.go", 15}
	if pos.BasePath() != exp.File {
		t.Fatalf(`Expect "%v" to be "%v"`, pos.BasePath(), exp.File)
	}
	if pos.Line != exp.Line &&
		pos.Line != 19 { // to pass go test -cover
		t.Fatalf(`Expect "%v" to be "%v"`, pos.Line, exp.Line)
	}
}
