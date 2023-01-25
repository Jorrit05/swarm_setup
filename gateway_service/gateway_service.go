package main

import (
	"encoding/json"
	"log"

	"github.com/Jorrit05/GoLib"
)

func mongo_test() {
	// mongo_test()
	// conn, err := GoLib.Connect("amqp://guest:guest@localhost:5672/")
	// if err != nil {
	// 	log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	// }
	// defer conn.Close()

	// inputQueue, err := GoLib.DeclareQueue("input_queue")
	// if err != nil {
	// 	log.Fatalf("Failed to declare input queue: %v", err)
	// }

	// messages, err := GoLib.Consume(inputQueue.Name)
	// if err != nil {
	// 	log.Fatalf("Failed to register consumer: %v", err)
	// }

	// for msg := range messages {
	// 	log.Printf("Received message: %v", string(msg.Body))
	// 	if msg.Priority == 3 {
	// 		publish(inputQueue.Name)
	// 	}

	// 	msg.Headers.Validate()
	// }
}

type ExampleMessage struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func publish(queueName string) {
	// conn, err := GoLib.Connect()
	// if err != nil {
	// 	log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	// }
	// defer conn.Close()

	// outputQueue, err := GoLib.DeclareQueue("output_queue")
	// if err != nil {
	// 	log.Fatalf("Failed to declare output queue: %v", err)
	// }

	exampleMessage := ExampleMessage{Name: "John Smith", Age: 30}
	jsonMessage, err := json.Marshal(exampleMessage)
	if err != nil {
		log.Fatalf("Failed to marshal message: %v", err)
	}

	err = GoLib.Publish(queueName, jsonMessage)
	if err != nil {
		log.Fatalf("Failed to publish message: %v", err)
	}
	log.Printf("Sent message: %v", string(jsonMessage))
}
