package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/Jorrit05/GoLib"
)

var (
	counter            = 1
	serviceName string = "query_service"
	routingKey  string = GoLib.GetDefaultRoutingKey(serviceName)
)

func LastPartAfterSlash(s string) string {
	splitted := strings.Split(s, "/")
	return splitted[len(splitted)-1]
}

func main() {
	fmt.Println(LastPartAfterSlash("ThisIsString"))

	// Log to file
	f := GoLib.StartLog()
	defer f.Close()
	log.SetOutput(f)

	fmt.Println(routingKey)
	// ctx = context.Background()
	// Setup HTTP server
	// mux := http.NewServeMux()
	// mux.HandleFunc("/", handler)
	// go func() {
	// 	// 	log.Println("ListenAndServe: 1")

	// 	if err := http.ListenAndServe(":3000", mux); err != nil {
	// 		log.Fatalf("Error starting HTTP server: %s", err)
	// 	}
	// }()

	// select {}
}

func handler(w http.ResponseWriter, req *http.Request) {
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		log.Printf("handler: Error reading body: %v", err)
		http.Error(w, "handler: Error reading request body", http.StatusBadRequest)
		return
	}
	defer req.Body.Close()
	log.Printf("handler: 1: %s", string(body))

	// Create a channel to receive the response
	responseChan := make(chan string)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	go DoThing(ctx, responseChan)

	v := <-responseChan

	w.Write([]byte(v))
}

func DoThing(ctx context.Context, resp chan string) {

	if counter == 1 {
		counter += 1
		time.Sleep(10 * time.Second)
		resp <- "10"
	}
	if counter == 2 {
		counter += 1
		time.Sleep(5 * time.Second)
		resp <- "5"
	}
	if counter == 3 {
		counter += 1
		time.Sleep(2 * time.Second)
		resp <- "2"
	}
}
