package main

import (
	"log"
	"math/rand"
	"time"
	"unsafe"

	"github.com/ltovarm/Geoffrey_App/sensor/mosquitto"
)

import "C"

func sensorData() (int, unsafe.Pointer) {
	temperature := C.float(10 + 20*rand.Float32())
	return int(C.sizeof_float), unsafe.Pointer(&temperature)
}

func main() {
	sett := config()

	err := mosquitto.LibInit()
	if err != nil {
		log.Panic(err)
	}
	defer mosquitto.LibCleanup()

	mqtt, err := mosquitto.New("", true)
	if err != nil {
		log.Panic(err)
	}
	defer mqtt.Destroy()

	if len(sett.user) > 0 {
		err = mqtt.UsernamePwSet(sett.user, sett.password)
		if err != nil {
			log.Panic(err)
		}
	}

	mqtt.ConnectCallbackSet(func(m mosquitto.Mosquitto, p unsafe.Pointer, mid int) {})
	mqtt.PublishCallbackSet(func(m mosquitto.Mosquitto, p unsafe.Pointer, mid int) {})

	err = mqtt.Connect(sett.url, sett.port, 60)
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
		mid, err := mqtt.Publish(sett.topic, payloadlen, payload, 2, false)
		if err != nil {
			log.Panic(err)
		} else {
			log.Printf("Message %v has been published", mid)
		}
		time.Sleep(1 * time.Second)
	}
}
