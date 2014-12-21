// Copyright 2014, Hǎiliàng Wáng. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package core

import (
	"sort"
	"sync"
	"time"

	ext "h12.me/gspec/extension"
)

type collector struct {
	group *ext.TestGroup
	m     map[string]*ext.TestGroup
	mu    sync.Mutex
	r     ext.Reporter
	ext.Stats
}

func newCollector(r ext.Reporter) *collector {
	return &collector{
		m: make(map[string]*ext.TestGroup),
		r: r,
	}
}

func (c *collector) groupStart(g *ext.TestGroup, path Path) {
	c.mu.Lock()
	defer c.mu.Unlock()
	id := path.String()
	if c.m[id] != nil {
		return
	}
	c.Total++
	if len(path) == 1 { // root node
		c.group = g
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

func (c *collector) groupEnd(err error, path Path) {
	c.mu.Lock()
	defer c.mu.Unlock()
	id := path.String()
	if g, ok := c.m[id]; ok {
		g.Error = err
		if len(g.Children) == 0 {
			c.Ended++
			if err != nil {
				switch err.(type) {
				case *ext.PendingError:
					c.Pending++
				default:
					c.Failed++
				}
			}
		}
		c.progress(g)
	}
}

func (c *collector) sort() {
	sortTestGroup(c.group)
}

func (c *collector) setDuration(path Path, d time.Duration) {
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
	leaf        Path
	startTime   time.Time
	setDuration func(Path, time.Duration)
}

func (t *timer) start() {
	t.startTime = time.Now()
}

func (t *timer) end() {
	t.setDuration(t.leaf, time.Now().Sub(t.startTime))
}

// Sort sorts the elements by ID.
func sortTestGroup(s *ext.TestGroup) {
	sort.Sort(byID{s.Children})
	for _, c := range s.Children {
		sortTestGroup(c)
	}
}

// byID implements Less method of sort.Interface for sorting TestGroups by ID.
type byID struct{ ext.TestGroups }

// Len implements Len method of sort.Interface.
func (s byID) Len() int { return len(s.TestGroups) }

// Swap implements Swap method of sort.Interface.
func (s byID) Swap(i, j int) { s.TestGroups[i], s.TestGroups[j] = s.TestGroups[j], s.TestGroups[i] }

// Less implements Less method of sort.Interface.
func (s byID) Less(i, j int) bool { return s.TestGroups[i].ID < s.TestGroups[j].ID }
