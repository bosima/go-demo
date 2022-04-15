package main

import (
	"bufio"
	"context"
	"fmt"
	"grpchello/proto"
	"os"

	"github.com/asim/go-micro/plugins/client/grpc/v4"
	"go-micro.dev/v4"
)

func main() {
	// with this, you can call a gRPC service that is not registered
	// r := registry.NewRegistry()
	// r.Register(&registry.Service{
	// 	Name:    "grpchello.service",
	// 	Version: "1.0.0",
	// 	Nodes: []*registry.Node{
	// 		{
	// 			Id:      "grpchello.service-1",
	// 			Address: "127.0.0.1:8001",
	// 			Metadata: map[string]string{
	// 				"name": "grpchello",
	// 			},
	// 		},
	// 	},
	// })

	service := micro.NewService(
		micro.Client(grpc.NewClient()),
		//micro.Registry(r),
	)

	service.Init()
	client := proto.NewHelloService("grpchello.service", service.Client())

	rsp, err := client.Say(context.TODO(), &proto.SayRequest{Name: "BOSSMA"})
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(rsp)

	fmt.Println("Press Enter key to exit the program...")
	in := bufio.NewReader(os.Stdin)
	_, _, _ = in.ReadLine()
}
