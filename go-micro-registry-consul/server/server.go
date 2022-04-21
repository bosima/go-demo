package main

import (
	"context"
	"fmt"
	"log"
	"registry-consul/proto"

	"github.com/go-micro/plugins/v4/registry/consul"
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

	registry := consul.NewRegistry()

	rpcServer := server.NewServer(
		server.Name("registry-consul.service"),
		server.Address("0.0.0.0:8001"),
		server.Registry(registry),
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
