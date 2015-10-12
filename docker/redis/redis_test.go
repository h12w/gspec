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
	fmt.Println(s.Addr())
}
