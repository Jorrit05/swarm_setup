package main

import (
	"encoding/json"
	"fmt"

	"github.com/Jorrit05/GoLib"
	"github.com/docker/docker/client"
	amqp "github.com/rabbitmq/amqp091-go"
)

type CreateServicePayload struct {
	ImageName    string            `json:"image_name"`
	ImageVersion string            `json:"image_version"`
	EnvVars      map[string]string `json:"env_vars"`
	Networks     []string          `json:"networks"`
	Secrets      []string          `json:"secrets"`
	Volumes      map[string]string `json:"volumes"`
	Ports        map[string]string `json:"ports"`
}

type MicroServices struct {
	Services []CreateServicePayload
}

type DetachAttachServicePayload struct {
	ServiceName string `json:"service_name"`
	QueueName   string `json:"queue_name"`
}

type KillServicePayload struct {
	ServiceName string `json:"service_name"`
}

func handleCreateService(cli *client.Client, payload MicroServices) {
	fmt.Println("Handling Create Service")

	for _, microservice := range payload.Services {
		serviceSpec := GoLib.CreateServiceSpec(
			microservice.ImageName,
			microservice.ImageVersion,
			microservice.EnvVars,
			microservice.Networks,
			microservice.Secrets,
			microservice.Volumes,
			microservice.Ports,
			cli,
		)
		GoLib.CreateDockerService(cli, serviceSpec)
	}
}

func handleDetachService(payload DetachAttachServicePayload) {
	fmt.Println("Handling Detach Service")
	// Detach the service from the queue
}

func handleAttachService(payload DetachAttachServicePayload) {
	fmt.Println("Handling Attach Service")
	// Attach the service to the queue
}

func handleKillService(payload KillServicePayload) {
	fmt.Println("Handling Kill Service")
	// Kill the service
}

func startMessageLoop(
	messages <-chan amqp.Delivery,
	exchangeName string,
) {

	for msg := range messages {
		if exchangeName == "" {
			exchangeName = "topic_exchange"
		}

		switch msg.Type {
		case "CreateService":
			var payload MicroServices
			err := json.Unmarshal(msg.Body, &payload)
			if err != nil {
				log.Printf("Error decoding CreateServicePayload: %v", err)
				return
			}
			handleCreateService(dockerClient, payload)
		case "DetachService":
			// Handle DetachService
			// ...

		case "AttachService":
			// Handle AttachService
			// ...

		case "KillService":
			// Handle KillService
			// ...

		default:
			log.Printf("Unknown message type: %s", msg.Type)
		}

		// ... (Acknowledge the message)
	}
}
