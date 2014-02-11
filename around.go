package gspec

// DescFunc is the type of the function to define a test group with a
// descritpion and a closure.
type DescFunc func(description string, f func())

func (t *groupContext) Alias(name string) DescFunc {
	if name != "" {
		name += " "
	}
	return func(description string, f func()) {
		id := getFuncID(f)
		path := t.cur.slice()
		g := &TestGroup{
			ID:          id,
			Description: name + description,
		}
		t.group(id, func() {
			t.groupStart(g, path)
			terr := capturePanic(f)
			if terr == nil {
				if t.err != nil {
					terr = &TestError{Err: t.err} // TODO: fill other fields
					t.err = nil
				}
			}
			t.groupEnd(id, terr)
		})
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
