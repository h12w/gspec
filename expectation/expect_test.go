package expectation

import (
	"testing"
)

type testCase struct {
	actual, expected interface{}
	succeeded        bool
}

func TestEqual(t *testing.T) {
	cases := []testCase{
		{nil, nil, true},
		{nil, 1, false},
		{1, 1, true},
		{1, 2, false},
		{int32(1), int32(1), true},
		{1, int32(1), true},
		{1, "1", true},
		{[]int{1}, []int{1}, true},
	}

	for _, c := range cases {
		m, expect := mockExpect()
		expect(c.actual).Equal(c.expected)
		if (m.err != nil) == c.succeeded {
			t.Errorf("Equal test fails: %v", c)
		}
	}
}

func TestExpectTo(t *testing.T) {
	m, expect := mockExpect()
	expect(nil).To(func(actual, expected interface{}) *Error {
		return &Error{"x"}
	}, nil)
	if m.msg() != "x" {
		t.Errorf("Expect error message %v to be x", m.msg())
	}
}

type expectMock struct {
	err *Error
}

func mockExpect() (*expectMock, ExpectFunc) {
	m := &expectMock{}
	expect := Alias(func(err *Error) {
		m.err = err
	})
	return m, expect
}

func (m *expectMock) msg() string {
	if m.err != nil {
		return m.err.Error()
	}
	return ""
}
