package main

import (
	"log"

	"github.com/Jorrit05/GoLib"
)

func main() {
	conn, err := GoLib.Connect()
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	defer conn.Close()

	inputQueue, err := GoLib.DeclareQueue("input_queue")
	if err != nil {
		log.Fatalf("Failed to declare input queue: %v", err)
	}

	messages, err := GoLib.Consume(inputQueue.Name)
	if err != nil {
		log.Fatalf("Failed to register consumer: %v", err)
	}

	for msg := range messages {
		log.Printf("Received message: %v", string(msg.Body))
	}
}
