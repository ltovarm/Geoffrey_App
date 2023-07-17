package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/rs/cors"
)

type TemperatureData struct {
	Temperatures []float64 `json:"temperatures"`
}

func main() {
	// Start a timer to send temperature data every second
	go startTemperatureUpdates()

	// Initialize the router
	http.HandleFunc("/temperatures", func(w http.ResponseWriter, r *http.Request) {
		// Simulate the response with random temperature data
		temperatures := generateRandomTemperatures()

		data := TemperatureData{
			Temperatures: temperatures,
		}

		// Convert data to JSON
		jsonData, err := json.Marshal(data)
		if err != nil {
			http.Error(w, "Error converting to JSON", http.StatusInternalServerError)
			return
		}

		// Configure response headers
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		// Send the JSON data in the response
		fmt.Fprintf(w, string(jsonData))
	})

	// Configure CORS to allow requests from localhost:3000 (React)
	handler := cors.Default().Handler(http.DefaultServeMux)

	// Start the server on port 8080 with the CORS middleware enabled
	fmt.Println("Server running at http://localhost:8080")
	http.ListenAndServe(":8080", handler)
}

// Function to start the timer and generate random temperature data every second
func startTemperatureUpdates() {
	for {
		// Wait for one second
		time.Sleep(time.Second * 20)

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

// Function to send the updated temperatures to clients (not implemented in this code)
func sendTemperatureUpdates() {
	// Implement the logic to send the updated temperatures to connected clients.
	// For this example, we are just printing the updated temperatures to the console.
	temperatures := generateRandomTemperatures()
	fmt.Println("Updated temperatures:", temperatures)
}
