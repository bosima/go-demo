package main

import (
	"context"
	"log"

	"github.com/bosima/go-demo/go-micro-opentracing/config"
	"github.com/bosima/go-demo/go-micro-opentracing/zipkin"
	"github.com/go-micro/plugins/v4/wrapper/trace/opentracing"

	"go-micro.dev/v4"
)

type Hello struct {
}

func (h *Hello) Say(ctx context.Context, name *string, resp *string) error {
	*resp = "Hello " + *name
	return nil
}

func main() {
	tracer := zipkin.GetTracer(config.SERVICE_NAME, config.SERVICE_HOST)
	defer zipkin.Close()
	tracerHandler := opentracing.NewHandlerWrapper(tracer)

	service := micro.NewService(
		micro.Name(config.SERVICE_NAME),
		micro.Address(config.SERVICE_HOST),
		micro.WrapHandler(tracerHandler),
	)

	service.Init()

	micro.RegisterHandler(service.Server(), &Hello{})

	if err := service.Run(); err != nil {
		log.Println(err)
	}
}
