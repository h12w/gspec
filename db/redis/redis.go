package redis

import (
	"time"

	"github.com/garyburd/redigo/redis"
	"h12.me/gspec/db/docker"
)

const containerName = "gspec-db-redis-ac3bfb841b3c47378dfdecca51b23042"

type Redis struct {
	pool *redis.Pool
	c    *docker.Container
}

func New() (*Redis, error) {
	container, err := docker.Find(containerName)
	if err != nil {
		container, err = docker.New("--name="+containerName, "--detach=true", "--publish=8379:6379", "redis")
		if err != nil {
			return nil, err
		}
	}
	return &Redis{
		pool: &redis.Pool{
			MaxIdle:     3,
			IdleTimeout: 240 * time.Second,
			Dial: func() (redis.Conn, error) {
				c, err := redis.Dial("tcp", container.Addr.String())
				if err != nil {
					return nil, err
				}
				return c, err
			},
			TestOnBorrow: func(c redis.Conn, t time.Time) error {
				_, err := c.Do("PING")
				return err
			},
		},
		c: container,
	}, nil
}

func (s *Redis) Close() {
}

func (s *Redis) Addr() string {
	return s.c.Addr.String()
}
