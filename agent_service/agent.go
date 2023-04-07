package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sync"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/swarm"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/namesgenerator"

	"github.com/Jorrit05/GoLib"
	amqp "github.com/rabbitmq/amqp091-go"
)

var (
	serviceName   string = "gateway_service"
	routingKey    string = GoLib.GetDefaultRoutingKey(serviceName)
	outputChannel *amqp.Channel
	mutex         = &sync.Mutex{}
	requestMap    = make(map[string]*requestInfo)
)

type requestInfo struct {
	id       string
	response chan amqp.Delivery
}

func main() {
	// Log to file
	f := GoLib.StartLog()
	defer f.Close()
	log.SetOutput(f)

	// Connect to AMQ queue, declare own routingKey as queue
	messages, conn, channel, err := GoLib.SetupConnection(serviceName, routingKey, true)
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	defer conn.Close()

	go register()

	go func() {
		// Start listening for messages, this method keeps this method 'alive'
		GoLib.StartMessageLoop(placeholder, messages, channel, routingKey, "")
	}()

	// Create a new Docker client
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		log.Fatalf("Error creating Docker client: %v", err)
	}

	// Check if Swarm is active
	info, err := cli.Info(context.Background())
	if err != nil {
		log.Fatalf("Error getting Docker info: %v", err)
	}
	if !info.Swarm.ControlAvailable {
		log.Fatal("This node is not a swarm manager. The agent can only be run on a swarm manager.")
	}

	// Example service specification
	spec := swarm.ServiceSpec{
		Annotations: swarm.Annotations{
			Name: namesgenerator.GetRandomName(0), // Generate a random name for the service
		},
		TaskTemplate: swarm.TaskSpec{
			ContainerSpec: &swarm.ContainerSpec{
				Image: "nginx:latest",
			},
		},
		EndpointSpec: &swarm.EndpointSpec{
			Ports: []swarm.PortConfig{
				{
					Protocol:      swarm.PortConfigProtocolTCP,
					PublishedPort: 80,
					TargetPort:    80,
				},
			},
		},
	}

	// Create the service
	response, err := cli.ServiceCreate(context.Background(), spec, types.ServiceCreateOptions{})
	if err != nil {
		log.Fatalf("Error creating service: %v", err)
	}

	// Print the service ID
	fmt.Printf("Service created with ID: %s\n", response.ID)

	// Clean up and close the client
	if err := cli.Close(); err != nil {
		log.Fatalf("Error closing Docker client: %v", err)
	}
}

func placeholder(message amqp.Delivery) (amqp.Publishing, error) {
	var employeeSalary map[string]int
	employeeSalary["me"] = 2000

	jsonMessage, err := json.Marshal(employeeSalary)
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
