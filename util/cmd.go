package util

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

type Cmd struct {
	c      *exec.Cmd
	err    error
	errBuf bytes.Buffer
}

func Command(name string, arg ...string) *Cmd {
	cmd := Cmd{c: exec.Command(name, arg...)}
	cmd.c.Stderr = &cmd.errBuf
	return &cmd
}

func (c *Cmd) Output() []byte {
	r, err := c.c.Output()
	if err != nil {
		c.err = c.formatError(err)
		return nil
	}
	return bytes.TrimSpace(r)
}

func (c *Cmd) Err() error {
	return c.err
}

func (c *Cmd) Run() error {
	err := c.c.Run()
	if err != nil {
		c.err = c.formatError(err)
		return c.err
	}
	return nil
}

func (c *Cmd) formatError(err error) error {
	return fmt.Errorf("%s (%s): %s", strings.Join(c.c.Args, " "), err.Error(), c.errBuf.String())
}

func (c *Cmd) String() string {
	return c.c.Path + " " + strings.Join(c.c.Args, " ")
}

func CmdExists(file string) bool {
	_, err := exec.LookPath(file)
	return err == nil
}
