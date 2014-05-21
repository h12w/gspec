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

func (s serial) String() string {
	return fmt.Sprint(int(s))
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

func (s *serialStack) push(i serial) {
	s.path = append(s.path, i)
}
func (s *serialStack) pop() (i serial) {
	if len(s.path) == 0 {
		panic("call pop when serialStack is empty.")
	}
	s.path, i = s.path[:len(s.path)-1], s.path[len(s.path)-1]
	return
}

// pathQueue implements a thread safe queue for path (test group destination).
type pathQueue struct {
	a  []path
	mu sync.Mutex
}

func (q *pathQueue) count() int {
	q.mu.Lock()
	defer q.mu.Unlock()
	return len(q.a)
}

func (q *pathQueue) enqueue(p path) {
	q.mu.Lock()
	defer q.mu.Unlock()
	q.a = append(q.a, p.clone())
}

func (q *pathQueue) dequeue() (p path) {
	q.mu.Lock()
	defer q.mu.Unlock()
	p, q.a = q.a[0], q.a[1:]
	return
}
