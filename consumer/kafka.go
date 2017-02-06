package consumer

import (
	"log"
	"os"

	"github.com/wvanbergen/kafka/consumergroup"
	"github.com/wvanbergen/kazoo-go"
	"github.com/Shopify/sarama"

	"github.com/xunyu/config"
)

func StartConsumer(cfg config.Config) <-chan *sarama.ConsumerMessage {
	out := make(chan *sarama.ConsumerMessage, 0)
	var zookeeperNodes []string
	consumerCfg := cfg.ConsumerConfig
	cgConfig := consumergroup.NewConfig()
	zookeeperNodes, cgConfig.Zookeeper.Chroot = kazoo.ParseConnectionString(consumerCfg.Zookeeper)
	consumer, consumerErr := consumergroup.JoinConsumerGroup(consumerCfg.GroupId, consumerCfg.Topics, zookeeperNodes, cgConfig)

	go func() {
		if consumerErr != nil {
			panic(consumerErr)
		}

		defer func() {
			if err := consumer.Close(); err != nil {
				log.Fatalln(err)
			}
		}()

		for {
			select {
			case msg := <-consumer.Messages():
				out <- msg
				consumer.CommitUpto(msg)
			}
		}
	}()
}
