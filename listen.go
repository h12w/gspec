package gspec

import (
	"sync"
)

type groupListeners struct {
	a  []GroupListener
	mu sync.Mutex
}

func (ls *groupListeners) add(l Listener) {
	ls.a = append(ls.a, l)
}

func (ls *groupListeners) GroupStart(g *TestGroup, path []FuncId) {
	ls.mu.Lock()
	defer ls.mu.Unlock()
	for _, l := range ls.a {
		l.GroupStart(g, path)
	}
}

func (ls *groupListeners) GroupEnd(g *TestGroup, path []FuncId) {
	ls.mu.Lock()
	defer ls.mu.Unlock()
	for _, l := range ls.a {
		l.GroupEnd(g, path)
	}
}

type GroupListener interface {
	GroupStart(g *TestGroup, path []FuncId)
	GroupEnd(g *TestGroup, path []FuncId)
}

type Listener interface {
	Start()
	End(groups []*TestGroup)
	GroupListener
}

type TestGroup struct {
	Id          FuncId
	Description string
	Error       *TestError
	Children    []*TestGroup
}

type TestError struct {
	Err  error
	File string
	Line int
}

type DescFunc func(description string, f func())

func (t *G) Alias(name string) DescFunc {
	return func(description string, f func()) {
		id := getFuncId(f)
		g := &TestGroup{
			Id:          id,
			Description: name + " " + description,
		}
		path := t.cur.slice()
		t.GroupStart(g, path)

		if t.Group(f) {
			// g.Error =
			t.GroupEnd(g, path)
		}
	}
}

func (t *G) Alias2(n1, n2 string) (_, _ DescFunc) {
	return t.Alias(n1), t.Alias(n2)
}

func (t *G) Alias3(n1, n2, n3 string) (_, _, _ DescFunc) {
	return t.Alias(n1), t.Alias(n2), t.Alias(n3)
}

type treeCollector struct {
	groups []*TestGroup
	m      map[FuncId]*TestGroup
	mu     sync.Mutex
}

func newTreeCollector() treeCollector {
	return treeCollector{m: make(map[FuncId]*TestGroup)}
}

func (c *treeCollector) AddGroup(g *TestGroup, path []FuncId) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.m[g.Id] != nil {
		return
	}
	if len(path) > 0 {
		parentId := path[len(path)-1]
		parent := c.m[parentId] // must exists
		parent.Children = append(parent.Children, g)
	} else {
		c.groups = append(c.groups, g)
	}
	c.m[g.Id] = g
}
