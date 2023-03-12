package main

import (
	"flag"
	"log"
	"net/url"
)

type settings struct {
	url      *url.URL
	user     string
	password string
	topic    string
}

func settingsDefault() settings {
	localhostUrl, _ := url.Parse("http://localhost:1883")
	return settings{
		url:      localhostUrl,
		user:     "",
		password: "",
		topic:    "",
	}
}

func settingsFlags(s *settings) {
	var localhostUrl string
	flag.StringVar(&localhostUrl, "host", s.url.String(), "MQTT host broker")
	flag.StringVar(&s.user, "user", s.user, "MQTT broker user")
	flag.StringVar(&s.password, "password", s.password, "MQTT broker password")
	flag.StringVar(&s.topic, "topic", s.topic, "MQTT topic")
	flag.Parse()

	urlHost, err := url.Parse(localhostUrl)
	if err != nil {
		log.Panic(err)
	}
	s.url = urlHost
}

func config() settings {
	sett := settingsDefault()
	settingsFlags(&sett)

	return sett
}
