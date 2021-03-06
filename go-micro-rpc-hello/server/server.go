package main

import (
	"context"
	"fmt"
	"log"
	"rpchello/proto"

	"go-micro.dev/v4"
	"go-micro.dev/v4/server"
)

type Hello struct{}

func (s *Hello) Say(ctx context.Context, req *proto.SayRequest, rsp *proto.SayResponse) error {
	fmt.Println("request:", req.Name)
	rsp.Message = "Hello " + req.Name
	return nil
}

func main() {

	rpcServer := server.NewServer(
		server.Name("rpchello.service"),
		server.Address("0.0.0.0:8001"),
	)

	// Register Handlers
	proto.RegisterHelloHandler(rpcServer, &Hello{})

	service := micro.NewService(
		micro.Server(rpcServer),
	)

	// optionally setup command line usage
	service.Init()

	// Run server
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
