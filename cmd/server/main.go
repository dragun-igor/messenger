package main

import (
	"context"
	"log"

	"github.com/dragun-igor/messenger/config"
	"github.com/dragun-igor/messenger/internal/server"
)

func main() {
	config := config.Get()
	server, err := server.NewServer(context.Background(), config)
	if err != nil {
		log.Fatalln(err)
	}
	if err := server.Serve(); err != nil {
		log.Fatalln(err)
	}
}
