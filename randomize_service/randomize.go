package main

import (
	"encoding/json"
	"log"

	"github.com/Jorrit05/GoLib"
	amqp "github.com/rabbitmq/amqp091-go"
)

var (
	serviceName string = "randomize_service"
	routingKey  string = GoLib.GetDefaultRoutingKey(serviceName)
)

func main() {
	// Log to file
	f := GoLib.StartLog()
	defer f.Close()
	log.SetOutput(f)

	// Connect to AMQ queue, declare own routingKey as queue, start listening for messages
	messages, conn, channel, err := GoLib.SetupConnection(serviceName, routingKey, true)
	if err != nil {
		log.Fatalf("Failed to setup proper connection to RabbitMQ: %v", err)
	}
	defer conn.Close()

	GoLib.StartMessageLoop(anonymize, messages, channel, routingKey, "")
}

type SkillQuery struct {
	PersonId  int    `json:"person_id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Sex       string `json:"sex"`
	// DriversLicense string `json:"drivers_license"`
	// Programming    string `json:"programming"`
}

func anonymize(message amqp.Delivery) (amqp.Publishing, error) {
	var skillQueries []SkillQuery

	err := json.Unmarshal(message.Body, &skillQueries)
	if err != nil {
		log.Printf("Error unmarshaling JSON:", err)
		return amqp.Publishing{}, err
	}

	// Anonymise last name
	for i := range skillQueries {
		skillQueries[i].LastName = "anonymized"
	}

	jsonMessage, err := json.Marshal(skillQueries)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		return amqp.Publishing{}, err
	}

	return amqp.Publishing{
		Body:          jsonMessage,
		Type:          "application/json",
		CorrelationId: message.CorrelationId,
	}, nil
}
