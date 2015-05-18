package mongo

import (
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

func New() (*Session, error) {
	container, err := docker.New("-d", "-P", "mongo:latest")
	if err != nil {
		return nil, err
	}
	session, err := mgo.Dial(container.Addr.String())
	if err != nil {
		return nil, err
	}
	return &Session{mgoSession{session}, container}, nil
}

func (s *Session) Copy() *mgo.Session {
	return s.Session.Copy()
}

func (s *Session) Close() {
	s.Session.Close()
	s.c.Close()
}

func (s *Session) Addr() string {
	return s.c.Addr.String()
}
