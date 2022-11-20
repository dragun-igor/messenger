package main

import (
	"flag"
	"log"

	"github.com/dragun-igor/messenger/internal/client"
)

var (
	ghost = flag.String("ghost", "localhost", "GRPC host")
	gport = flag.String("gport", "50051", "GRPC port")
	phost = flag.String("phost", "localhost", "Prometheus host")
	pport = flag.String("pport", "9094", "Prometheus port")
)

func main() {
	flag.Parse()
	client := client.NewClient(*phost, *pport)
	if err := client.Serve(*ghost, *gport); err != nil {
		log.Fatalln(err)
	}
}
