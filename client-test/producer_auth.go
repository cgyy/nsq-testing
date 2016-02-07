package main

import (
	"github.com/bitly/go-nsq"
	"log"
	"os"
)

func main() {
	config := nsq.NewConfig()
	// Get the key from osin oauth2.0 server
	// user: test
	// password: test
	config.AuthSecret = "Cfs1G9zZTuC2iQLlpiwIjA"

	//cfg.TlsV1 = true
	//cfg.AuthSecret = "$ACCESS_TOKEN"
	//cfg.MaxInFlight = 1000
	//c := nsq.NewConsumer(topic, channel, cfg)
	//c.SetHandler(....)
	//lookup := "https://api-ssl.bitly.com/v3/nsq/lookup?access_token=$ACCESS_TOKEN"
	//c.ConnectToNSQLookupd(lookup)
	//<- c.StopChan

	w, er1 := nsq.NewProducer("127.0.0.1:4150", config)
	if er1 != nil {
		log.Fatal("Error :", er1)
	}

	err := w.Publish("auth_test", []byte("test"))
	if err != nil {
		log.Panic("Could not connect")
	}

	w.Stop()
}
