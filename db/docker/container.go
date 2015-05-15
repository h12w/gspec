package docker

import (
	"bytes"
	"fmt"
	"log"
	"net"
	"os/exec"
	"strconv"
	"time"
)

type Container struct {
	ID   string
	Addr *net.TCPAddr
}

func newContainer(id string) (*Container, error) {
	c := &Container{ID: id}
	var err error
	c.Addr, err = c.addr()
	if err != nil {
		return nil, err
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
	out, err := exec.Command("docker", "port", c.ID).Output()
	if err != nil {
		return 0, err
	}
	tok := bytes.Split(out, []byte(":"))
	if len(tok) == 2 {
		return strconv.Atoi(string(bytes.TrimSpace(tok[1])))
	}
	return 0, fmt.Errorf("fail to parse port from %s", string(out))
}

func (c *Container) Kill() error {
	return exec.Command("docker", "kill", c.ID).Run()
}

func (c *Container) Remove() error {
	return exec.Command("docker", "rm", c.ID).Run()
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

// lookup retrieves the ip address of the container, and tries to reach
// before timeout the tcp address at this ip and given port.
func (c *Container) lookup(port int, timeout time.Duration) (ip string, err error) {
	ip, err = c.ip()
	if err != nil {
		err = fmt.Errorf("error getting IP: %v", err)
		return
	}
	addr := fmt.Sprintf("%s:%d", ip, port)
	err = awaitReachable(addr, timeout)
	return
}
