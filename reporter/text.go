// Copyright 2014, Hǎiliàng Wáng. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package reporter // import "h12.me/gspec/reporter"

import (
	"fmt"
	"io"
	"strings"

	ge "h12.me/gspec/errors"
	ext "h12.me/gspec/extension"
)

const gspecPath = "h12.me/gspec"

// NewTextReporter creates and initialize a new text reporter using w to write
// the output.
func NewTextReporter(w io.Writer, verbose bool) ext.Reporter {
	return &textReporter{w: w, verbose: verbose}
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
	w       io.Writer
	verbose bool
}

func (l *textReporter) End(root *ext.TestGroup) {
	mid := make(map[string]bool)
	if root.Error != nil {
		writeTestGroups(l.w, ext.TestGroups{root}, mid)
	}
	for _, g := range root.Children {
		completed := g.For(func(path ext.TestGroups) bool {
			last := path[len(path)-1]
			if l.verbose || last.Error != nil {
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
	}
	if l.Stats.Pending > 0 {
		fmt.Fprintf(l.w, ">>> PENDING COUNT: %d of %d.\n", l.Stats.Pending, l.Stats.Total)
	}
	fmt.Fprintf(l.w, ">>> TOTAL: %d.\n", l.Stats.Total)
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
			if isPending(g.Error) {
				sym = "p"
			} else {
				sym = "F"
			}
		}
		fmt.Fprint(p.w, sym)
	}
	p.Stats = *s
}

func (p *textProgresser) End(group *ext.TestGroup) {
	fmt.Fprintln(p.w, "$")
}

type dummyReporter struct{}

func (dummyReporter) Start()                              {}
func (dummyReporter) End(*ext.TestGroup)                  {}
func (dummyReporter) Progress(*ext.TestGroup, *ext.Stats) {}

// Write writes TestGroups from root to leaf.
func writeTestGroups(w io.Writer, gs ext.TestGroups, mid map[string]bool) bool {
	for i, g := range gs {
		indent := strings.Repeat("    ", i)
		if !mid[g.ID] {
			fmt.Fprintln(w, indent+g.Description)
			mid[g.ID] = true
		}
		if g.Error != nil {
			if panicError, ok := g.Error.(*ext.PanicError); ok {
				writePanicError(w, panicError)
				printFocusInstruction(w, indent, g.ID)
				//				fmt.Fprintf(w, string(panicError.SS))
				fmt.Fprintln(w, ">>> Stop printing more errors due to a panic.")
				return false
			}
			fmt.Fprintln(w, ge.Indent(g.Error.Error(), indent+"  "))
			if !isPending(g.Error) {
				printFocusInstruction(w, indent, g.ID)
			}
		}
	}
	return true
}

func printFocusInstruction(w io.Writer, indent, id string) {
	fmt.Fprintf(w, ge.Indent(`  (use "go test -focus %s" to run the test case only.)`, indent), id)
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
	if g.Error != nil && !isPending(g.Error) {
		r.t.Fail()
	}
}

func isPending(err error) bool {
	switch err.(type) {
	case *ext.PendingError:
		return true
	}
	return false
}
