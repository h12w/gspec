package mysql

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"h12.me/gspec/db/docker"
)

type Session struct {
	ConnStr string
	*sql.DB
	c *docker.Container
}

func New() (*Session, error) {
	container, err := docker.New("--detach=true", "--publish-all=true", "h12w/mariadb:latest")
	if err != nil {
		return nil, err
	}
	time.Sleep(time.Second)
	password := ""
	connStr := fmt.Sprintf("root:%s@tcp(%s)/", password, container.Addr.String())
	x, err := sql.Open("mysql", connStr)
	if err != nil {
		container.Close()
		return nil, err
	}
	return &Session{
		ConnStr: connStr,
		DB:      x,
		c:       container,
	}, nil
}

func (s *Session) Close() {
	if s.DB != nil {
		s.DB.Close()
		s.DB = nil
	}
	s.c.Close()
}

func (s *Session) Addr() string {
	return s.c.Addr.String()
}