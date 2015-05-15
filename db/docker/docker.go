package docker

import (
	"fmt"
	"log"
	"net"
	"strings"
	"time"
)

func New(image string, args ...string) (*Container, error) {
	if err := initImage(image); err != nil {
		return nil, err
	}
	id, err := run(image, args)
	if err != nil {
		return nil, err
	}
	return newContainer(id)
}

func initImage(image string) error {
	if err := initDocker(); err != nil {
		return err
	}
	if ok, err := haveImage(image); !ok || err != nil {
		if err != nil {
			return err
		}
		if err := pull(image); err != nil {
			return err
		}
	}
	return nil
}

func haveImage(name string) (bool, error) {
	cmd := command("docker", "images", "--no-trunc")
	out := cmd.Output()
	if cmd.Err() != nil {
		return false, cmd.Err()
	}
	return strings.Contains(out, name), nil
}

func run(image string, args []string) (string, error) {
	args = append([]string{"run"}, args...)
	args = append(args, image)
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
