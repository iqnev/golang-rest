package main

import (
	"fmt"
	"net"
	"os"

	"github.com/hashicorp/go-hclog"
	"github.com/iqnev/golang-rest/currency/server"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	protos "github.com/iqnev/golang-rest/currency/protos/currency"
)

func main() {

	log := hclog.Default()

	gc := grpc.NewServer()

	cs := server.NewCurrency(log)

	protos.RegisterCurrencyServer(gc, cs)

	reflection.Register(gc)

	l, err := net.Listen("tcp", fmt.Sprintf(":%d", 8989))
	if err != nil {
		log.Error("Unable to create listener", "error", err)
		os.Exit(1)
	}

	// listen for requests
	gc.Serve(l)
}
