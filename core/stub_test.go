// Copyright 2014, Hǎiliàng Wáng. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package core

import (
	"fmt"
	"io/ioutil"
	"reflect"
	"runtime"
	"sort"
	"strings"
	"sync"

	"github.com/hailiang/gspec/errors"
	. "github.com/hailiang/gspec/extension"
	. "github.com/hailiang/gspec/reporter"
	ogdl "github.com/ogdl/flow"
)

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	errors.Sprint = ogdlPrint
}

func ogdlPrint(v interface{}) string {
	buf, _ := ogdl.MarshalIndent(v, "    ", "    ")
	typ := ""
	if v != nil {
		typ = reflect.TypeOf(v).String() + "\n"
	}
	return "\n" +
		typ +
		string(buf) +
		"\n"
}

var (
	globalController = NewController(&Config{}, NewTextReporter(ioutil.Discard))
)

func runCon(f ...TestFunc) {
	globalController.Start(true, f...)
}

func runSeq(f ...TestFunc) {
	globalController.Start(false, f...)
}

type SS struct {
	ss []string
	mu sync.Mutex
}

func NewSS() *SS {
	return &SS{}
}

func (c *SS) Send(s string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.ss = append(c.ss, s)
}

func (c *SS) Slice() []string {
	return c.ss
}

func (c *SS) Unsorted() []string {
	return c.ss
}

func (c *SS) Sorted() []string {
	sort.Strings(c.ss)
	return c.ss
}

func (c *SS) Equal(ss []string) bool {
	return reflect.DeepEqual(c.Unsorted(), ss)
}

func (c *SS) EqualSorted(ss []string) bool {
	return reflect.DeepEqual(c.Sorted(), ss)
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

func clearGroupForTest(gs TestGroups) {
	for i := range gs {
		gs[i].For(func(gs TestGroups) bool {
			for j := range gs {
				gs[j].ID = ""
				gs[j].Duration = 0
			}
			return true
		})
	}
}

func p(v ...interface{}) error {
	_, err := fmt.Println(v...)
	return err
}
