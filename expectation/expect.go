package expectation

// Actual provides checking methods for an actual value in an expectation.
type Actual struct {
	v    interface{}
	fail func(*Error)
}

// To is a general method for checking an expectation.
func (a *Actual) To(check Checker, expected interface{}) {
	if err := check(a.v, expected); err != nil {
		a.fail(err)
	}
}

// An Error is returned when the expectation fails.
type Error struct {
	Msg string
}

func (e *Error) Error() string {
	return e.Msg
}

// ExpectFunc is the type of function that returns an Actual object given an
// actual value.
type ExpectFunc func(actual interface{}) *Actual

// Alias registers a fail function and returns an ExpectFunc.
func Alias(fail func(*Error)) ExpectFunc {
	return func(actual interface{}) *Actual {
		return &Actual{actual, fail}
	}
}

// T is a subset of testing.T used in this package.
type T interface {
	Error(...interface{})
}

// AliasForT registers T as the fail handler and returns an ExpectFunc.
func AliasForT(t T) ExpectFunc {
	return Alias(func(err *Error) {
		t.Error(err.Msg)
	})
}
