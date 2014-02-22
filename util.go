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

// funcID is an ID unique for each function (closure)
type funcID struct {
	p   uintptr // pointer value of function
	ver int     // allow one function run multiple times with unique ID each time
}

func getFuncID(f interface{}) funcID {
	return funcID{reflect.ValueOf(f).Pointer(), 0}
}

func (id funcID) version(ver int) funcID {
	return funcID{id.p, ver}
}

func imin(a, b int) int {
	if a < b {
		return a
	}
	return b
}
