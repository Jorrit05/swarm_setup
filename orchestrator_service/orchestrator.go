package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/Jorrit05/GoLib"
	"github.com/docker/docker/client"
	clientv3 "go.etcd.io/etcd/client/v3"
)

var (
	serviceName string = "orchestrator_service"
	routingKey  string

	log, logFile                = GoLib.InitLogger(serviceName)
	dockerClient *client.Client = GoLib.GetDockerClient()

	externalRoutingKey  string
	externalServiceName string
	etcdClient          *clientv3.Client = GoLib.GetEtcdClient()
	agentConfig         GoLib.EnvironmentConfig
)

func main() {
	defer logFile.Close()
	routingKey = GoLib.GetDefaultRoutingKey(serviceName)

	// Register a yaml file of available microservices in etcd.
	processedServices, err := GoLib.SetMicroservicesEtcd(&GoLib.EtcdClientWrapper{Client: etcdClient}, "/var/log/stack-files/microservices.yml")
	if err != nil {
		log.Fatalf("Error setting microservices in etcd: %v", err)
	}

	for serviceName, _ := range processedServices {
		log.Infof("serviceName added to etcd, %s", serviceName)
	}

	http.HandleFunc("/put", putHandler(etcdClient))
	http.HandleFunc("/get", getHandler(etcdClient))

	log.Println("Starting server on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func putHandler(client *clientv3.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		key := r.URL.Query().Get("key")
		value := r.URL.Query().Get("value")

		if key == "" || value == "" {
			http.Error(w, "Both key and value parameters are required", http.StatusBadRequest)
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()

		_, err := client.Put(ctx, key, value)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		fmt.Fprintf(w, "Key '%s' set to '%s'\n", key, value)
	}
}

func getHandler(client *clientv3.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		key := r.URL.Query().Get("key")

		if key == "" {
			http.Error(w, "Key parameter is required", http.StatusBadRequest)
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()

		resp, err := client.Get(ctx, key)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if len(resp.Kvs) == 0 {
			http.Error(w, "Key not found", http.StatusNotFound)
			return
		}

		fmt.Fprintf(w, "Key '%s' has value '%s'\n", key, string(resp.Kvs[0].Value))
	}
}
