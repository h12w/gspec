package docker

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"net"
	"os/exec"
	"strings"
	"time"
)

type T interface {
	Fatalf(format string, args ...interface{})
}

func New(t T, image string, timeout time.Duration, start func() (string, error)) (*Container, error) {
	checkStatus(t, image)
	id, err := start()
	if err != nil {
		return nil, err
	}
	return newContainer(id)
}

func checkStatus(t T, image string) {
	if !haveDocker() {
		t.Fatalf("'docker' command not found")
	}
	if ok, err := haveImage(image); !ok || err != nil {
		if err != nil {
			t.Fatalf("Error running docker to check for %s: %v", image, err)
		}
		log.Printf("Pulling docker image %s ...", image)
		if err := Pull(image); err != nil {
			t.Fatalf("Error pulling %s: %v", image, err)
		}
	}
}

func haveImage(name string) (ok bool, err error) {
	out, err := exec.Command("docker", "images", "--no-trunc").Output()
	if err != nil {
		log.Println(err)
		return
	}
	return bytes.Contains(out, []byte(name)), nil
}

func Run(args ...string) (containerID string, err error) {
	cmd := exec.Command("docker", append([]string{"run"}, args...)...)
	var stdout, stderr bytes.Buffer
	cmd.Stdout, cmd.Stderr = &stdout, &stderr
	if err = cmd.Run(); err != nil {
		err = fmt.Errorf("%v%v", stderr.String(), err)
		return
	}
	containerID = strings.TrimSpace(stdout.String())
	if containerID == "" {
		return "", errors.New("unexpected empty output from `docker run`")
	}
	return
}

func Pull(image string) error {
	out, err := exec.Command("docker", "pull", image).CombinedOutput()
	if err != nil {
		err = fmt.Errorf("%v: %s", err, out)
	}
	return err
}

func AwaitReachable(addr string, maxWait time.Duration) error {
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
