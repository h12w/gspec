// Copyright 2014, Hǎiliàng Wáng. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package core

import (
	"sync"
	"time"

	ext "github.com/hailiang/gspec/extension"
)

type collector struct {
	groups ext.TestGroups
	m      map[string]*ext.TestGroup
	mu     sync.Mutex
	r      ext.Reporter
	ext.Stats
}

func newCollector(r ext.Reporter) *collector {
	return &collector{
		m: make(map[string]*ext.TestGroup),
		r: r,
	}
}

func (c *collector) groupStart(g *ext.TestGroup, path path) {
	c.mu.Lock()
	defer c.mu.Unlock()
	id := path.String()
	if c.m[id] != nil {
		return
	}
	c.Total++
	if len(path) == 1 { // root node
		c.groups = append(c.groups, g)
	} else {
		parentID := path[:len(path)-1].String()
		parent := c.m[parentID] // must exists
		if len(parent.Children) == 0 {
			c.Total--
		}
		parent.Children = append(parent.Children, g)
	}
	c.m[id] = g
	c.progress(g)
}

func (c *collector) groupEnd(err error, path path) {
	c.mu.Lock()
	defer c.mu.Unlock()
	id := path.String()
	if g, ok := c.m[id]; ok {
		g.Error = err
		if len(g.Children) == 0 {
			c.Ended++
			if err != nil {
				c.Failed++
			}
		}
		c.progress(g)
	}
}

func (c *collector) setDuration(path path, d time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()
	id := path.String()
	if g, ok := c.m[id]; ok {
		g.Duration = d
	}
}

func (c *collector) progress(g *ext.TestGroup) {
	c.r.Progress(g, &c.Stats)
}

// timer implements the start and end method of S, and is responsible for
// measuring test time.
type timer struct {
	leaf        path
	startTime   time.Time
	setDuration func(path, time.Duration)
}

func (t *timer) start() {
	t.startTime = time.Now()
}

func (t *timer) end() {
	t.setDuration(t.leaf, time.Now().Sub(t.startTime))
}
