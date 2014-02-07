package gspec

import (
	"sync"
)

type listener interface {
	groupStart(g *TestGroup, path []FuncId)
	groupEnd(id FuncId, err *TestError)
}

type TestGroup struct {
	Id          FuncId
	Description string
	Error       *TestError
	//	Parent      *TestGroup
	Children []*TestGroup
}

type TestError struct {
	Err  interface{}
	File string
	Line int
}

type treeListener struct {
	groups []*TestGroup
	m      map[FuncId]*TestGroup
	r      Reporter
	mu     sync.Mutex
	Stats
}

func newTreeListener(r Reporter) *treeListener {
	return &treeListener{
		m: make(map[FuncId]*TestGroup),
		r: r}
}

func (c *treeListener) groupStart(g *TestGroup, path []FuncId) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.m[g.Id] != nil {
		return
	}
	c.Total++
	if len(path) == 0 {
		c.groups = append(c.groups, g)
	} else {
		parentId := path[len(path)-1]
		parent := c.m[parentId] // must exists
		if len(parent.Children) == 0 {
			c.Total--
		}
		parent.Children = append(parent.Children, g)
		//	g.Parent = parent
	}
	c.m[g.Id] = g
	c.progress(g)
}

func (c *treeListener) groupEnd(id FuncId, err *TestError) {
	c.mu.Lock()
	defer c.mu.Unlock()
	g := c.m[id]
	g.Error = err
	if len(g.Children) == 0 {
		c.Ended++
	}
	c.progress(g)
}

func (c *treeListener) progress(g *TestGroup) {
	c.r.Progress(g, &c.Stats)
}
