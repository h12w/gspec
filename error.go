package gspec

func (t *specImpl) run(f func()) *TestError {
	return capturePanic(f)
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
