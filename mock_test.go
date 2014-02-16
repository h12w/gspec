package gspec

import (
	"encoding/json"
	exp "github.com/hailiang/gspec/expectation"
	"io/ioutil"
	"runtime"
	"sort"
	"strings"
	"sync"
)

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	exp.Sprint = func(v interface{}) string {
		buf, _ := json.MarshalIndent(v, "    ", "  ")
		return string(buf)
	}
}

var (
	globalScheduler = NewScheduler(NewTextReporter(ioutil.Discard))
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

type MockReporter struct {
	mu     sync.Mutex
	groups []*TestGroup
}

func (l *MockReporter) Start() {
	l.mu.Lock()
	defer l.mu.Unlock()
}

func (l *MockReporter) End(groups []*TestGroup) {
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.groups == nil {
		l.groups = groups
	} else {
		panic("End should be only called once.")
	}
}

func (l *MockReporter) Progress(g *TestGroup, s *Stats) {
	l.mu.Lock()
	defer l.mu.Unlock()
}
