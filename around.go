package gspec

// DescFunc is the type of the function to define a test group with a
// descritpion and a closure.
type DescFunc func(description string, f func())

func (t *groupContext) Group(f func()) {
	alias := t.Alias("")
	alias("", f)
}

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
			t.groupEnd(id, capturePanic(f))
		})
	}
}

func (t *groupContext) Alias2(n1, n2 string) (_, _ DescFunc) {
	return t.Alias(n1), t.Alias(n2)
}

func (t *groupContext) Alias3(n1, n2, n3 string) (_, _, _ DescFunc) {
	return t.Alias(n1), t.Alias(n2), t.Alias(n3)
}

func capturePanic(f func()) (terr *TestError) {
	defer func() {
		if err := recover(); err != nil {
			terr = &TestError{
				Err:  err,
				File: "",
				Line: 0,
			}
		}
	}()
	f()
	return
}
