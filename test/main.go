package main

import (
	"context"
	"encoding/json"
	"time"

	"github.com/Jorrit05/GoLib"
	clientv3 "go.etcd.io/etcd/client/v3"
)

var (
	serviceName string = "test_service"

	log, logFile = GoLib.InitLogger(serviceName)
	etcdClient   *clientv3.Client
)

// func (c CreateServicePayload) String() string {
// 	var sb strings.Builder

// 	sb.WriteString(fmt.Sprintf("ImageName: %s\n", c.ImageName))
// 	sb.WriteString(fmt.Sprintf("Tag: %s\n", c.Tag))
// 	sb.WriteString("EnvVars:\n")
// 	for k, v := range c.EnvVars {
// 		sb.WriteString(fmt.Sprintf("  %s: %s\n", k, v))
// 	}
// 	sb.WriteString(fmt.Sprintf("Networks: %v\n", c.Networks))
// 	sb.WriteString(fmt.Sprintf("Secrets: %v\n", c.Secrets))
// 	sb.WriteString("Volumes:\n")
// 	for k, v := range c.Volumes {
// 		sb.WriteString(fmt.Sprintf("  %s: %s\n", k, v))
// 	}
// 	sb.WriteString("Ports:\n")
// 	for k, v := range c.Ports {
// 		sb.WriteString(fmt.Sprintf("  %s: %s\n", k, v))
// 	}
// 	sb.WriteString(fmt.Sprintf("Deploy: \n"))
// 	sb.WriteString(fmt.Sprintf("  Replicas: %d\n", c.Deploy.Replicas))
// 	sb.WriteString(fmt.Sprintf("  Placement: \n"))
// 	sb.WriteString(fmt.Sprintf("    Constraints: %v\n", c.Deploy.Placement.Constraints))
// 	sb.WriteString(fmt.Sprintf("  Resources: \n"))
// 	sb.WriteString(fmt.Sprintf("    Reservations: \n"))
// 	sb.WriteString(fmt.Sprintf("      Memory: %s\n", c.Deploy.Resources.Reservations.Memory))
// 	sb.WriteString(fmt.Sprintf("    Limits: \n"))
// 	sb.WriteString(fmt.Sprintf("      Memory: %s\n", c.Deploy.Resources.Limits.Memory))

//		return sb.String()
//	}

// T is the struct type to be saved.
// target is an instance of the struct.
// etcdClient is an instance of the etcd client.
// key is the etcd key where the value will be stored.
func SaveStructToEtcd[T any](etcdClient *clientv3.Client, key string, target T) error {
	// Marshal the target struct into a JSON representation
	jsonRep, err := json.Marshal(target)
	if err != nil {
		log.Errorf("failed to marshal struct: %v", err)
		return err
	}

	// Save the JSON representation to the etcd key-value store
	_, err = etcdClient.Put(context.Background(), key, string(jsonRep))
	if err != nil {
		log.Errorf("failed to save struct to etcd: %v", err)
		return err
	}

	return nil
}

func main() {
	etcdClient, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"localhost:2379"},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		log.Fatal(err)
	}
	type SampleStruct struct {
		Name  string `json:"name"`
		Value int    `json:"value"`
	}

	sample := SampleStruct{Name: "test", Value: 42}
	err = SaveStructToEtcd(etcdClient, "/sample/test", sample)
	if err != nil {
		log.Fatalf("failed to save struct to etcd: %v", err)
	}
	// registerJSONArray[GoLib.ArcheType](archetypesJSON, &GoLib.ArcheTypes{}, etcdClient, "/reasoner/archetype_config")

	// 	// var err error = nil
	// 	payload := CreateServicePayload{
	// 		ImageName: "my-image",
	// 		Tag:       "latest",
	// 		EnvVars:   map[string]string{"ENV1": "value1", "ENV2": "value2"},
	// 		Networks:  []string{"network1", "network2"},
	// 		Secrets:   []string{"secret1", "secret2"},
	// 		Volumes:   map[string]string{"volume1": "/path1", "volume2": "/path2"},
	// 		Ports:     map[string]string{"8080": "80"},
	// 		Deploy: Deploy{
	// 			Replicas:  2,
	// 			Placement: Placement{Constraints: []string{"node.role == worker"}},
	// 			Resources: Resources{
	// 				Reservations: Resource{Memory: "100M"},
	// 				Limits:       Resource{Memory: "200M"},
	// 			},
	// 		},
	// 	}

	// 	jsonData, err := json.Marshal(payload)
	// 	if err != nil {
	// 		fmt.Printf("Error marshaling payload to JSON:", err)
	// 	}

	// 	fmt.Printf(string(jsonData))

	// defer logFile.Close()
	// mux := http.NewServeMux()
	// mux.HandleFunc("/", handler)
	// go func() {
	// 	fmt.Println("ListenAndServe: 1")

	// 	if err := http.ListenAndServe(":3000", mux); err != nil {

	// 		log.Fatalf("Error starting HTTP server: %s", err)
	// 	}
	// }()
	// fmt.Println("3")
	// select {}
}
