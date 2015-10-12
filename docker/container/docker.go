package container

import (
	"bytes"
	"fmt"
	"strconv"

	"h12.me/gspec/util"
)

func dockerRun(args []string) (string, error) {
	args = append([]string{"run"}, args...)
	cmd := util.Command("docker", args...)
	containerID := string(cmd.Output())
	if cmd.Err() != nil {
		return "", cmd.Err()
	}
	return containerID, nil
}

func dockerPS(filter string) (string, error) {
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

func (c Container) Log() string {
	return string(util.Command("docker", "logs", c.ID).Output())
}
