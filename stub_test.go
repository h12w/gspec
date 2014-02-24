// Copyright 2014, Hǎiliàng Wáng. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gspec

import (
	"encoding/json"
	"encoding/xml"
	"io/ioutil"
	"runtime"
	"sort"
	"strings"
	"sync"

	"github.com/davecgh/go-spew/spew"
	"github.com/hailiang/gspec/errors"
)

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	errors.Sprint = jsonPrint
}

func dumpPrint(v interface{}) string {
	spew.Config.Indent = "    "
	return "\n" + spew.Sdump(v)
}

func jsonPrint(v interface{}) string {
	buf, _ := json.MarshalIndent(v, "    ", "  ")
	return "\n    " + string(buf) + "\n"
}

func xmlPrint(v interface{}) string {
	buf, _ := xml.MarshalIndent(v, "    ", "    ")
	return "\n" + string(buf) + "\n"
}

var (
	globalScheduler = NewScheduler(&TStub{}, NewTextReporter(ioutil.Discard))
)

func Run(f ...TestFunc) {
	globalScheduler.Start(false, f...)
}

func RunSeq(f ...TestFunc) {
	globalScheduler.Start(true, f...)
}

type SChan struct {
	ch chan string
	ss []string
	wg sync.WaitGroup
}

func NewSChan() *SChan {
	return &SChan{ch: make(chan string)}
}

func (c *SChan) Send(s string) {
	c.wg.Add(1)
	go func() {
		c.ch <- s
		c.wg.Done()
	}()
}

func (c *SChan) Slice() []string {
	return c.ss
}

func (c *SChan) receiveAll() {
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		for s := range c.ch {
			c.ss = append(c.ss, s)
		}
		wg.Done()
	}()
	c.wg.Wait()
	close(c.ch)
	wg.Wait()
}

func (c *SChan) EqualSorted(ss []string) bool {
	c.receiveAll()
	sort.Strings(c.ss)
	return c.equal(ss)
}

func (c *SChan) equal(ss []string) bool {
	if len(ss) != len(c.ss) {
		return false
	}
	for i := range ss {
		if ss[i] != c.ss[i] {
			return false
		}
	}
	return true
}

func sortBytes(s string) string {
	bs := []byte(strings.TrimSpace(s))
	sort.Sort(Bytes(bs))
	return string(bs)
}

type Bytes []byte

func (bs Bytes) Len() int {
	return len(bs)
}

func (bs Bytes) Swap(i, j int) {
	bs[i], bs[j] = bs[j], bs[i]
}

func (bs Bytes) Less(i, j int) bool {
	return bs[i] < bs[j]
}

type ReporterStub struct {
	mu     sync.Mutex
	groups TestGroups
}

func (l *ReporterStub) Start() {
	l.mu.Lock()
	defer l.mu.Unlock()
}

func (l *ReporterStub) End(groups TestGroups) {
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.groups == nil {
		l.groups = groups
	} else {
		panic("End should be only called once.")
	}
}

func (l *ReporterStub) Progress(g *TestGroup, s *Stats) {
	l.mu.Lock()
	defer l.mu.Unlock()
}

type TStub struct {
	s string
}

func (m *TStub) Fail() {
	m.s += "Fail."
}

func (m *TStub) Parallel() {
}
