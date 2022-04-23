package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"registry-consul/proto"
	"time"

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

	registry := consul.NewRegistry(
	//registry.Addrs("127.0.0.1:8500"),
	)

	regCheckFunc := func(ctx context.Context) error {
		fmt.Println(time.Now().Format("2006-01-02 15:04:05") + " do register check")
		if 1+1 == 2 {
			return nil
		}
		return errors.New("this not earth")
	}

	rpcServer := server.NewServer(
		server.Name("registry-consul.service"),
		server.Address("0.0.0.0:8001"),
		server.Registry(registry),
		server.RegisterCheck(regCheckFunc),
		server.RegisterInterval(10*time.Second),
		server.RegisterTTL(20*time.Second),
	)

	proto.RegisterHelloHandler(rpcServer, &Hello{})

	service := micro.NewService(
		micro.Server(rpcServer),
	)

	// This will override RegisterInterval and RegisterTTL,
	// so we don't use it
	//service.Init()

	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
