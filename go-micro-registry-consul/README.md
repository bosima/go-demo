## Server

	registry := consul.NewRegistry()

	rpcServer := server.NewServer(
		...
		server.Registry(registry),
	)

## Client

	registry := consul.NewRegistry()

	service := micro.NewService(
		...
		micro.Registry(registry),
	)