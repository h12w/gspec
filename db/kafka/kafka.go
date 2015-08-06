package kafka

import (
	"math/rand"
	"strconv"

	"h12.me/gspec/util"
)

type ZooKeeper struct {
	Addr string
}

func (zk *ZooKeeper) NewTopic() (string, error) {
	topic := "topic_" + strconv.Itoa(rand.Int())
	return topic, util.Command("kafka-topics.sh", "--zookeeper", zk.Addr, "--create", "--topic", topic, "--partitions", "1", "--eplication-factor", "1").Run()
}

func (zk *ZooKeeper) DeleteTopic(topic string) error {
	return util.Command("kafka-topics.sh", "--zookeeper", zk.Addr, "--delete", "--topic", topic).Run()
}
