package mysql

import "testing"

func TestMysql(t *testing.T) {
	s, err := New()
	if err != nil {
		t.Fatal(err)
	}
	defer s.Close()
	if _, err := s.Exec("CREATE DATABASE abc"); err != nil {
		t.Fatal(err)
	}
}
