package docker

import (
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

func run(args []string) (string, error) {
	args = append([]string{"run"}, args...)
	cmd := command("docker", args...)
	containerID := string(cmd.Output())
	if cmd.Err() != nil {
		return "", cmd.Err()
	}
	return containerID, nil
}
