package kafka

import (
	"fmt"
	"testing"
)

func TestIt(t *testing.T) {
	zk := ZooKeeper{"192.168.59.103:2181"}
	topic, err := zk.NewTopic(2)
	fmt.Println(topic)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(zk.DescribeTopic(topic))
	zk.DeleteTopic(topic)
}
