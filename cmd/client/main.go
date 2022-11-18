package main

import (
	"flag"
	"log"

	"github.com/dragun-igor/messenger/internal/client"
)

var name = flag.String("name", "Igor", "User name")

func main() {
	flag.Parse()
	client := client.NewClient("localhost:50051", *name)
	if err := client.Serve(); err != nil {
		log.Fatalln(err)
	}
}
