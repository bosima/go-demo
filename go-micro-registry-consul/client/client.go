package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"registry-consul/proto"

	"github.com/go-micro/plugins/v4/registry/consul"
	"go-micro.dev/v4"
	"go-micro.dev/v4/client"
)

func main() {

	registry := consul.NewRegistry()

	service := micro.NewService(
		micro.Client(client.NewClient()),
		micro.Registry(registry),
	)

	service.Init()
	client := proto.NewHelloService("registry-consul.service", service.Client())

	rsp, err := client.Say(context.TODO(), &proto.SayRequest{Name: "BOSSMA"})
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(rsp)

	fmt.Println("Press Enter key to exit the program...")
	in := bufio.NewReader(os.Stdin)
	_, _, _ = in.ReadLine()
}
