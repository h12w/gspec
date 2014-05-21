// Copyright 2014, Hǎiliàng Wáng. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package core

import (
	"fmt"
	"strings"
	"sync"
)

// serial is an integer serial number of the execution of a test closure
type serial int

func (f serial) String() string {
	return fmt.Sprint(int(f))
}

func parseSerial(s string) (f serial, _ error) {
	n, err := fmt.Sscanf(s, "%d", &f)
	if n == 1 {
		return f, nil
	}
	return 0, err
}

type path []serial

func (p path) clone() path {
	return append(path{}, p...)
}

func (p path) onPath(dst path) bool {
	last := imin(len(p), len(dst)) - 1
	for i := 0; i <= last; i++ {
		if p[i] != dst[i] {
			return false
		}
	}
	return true
}
func imin(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// Set implements the String method of flag.Value interface.
func (p path) String() string {
	ss := make([]string, len(p))
	for i := range ss {
		ss[i] = p[i].String()
	}
	return strings.Join(ss, "/")
}

// Set implements the Set method of flag.Value interface.
func (p *path) Set(s string) (err error) {
	ss := strings.Split(s, "/")
	*p = make(path, len(ss))
	for i := range ss {
		(*p)[i], err = parseSerial(ss[i])
		if err != nil {
			*p = nil
			return
		}
	}
	return
}

type serialStack struct {
	path
}

func (p *serialStack) push(i serial) {
	p.path = append(p.path, i)
}
func (p *serialStack) pop() (i serial) {
	if len(p.path) == 0 {
		panic("call pop when serialStack is empty.")
	}
	p.path, i = p.path[:len(p.path)-1], p.path[len(p.path)-1]
	return
}

// pathQueue implements a thread safe queue for path (test group destination).
type pathQueue struct {
	a  []path
	mu sync.Mutex
}

func (f *pathQueue) count() int {
	f.mu.Lock()
	defer f.mu.Unlock()
	return len(f.a)
}

func (f *pathQueue) enqueue(p path) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.a = append(f.a, p.clone())
}

func (f *pathQueue) dequeue() (p path) {
	f.mu.Lock()
	defer f.mu.Unlock()
	p, f.a = f.a[0], f.a[1:]
	return
}
