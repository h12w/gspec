package util

import (
	"fmt"
	"net"
	"time"
)

func AwaitReachable(addr *net.TCPAddr, maxWait time.Duration) error {
	done := time.Now().Add(maxWait)
	for time.Now().Before(done) {
		c, err := net.DialTCP("tcp", nil, addr)
		if err == nil {
			c.Close()
			return nil
		}
		time.Sleep(100 * time.Millisecond)
	}
	return fmt.Errorf("%v unreachable for %v", addr, maxWait)
}
