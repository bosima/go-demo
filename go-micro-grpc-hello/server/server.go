package main

import (
	"context"
	"fmt"
	"grpchello/proto"
	"log"
	"time"

	"github.com/asim/go-micro/plugins/server/grpc/v4"
	"go-micro.dev/v4"
)

type Hello struct{}

func (s *Hello) Say(ctx context.Context, req *proto.SayRequest, rsp *proto.SayResponse) error {
	fmt.Println("request:", req.Name)
	rsp.Message = "Hello " + req.Name
	return nil
}

func main() {

	grpcServer := grpc.NewServer()

	service := micro.NewService(
		micro.RegisterTTL(time.Second*30),
		micro.RegisterInterval(time.Second*10),
		micro.Server(grpcServer),
		micro.Name("grpchello.service"),
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
