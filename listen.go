// Copyright 2014, Hǎiliàng Wáng. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gspec

import (
	"sync"
)

type listener struct {
	groups TestGroups
	m      map[funcID]*TestGroup
	mu     sync.Mutex
	Reporter
	Stats
}

func newListener(r Reporter) *listener {
	return &listener{
		m:        make(map[funcID]*TestGroup),
		Reporter: r}
}

func (l *listener) groupStart(g *TestGroup, path path) {
	l.mu.Lock()
	defer l.mu.Unlock()
	id := path[len(path)-1]
	if l.m[id] != nil {
		return
	}
	l.Total++
	if len(path) == 1 { // root node
		l.groups = append(l.groups, g)
	} else {
		parentID := path[len(path)-2]
		parent := l.m[parentID] // must exists
		if len(parent.Children) == 0 {
			l.Total--
		}
		parent.Children = append(parent.Children, g)
	}
	l.m[id] = g
	l.progress(g)
}

func (l *listener) groupEnd(err error, id funcID) {
	l.mu.Lock()
	defer l.mu.Unlock()
	g := l.m[id]
	g.Error = err
	if len(g.Children) == 0 {
		l.Ended++
	}
	l.progress(g)
}

func (l *listener) progress(g *TestGroup) {
	l.Progress(g, &l.Stats)
}
