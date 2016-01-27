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

	file, er2 := os.Open("h2.png")
	if er2 != nil {
		log.Fatal("Open File error")
	}

	// Get file size
	Size, _ := file.Seek(0, 2)
	file.Seek(0, 0)
	log.Print("File Size =", Size)
	buf := make([]byte, Size)
	file.Read(buf)
	//log.Print(buf)

	//err := w.Publish("write_test", []byte("test"))
	err := w.Publish("write_file", buf)
	if err != nil {
		log.Panic("Could not connect")
	}

	w.Stop()
}
