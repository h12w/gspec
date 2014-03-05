package core

import (
	"flag"
)

var globalConfig config

func init() {
	flag.Var(&globalConfig.focus, "focus", "test case id to select one test case to run")
}

type config struct {
	focus path
}
