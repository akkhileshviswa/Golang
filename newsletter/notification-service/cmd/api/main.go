package main

import (
	"fmt"
	"log"
	"mail-service/event"
	"math"
	"net/http"
	"os"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

const webPort = "8084"

func main() {
	rabbitConn, err := rabbitConnect()
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	defer rabbitConn.Close()

	log.Println("Listening for and consuming RabbitMQ messages...")

	consumer, err := event.NewConsumer(rabbitConn)
	if err != nil {
		panic(err)
	}

	err = consumer.Listen([]string{"Mail"}, rabbitConn)
	if err != nil {
		log.Println(err)
	}

	log.Println("Starting mail service on port", webPort)

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: Routes(),
	}

	err = srv.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}
}

// This function is used to connect to rabbit mq
func rabbitConnect() (*amqp.Connection, error) {
	var counts int64
	var backOff = 1 * time.Second
	var connection *amqp.Connection

	for {
		c, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
		if err != nil {
			fmt.Println("RabbitMQ not yet ready...")
			counts++
		} else {
			log.Println("Connected to RabbitMQ!")
			connection = c
			break
		}

		if counts > 5 {
			fmt.Println(err)
			return nil, err
		}

		backOff = time.Duration(math.Pow(float64(counts), 2)) * time.Second
		log.Println("backing off...")
		time.Sleep(backOff)
		continue
	}

	return connection, nil
}
