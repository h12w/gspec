package gspec

import (
	"fmt"
	"io"
)

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

type Reporter interface {
	Start()
	End(groups []*TestGroup)
	Progress(g *TestGroup, s *Stats)
}

/*
passed failed skipped pending
*/
type Stats struct {
	Total int
	Ended int
}

type TextReporter struct {
	Stats
	w io.Writer
}

func NewTextReporter(w io.Writer) *TextReporter {
	return &TextReporter{w: w}
}

func (l *TextReporter) Start() {
	l.Stats = Stats{}
}

func (l *TextReporter) End(groups []*TestGroup) {
	fmt.Fprintln(l.w, "")
}

func (l *TextReporter) Progress(g *TestGroup, s *Stats) {
	if s.Ended > l.Ended {
		sym := "."
		if g.Error != nil {
			sym = "F"
		}
		fmt.Fprint(l.w, sym)
	}
	l.Stats = *s
}
