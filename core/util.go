// Copyright 2014, Hǎiliàng Wáng. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package core

import (
	"reflect"
)

// funcID is an ID unique for each function
type funcID struct {
	p     uintptr // address of function
	count int     // allow one function run multiple times with unique ID each time
}

func getFuncAddress(f interface{}) uintptr {
	return reflect.ValueOf(f).Pointer()
}

type funcCounter map[uintptr]int

func newFuncCounter() funcCounter {
	return make(funcCounter)
}

// funcID method gets the function ID and increment the count by one.
//
// During the running of a test case, each node gets executed exactly once from
// the root to the leaf, so the only case that causes multiple running is a
// loop in table driven test.
//
// Once a test case on the leaf is finished, other test cases run with newly
// created context (specImpl), and this counter is zeroed automatically.
func (c funcCounter) funcID(f func()) funcID {
	p := getFuncAddress(f)
	count := c[p]
	c[p]++
	return funcID{p, count}
}

func imin(a, b int) int {
	if a < b {
		return a
	}
	return b
}
