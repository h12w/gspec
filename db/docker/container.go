package docker

import (
	"bytes"
	"fmt"
	"log"
	"net"
	"strconv"
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
	out := command("docker", "port", c.ID).Output()
	tok := bytes.Split(out, []byte(":"))
	if len(tok) == 2 {
		return strconv.Atoi(string(bytes.TrimSpace(tok[1])))
	}
	return 0, fmt.Errorf("fail to parse port from %s", string(out))
}

func (c *Container) Kill() error {
	return command("docker", "kill", c.ID).Run()
}

func (c *Container) Remove() error {
	return command("docker", "rm", c.ID).Run()
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
	return string(command("docker", "logs", c.ID).Output())
}