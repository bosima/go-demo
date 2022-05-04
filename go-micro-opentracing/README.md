## Run a Zipkin

    docker run -d -p 9411:9411 openzipkin/zipkin

## Run demo

    go run main.go
    go run client/main.go

## View Zipkin

http://127.0.0.1:9411/zipkin