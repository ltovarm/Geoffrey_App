package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"time"
)

const (
	CONN_HOST   = "172.29.108.252"
	CONN_PORT   = "9090"
	CONN_TYPE   = "tcp"
	SUCCSESSFUL = 1
)

var conn_close = false

func main() {
	// Listen for incoming connections.
	l, err := net.Listen(CONN_TYPE, CONN_HOST+":"+CONN_PORT)
	if err != nil {
		fmt.Println("Error listening:", err.Error())
		os.Exit(1)
	}

	// Close the listener when the application closes.
	defer l.Close()
	fmt.Println("Listening on " + CONN_HOST + ":" + CONN_PORT)

	for {
		// Listen for an incoming connection.
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting: ", err.Error())
			os.Exit(1)
		}
		fmt.Println(("state 0.\n"))
		// Handle connections in a new goroutine.
		go handleRequest(conn, l)
	}
}

// Handles incoming requests.
func handleRequest(conn net.Conn, l net.Listener) {
	// Make a buffer to hold incoming data.
	buf := make([]byte, 1024)
	// Read the incoming connection into the buffer.
	_, err := conn.Read(buf)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Msg:", string(buf))
	if strings.Contains(string(buf), "exit") {
		l.Close()
	}
	// Send a response back to person contacting us.
	time := time.Now().Format("Monday, 02-Jan-06 15:04:05 MST")
	conn.Write([]byte("Message received.\n"))
	conn.Write([]byte(time))
	conn.Write([]byte("\n"))

	// Close the connection when you're done with it.
	conn.Close()
}
