package main

import (
	"github.com/Shopify/sarama"
	"log"
	"os"
	"os/signal"
	"time"
)

func main() {
	consumer, err := sarama.NewConsumer([]string{"192.168.20.217:9092", "192.168.21.195:9092", "192.168.22.49:9092", "192.168.22.152:9092", "192.168.24.113:9092"}, nil)
	if err != nil {
		panic(err)
	}

	defer func() {
		if err := consumer.Close(); err != nil {
			log.Fatalln(err)
		}
	}()

	partitionConsumer, err := consumer.ConsumePartition("zabbix-metrics-switcher", 0, sarama.OffsetOldest)
	if err != nil {
		panic(err)
	}

	defer func() {
		if err := partitionConsumer.Close(); err != nil {
			log.Fatalln(err)
		}
	}()

	// Trap SIGINT to trigger a shutdown.
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)

	consumed := 0
	start := time.Now()
ConsumerLoop:
	for {
		select {
		case <-partitionConsumer.Messages():
			consumed++
			if consumed%10000 == 0 {
				end := time.Since(start).Seconds()
				log.Printf("%d, %.2f\n", consumed, float64(consumed)/end)
			}
		case <-signals:
			break ConsumerLoop
		}
	}

	log.Printf("Consumed: %d\n", consumed)
}
