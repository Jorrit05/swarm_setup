package main

import (
	"context"
	"encoding/json"
	"log"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
)

type AgentConfig struct {
	IP        string `json:"ip"`
	Port      int    `json:"port"`
	OtherData string `json:"other_data"`
}

func register() {
	// Connect to etcd
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"etcd1:2379", "etcd2:2379", "etcd3:2379"},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		log.Fatal(err)
	}
	defer cli.Close()

	// Create a lease with a 10-second TTL
	lease, err := cli.Grant(context.Background(), 10)
	if err != nil {
		log.Fatal(err)
	}

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

	// Write agent information to etcd with the lease attached
	_, err = cli.Put(context.Background(), "/agents/agent1", string(configData), clientv3.WithLease(lease.ID))
	if err != nil {
		log.Fatal(err)
	}

	// Keep the lease alive by refreshing it periodically
	leaseKeepAlive, err := cli.KeepAlive(context.Background(), lease.ID)
	if err != nil {
		log.Fatal(err)
	}

	// Periodically refresh the lease
	for range leaseKeepAlive {
		log.Println("Lease refreshed")
	}
}
