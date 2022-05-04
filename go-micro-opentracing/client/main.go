package main

import (
	"context"
	"log"
	"time"

	"github.com/bosima/go-demo/go-micro-opentracing/config"
	"github.com/bosima/go-demo/go-micro-opentracing/zipkin"
	"github.com/go-micro/plugins/v4/wrapper/trace/opentracing"
	"go-micro.dev/v4"
)

func main() {
	tracer := zipkin.GetTracer(config.CLIENT_NAME, config.CLIENT_HOST)
	defer zipkin.Close()
	tracerClient := opentracing.NewClientWrapper(tracer)

	service := micro.NewService(
		micro.Name(config.CLIENT_NAME),
		micro.Address(config.CLIENT_HOST),
		micro.WrapClient(tracerClient),
	)

	client := service.Client()

	go func() {
		for {
			<-time.After(time.Second)
			result := new(string)
			request := client.NewRequest(config.SERVICE_NAME, "Hello.Say", "FireflySoft")
			err := client.Call(context.TODO(), request, result)
			if err != nil {
				log.Println(err)
				continue
			}
			log.Println(*result)
		}
	}()

	if err := service.Run(); err != nil {
		log.Println(err)
	}
}
