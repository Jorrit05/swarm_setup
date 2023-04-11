package main

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/Jorrit05/GoLib"
	"github.com/docker/docker/client"
	amqp "github.com/rabbitmq/amqp091-go"
	clientv3 "go.etcd.io/etcd/client/v3"
)

var (
	serviceName         string
	log, logFile                       = GoLib.InitLogger(serviceName)
	dockerClient        *client.Client = GoLib.GetDockerClient()
	routingKey          string
	externalRoutingKey  string
	externalServiceName string
	etcdClient          *clientv3.Client = GoLib.GetEtcdClient()
	agentConfig         EnvironmentConfig
)

type EnvironmentConfig struct {
	Name             string    `json:"name"`
	ActiveServices   *[]string `json:"services"`
	ActiveSince      *time.Time
	ConfigUpdated    *time.Time
	RoutingKeyOutput string
	RoutingKeyInput  string
	InputQueueName   string
	ServiceName      string
}

func main() {
	defer logFile.Close()
	defer etcdClient.Close()

	// Because there will be several agents running in this test setup add (and register) a guid for uniqueness
	serviceName = fmt.Sprintf("agent-%s_service", GoLib.GenerateGuid(1))
	routingKey = GoLib.GetDefaultRoutingKey(serviceName)
	externalRoutingKey = fmt.Sprintf("%s-input", routingKey)
	externalServiceName = fmt.Sprintf("%s-input", serviceName)

	// Define a WaitGroup
	var wg sync.WaitGroup
	wg.Add(2)

	// Connect to AMQ queue, declare own routingKey as queue
	_, conn, channel, err := GoLib.SetupConnection(serviceName, routingKey, false)
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	defer conn.Close()

	// The agent will need an external queue to receive requests from the orchestrator
	messages := registerExternalQueue(channel)

	// Register all basic info of this agents environment
	registerAgent()

	// Start listening for messages, this method keeps this method 'alive'
	go func() {
		startMessageLoop(messages, "")
		wg.Done() // Decrement the WaitGroup counter when the goroutine finishes
	}()

	// Wait for both goroutines to finish
	wg.Wait()
}

func registerAgent() {
	// Prepare agent configuration data

	agentConfig = EnvironmentConfig{
		Name:             os.Getenv("HOSTNAME"),
		RoutingKeyOutput: routingKey,
		ServiceName:      serviceName,
		InputQueueName:   externalServiceName,
		RoutingKeyInput:  externalRoutingKey,
	}

	// Serialize agent configuration data as JSON
	configData, err := json.Marshal(agentConfig)
	if err != nil {
		log.Fatal(err)
	}

	go GoLib.CreateEtcdLeaseObject(etcdClient, fmt.Sprintf("/agents/%s", agentConfig.Name), string(configData))
}

func updateAgent() {

	// Update the ActiveSince field
	now := time.Now()
	agentConfig.ConfigUpdated = &now

	// Serialize agent configuration data as JSON
	configData, err := json.Marshal(agentConfig)
	if err != nil {
		log.Fatal(err)
	}

	go GoLib.CreateEtcdLeaseObject(etcdClient, fmt.Sprintf("/agents/%s", agentConfig.Name), string(configData))

}

func registerExternalQueue(channel *amqp.Channel) <-chan amqp.Delivery {

	queue, err := GoLib.DeclareQueue(externalServiceName, channel)
	if err != nil {
		log.Fatalf("Failed to declare queue: %v", err)
	}

	// Bind queue to "topic_exchange"
	// TODO: Make "topic_exchange" flexible?
	if err := channel.QueueBind(
		queue.Name,         // name
		externalRoutingKey, // key
		"topic_exchange",   // exchange
		false,              // noWait
		nil,                // args
	); err != nil {
		log.Fatalf("Queue Bind: %s", err)
	}

	// Start listening to queue
	messages, err := GoLib.Consume(externalServiceName, channel)
	if err != nil {
		log.Fatalf("Failed to register consumer: %v", err)

	} else {
		log.Printf("Registered consumer: %s", externalServiceName)
	}

	return messages
}
