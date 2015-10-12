package container

import (
	"bytes"
	"encoding/json"
	"errors"
)

func (c *Container) ip() (string, error) {
	out, err := exec.Command("docker", "inspect", c.ID).Output()
	if err != nil {
		return "", err
	}
	type networkSettings struct {
		IPAddress string
	}
	type container struct {
		NetworkSettings networkSettings
	}
	var cs []container
	if err := json.NewDecoder(bytes.NewReader(out)).Decode(&cs); err != nil {
		return "", err
	}
	if len(cs) == 0 {
		return "", errors.New("no output from docker inspect")
	}
	if ip := cs[0].NetworkSettings.IPAddress; ip != "" {
		return ip, nil
	}
	return "", errors.New("could not find an IP. Not running?")
}

func initDocker() error {
	if !(cmdExists("boot2docker") && cmdExists("docker")) {
		return errors.New("docker not installed")
	}
	return nil
}
