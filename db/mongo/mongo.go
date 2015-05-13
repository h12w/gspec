package mongo

import (
	"time"

	"h12.me/gspec/db/docker"
)

func New(t docker.T) (*docker.Container, error) {
	return docker.New(t, "mongo", 10*time.Second, func() (string, error) {
		return docker.Run("-d", "-P", "mongo")
	})
}
