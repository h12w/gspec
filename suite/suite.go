// Copyright 2014, Hǎiliàng Wáng. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package suite

import (
	"os"
	"github.com/hailiang/gspec"
)

var (
	Reporter = gspec.NewTextReporter(os.Stdout)
	testFunctions []gspec.TestFunc
)

type T interface {
	Fail()
}

func Add(fs ...gspec.TestFunc) {
	testFunctions = append(testFunctions, fs...)
}

func Run(t T, sequential bool) {
	s := gspec.NewScheduler(Reporter)
	s.Start(sequential, testFunctions...)
}


