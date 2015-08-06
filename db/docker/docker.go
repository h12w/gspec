package docker

import (
	"fmt"
	"time"

	"h12.me/gspec/util"
)

func New(args ...string) (*Container, error) {
	if err := initDocker(); err != nil {
		return nil, err
	}
	id, err := run(args)
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
	id, err := ps("name=" + name)
	if err != nil {
		return nil, err
	}
	return newContainer(id)
}

func run(args []string) (string, error) {
	args = append([]string{"run"}, args...)
	cmd := util.Command("docker", args...)
	containerID := string(cmd.Output())
	if cmd.Err() != nil {
		return "", cmd.Err()
	}
	return containerID, nil
}

func ps(filter string) (string, error) {
	args := []string{
		"ps",
		"--no-trunc=true",
		"--quiet=true",
		"--filter=" + filter,
	}
	cmd := util.Command("docker", args...)
	containerID := string(cmd.Output())
	if cmd.Err() != nil {
		return "", cmd.Err()
	}
	if containerID == "" {
		return "", fmt.Errorf("cannot find container with condition %s", filter)
	}
	return containerID, nil
}
