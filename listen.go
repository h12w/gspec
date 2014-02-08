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

func (l *treeListener) setReporter(r Reporter) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.r = r
}

func (l *treeListener) groupStart(g *TestGroup, path []FuncId) {
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.m[g.Id] != nil {
		return
	}
	l.Total++
	if len(path) == 0 {
		l.groups = append(l.groups, g)
	} else {
		parentId := path[len(path)-1]
		parent := l.m[parentId] // must exists
		if len(parent.Children) == 0 {
			l.Total--
		}
		parent.Children = append(parent.Children, g)
		//	g.Parent = parent
	}
	l.m[g.Id] = g
	l.progress(g)
}

func (l *treeListener) groupEnd(id FuncId, err *TestError) {
	l.mu.Lock()
	defer l.mu.Unlock()
	g := l.m[id]
	g.Error = err
	if len(g.Children) == 0 {
		l.Ended++
	}
	l.progress(g)
}

func (l *treeListener) progress(g *TestGroup) {
	l.r.Progress(g, &l.Stats)
}
