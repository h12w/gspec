package gspec

import (
	"sync"
)

var (
	NewCollector = func() collector {
		return &treeCollector{}
	}
)

type CollectorBuilder func() collector

type collector interface {
	Start()
	GroupStart(g *TestGroup, path []FuncId)
	GroupEnd(id FuncId, err *TestError)
	End()
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
		t.GroupStart(
			&TestGroup{
				Id:          id,
				Description: name + " " + description,
			},
			append([]FuncId{}, t.cur.a...),
		)

		if t.Group(f) {
			t.GroupEnd(id, nil)
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
	Groups []*TestGroup
	m      map[FuncId]*TestGroup
	mu     sync.Mutex
}

func (c *treeCollector) GroupStart(g *TestGroup, path []FuncId) {
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
		c.Groups = append(c.Groups, g)
	}
	c.m[g.Id] = g
}

func (c *treeCollector) GroupEnd(id FuncId, err *TestError) {
	c.m[id].Error = err
}

func (c *treeCollector) Start() {
	c.m = make(map[FuncId]*TestGroup)
}

func (c *treeCollector) End() {
}
