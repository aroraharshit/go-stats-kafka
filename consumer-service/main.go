package main

import (
	"fmt"
	"log"

	"github.com/IBM/sarama"
)

const (
	brokerAddress = "localhost:9092"
	topic         = "systemstats"
)


func main() {
	brokers := []string{"localhost:9092"}
	config := sarama.NewConfig()
	config.Consumer.Offsets.Initial = sarama.OffsetNewest

	consumer, err := sarama.NewConsumer(brokers, config)
	if err != nil {
		log.Fatalf("Error creating Kafka consumer %v", err)
	}

	defer consumer.Close()

	partitionConsumer, err := consumer.ConsumePartition("systemstats", 0, sarama.OffsetNewest)
	if err != nil {
		log.Fatalf("Error creating partition consumer %v", err)
	}

	defer partitionConsumer.Close()

	fmt.Println("Consumer started...")

	for msg := range partitionConsumer.Messages() {
		fmt.Printf("Recieved from Kafka %s\n", string(msg.Value))
	}
}
