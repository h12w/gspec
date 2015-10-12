package container

import (
	"bytes"
	"fmt"
	"regexp"
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
		"--all=true",
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

var rxPortMapping = regexp.MustCompile(`([0-9]+)/.* -> .*:([0-9]+)`)

func (c *Container) ports() (map[int]int, error) {
	result := make(map[int]int)
	cmd := util.Command("docker", "port", c.ID)
	out := cmd.Output()
	lines := bytes.Split(out, []byte("\n"))
	for _, line := range lines {
		m := rxPortMapping.FindSubmatch(line)
		if len(m) != 3 {
			return nil, fmt.Errorf("unexpected input %s", string(line))
		}
		from, err := strconv.Atoi(string(m[1]))
		if err != nil {
			return nil, err
		}
		to, err := strconv.Atoi(string(m[2]))
		if err != nil {
			return nil, err
		}
		result[from] = to
	}
	return result, nil
}

func dockerStart(id string) error {
	return util.Command("docker", "start", id).Run()
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
