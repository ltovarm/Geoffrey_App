package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"

	"github.com/gorilla/websocket"
	query "github.com/ltovarm/Geoffrey_App/BackEnd/internal/query"

	_ "github.com/lib/pq"

	"github.com/streadway/amqp"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Permitir conexiones desde cualquier origen (CORS)
	},
}

// Estructura para manejar los clientes WebSocket
type WebSocketClients struct {
	clients map[*websocket.Conn]bool
	mutex   sync.Mutex
}

var wsClients = WebSocketClients{
	clients: make(map[*websocket.Conn]bool),
}

func main() {
	// Define RabbitMQ server URL.
	amqpServerURL := os.Getenv("AMQP_SERVER_URL")
	// Create a new RabbitMQ connection.
	connectRabbitMQ, err := amqp.Dial("amqpServerURL")
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
	// Establecer conexión WebSocket para enviar datos al frontend
	http.HandleFunc("/ws", handleWebSocket)
	fmt.Println("Server running at http://localhost:8080")
	http.ListenAndServe(":8080", nil)

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

		insertJsonToTable(data, "temperatures")
		// Procesar el mensaje
		log.Printf("Mensaje enviado: %v", data)
		// Send data to all connected clients
		for client := range wsClients.clients {
			err := client.WriteMessage(websocket.TextMessage, msg.Body)
			if err != nil {
				fmt.Println("Error sending data to client:", err)
			}
		}
	}
}

func insertJsonToTable(data map[string]interface{}, sqltable string) {

	// Set-up connection
	my_db := query.NewDb()
	if err := my_db.ConnectToDatabaseFromEnvVar(); err != nil {
		log.Fatalf("Error connecting to db: %s\n", err)
	}
	// Insert into table
	if err := my_db.SendDataAsJSON(data, sqltable); err != nil {
		log.Fatalf("Error inserting data to db: %s\n", err)
	}
	my_db.CloseDatabase()
}

// Función para manejar WebSocket connections
func handleWebSocket(w http.ResponseWriter, r *http.Request) {
	// Upgrade the HTTP connection to WebSocket
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatalf("Error al actualizar la conexión a WebSocket: %s", err)
		return
	}
	defer conn.Close()

	// Agregar el cliente WebSocket a la lista de clientes
	wsClients.mutex.Lock()
	wsClients.clients[conn] = true
	wsClients.mutex.Unlock()

	// Leer mensajes entrantes desde el cliente (no utilizado en este ejemplo)
	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			// Eliminar el cliente WebSocket de la lista cuando se desconecta
			wsClients.mutex.Lock()
			delete(wsClients.clients, conn)
			wsClients.mutex.Unlock()
			break
		}
	}
}
