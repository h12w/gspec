// Copyright 2014, Hǎiliàng Wáng. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gspec

import (
	"fmt"
	"io"
	"sync"
)

// A TestGroup contains a test group's related data.
type TestGroup struct {
	Description string
	Error       error
	Children    []*TestGroup
}

// Stats contains statistics of tests running.
type Stats struct {
	Total int
	Ended int
}

// Reporter is a interface to accept events from tests running.
type Reporter interface {
	Start()
	End(groups []*TestGroup)
	Progress(g *TestGroup, s *Stats)
}

// NewTextReporter creates and initialize a new text reporter using w to write
// the output.
func NewTextReporter(w io.Writer) Reporter {
	return &textReporter{w: w}
}

// NewTextProgresser creates and initialize a new text progresser using w to
// write the output.
func NewTextProgresser(w io.Writer) Reporter {
	return &textProgresser{w: w}
}

// broadcaster syncs, receives and broadcasts all messages via Reporter interface
type broadcaster struct {
	t  T
	a  []Reporter
	mu sync.Mutex
}

func newBroadcaster(t T, reporters []Reporter) broadcaster {
	return broadcaster{t: t, a: reporters}
}

func (b *broadcaster) Start() {
	b.mu.Lock()
	defer b.mu.Unlock()
	for _, r := range b.a {
		r.Start()
	}
}

func (b *broadcaster) End(groups []*TestGroup) {
	b.mu.Lock()
	defer b.mu.Unlock()
	for _, r := range b.a {
		r.End(groups)
	}
}

func (b *broadcaster) Progress(g *TestGroup, stats *Stats) {
	b.mu.Lock()
	defer b.mu.Unlock()
	if g.Error != nil {
		b.t.Fail()
	}
	for _, r := range b.a {
		r.Progress(g, stats)
	}
}

// TextReporter implements a simple plain text CLI reporter.
type textReporter struct {
	dummyReporter
	w io.Writer
}

func (l *textReporter) Start() {
}

func (l *textReporter) End(groups []*TestGroup) {
	fmt.Fprintln(l.w, "")
}

type textProgresser struct {
	dummyReporter
	Stats
	w io.Writer
}

func (p *textProgresser) Progress(g *TestGroup, s *Stats) {
	if s.Ended > p.Ended {
		sym := "."
		if g.Error != nil {
			sym = "F"
		}
		fmt.Fprint(p.w, sym)
	}
	p.Stats = *s
}

type dummyReporter struct{}

func (dummyReporter) Start()                      {}
func (dummyReporter) End([]*TestGroup)            {}
func (dummyReporter) Progress(*TestGroup, *Stats) {}
