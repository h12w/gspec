package core

import (
	"testing"
)

/*
Scenario: run a specified test case
	Given a series of test cases
	When the focus argument is set to a specific path
	Then only the corresponding test case is executed
*/

func focusFuncs(ch *SS) []TestFunc {
	return []TestFunc{
		func(s S) {
			do := aliasGroup(s)
			do(func() {
				ch.Send("x")
				do(func() {
					ch.Send("a")
				})
				do(func() {
					ch.Send("b")
				})
			})
		},
		func(s S) {
			do := aliasGroup(s)
			do(func() {
				ch.Send("y")
				do(func() {
					ch.Send("c")
				})
				do(func() {
					ch.Send("d")
				})
			})
		},
	}
}

type focusTestCase struct {
	path Path
	exp  []string
}

var focusTestCases = []focusTestCase{
	{Path{0, 0, 0}, []string{"x", "a"}},
	{Path{0, 0, 1}, []string{"x", "b"}},
	{Path{0, 1, 0}, []string{"y", "c"}},
	{Path{0, 1, 1}, []string{"y", "d"}},
}

func TestRunFocus(t *testing.T) {
	for _, c := range focusTestCases {
		ch := NewSS()
		NewController(&ReporterStub{}).Start(c.path, true, focusFuncs(ch)...)
		if !ch.Equal(c.exp) {
			t.Fatalf("Wrong execution of a closure test, got %v.", ch.Slice())
		}
	}
}
