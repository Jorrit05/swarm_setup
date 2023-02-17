package main

import (
	"context"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/Jorrit05/GoLib"
	"github.com/google/uuid"
	amqp "github.com/rabbitmq/amqp091-go"
)

var (
	serviceName   string = "gateway_service"
	routingKey    string = GoLib.GetDefaultRoutingKey(serviceName)
	outputChannel *amqp.Channel
	// responseChannels = make(map[string]*requestInfo)
	mutex      = &sync.Mutex{}
	requestMap = make(map[string]*requestInfo)
	// main_ctx   context.Context
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
	_, conn, channel, err := GoLib.SetupConnection(serviceName, routingKey, false)
	if err != nil {
		log.Fatalf("Failed to setup proper connection to RabbitMQ: %v", err)
	}

	outputChannel = channel
	defer conn.Close()
	// main_ctx = context.Background()
	// Setup HTTP server
	mux := http.NewServeMux()
	mux.HandleFunc("/", handler)
	go func() {
		// 	log.Println("ListenAndServe: 1")

		if err := http.ListenAndServe(":3000", mux); err != nil {
			log.Fatalf("Error starting HTTP server: %s", err)
		}
	}()

	// Start a separate go routine to handle reply messages
	go startReplyLoop()

	select {}
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

	// Generate a unique identifier for the request
	requestID := uuid.New().String()

	// Create a channel to receive the response
	responseChan := make(chan amqp.Delivery)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	// req = req.WithContext(ctx)

	// Store the request information in the map
	mutex.Lock()
	requestMap[requestID] = &requestInfo{id: requestID, response: responseChan}
	mutex.Unlock()

	// Send the message to the start queue
	convertedAmqMessage := amqp.Publishing{
		// DeliveryMode: amqp.Persistent,
		Timestamp:     time.Now(),
		ContentType:   "application/json",
		CorrelationId: requestID,
		Body:          body,
		// Headers:       amqp.Table{"context": json.Marshal()},
	}
	log.Printf("handler: 3, %s", routingKey)

	if err := Publish(ctx, outputChannel, routingKey, convertedAmqMessage, ""); err != nil {
		log.Println("handler: 4")
		log.Printf("Handler: Error publishing: %s", err)
	}

	// Wait for the response from the response channel
	select {
	case msg := <-responseChan:
		log.Printf("handler: 5, msg received: %s", msg.Body)
		w.Write(msg.Body)
	case <-ctx.Done():
		log.Println("handler: 6, context timed out")
		http.Error(w, "handler: Request timed out", http.StatusRequestTimeout)
	}
}

func startReplyLoop() {
	log.Println("startReplyLoop: 0")
	// Connect to AMQ queue, declare own routingKey as queue, start listening for messages
	messages, conn, _, err := GoLib.SetupConnection("reply_service", "service.reply", true)
	if err != nil {
		log.Fatalf("Failed to setup proper connection to RabbitMQ: %v", err)
	}
	defer conn.Close()

	for msg := range messages {
		log.Println("startReplyLoop: 1")

		mutex.Lock()
		info, exists := requestMap[msg.CorrelationId]
		mutex.Unlock()
		if exists {
			log.Println("startReplyLoop: msg exists")

			info.response <- msg
			close(info.response)
			mutex.Lock()
			delete(requestMap, msg.CorrelationId)
			mutex.Unlock()
		} else {
			log.Println("startReplyLoop: msg does not exists")

		}
	}
}

func Publish(ctx context.Context, chann *amqp.Channel, routingKey string, message amqp.Publishing, exchangeName string) error {
	if exchangeName == "" {
		exchangeName = "topic_exchange"
	}
	log.Printf("Publish: 1 %s", routingKey)

	err := chann.PublishWithContext(ctx, exchangeName, routingKey, false, false, message)
	if err != nil {
		log.Printf("Publish: 2 %s", err)
		return err
	}
	log.Println("Publish: 3")

	return nil
}
