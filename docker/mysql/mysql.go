package mysql

import (
	"database/sql"
	"fmt"
	"math/rand"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
	"h12.me/gspec/docker/container"
)

const containerName = "gspec-db-mysql-f762b7f19a06403cb27bc8ab5f735840"
const password = "1234"

type Database struct {
	DBName  string
	ConnStr string
	*sql.DB
	c *container.Container
}

func New() (*Database, error) {
	c, err := container.Find(containerName)
	if err != nil {
		c, err = container.New("--name="+containerName, "--detach=true", "--publish=3306:3306", "--env=MYSQL_ROOT_PASSWORD="+password, "mysql:latest")
		if err != nil {
			return nil, err
		}
	}

	connStr := fmt.Sprintf("root:%s@tcp(%s)/", password, c.Addr.String())
	x, err := sql.Open("mysql", connStr)
	if err != nil {
		c.Close()
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
		ConnStr: connStr + dbName,
		DBName:  dbName,
		DB:      x,
		c:       c,
	}, nil
}

func (s *Database) Close() {
	s.DB.Exec("DROP DATABASE " + s.DBName)
	if s.DB != nil {
		s.DB.Close()
		s.DB = nil
	}
}

func (s *Database) Addr() string {
	return s.c.Addr.String()
}
