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
	s, err := server.New(config)
	if err != nil {
		log.Fatalln(err)
	}
	if err := s.Serve(ctx); err != nil {
		log.Fatalln(err)
	}
}
