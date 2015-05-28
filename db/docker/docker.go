package docker

import (
	"fmt"
	"net"
	"time"
)

func New(args ...string) (*Container, error) {
	if err := initDocker(); err != nil {
		return nil, err
	}
	id, err := run(args)
	if err != nil {
		return nil, err
	}
	return newContainer(id)
}

func run(args []string) (string, error) {
	args = append([]string{"run"}, args...)
	cmd := command("docker", args...)
	containerID := string(cmd.Output())
	if cmd.Err() != nil {
		return "", cmd.Err()
	}
	return containerID, nil
}

func awaitReachable(addr string, maxWait time.Duration) error {
	done := time.Now().Add(maxWait)
	for time.Now().Before(done) {
		c, err := net.Dial("tcp", addr)
		if err == nil {
			c.Close()
			return nil
		}
		time.Sleep(100 * time.Millisecond)
	}
	return fmt.Errorf("%v unreachable for %v", addr, maxWait)
}
