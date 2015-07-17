package redis

import (
	"fmt"
	"testing"
)

func TestMysql(t *testing.T) {
	s, err := New()
	if err != nil {
		t.Fatal(err)
	}
	defer s.Close()
	fmt.Println(s.Addr())
}
