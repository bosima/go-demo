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

	service := micro.NewService(
		micro.Client(grpc.NewClient()),
	)

	service.Init()
	client := proto.NewHelloService("grpchello.service", service.Client())

	rsp, err := client.Say(context.TODO(), &proto.SayRequest{Name: "BOSSMA"})
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(rsp)

	fmt.Println("按回车键退出程序...")
	in := bufio.NewReader(os.Stdin)
	_, _, _ = in.ReadLine()
}
