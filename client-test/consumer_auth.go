package main

import (
	"log"
	"os"
	"sync"

	"github.com/bitly/go-nsq"
)

func main() {

	wg := &sync.WaitGroup{}
	wg.Add(1)

	config := nsq.NewConfig()
	//
	//cfg.TlsV1 = true
	//cfg.AuthSecret = "$ACCESS_TOKEN"
	//cfg.MaxInFlight = 1000
	//c := nsq.NewConsumer(topic, channel, cfg)
	//c.AddHandler(....)
	//lookup := "https://api-ssl.bitly.com/v3/nsq/lookup?access_token=$ACCESS_TOKEN"
	//c.ConnectToNSQLookupd(lookup)
	//<- c.StopChan
	//
	q, _ := nsq.NewConsumer("auth_test", "ch", config)
	q.AddHandler(nsq.HandlerFunc(func(message *nsq.Message) error {
		log.Printf("Got a message: %v", message)
		//
		wg.Done()
		return nil
	}))
	// connect to lookupd http.
	err := q.ConnectToNSQLookupd("127.0.0.1:4161")
	if err != nil {
		log.Panic("Could not connect")
	}
	wg.Wait()

}
