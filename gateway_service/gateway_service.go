package main

import (
	"lib/queues"
	"log"
)

func main() {
	conn, err := queues.Connect()
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	defer conn.Close()

	inputQueue, err := queues.DeclareQueue("input_queue")
	if err != nil {
		log.Fatalf("Failed to declare input queue: %v", err)
	}

	messages, err := queues.Consume(inputQueue.Name)
	if err != nil {
		log.Fatalf("Failed to register consumer: %v", err)
	}

	for msg := range messages {
		log.Printf("Received message: %v", string(msg.Body))
	}
}
