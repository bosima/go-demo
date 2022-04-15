package main

import (
	"context"
	"fmt"
	"log"
	"rpchello/proto"
	"time"

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

	rpcServer := server.NewServer()

	service := micro.NewService(
		micro.RegisterTTL(time.Second*30),
		micro.RegisterInterval(time.Second*10),
		micro.Server(rpcServer),
		micro.Name("rpchello.service"),
		micro.Address("0.0.0.0:8001"),
	)

	// optionally setup command line usage
	service.Init()

	// Register Handlers
	proto.RegisterHelloHandler(service.Server(), &Hello{})

	// Run server
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
