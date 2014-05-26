// Copyright 2014, Hǎiliàng Wáng. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package expectation

import (
	"testing"
)

var stringTestCases = []expectTestCase{
	{`expect("ab").HasPrefix("a")`, func(expect ExpectFunc) { expect("ab").HasPrefix("a") }, true},
	{`expect("ab").HasPrefix("c")`, func(expect ExpectFunc) { expect("ab").HasPrefix("c") }, false},
	{`expect(1).HasPrefix("c")`, func(expect ExpectFunc) { expect(1).HasPrefix("c") }, false},
	{`expect("c").HasPrefix(1)`, func(expect ExpectFunc) { expect("c").HasPrefix(1) }, false},

	{`expect("ab").HasSuffix("b")`, func(expect ExpectFunc) { expect("ab").HasSuffix("b") }, true},
	{`expect("ab").HasSuffix("c")`, func(expect ExpectFunc) { expect("ab").HasSuffix("c") }, false},
	{`expect("c").HasSuffix(1)`, func(expect ExpectFunc) { expect("c").HasSuffix(1) }, false},

	{`expect("abcd").Contains("bc")`, func(expect ExpectFunc) { expect("abcd").Contains("bc") }, true},
	{`expect("abcd").Contains("cb")`, func(expect ExpectFunc) { expect("abcd").Contains("cb") }, false},
	{`expect("abcd").Contains(1)`, func(expect ExpectFunc) { expect("abcd").Contains(1) }, false},
}

func TestStringExpectations(t *testing.T) {
	testExpectations(t, stringTestCases)
}
