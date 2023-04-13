package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/Jorrit05/GoLib"
	clientv3 "go.etcd.io/etcd/client/v3"
)

var (
	serviceName string = "test_service"

	log, logFile = GoLib.InitLogger(serviceName)
	etcdClient   *clientv3.Client
)

func main() {
	var err error = nil
	etcdClient, err = clientv3.New(clientv3.Config{
		Endpoints:   []string{"localhost:2379"},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		log.Fatal(err)
	}

	defer logFile.Close()
	mux := http.NewServeMux()
	mux.HandleFunc("/", handler)
	go func() {
		fmt.Println("ListenAndServe: 1")

		if err := http.ListenAndServe(":3000", mux); err != nil {

			log.Fatalf("Error starting HTTP server: %s", err)
		}
	}()
	fmt.Println("3")
	select {}
}
