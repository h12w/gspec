package docker

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

type cmd struct {
	c      *exec.Cmd
	err    error
	errBuf bytes.Buffer
}

func command(name string, arg ...string) *cmd {
	cmd := cmd{c: exec.Command(name, arg...)}
	cmd.c.Stderr = &cmd.errBuf
	return &cmd
}

func (c *cmd) Output() []byte {
	r, err := c.c.Output()
	if err != nil {
		c.err = c.formatError(err)
		return nil
	}
	return bytes.TrimSpace(r)
}

func (c *cmd) Err() error {
	return c.err
}

func (c *cmd) Run() error {
	err := c.c.Run()
	if err != nil {
		c.err = c.formatError(err)
	}
	return nil
}

func (c *cmd) formatError(err error) error {
	return fmt.Errorf("%s (%s): %s", strings.Join(c.c.Args, " "), err.Error(), c.errBuf.String())
}

func (c *cmd) String() string {
	return c.c.Path + " " + strings.Join(c.c.Args, " ")
}

func cmdExists(file string) bool {
	_, err := exec.LookPath(file)
	return err == nil
}
