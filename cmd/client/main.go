package main

import (
	"flag"
	"log"

	"github.com/dragun-igor/messenger/internal/client"
)

func main() {
	flag.Parse()
	client := client.NewClient()
	if err := client.Serve(); err != nil {
		log.Fatalln(err)
	}
}
