package docker

import (
	"bytes"
	"fmt"
	"log"
	"net"
	"strconv"

	"h12.me/gspec/util"
)

type Container struct {
	ID   string
	Addr *net.TCPAddr
}

func newContainer(id string) (_ *Container, err error) {
	c := &Container{ID: id}
	c.Addr, err = c.addr()
	if err != nil {
		log := c.Log()
		c.Close()
		return nil, fmt.Errorf("%s: %s", err.Error(), log)
	}
	return c, nil
}

func (c *Container) addr() (*net.TCPAddr, error) {
	ip, err := c.ip()
	if err != nil {
		return nil, err
	}
	port, err := c.port()
	if err != nil {
		return nil, err
	}
	return &net.TCPAddr{
		IP:   net.ParseIP(ip),
		Port: port,
	}, nil
}

func (c *Container) port() (int, error) {
	cmd := util.Command("docker", "port", c.ID)
	out := cmd.Output()
	tok := bytes.Split(out, []byte(":"))
	if len(tok) == 2 {
		return strconv.Atoi(string(bytes.TrimSpace(tok[1])))
	}
	return 0, fmt.Errorf("fail to parse port from %s, cmd: %v, id: %s\n", string(out), cmd, c.ID)
}

func (c *Container) Kill() error {
	return util.Command("docker", "kill", c.ID).Run()
}

func (c *Container) Remove() error {
	return util.Command("docker", "rm", c.ID).Run()
}

// KillRemove calls Kill on the container, and then Remove if there was
// no error. It logs any error to t.
func (c *Container) Close() {
	if err := c.Kill(); err != nil {
		log.Println(err)
	}
	if err := c.Remove(); err != nil {
		log.Println(err)
	}
}

func (c Container) Log() string {
	return string(util.Command("docker", "logs", c.ID).Output())
}
