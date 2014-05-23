// Copyright 2014, Hǎiliàng Wáng. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package reporter

import (
	"fmt"
	"io"
	"strings"

	"github.com/hailiang/gspec/errors"
	ext "github.com/hailiang/gspec/extension"
)

const gspecPath = "github.com/hailiang/gspec"

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
	ext.Stats
	w io.Writer
}

func (l *textReporter) End(groups ext.TestGroups) {
	mid := make(map[string]bool)
	for _, g := range groups {
		completed := g.For(func(path ext.TestGroups) bool {
			last := path[len(path)-1]
			if last.Error != nil {
				if !writeTestGroups(l.w, path, mid) {
					return false
				}
			}
			return true
		})
		if !completed {
			break
		}
	}
	if l.Stats.Failed > 0 {
		fmt.Fprintf(l.w, ">>> FAIL COUNT: %d of %d.\n", l.Stats.Failed, l.Stats.Total)
	} else {
		fmt.Fprintf(l.w, ">>> TOTAL: %d.\n", l.Stats.Total)
	}
}

func (l *textReporter) Progress(g *ext.TestGroup, s *ext.Stats) {
	l.Stats = *s
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

// Write writes TestGroups from root to leaf.
func writeTestGroups(w io.Writer, gs ext.TestGroups, mid map[string]bool) bool {
	for i, g := range gs {
		indent := strings.Repeat("    ", i)
		if mid[g.ID] {
			fmt.Fprintln(w, "")
		} else {
			fmt.Fprintln(w, indent+g.Description)
			mid[g.ID] = true
		}
		if g.Error != nil {
			if panicError, ok := g.Error.(*ext.PanicError); ok {
				writePanicError(w, panicError)
				fmt.Fprintf(w, errors.Indent("(Focus mode: go test -focus %s)", indent), g.ID)
				//				fmt.Fprintf(w, string(panicError.SS))
				fmt.Fprintln(w, ">>> Stop printing more errors due to a panic.")
				return false
			}
			fmt.Fprintln(w, errors.Indent(g.Error.Error(), indent+"  "))
			fmt.Fprintf(w, errors.Indent("(Focus mode: go test -focus %s)", indent), g.ID)
		}
	}
	return true
}

func writePanicError(w io.Writer, e *ext.PanicError) {
	fmt.Fprint(w, "panic: ")
	fmt.Fprintln(w, e.Err.Error())
	for _, f := range e.Stack {
		if strings.Contains(f.File, gspecPath) {
			fmt.Fprintln(w, "    ......")
			break
		}
		fmt.Fprint(w, "    ")
		fmt.Fprintln(w, f.Name)
		fmt.Fprint(w, "        ")
		fmt.Fprintf(w, "%s:%d\n", f.File, f.Line)
	}
}

// T is an interface that allows a testing.T to be passed.
type T interface {
	Fail()
}

type failReporter struct {
	t T
	dummyReporter
}

// NewFailReporter creates and initializes a reporter that calls T.Fail when
// any test error occurs.
func NewFailReporter(t T) ext.Reporter {
	return &failReporter{t: t}
}

func (r *failReporter) Progress(g *ext.TestGroup, stats *ext.Stats) {
	if g.Error != nil {
		r.t.Fail()
	}
}
