package docker

import (
	"bytes"
	"os/exec"
)

func (c Container) ip() (string, error) {
	out, err := exec.Command("boot2docker", "ip").Output()
	if err != nil {
		return "", err
	}
	return string(bytes.TrimSpace(out)), nil
}

func haveDocker() bool {
	_, err := exec.LookPath("docker")
	if err != nil {
		return false
	}
	_, err = exec.LookPath("boot2docker")
	return err == nil
}
