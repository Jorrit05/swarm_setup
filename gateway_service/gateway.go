package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/Jorrit05/GoLib"
	amqp "github.com/rabbitmq/amqp091-go"
)

var routingKey string = os.Getenv("ROUTING_KEY")
var serviceName string = "gateway_service"

func main() {
	// Log to file
	f := GoLib.StartLog()
	defer f.Close()
	log.SetOutput(f)

	http.HandleFunc("/", handler)
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Error starting HTTP server: %s", err)
	}

}

func placeholder(message amqp.Delivery) (amqp.Publishing, error) {
	return amqp.Publishing{}, nil
}

func handler(w http.ResponseWriter, r *http.Request) {

	// Connect to AMQ queue, declare own routingKey as queue, start listening for messages
	messages, conn, channel, err := GoLib.SetupConnection(serviceName, routingKey)
	if err != nil {
		log.Fatalf("Failed to setup proper connection to RabbitMQ: %v", err)
	}
	defer conn.Close()
	GoLib.StartMessageLoop(placeholder, messages, channel, routingKey, "")

	fmt.Fprintf(w, "Response sent to reply-queue")
}
