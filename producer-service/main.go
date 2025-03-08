package main

import (
	"fmt"
	"io"
	"net/http"

	"github.com/IBM/sarama"
)

const (
	brokerAddress = "localhost:9092"
)

var producer sarama.SyncProducer

func init() {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true

	var err error
	producer, err = sarama.NewSyncProducer([]string{brokerAddress}, config)
	if err != nil {
		panic(err)
	}
}

func publishToKafka(stats []byte) {
	message := &sarama.ProducerMessage{
		Topic: "systemstats",
		Value: sarama.ByteEncoder(stats),
	}

	_, _, err := producer.SendMessage(message)
	if err != nil {
		fmt.Println("Error sending message to Kafka", err)
	} else {
		fmt.Println("Published to Kafka", string(stats))
	}
}

func statsHandler(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	publishToKafka(body)
	w.WriteHeader(http.StatusOK)
}

func main(){
	http.HandleFunc("/api/producer/v1/stats",statsHandler)
	fmt.Println("Producer running on port 5000")
	http.ListenAndServe(":5000",nil)
}
