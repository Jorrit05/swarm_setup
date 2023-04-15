package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"sync"

	"github.com/Jorrit05/GoLib"
	amqp "github.com/rabbitmq/amqp091-go"

	_ "github.com/go-sql-driver/mysql"
)

var (
	log, logFile        = GoLib.InitLogger(serviceName)
	serviceName  string = "query_service"
	routingKey   string = GoLib.GetDefaultRoutingKey(serviceName)
	db           *sql.DB
)

func main() {
	defer logFile.Close()
	defer GoLib.HandlePanicAndFlushLogs(log, logFile)

	// Define a WaitGroup
	var wg sync.WaitGroup
	wg.Add(1)

	// Open a database connection
	connectionString, _ := GoLib.GetSQLConnectionString()

	// Error specified separately because DB is 'global' and shouldn't be overwritten
	var err error
	db, err = sql.Open("mysql", connectionString)
	if err != nil {
		log.Fatalf("Error opening database:", err)
	}
	defer db.Close()

	// Connect to AMQ queue, declare own routingKey as queue
	messages, conn, channel, err := GoLib.SetupConnection(serviceName, routingKey, true)
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	defer conn.Close()

	go func() {
		GoLib.StartMessageLoop(doQuery, messages, channel, routingKey, "")
		wg.Done() // Decrement the WaitGroup counter when the goroutine finishes
	}()

	// Wait for both goroutines to finish
	wg.Wait()
	// Start listening for messages, this method keeps this method 'alive'
}

type SkillQuery struct {
	PersonId  int    `json:"person_id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Sex       string `json:"sex"`
	// DriversLicense string `json:"drivers_license"`
	// Programming    string `json:"programming"`
}

func doQuery(message amqp.Delivery) (amqp.Publishing, error) {

	var data interface{}

	if err := json.Unmarshal(message.Body, &data); err != nil {
		log.Printf("Error unmarshaling JSON: %s", err)
		return amqp.Publishing{}, err
	}

	// Access values using type assertions
	query, ok := data.(map[string]interface{})["query"]
	if !ok {
		log.Println("Error accessing key 'query'")
		return amqp.Publishing{}, errors.New("issue accessing json key: 'query'")
	}

	queryString, ok := query.(string)
	if !ok {
		log.Println("queryString is not a string")
		return amqp.Publishing{}, errors.New("queryString is not a string")
	}

	if db == nil {
		log.Println("Error: db is not initialized")
		return amqp.Publishing{}, errors.New("db is not initialized")
	}

	rows, err := db.Query(queryString)
	if err != nil {
		log.Printf("Error executing query:", err)
		return amqp.Publishing{}, err
	}
	defer rows.Close()

	var skillQueries []SkillQuery

	for rows.Next() {
		var first_name string
		var last_name string
		var sex string
		var person_id int
		if err := rows.Scan(&first_name, &last_name, &sex, &person_id); err != nil {
			log.Printf("Error scanning row:", err)
			return amqp.Publishing{}, err
		}

		skillQueries = append(skillQueries, SkillQuery{
			PersonId:  person_id,
			FirstName: first_name,
			LastName:  last_name,
			Sex:       sex,
		})
	}

	jsonMessage, err := json.Marshal(skillQueries)
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
