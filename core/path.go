// Copyright 2014, Hǎiliàng Wáng. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package core

import (
	"fmt"
	"strings"
	"sync"
)

// Serial is an integer serial number of the execution of a test closure.
type Serial int

// String converts a Serial to a string.
func (s Serial) String() string {
	return fmt.Sprint(int(s))
}

func parseSerial(s string) (f Serial, _ error) {
	n, err := fmt.Sscanf(s, "%d", &f)
	if n == 1 {
		return f, nil
	}
	return 0, err
}

// Path represents a path from the root to the leaf of nested test groups.
type Path []Serial

func (p Path) clone() Path {
	return append(Path{}, p...)
}

func (p Path) onPath(dst Path) bool {
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
func (p Path) String() string {
	ss := make([]string, len(p))
	for i := range ss {
		ss[i] = p[i].String()
	}
	return strings.Join(ss, "/")
}

// Set implements the Set method of flag.Value interface.
func (p *Path) Set(s string) (err error) {
	ss := strings.Split(s, "/")
	*p = make(Path, len(ss))
	for i := range ss {
		(*p)[i], err = parseSerial(ss[i])
		if err != nil {
			*p = nil
			return
		}
	}
	return
}

// serialStack implements stack operations on Path.
type serialStack struct {
	Path
}

func (s *serialStack) push(i Serial) {
	s.Path = append(s.Path, i)
}

func (s *serialStack) pop() (i Serial) {
	if len(s.Path) == 0 {
		panic("call pop when serialStack is empty.")
	}
	s.Path, i = s.Path[:len(s.Path)-1], s.Path[len(s.Path)-1]
	return
}

// pathQueue implements a thread safe queue for path (test group destination).
type pathQueue struct {
	a  []Path
	mu sync.Mutex
}

func (q *pathQueue) count() int {
	q.mu.Lock()
	defer q.mu.Unlock()
	return len(q.a)
}

func (q *pathQueue) enqueue(p Path) {
	q.mu.Lock()
	defer q.mu.Unlock()
	q.a = append(q.a, p.clone())
}

func (q *pathQueue) dequeue() (p Path) {
	q.mu.Lock()
	defer q.mu.Unlock()
	p, q.a = q.a[0], q.a[1:]
	return
}
