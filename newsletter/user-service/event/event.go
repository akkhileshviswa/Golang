package event

import amqp "github.com/rabbitmq/amqp091-go"

// This function is used to return the channel with required parameters for exchange.
func declareExchange(ch *amqp.Channel) error {
	return ch.ExchangeDeclare(
		"insert_log", // name
		"topic",      // type
		true,         // durable?
		false,        // auto-deleted?
		false,        // internal?
		false,        // no-wait?
		nil,          // arguements?
	)
}

// This function is used to return the channel with required parameters for exchange.
func declareRandomQueue(ch *amqp.Channel) (amqp.Queue, error) {
	return ch.QueueDeclare(
		"",    // name?
		false, // durable?
		false, // delete when unused?
		true,  // exclusive?
		false, // no-wait?
		nil,   // arguments?
	)
}
