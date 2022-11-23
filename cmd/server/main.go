package main

import (
	"context"
	"log"

	"github.com/dragun-igor/messenger/config"
	"github.com/dragun-igor/messenger/internal/server"
)

func main() {
	config := config.Get()
	ctx := context.Background()
	server, err := server.New(ctx, config)
	if err != nil {
		log.Fatalln(err)
	}
	if err := server.Serve(ctx); err != nil {
		log.Fatalln(err)
	}
}
