// Copyright 2014, Hǎiliàng Wáng. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gspec

import (
	"fmt"
	"reflect"
)

func p(v ...interface{}) error {
	_, err := fmt.Println(v...)
	return err
}

// FuncID is an ID unique for each function (closure)
type FuncID uintptr

func getFuncID(f interface{}) FuncID {
	return FuncID(reflect.ValueOf(f).Pointer())
}

func imin(a, b int) int {
	if a < b {
		return a
	}
	return b
}
