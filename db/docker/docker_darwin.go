package docker

import (
	"bytes"
	"errors"
	"log"
	"os"
	"strings"
)

func (c Container) ip() (string, error) {
	out := command("boot2docker", "ip").Output()
	return string(bytes.TrimSpace(out)), nil
}

func initDocker() error {
	if !(cmdExists("boot2docker") && cmdExists("docker")) {
		return errors.New("docker not installed")
	}
	if status := string(command("boot2docker", "status").Output()); status != "running" {
		log.Println("boot2docker start ...")
		if err := command("boot2docker", "start").Run(); err != nil {
			return err
		}
	}
	if os.Getenv("DOCKER_HOST") == "" || os.Getenv("DOCKER_CERT_PATH") == "" || os.Getenv("DOCKER_TLS_VERIFY") == "" {
		for _, line := range strings.Split(string(command("boot2docker", "shellinit").Output()), "\n") {
			line = strings.TrimSpace(line)
			if line == "" {
				continue
			}
			line = strings.TrimPrefix(line, "set -x ")
			keyVal := strings.Split(line, " ")
			os.Setenv(keyVal[0], keyVal[1])
		}
	}
	return nil
}
