package main

import (
	"flag"
	"log"
	"os"

	"github.com/streadway/amqp"
)

type settings struct {
	brokerHost string
	queue      string
}

func main() {

	var sett settings
	flag.StringVar(&sett.brokerHost, "broker", os.Getenv("AMQP_SERVER_URL"), "Broker host")
	flag.StringVar(&sett.queue, "queue", "mqtt", "Queue")
	flag.Parse()

	// Create a new RabbitMQ connection.
	connectRabbitMQ, err := amqp.Dial(sett.brokerHost)
	if err != nil {
		panic(err)
	}
	defer connectRabbitMQ.Close()

	// Opening a channel to our RabbitMQ instance over
	// the connection we have already established.
	channelRabbitMQ, err := connectRabbitMQ.Channel()
	if err != nil {
		panic(err)
	}
	defer channelRabbitMQ.Close()

	// Subscribing to QueueService1 for getting messages.
	messages, err := channelRabbitMQ.Consume(
		sett.queue, // queue name
		"",         // consumer
		true,       // auto-ack
		false,      // exclusive
		false,      // no local
		false,      // no wait
		nil,        // arguments
	)
	if err != nil {
		log.Println(err)
	}

	// Make a channel to receive messages into infinite loop.
	forever := make(chan bool)

	go func() {
		for message := range messages {
			// For example, show received message in a console.
			// TO-DO: Sender the data to process data container
			log.Printf(" > Received message: %s\n", message.Body)
		}
	}()

	<-forever
}
