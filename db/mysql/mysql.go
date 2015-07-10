package mysql

import (
	"database/sql"
	"fmt"
	"math/rand"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
	"h12.me/gspec/db/docker"
)

const containerName = "gspec-db-mysql-f762b7f19a06403cb27bc8ab5f735840"

type Database struct {
	Name    string
	ConnStr string
	DBName  string
	*sql.DB
	c *docker.Container
}

func New() (*Database, error) {
	container, err := docker.Find(containerName)
	if err != nil {
		container, err = docker.New("--name="+containerName, "--detach=true", "--publish=3306:3306", "h12w/mysql:latest")
		if err != nil {
			return nil, err
		}
	}

	password := ""
	connStr := fmt.Sprintf("root:%s@tcp(%s)/", password, container.Addr.String())
	x, err := sql.Open("mysql", connStr)
	if err != nil {
		container.Close()
		return nil, err
	}
	dbName := "db_" + strconv.Itoa(rand.Int())
	if _, err := x.Exec("CREATE DATABASE " + dbName); err != nil {
		return nil, err
	}
	if _, err := x.Exec("USE " + dbName); err != nil {
		return nil, err
	}
	return &Database{
		Name:    dbName,
		ConnStr: connStr + dbName,
		DBName:  dbName,
		DB:      x,
		c:       container,
	}, nil
}

func (s *Database) Close() {
	s.DB.Exec("DROP DATABASE " + s.Name)
	if s.DB != nil {
		s.DB.Close()
		s.DB = nil
	}
}

func (s *Database) Addr() string {
	return s.c.Addr.String()
}
