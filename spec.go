package gspec

// TestFunc is the type of the function prepared to run in a goroutine for each
// test case.
type TestFunc func(S)

// S (short for "spec") provides the interface for writing tests and internally
// holds an object that contains minimal context needed to pass into a testing
// goroutine.
type S interface {
	Alias(name string) DescFunc
	Fail(err error)
}

// specImpl implements "S" interface.
type specImpl struct {
	*group
	*listener
	err error
}

// DescFunc is the type of the function to define a test group with a
// descritpion and a closure.
type DescFunc func(description string, f func())

func newSpec(g *group, l *listener) S {
	return &specImpl{g, l, nil}
}

func (t *specImpl) Alias(name string) DescFunc {
	if name != "" {
		name += " "
	}
	return func(description string, f func()) {
		t.visit(getFuncID(f), func() {
			t.groupStart(&TestGroup{Description: name + description}, t.current())
			terr := t.run(f)
			t.groupEnd(terr, getFuncID(f))
		})
	}
}
