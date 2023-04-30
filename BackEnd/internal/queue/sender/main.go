package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/streadway/amqp"
)

type Message struct {
	Id        int     `json:"id"`
	Temp      float32 `json:"temp"`
	Timestamp int64   `json:"ts"`
}

func main() {
	// Define RabbitMQ server URL.
	amqpServerURL := os.Getenv("AMQP_SERVER_URL")

	// Create a new RabbitMQ connection.
	connectRabbitMQ, err := amqp.Dial(amqpServerURL)
	if err != nil {
		panic(err)
	}
	defer connectRabbitMQ.Close()

	// Let's start by opening a channel to our RabbitMQ
	// instance over the connection we have already
	// established.
	channelRabbitMQ, err := connectRabbitMQ.Channel()
	if err != nil {
		panic(err)
	}
	defer channelRabbitMQ.Close()

	// With the instance and declare Queues that we can publish and subscribe to.
	// TO-DO: Queues must be declare from xml or json.
	q, err := channelRabbitMQ.QueueDeclare(
		"QueueService1", // queue name
		true,            // durable
		false,           // auto delete
		false,           // exclusive
		false,           // no wait
		nil,             // arguments
	)
	if err != nil {
		panic(err)
	}

	// Create a new Fiber instance.
	app := fiber.New()

	// Add middleware.
	app.Use(
		logger.New(), // add simple logger
	)

	// Add route for send message to Service 1.
	app.Get("/send", func(c *fiber.Ctx) error {
		// Create a message to publish.
		messageRcv := amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(c.Query("msg")),
		}
		// Process msg received
		valueStr := string(messageRcv.Body)
		value, err := strconv.ParseFloat(valueStr, 32)
		if err != nil {
			fmt.Println("El string no es un float válido.")
		}

		// Create message in JSON format
		message := Message{Id: 1, Temp: float32(value), Timestamp: time.Now().Unix()}
		body, err := json.Marshal(message)
		if err != nil {
			log.Fatalf("Error al crear mensaje JSON: %v", err)
		}

		// Publicar mensaje en la cola
		err = channelRabbitMQ.Publish(
			"",     // Intercambio vacío
			q.Name, // Nombre de la cola
			false,  // No esperar confirmación
			false,  // No requerir confirmación de mandatarios
			amqp.Publishing{
				ContentType: "application/json",
				Body:        body,
			},
		)
		if err != nil {
			log.Fatalf("Error al publicar mensaje: %v", err)
		}
		return nil
	})

	// Start Fiber API server.
	log.Fatal(app.Listen(":3000"))
}
