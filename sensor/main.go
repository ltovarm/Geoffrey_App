package main

import (
	"log"

	"github.com/ltovarm/Geoffrey_App/sensor/mqtt"
)

func main() {
	err := mqtt.LibInit()
	if err != nil {
		log.Panic(err)
	}

	err = mqtt.LibCleanup()
	if err != nil {
		log.Panic(err)
	}
}
