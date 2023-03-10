package main

import (
	"log"
	"net/url"

	"github.com/ltovarm/Geoffrey_App/sensor/mosquitto"
)

func main() {
	err := mosquitto.LibInit()
	if err != nil {
		log.Panic(err)
	}
	defer mosquitto.LibCleanup()

	mqtt, err := mosquitto.New("mqtt", true)
	if err != nil {
		log.Panic(err)
	}
	defer mqtt.Destroy()

	err = mqtt.UsernamePwSet("mqtt", "mqtt")
	if err != nil {
		log.Panic(err)
	}

	u, err := url.Parse("http://host.docker.internal:5672")
	if err != nil {
		log.Panic(err)
	}

	err = mqtt.Connect(u, 10)
	if err != nil {
		log.Panic(err)
	}
	defer mqtt.Disconnect()

}
