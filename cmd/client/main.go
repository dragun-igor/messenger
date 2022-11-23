package main

import (
	"context"
	"flag"
	"log"

	"github.com/dragun-igor/messenger/internal/client"
)

var (
	grpcAddr = flag.String("gaddr", "localhost:50051", "GRPC host")
	promAddr = flag.String("paddr", "localhost:9094", "Prometheus host")
)

func main() {
	flag.Parse()
	ctx := context.Background()
	c, err := client.New(*grpcAddr, *promAddr)
	if err != nil {
		log.Fatalln(err)
	}
	if err := c.Serve(ctx); err != nil {
		log.Fatalln(err)
	}
}
