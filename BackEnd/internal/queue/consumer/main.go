package main

import (
	"encoding/json"
	"log"
	"os"

	_ "github.com/lib/pq"
	query "github.com/ltovarm/Geoffrey_App/BackEnd/internal/query/queries"
	"github.com/streadway/amqp"
)

func insertJsonToTable(data map[string]interface{}) {

	// Set-up connection
	my_db := query.NewDb()
	if err := my_db.ConnectToDatabaseFromEnvVar(); err != nil {
		log.Fatalf("Error connecting to db: %s\n", err)
	}
	// Insert into table
	if err := my_db.SendDataAsJSON(data, "tempeatures"); err != nil {
		log.Fatalf("Error inserting data to db: %s\n", err)
	}

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

	// Configuración del canal
	channelRabbitMQ, err := connectRabbitMQ.Channel()
	if err != nil {
		log.Fatalf("Error al abrir un canal: %s", err)
	}
	defer channelRabbitMQ.Close()

	// Configuración de la cola
	q, err := channelRabbitMQ.QueueDeclare(
		"QueueService1", // nombre de la cola
		true,            // durabilidad
		false,           // autoeliminación
		false,           // exclusividad
		false,           // no espera
		nil,             // argumentos
	)
	if err != nil {
		log.Fatalf("Error al declarar la cola: %s", err)
	}

	// Consumir mensajes
	msgs, err := channelRabbitMQ.Consume(
		q.Name, // nombre de la cola
		"",     // etiqueta del consumidor
		true,   // autoack
		false,  // exclusividad
		false,  // no espera
		false,  // no local
		nil,    // argumentos
	)
	if err != nil {
		log.Fatalf("Error al registrar el consumidor: %s", err)
	}

	// Esperar por los mensajes
	for msg := range msgs {
		// Decodificar el JSON del mensaje
		var data map[string]interface{}
		if err := json.Unmarshal(msg.Body, &data); err != nil {
			log.Printf("Error al decodificar el JSON: %s", err)
		}
		// Procesar el mensaje
		log.Printf("Mensaje recibido: %v", data)

		insertJsonToTable(data)
		// Procesar el mensaje
		log.Printf("Mensaje enviado: %v", data)
	}
}
