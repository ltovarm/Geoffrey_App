package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/rs/cors"
)

type TemperatureData struct {
	Temperatures []float64 `json:"temperatures"`
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Permitimos conexiones desde cualquier origen (CORS)
	},
}

func main() {
	// Start a timer to send temperature data every 20 seconds
	go startTemperatureUpdates()

	// Handle WebSocket connections
	http.HandleFunc("/ws", handleWebSocket)

	// Configure CORS to allow requests from localhost:3000 (React)
	handler := cors.Default().Handler(http.DefaultServeMux)

	// Start the server on port 8080 with the CORS middleware enabled
	fmt.Println("Server running at http://localhost:8080")
	http.ListenAndServe(":8080", handler)
}

// Function to start the timer and generate random temperature data every 20 seconds
func startTemperatureUpdates() {
	for {
		// Wait for 20 seconds
		time.Sleep(time.Second * 5)

		// Generate random temperatures and send them to clients
		sendTemperatureUpdates()
	}
}

// Function to generate random temperatures between 0.0 and 40.0
func generateRandomTemperatures() []float64 {
	temperatures := make([]float64, 5)
	for i := 0; i < 5; i++ {
		temperatures[i] = rand.Float64() * 40.0
	}
	return temperatures
}

// Function to send the updated temperatures to connected clients
func sendTemperatureUpdates() {
	// Generate random temperatures
	temperatures := generateRandomTemperatures()

	// Convert data to JSON
	data := TemperatureData{
		Temperatures: temperatures,
	}
	jsonData, err := json.Marshal(data)
	if err != nil {
		fmt.Println("Error converting to JSON:", err)
		return
	}

	// Send data to all connected clients
	for client := range clients {
		err := client.WriteMessage(websocket.TextMessage, jsonData)
		if err != nil {
			fmt.Println("Error sending data to client:", err)
		}
	}

	// Show temperatures in the terminal
	fmt.Println("Updated temperatures:", temperatures)
}

// WebSocket clients
var clients = make(map[*websocket.Conn]bool)

// Function to handle WebSocket connections
func handleWebSocket(w http.ResponseWriter, r *http.Request) {
	// Upgrade the HTTP connection to WebSocket
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("Error upgrading connection to WebSocket:", err)
		return
	}
	defer conn.Close()

	// Add the client connection to the map
	clients[conn] = true

	// Read incoming messages from the client (not used in this example)
	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			// Remove the client connection from the map when disconnected
			delete(clients, conn)
			break
		}
	}
}
