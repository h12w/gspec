// Copyright 2014, Hǎiliàng Wáng. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package expectation

import (
	"testing"
)

func TestHasPrefix(t *testing.T) {
	{
		m, expect := mockExpect()
		expect("ab").HasPrefix("a")
		if m.err != nil {
			t.Errorf("HasPrefix test: ab has prefix a but returns an error.")
		}
	}
	{
		m, expect := mockExpect()
		expect("ab").HasPrefix("c")
		if m.err == nil {
			t.Errorf("HasPrefix test: ab does not have prefix c but returns no error.")
		}
	}
	{
		m, expect := mockExpect()
		expect(1).HasPrefix("c")
		if m.err == nil {
			t.Errorf("HasPrefix test: non-string value should cause an error.")
		}
	}
	{
		m, expect := mockExpect()
		expect("ab").HasPrefix(1)
		if m.err == nil {
			t.Errorf("HasPrefix test: non-string value should cause an error.")
		}
	}
}

func TestHasSuffix(t *testing.T) {
	{
		m, expect := mockExpect()
		expect("ab").HasSuffix("b")
		if m.err != nil {
			t.Errorf("HasSuffix test: ab has suffix b but returns an error.")
		}
	}
	{
		m, expect := mockExpect()
		expect("ab").HasSuffix("c")
		if m.err == nil {
			t.Errorf("HasSuffix test: ab does not have suffix c but returns no error.")
		}
	}
	{
		m, expect := mockExpect()
		expect(1).HasSuffix("c")
		if m.err == nil {
			t.Errorf("HasSuffix test: non-string value should cause an error.")
		}
	}
	{
		m, expect := mockExpect()
		expect("ab").HasSuffix(1)
		if m.err == nil {
			t.Errorf("HasSuffix test: non-string value should cause an error.")
		}
	}
}
