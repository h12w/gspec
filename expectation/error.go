// Copyright 2014, Hǎiliàng Wáng. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package expectation

// An Error is returned when the expectation fails.
type Error struct {
	Msg string
}

func (e *Error) Error() string {
	return e.Msg
}
