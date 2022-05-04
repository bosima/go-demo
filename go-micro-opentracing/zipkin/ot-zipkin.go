package zipkin

import (
	"log"

	"github.com/bosima/go-demo/go-micro-opentracing/config"
	"github.com/opentracing/opentracing-go"
	zipkinot "github.com/openzipkin-contrib/zipkin-go-opentracing"
	"github.com/openzipkin/zipkin-go"
	"github.com/openzipkin/zipkin-go/reporter"
	zipkinhttp "github.com/openzipkin/zipkin-go/reporter/http"
)

var zipkinReporter reporter.Reporter

func GetTracer(serviceName string, host string) opentracing.Tracer {
	// set up a span reporter
	zipkinReporter = zipkinhttp.NewReporter(config.ZIPKIN_SERVER_URL)

	// create our local service endpoint
	endpoint, err := zipkin.NewEndpoint(serviceName, host)
	if err != nil {
		log.Fatalf("unable to create local endpoint: %+v\n", err)
	}

	// initialize our tracer
	nativeTracer, err := zipkin.NewTracer(zipkinReporter, zipkin.WithLocalEndpoint(endpoint))
	if err != nil {
		log.Fatalf("unable to create tracer: %+v\n", err)
	}

	// use zipkin-go-opentracing to wrap our tracer
	tracer := zipkinot.Wrap(nativeTracer)
	opentracing.InitGlobalTracer(tracer)
	return tracer
}

func Close() {
	zipkinReporter.Close()
}
