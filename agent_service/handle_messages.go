package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/Jorrit05/GoLib"
	"github.com/docker/docker/client"
	amqp "github.com/rabbitmq/amqp091-go"
)

func handleCreateService(cli *client.Client, payload []GoLib.CreateServicePayload) {

	// Get some values from etcd
	for _, microservice := range payload {
		serviceSpec := GoLib.CreateServiceSpec(
			microservice.ImageName,
			microservice.Tag,
			microservice.EnvVars,
			microservice.Networks,
			microservice.Secrets,
			microservice.Volumes,
			microservice.Ports,
			cli,
		)
		GoLib.CreateDockerService(cli, serviceSpec)

		jsonData, err := json.Marshal(microservice)
		if err != nil {
			log.Warn("Error marshaling payload to JSON:", err)
		}

		_, err = etcdClient.Put(context.Background(), fmt.Sprintf("%s/%s", msEtcdPath, microservice.ImageName), string(jsonData))
		if err != nil {
			log.Fatalf("Failed creating an item in etcd: %s", err)
		}
	}
}

func handleDetachService(payload GoLib.DetachAttachServicePayload) {
	fmt.Println("Handling Detach Service")
	// Detach the service from the queue
}

func handleAttachService(payload GoLib.DetachAttachServicePayload) {
	fmt.Println("Handling Attach Service")
	// Attach the service to the queue
}

func handleKillService(payload GoLib.KillServicePayload) {
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
		case "datarequest":

		case "CreateService":
			var payload []GoLib.CreateServicePayload
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
