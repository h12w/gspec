package gspec

type DescFunc func(description string, f func())

func (t *G) Group(f func()) {
	alias := t.Alias("")
	alias("", f)
}

func (t *G) Alias(name string) DescFunc {
	if name != "" {
		name += " "
	}
	return func(description string, f func()) {
		id := getFuncId(f)
		path := t.cur.slice()
		g := &TestGroup{
			Id:          id,
			Description: name + description,
		}
		t.group(id, func() {
			t.groupStart(g, path)
			t.groupEnd(id, returnOnPanic(f))
		})
	}
}

func (t *G) Alias2(n1, n2 string) (_, _ DescFunc) {
	return t.Alias(n1), t.Alias(n2)
}

func (t *G) Alias3(n1, n2, n3 string) (_, _, _ DescFunc) {
	return t.Alias(n1), t.Alias(n2), t.Alias(n3)
}

func returnOnPanic(f func()) (terr *TestError) {
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
