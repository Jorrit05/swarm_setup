package main

import (
	"context"
	"encoding/json"
	"log"
	"os"

	"github.com/Jorrit05/GoLib"
	amqp "github.com/rabbitmq/amqp091-go"
)

var routingKey string = os.Getenv("ROUTING_KEY")

func startLog() *os.File {
	f, err := os.OpenFile("log.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	return f
}

func main() {
	f := startLog()
	defer f.Close()
	log.SetOutput(f)

	// Connect to AMQ queue, declare own routingKey
	conn, channel, err := GoLib.SetupConnection("randomize_service", routingKey)
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	defer conn.Close()

	// Start listening to queue defined by environment var INPUT_QUEUE
	messages, err := GoLib.Consume(os.Getenv("INPUT_QUEUE"))
	if err != nil {
		log.Fatalf("Failed to register consumer: %v", err)
	}

	// Message loop stays alive
	for msg := range messages {
		log.Printf("Received message: %v", string(msg.Body))
		anonymizedMsg := anonymize(msg)

		err := channel.PublishWithContext(context.Background(), "topic_exchange", routingKey, false, false, anonymizedMsg)
		if err != nil {
			log.Fatalf("Error publishing message: %v", err)
		}
	}
}

type SkillQuery struct {
	PersonId  int    `json:"person_id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Sex       string `json:"sex"`
	// DriversLicense string `json:"drivers_license"`
	// Programming    string `json:"programming"`
}

func anonymize(message amqp.Delivery) amqp.Publishing {
	var skillQueries []SkillQuery

	err := json.Unmarshal(message.Body, &skillQueries)
	if err != nil {
		log.Fatalf("Error unmarshaling JSON:", err)
	}

	// Anonymise last name
	for i := range skillQueries {
		skillQueries[i].LastName = "anonymized"
	}

	jsonMessage, _ := json.Marshal(skillQueries)

	return amqp.Publishing{
		Body: jsonMessage,
		Type: "text/json",
	}
}
