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

func TestRunFocusA(t *testing.T) {
	globalConfig.focus = path{0, 0}
	ch := NewSS()
	runCon(focusFuncs(ch)...)
	if exp := []string{"a"}; !ch.EqualSorted(exp) {
		t.Fatalf("Wrong execution of a closure test, got %v.", ch.Slice())
	}
	globalConfig.focus = path{}
}

func TestRunFocusB(t *testing.T) {
	globalConfig.focus = path{0, 1}
	ch := NewSS()
	runCon(focusFuncs(ch)...)
	if exp := []string{"b"}; !ch.EqualSorted(exp) {
		t.Fatalf("Wrong execution of a closure test, got %v.", ch.Slice())
	}
	globalConfig.focus = path{}
}

func TestRunFocusC(t *testing.T) {
	globalConfig.focus = path{1, 0}
	ch := NewSS()
	runCon(focusFuncs(ch)...)
	if exp := []string{"c"}; !ch.EqualSorted(exp) {
		t.Fatalf("Wrong execution of a closure test, got %v.", ch.Slice())
	}
	globalConfig.focus = path{}
}

func TestRunFocusD(t *testing.T) {
	globalConfig.focus = path{1, 1}
	ch := NewSS()
	runCon(focusFuncs(ch)...)
	if exp := []string{"d"}; !ch.EqualSorted(exp) {
		t.Fatalf("Wrong execution of a closure test, got %v.", ch.Slice())
	}
	globalConfig.focus = path{}
}
