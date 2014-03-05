// Copyright 2014, Hǎiliàng Wáng. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package core

import (
	"fmt"
	"reflect"
	"runtime"
)

// funcID is an ID unique for each function
type funcID struct {
	p     uintptr // address of function
	count int     // allow one function run multiple times with unique ID each time
}

func (id funcID) valid() bool {
	return isValidFuncAddress(id.p)
}

func (id funcID) String() string {
	if id.count == 0 {
		return fmt.Sprintf("%X", id.p)
	}
	return fmt.Sprintf("%X-%d", id.p, id.count)
}

func parseFuncID(s string) (id funcID, _ error) {
	n, err := fmt.Sscanf(s, "%X-%d", &id.p, &id.count)
	if n == 2 {
		return id, nil
	}
	n, err = fmt.Sscanf(s, "%X", &id.p)
	if n == 1 {
		return id, nil
	}
	return id, err
}

func getFuncAddress(f interface{}) uintptr {
	return reflect.ValueOf(f).Pointer()
}

func isValidFuncAddress(p uintptr) bool {
	return runtime.FuncForPC(p) != nil
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
