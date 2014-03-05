package core

import (
	"io/ioutil"
	"testing"

	exp "github.com/hailiang/gspec/expectation"
	. "github.com/hailiang/gspec/reporter"
)

func TestFocus(t *testing.T) {
	expect := exp.Alias(exp.TFail(t))
	c := config{focus: path{funcID{0x12, 0}}}
	_, err := c.dst()
	expect(err).NotEqual(nil)

	f := func() {}
	p := getFuncAddress(f)
	c = config{focus: path{funcID{p, 0}}}
	dst, err := c.dst()
	expect(err).Equal(nil)
	expect(dst).Equal(path{funcID{p, 0}})
}

func TestGlobalFocus(t *testing.T) {
	expect := exp.Alias(exp.TFail(t))
	globalConfig.focus = path{funcID{0x12, 0}}
	defer func() {
		globalConfig.focus = path{}
	}()

	s := NewScheduler(&TStub{}, NewTextReporter(ioutil.Discard))
	err := s.Start(true, func(S) {})
	expect(err).NotEqual(nil)
}
