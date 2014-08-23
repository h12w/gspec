// Copyright 2014, Hǎiliàng Wáng. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package expectation

import (
	"bytes"
	"errors"
	"os"
	"testing"
)

func TestSetOutput(t *testing.T) {
	if std.out != os.Stdout {
		t.Errorf("Default output should be stdout.")
	}
	SetOutput(os.Stderr)
	if std.out != os.Stderr {
		t.Errorf("Output should be changed to stderr.")
	}
}

func TestTFail(t *testing.T) {
	b := new(bytes.Buffer)
	SetOutput(b)
	defer SetOutput(os.Stdout)
	mt := mockT()

	f := TFail(mt.Fail)
	f(errors.New("a"))
	if *mt != (tMock{failCnt: 1, failNowCnt: 0}) {
		t.Errorf("Only Fail should be called once: %v", mt)
	}
	if b.String() != "a\n" {
		t.Errorf(`"a\n" should be printed but got "%s".`, b.String())
	}
}

func TestTFailNow(t *testing.T) {
	b := new(bytes.Buffer)
	SetOutput(b)
	defer SetOutput(os.Stdout)
	mt := mockT()

	f := TFail(mt.FailNow)
	f(errors.New("a"))
	if *mt != (tMock{failCnt: 0, failNowCnt: 1}) {
		t.Errorf("Only FailNow should be called once: %v", mt)
	}
	if b.String() != "a\n" {
		t.Errorf(`"a\n" should be printed but got "%s".`, b.String())
	}
}

type tMock struct {
	failCnt    int
	failNowCnt int
}

func mockT() *tMock {
	return &tMock{}
}

func (t *tMock) Fail() {
	t.failCnt++
}

func (t *tMock) FailNow() {
	t.failNowCnt++
}
