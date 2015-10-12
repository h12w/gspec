package container

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"h12.me/gspec/util"
)

type Container struct {
	ID    string
	IP    string
	Ports map[int]int
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
	if err := util.AwaitReachable(c.anyAddr(), 30*time.Second); err != nil {
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
	if err = dockerStart(id); err != nil {
		return nil, err
	}
	c, err := newContainer(id)
	if err != nil {
		return nil, err
	}
	return c, nil
}

func newContainer(id string) (_ *Container, err error) {
	c := &Container{ID: id}
	c.IP, err = c.ip()
	if err != nil {
		log := c.Log()
		c.Close()
		return nil, fmt.Errorf("%s: %s", err.Error(), log)
	}
	c.Ports, err = c.ports()
	if err != nil {
		log := c.Log()
		c.Close()
		return nil, fmt.Errorf("%s: %s", err.Error(), log)
	}
	return c, nil
}

func (c *Container) anyAddr() string {
	for _, port := range c.Ports {
		return c.IP + ":" + strconv.Itoa(port)
	}
	return ""
}

func (c *Container) Addr(port int) string {
	return c.IP + ":" + strconv.Itoa(c.Ports[port])
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
