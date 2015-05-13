package mongo

import (
	"time"

	"gopkg.in/mgo.v2"
	"h12.me/gspec/db/docker"
)

type Session struct {
	mgoSession
	c *docker.Container
}
type mgoSession struct {
	*mgo.Session
}

func (s *Session) Close() {
	s.Session.Close()
	s.c.Close()
}

func New(t docker.T) (*Session, error) {
	container, err := docker.New(t, "mongo", 10*time.Second, func() (string, error) {
		return docker.Run("-d", "-P", "mongo")
	})
	if err != nil {
		return nil, err
	}

	session, err := mgo.Dial(container.Addr.String())
	if err != nil {
		return nil, err
	}
	return &Session{mgoSession{session}, container}, nil
}
