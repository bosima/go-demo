package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"rpchello/proto"

	"go-micro.dev/v4"
	"go-micro.dev/v4/client"
)

func main() {

	service := micro.NewService(
		micro.Client(client.NewClient()),
	)

	service.Init()
	client := proto.NewHelloService("rpchello.service", service.Client())

	rsp, err := client.Say(context.TODO(), &proto.SayRequest{Name: "BOSSMA"})
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(rsp)

	fmt.Println("Press Enter key to exit the program...")
	in := bufio.NewReader(os.Stdin)
	_, _, _ = in.ReadLine()
}
