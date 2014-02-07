// Copyright 2014, Hǎiliàng Wáng. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gspec

import (
	"io"
	"os"
	"sync"
)

type RootFunc func(g *G)

// Run with default option
func Run(f RootFunc) {
	(&Runner{}).Run(f)
}

// Runner is a struct of options to configure how to run tests.
type Runner struct {
	Sequential bool
	Output     io.Writer
	Reporter   Reporter
}

func (r *Runner) Run(f RootFunc) {
	r.setDefault()
	r.setFlag() // flags have higher priority
	if r.Sequential {
		newSequentialRunner(f, r.Reporter).start()
	} else {
		newConcurrentRunner(f, r.Reporter).start()
	}
}

func (r *Runner) setDefault() {
	if r.Output == nil {
		r.Output = os.Stdout
	}
	if r.Reporter == nil {
		r.Reporter = NewTextReporter(r.Output)
	}
}

// TODO:
func (r *Runner) setFlag() {
}

type concurrentRunner struct {
	*sequentialRunner
	wg sync.WaitGroup
}

func newConcurrentRunner(f RootFunc, l Reporter) *concurrentRunner {
	r := &concurrentRunner{sequentialRunner: newSequentialRunner(f, l)}
	r.self = r
	return r
}

func (r *concurrentRunner) start() {
	defer func() {
		r.wg.Wait()
		r.treeListener.r.End(r.groups)
	}()
	r.treeListener.r.Start()
	r.run(path{})
}

func (r *concurrentRunner) run(p path) {
	r.wg.Add(1) // no need to lock
	go func() {
		defer r.wg.Done()
		r.sequentialRunner.run(p)
	}()
}

type sequentialRunner struct {
	f    RootFunc
	self scheduler
	*treeListener
}

func newSequentialRunner(f RootFunc, reporter Reporter) *sequentialRunner {
	r := &sequentialRunner{f, nil, newTreeListener(reporter)}
	r.self = r
	return r
}

func (r *sequentialRunner) start() {
	r.treeListener.r.Start()
	defer func() {
		r.treeListener.r.End(r.treeListener.groups)
	}()
	r.run(path{})
}

func (r *sequentialRunner) run(p path) {
	r.f(newG(p, r.self))
}
