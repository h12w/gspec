package mongo

import (
	"math/rand"
	"strconv"

	"gopkg.in/mgo.v2"
	"h12.me/gspec/docker/container"
)

const (
	containerName = "gspec-db-mongo-79cb399e9230494cb475d8461a0183c7"
	internalPort  = 27017
)

type Mongo struct {
	DBName  string
	ConnStr string
	*mgo.Session
	*mgo.Database
	c *container.Container
}

func New() (*Mongo, error) {
	c, err := container.Find(containerName)
	if err != nil {
		c, err = container.New("--name="+containerName, "--detach=true", "--publish=27017:27017", "mongo:latest")
		if err != nil {
			return nil, err
		}
	}
	connStr := "mongodb://" + c.Addr(internalPort)
	session, err := mgo.Dial(connStr)
	if err != nil {
		c.Close()
		return nil, err
	}
	dbName := "db_" + strconv.Itoa(rand.Int())
	db := session.DB(dbName)
	return &Mongo{
		DBName:   dbName,
		ConnStr:  connStr,
		Session:  session,
		Database: db,
		c:        c,
	}, nil
}

func (s *Mongo) Close() {
	if s.Database != nil {
		s.Database.DropDatabase()
		s.Database = nil
	}
	if s.Session != nil {
		s.Session.Close()
		s.Session = nil
	}
}

func (s *Mongo) Addr() string {
	return s.c.Addr(internalPort)
}
