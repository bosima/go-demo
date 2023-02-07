package main

import (
	"github.com/bosima/go/leader-election/consul"
	"github.com/bosima/ylog"
)

func main() {
	serviceId := "test-service-1"
	serviceName := "test-service"

	go func() {
		election := consul.NewConsulLeaderElection(serviceId, serviceName, 10, electResultHandler)
		err := election.Run()
		if err != nil {
			ylog.Error(err)
		}
	}()

	ylog.Info("Elect Started.")

	// here do your business

	select {}
}

func electResultHandler(result bool) {
	ylog.Info("Elect Result:", result)
}
