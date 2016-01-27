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

	file, e1 := os.Create("test.png")
	if e1 != nil {
		log.Fatal("file create error")
	}

	config := nsq.NewConfig()
	//
	q, _ := nsq.NewConsumer("write_file", "ch", config)
	q.AddHandler(nsq.HandlerFunc(func(message *nsq.Message) error {
		log.Print("Body Size = ", len(message.Body))
		write_size, e2 := file.Write(message.Body)
		if e2 != nil {
			log.Print("Error write")
		}
		log.Print("write size", write_size)

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
