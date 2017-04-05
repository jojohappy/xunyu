package kafka

import (
	"strings"

	"github.com/xunyu/common"
	"github.com/xunyu/config"
	"github.com/xunyu/lib/log"

	"github.com/Shopify/sarama"
	"github.com/wvanbergen/kafka/consumergroup"
	"github.com/wvanbergen/kazoo-go"
)

type KafkaConfig struct {
	GroupId   string `config:"groupId"`
	Topics    string `config:"topics"`
	Zookeeper string `config:"zookeeper"`
	Initial   int64  `config:"initial"`
}

type kafka struct {
	common.PluginPrototype
	done   chan struct{}
	config KafkaConfig
}

var (
	defaultConfig = KafkaConfig{
		GroupId: "xunyu-kafka",
		Initial: sarama.OffsetOldest,
	}
)

func init() {
	common.RegisterInputPlugin("kafka", New)
}

func New(config *config.Config) (common.Pluginer, error) {
	k := &kafka{
		config: defaultConfig,
		done:   make(chan struct{}),
	}
	if err := k.init(config); nil != err {
		return nil, err
	}
	return k, nil
}

func (k *kafka) init(config *config.Config) error {
	if err := config.Assemble(&k.config); nil != err {
		return err
	}

	log.Debug("config of inputs kafka is %v", k.config)

	if _, err := k.newKafkaConfig(); nil != err {
		return err
	}

	return nil
}

func (k *kafka) newKafkaConfig() (*consumergroup.Config, error) {
	cfg := consumergroup.NewConfig()
	if k.config.Initial != 0 {
		cfg.Offsets.Initial = k.config.Initial
	}

	return cfg, nil
}

func newKafkaClient(cfg *consumergroup.Config, zookeeper string, topics string, groupId string) (*consumergroup.ConsumerGroup, error) {
	var zookeepers []string
	zookeepers, cfg.Zookeeper.Chroot = kazoo.ParseConnectionString(zookeeper)
	tp := strings.Split(topics, ",")
	return consumergroup.JoinConsumerGroup(groupId, tp, zookeepers, cfg)
}

func (k *kafka) Start() <-chan common.DataInter {
	out := make(chan common.DataInter, 1)
	cfg, err := k.newKafkaConfig()

	if nil != err {
		log.Error("error on creating config of inputs kafka: %s", err)
		close(out)
		return out
	}

	consumer, err := newKafkaClient(cfg, k.config.Zookeeper, k.config.Topics, k.config.GroupId)

	if nil != err {
		log.Error("error on creating kafka client: %s", err)
		close(out)
		return out
	}

	go func() {
		defer close(out)
		for {
			select {
			case msg := <-consumer.Messages():
				out <- string(msg.Value)
				consumer.CommitUpto(msg)
			case <-k.done:
				if err := consumer.Close(); err != nil {
					log.Error("error closing the consumer %s", err)
				}
				return
			}
		}
	}()
	return out
}

func (k *kafka) Close() {
	close(k.done)
}
