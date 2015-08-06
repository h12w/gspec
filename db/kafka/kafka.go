package kafka

import (
	"errors"
	"math/rand"
	"strconv"
	"time"

	"h12.me/gspec/util"
)

const (
	kafkaTopicsCmd = "kafka-topics.sh"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

type ZooKeeper struct {
	Addr string
}

func (zk *ZooKeeper) NewTopic() (string, error) {
	if !util.CmdExists(kafkaTopicsCmd) {
		return "", errors.New(kafkaTopicsCmd + " not found in path, please install Kafka first")
	}
	topic := "topic_" + strconv.Itoa(rand.Int())
	return topic, util.Command("kafka-topics.sh", "--zookeeper", zk.Addr, "--create", "--topic", topic, "--partitions", "1", "--eplication-factor", "1").Run()
}

func (zk *ZooKeeper) DeleteTopic(topic string) error {
	return util.Command("kafka-topics.sh", "--zookeeper", zk.Addr, "--delete", "--topic", topic).Run()
}
