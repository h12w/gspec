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

func (zk *ZooKeeper) NewTopic(partition int) (string, error) {
	if !util.CmdExists(kafkaTopicsCmd) {
		return "", errors.New(kafkaTopicsCmd + " not found in path, please install Kafka first")
	}
	topic := "topic_" + strconv.Itoa(rand.Int())
	return topic, util.Command(kafkaTopicsCmd, "--zookeeper", zk.Addr, "--create", "--topic", topic, "--partitions", strconv.Itoa(partition), "--replication-factor", "1").Run()
}

func (zk *ZooKeeper) DeleteTopic(topic string) error {
	return util.Command(kafkaTopicsCmd, "--zookeeper", zk.Addr, "--delete", "--topic", topic).Run()
}

func (zk *ZooKeeper) DescribeTopic(topic string) string {
	return string(util.Command(kafkaTopicsCmd, "--zookeeper", zk.Addr, "--describe", "--topic", topic).Output())
}
