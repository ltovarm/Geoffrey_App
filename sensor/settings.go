package main

import (
	"flag"
)

type settings struct {
	url       string
	port      int
	keepalive int
	user      string
	password  string
	topic     string
}

func settingsDefault() settings {
	return settings{
		url:       "localhost",
		port:      1883,
		keepalive: 60,
		user:      "",
		password:  "",
		topic:     "",
	}
}

func settingsFlags(s *settings) {
	flag.IntVar(&s.port, "port", s.port, "MQTT port")
	flag.IntVar(&s.keepalive, "keepalive", s.keepalive, "MQTT broker keepalive")
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
