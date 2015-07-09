package util

import (
	"fmt"
	"net"
	"time"
)

func AwaitReachable(addr string, maxWait time.Duration) error {
	done := time.Now().Add(maxWait)
	for time.Now().Before(done) {
		c, err := net.Dial("tcp", addr)
		if err == nil {
			c.Close()
			return nil
		}
		time.Sleep(100 * time.Millisecond)
	}
	return fmt.Errorf("%v unreachable for %v", addr, maxWait)
}
