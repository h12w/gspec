// Copyright 2014, Hǎiliàng Wáng. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package reporter

import (
	"fmt"
	"io"

	ext "github.com/hailiang/gspec/extension"
)

// NewTextReporter creates and initialize a new text reporter using w to write
// the output.
func NewTextReporter(w io.Writer) ext.Reporter {
	return &textReporter{w: w}
}

// NewTextProgresser creates and initialize a new text progresser using w to
// write the output.
func NewTextProgresser(w io.Writer) ext.Reporter {
	return &textProgresser{w: w}
}

// TextReporter implements a simple plain text CLI reporter.
type textReporter struct {
	dummyReporter
	w io.Writer
}

func (l *textReporter) End(groups ext.TestGroups) {
	for _, g := range groups {
		g.For(func(path ext.TestGroups) {
			last := path[len(path)-1]
			if last.Error != nil {
				path.Write(l.w)
			}
		})
	}
}

type textProgresser struct {
	ext.Stats
	w io.Writer
}

func (p *textProgresser) Start() {
	fmt.Fprint(p.w, "^")
}

func (p *textProgresser) Progress(g *ext.TestGroup, s *ext.Stats) {
	if s.Ended > p.Ended {
		sym := "."
		if g.Error != nil {
			sym = "F"
		}
		fmt.Fprint(p.w, sym)
	}
	p.Stats = *s
}

func (p *textProgresser) End(groups ext.TestGroups) {
	fmt.Fprintln(p.w, "$")
}

type dummyReporter struct{}

func (dummyReporter) Start()                              {}
func (dummyReporter) End(ext.TestGroups)                  {}
func (dummyReporter) Progress(*ext.TestGroup, *ext.Stats) {}
