package main

import (
	"encoding/json"
	"sync"

	"github.com/Jorrit05/GoLib"
	"github.com/docker/docker/client"
	amqp "github.com/rabbitmq/amqp091-go"
	clientv3 "go.etcd.io/etcd/client/v3"
)

var (
	serviceName  string           = "agent_service"
	log, logFile                  = GoLib.InitLogger(serviceName)
	cli          *client.Client   = GoLib.GetDockerClient()
	routingKey   string           = GoLib.GetDefaultRoutingKey(serviceName)
	etcdClient   *clientv3.Client = GoLib.GetEtcdClient()
)

type AgentConfig struct {
	IP        string `json:"ip"`
	Port      int    `json:"port"`
	OtherData string `json:"other_data"`
}

func main() {
	defer logFile.Close()
	defer etcdClient.Close()

	// Define a WaitGroup
	var wg sync.WaitGroup
	wg.Add(2)

	// Prepare agent configuration data
	agentConfig := AgentConfig{
		IP:        "10.0.0.1",
		Port:      8080,
		OtherData: "example_data",
	}

	// Serialize agent configuration data as JSON
	configData, err := json.Marshal(agentConfig)
	if err != nil {
		log.Fatal(err)
	}

	go GoLib.CreateEtcdLeaseObject(etcdClient, "agent1", string(configData))

	// Connect to AMQ queue, declare own routingKey as queue
	messages, conn, channel, err := GoLib.SetupConnection(serviceName, routingKey, true)
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	defer conn.Close()

	// Start listening for messages, this method keeps this method 'alive'
	go func() {
		GoLib.StartMessageLoop(placeholder, messages, channel, routingKey, "")
		wg.Done() // Decrement the WaitGroup counter when the goroutine finishes
	}()

	// // Example service specification
	// envVars := map[string]string{
	// 	"INPUT_QUEUE":       "query_service",
	// 	"AMQ_PASSWORD_FILE": "/run/secrets/rabbitmq_user",
	// 	"AMQ_USER":          "normal_user",
	// }

	// networks := []string{
	// 	"appnet",
	// }

	// secrets := []string{
	// 	"rabbitmq_user",
	// }

	// volumes := map[string]string{
	// 	fmt.Sprintf("/var/log/thesis_logs/%s_log.txt", serviceName): "/app/log.txt",
	// }

	// // ports := map[string]string{"80": "80"}

	// spec := GoLib.CreateServiceSpec("anonymize_service", "", envVars, networks, secrets, volumes, nil, cli)
	// defer cli.Close()

	// Create Docker Service in a separate goroutine
	// go func() {
	// 	GoLib.CreateDockerService(cli, spec)
	// 	wg.Done() // Decrement the WaitGroup counter when the goroutine finishes
	// }()

	// Wait for both goroutines to finish
	wg.Wait()
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
