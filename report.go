package gspec

import (
	"fmt"
	"io"
	"sync"
)

// A TestGroup contains a test group's related data.
type TestGroup struct {
	ID          FuncID
	Description string
	Error       *TestError
	//	Parent      *TestGroup
	Children []*TestGroup
}

// TestError contains the error information of a finished test group.
type TestError struct {
	Err  interface{}
	File string
	Line int
}

// Reporter is a interface to accept events from tests running.
type Reporter interface {
	Start()
	End(groups []*TestGroup)
	Progress(g *TestGroup, s *Stats)
}

/*
passed failed skipped pending
*/

// Stats contains statistics of tests running.
type Stats struct {
	Total int
	Ended int
}

// TextReporter implements a simple plain text CLI reporter.
type textReporter struct {
	Stats
	w  io.Writer
	mu sync.Mutex
}

// NewTextReporter creates and initialize a new TextReporter using w to write
// the output.
func NewTextReporter(w io.Writer) Reporter {
	return &textReporter{w: w}
}

func (l *textReporter) Start() {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.Stats = Stats{}
}

func (l *textReporter) End(groups []*TestGroup) {
	l.mu.Lock()
	defer l.mu.Unlock()
	fmt.Fprintln(l.w, "")
}

func (l *textReporter) Progress(g *TestGroup, s *Stats) {
	l.mu.Lock()
	defer l.mu.Unlock()
	if s.Ended > l.Ended {
		sym := "."
		if g.Error != nil {
			sym = "F"
		}
		fmt.Fprint(l.w, sym)
	}
	l.Stats = *s
}
