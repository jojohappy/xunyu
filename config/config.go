package config

import (
)

type ElasticConfig struct {
	Host string `json:"host"`
}

type KafkaConfig struct {
	Broker string `json:"broker"`
}

type WebsocketConfig struct {
	Uri string `json:"uri"`
}

type ZabbixConfig struct {
	Url string `json:"url"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type EtcdConfig struct {
	Host string `json:"host"`
}

type ConsumerConfig struct {
	GroupId string `json:"groupId"`
	Topics []string `json:"Topics"`
	Zookeeper string `json:"zookeeper"`
}

type Config struct {
	Elastic *ElasticConfig `json:"elastic"`
	Kafka *KafkaConfig `json:"kafka"`
	Websocket *WebsocketConfig `json:"websocket"`
	Zabbix *ZabbixConfig `json:"zabbix"`
	Etcd *EtcdConfig `json:"etcd"`
	Consumers *ConsumerConfig `json:"consumers"`
}

func initConfig() {

}
