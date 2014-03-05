package core

import (
	"flag"
	"fmt"
)

var globalConfig config

func init() {
	flag.Var(&globalConfig.focus, "focus", "test case id to select one test case to run")
}

type config struct {
	focus path
}

func (c *config) dst() (path, error) {
	if len(c.focus) > 0 {
		if !c.focus.valid() {
			return path{}, fmt.Errorf("\n%v is not a valid test case ID.\n", c.focus.String())
		}
		return c.focus, nil
	}
	return path{}, nil
}
