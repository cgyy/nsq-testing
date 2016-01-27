package main

import (
	"github.com/bitly/go-nsq"
	"log"
	"os"
)

func main() {
	config := nsq.NewConfig()
	w, er1 := nsq.NewProducer("127.0.0.1:4150", config)
	if er1 != nil {
		log.Fatal("Error :", er1)
	}

	err := w.Publish("write_test", []byte("test"))
	if err != nil {
		log.Panic("Could not connect")
	}

	w.Stop()
}
