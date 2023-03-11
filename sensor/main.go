package main

import (
	"log"
	"math/rand"
	"net/url"
	"time"
	"unsafe"

	"github.com/ltovarm/Geoffrey_App/sensor/mosquitto"
)

import "C"

var temperature C.float = 0.0

func sensorData() (int, unsafe.Pointer) {

	temperature = C.float(10 + 20*rand.Float32())
	return int(C.sizeof_float), unsafe.Pointer(&temperature)
}

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

	u, err := url.Parse("http://host.docker.internal:1883")
	if err != nil {
		log.Panic(err)
	}

	err = mqtt.Connect(u, 10)
	if err != nil {
		log.Panic(err)
	}
	defer mqtt.Disconnect()

	err = mqtt.LoopStart()
	if err != nil {
		log.Panic(err)
	}

	for {
		payloadlen, payload := sensorData()
		mid, err := mqtt.Publish("mqtt", payloadlen, payload, 2, false)
		if err != nil {
			log.Print(err)
		}
		log.Print("Message ", mid, " has been published")
		time.Sleep(1 * time.Second)
	}
}
