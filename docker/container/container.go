package container

import (
	"fmt"
	"log"
	"net"
	"time"

	"h12.me/gspec/util"
)

type Container struct {
	ID   string
	Addr *net.TCPAddr
}

func New(args ...string) (*Container, error) {
	if err := initDocker(); err != nil {
		return nil, err
	}
	id, err := dockerRun(args)
	if err != nil {
		return nil, err
	}
	c, err := newContainer(id)
	if err != nil {
		return nil, err
	}
	if err := util.AwaitReachable(c.Addr.String(), 30*time.Second); err != nil {
		c.Close()
		return nil, err
	}
	return c, nil
}

func Find(name string) (*Container, error) {
	if err := initDocker(); err != nil {
		return nil, err
	}
	id, err := dockerPS("name=" + name)
	if err != nil {
		return nil, err
	}
	return newContainer(id)
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
