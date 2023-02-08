1. Install

```go get github.com/bosima/go-demo/leader-election```

2. Usage

```
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
```