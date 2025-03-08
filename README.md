This project demonstrates a real-time messaging microservices architecture in Go using Kafka. It consists of three microservices:
1) API Service (api/main.go) - Provides an endpoint to send messages to the producer.
2)Producer Service (producer/main.go) - Receives messages from the API service and publishes them to a Kafka topic.
3)Consumer Service (consumer/main.go) - Listens to Kafka and processes incoming messages.

Architecture
(API Service) → (Producer Service) → Kafka → (Consumer Service)

-)The API Service provides an HTTP endpoint (/publish) to send messages.
-)The Producer Service forwards the message to Kafka (messages topic).
-)The Consumer Service listens to Kafka in real-time and processes messages.

Technologies Used
-)Go (Golang)
-)Kafka
-)IBM Sarama (Kafka Go Client)
