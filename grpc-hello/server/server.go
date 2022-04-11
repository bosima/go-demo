package main

import (
	"context"
	"fmt"
	"grpchello/proto"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type server struct {
	proto.UnimplementedHelloServer
}

func (s *server) Say(ctx context.Context, req *proto.SayRequest) (*proto.SayResponse, error) {
	fmt.Println("request:", req.Name)
	return &proto.SayResponse{Message: "Hello " + req.Name}, nil
}

func main() {
	listen, err := net.Listen("tcp", ":8001")
	if err != nil {
		fmt.Printf("failed to listen: %v", err)
		return
	}
	s := grpc.NewServer()
	proto.RegisterHelloServer(s, &server{})
	reflection.Register(s)

	defer func() {
		s.Stop()
		listen.Close()
	}()

	fmt.Println("Serving 8001...")
	err = s.Serve(listen)
	if err != nil {
		fmt.Printf("failed to serve: %v", err)
		return
	}
}
