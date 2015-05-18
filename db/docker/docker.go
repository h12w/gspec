package docker

import (
	"fmt"
	"log"
	"net"
	"strings"
	"time"
)

func New(image string, args ...string) (*Container, error) {
	if err := initDocker(); err != nil {
		return nil, err
	}
	id, err := run(image, args)
	if err != nil {
		return nil, err
	}
	return newContainer(id)
}

func run(image string, args []string) (string, error) {
	args = append([]string{"run"}, args...)
	cmd := command("docker", args...)
	containerID := strings.TrimSpace(cmd.Output())
	if cmd.Err() != nil {
		return "", cmd.Err()
	}
	return containerID, nil
}

func pull(image string) error {
	log.Printf("docker pull %s ...\n", image)
	cmd := command("docker", "pull", image)
	if err := cmd.Run(); err != nil {
		return err
	}
	return nil
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
