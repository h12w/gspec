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
		id := getFuncID(f)
		path := append(t.current(), id)
		g := &TestGroup{
			Description: name + description,
		}
		t.run(id, func() {
			t.groupStart(g, path)
			terr := capturePanic(f)
			if terr == nil {
				if t.err != nil {
					terr = &TestError{Err: t.err} // TODO: fill other fields
					t.err = nil
				}
			}
			t.groupEnd(terr, id)
		})
	}
}

// Fail notify that the test case has failed with an error.
func (t *specImpl) Fail(err error) {
	if t.err != nil {
		t.err = err // only keeps the first failure.
	}
}

func capturePanic(f func()) (terr *TestError) {
	defer func() {
		if err := recover(); err != nil {
			terr = &TestError{
				Err:  err,
				File: "",
				Line: 0,
			}
			// TODO: print error, terminate all tests and exit
		}
	}()
	f()
	return
}
