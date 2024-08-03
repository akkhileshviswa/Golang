package event

import (
	"encoding/json"
	"fmt"
	"log"
	"mail-service/mail"
	"strconv"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Consumer struct {
	conn      *amqp.Connection
	queueName string
}

type LogPayload struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

// This function is used to declare a new consumer for consuming messages.
func NewConsumer(conn *amqp.Connection) (Consumer, error) {
	consumer := Consumer{
		conn:      conn,
		queueName: "insert_log",
	}

	err := consumer.setup()
	if err != nil {
		return Consumer{}, err
	}

	return consumer, nil
}

// This function is used to setup the channel for exchange.
func (consumer *Consumer) setup() error {
	channel, err := consumer.conn.Channel()
	if err != nil {
		return err
	}

	return declareExchange(channel)
}

// This function is used to listen to the topics which has been passed as argument.
func (consumer *Consumer) Listen(topics []string, conn *amqp.Connection) error {
	ch, err := consumer.conn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	q, err := declareRandomQueue(ch)
	if err != nil {
		return err
	}

	for _, s := range topics {
		ch.QueueBind(
			q.Name,
			s,
			"insert_log",
			false,
			nil,
		)

		if err != nil {
			return err
		}
	}

	messages, err := ch.Consume(q.Name, "", true, false, false, false, nil)
	if err != nil {
		return err
	}

	forever := make(chan bool)
	go func() {
		for d := range messages {
			var payload mail.MailPayload
			_ = json.Unmarshal(d.Body, &payload)
			go handlePayload(payload, conn)
		}
	}()

	fmt.Printf("Waiting for message [Exchange, Queue] [insert_log, %s]\n", q.Name)
	<-forever

	return nil
}

// This function handles which function should be called based on a switch statement.
func handlePayload(payload mail.MailPayload, conn *amqp.Connection) {
	switch payload.Name {
	case "AccountCreated":
		err := mail.SendMail(payload)
		if err != nil {
			log.Println(err)
		} else {
			emitter, err := NewEventEmitter(conn)
			if err != nil {
				log.Println(err)
			}

			j, err := json.MarshalIndent(&payload, "", "\t")
			if err != nil {
				log.Println(err)
			}

			err = emitter.Push(string(j), "UpdateRecord")
			if err != nil {
				log.Println(err)
			}

			logMessage := LogPayload{
				Name: "Update",
				Data: "Mail Sent for: " + payload.FirstName + "ID: " + strconv.Itoa(payload.ID),
			}

			msg, err := json.MarshalIndent(&logMessage, "", "\t")
			if err != nil {
				log.Println(err)
			}

			err = emitter.Push(string(msg), "log.INFO")
			if err != nil {
				log.Println(err)
			}
		}
	}
}
