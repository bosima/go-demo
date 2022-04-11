## Compile proto:

    protoc --go_out=. --go-grpc_out=. proto/hello.proto

## Build Command:

    go build server/server.go
    go build client/client.go







