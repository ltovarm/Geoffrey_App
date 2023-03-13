package main

import (
	"flag"
)

type settings struct {
	url      string
	port     int
	user     string
	password string
	topic    string
}

func settingsDefault() settings {
	return settings{
		url:      "localhost",
		port:     1883,
		user:     "",
		password: "",
		topic:    "",
	}
}

func settingsFlags(s *settings) {
	flag.IntVar(&s.port, "port", s.port, "MQTT port")
	flag.StringVar(&s.url, "host", s.url, "MQTT host broker")
	flag.StringVar(&s.user, "user", s.user, "MQTT broker user")
	flag.StringVar(&s.password, "password", s.password, "MQTT broker password")
	flag.StringVar(&s.topic, "topic", s.topic, "MQTT topic")
	flag.Parse()
}

func config() settings {
	sett := settingsDefault()
	settingsFlags(&sett)

	return sett
}
