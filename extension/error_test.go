// Copyright 2014, Hǎiliàng Wáng. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package extension

import (
	"errors"
	"testing"

	exp "h12.me/gspec/expectation"
)

/*
Story: Internal Tests
	Test PanicError creation
*/

func TestNewPanicError(t *testing.T) {
	expect := exp.Alias(exp.TFail(t.FailNow))
	e := NewPanicError(errors.New("a"), 0)
	expect(e.Error()).Equal(errors.New("a"))

	e = NewPanicError("b", 0)
	expect(e.Error()).Equal(errors.New("b"))

	e = NewPanicError(false, 0)
	expect(e.Error()).Equal(errors.New("false"))
}

func TestNewPendingError(t *testing.T) {
	expect := exp.Alias(exp.TFail(t.FailNow))
	e := NewPendingError()
	expect(e).IsType(&PendingError{})
	expect(e.Error()).Contains("pending")
}
