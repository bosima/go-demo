## Run Service
    go run main.go

## Client Request

Get http://127.0.0.1:8002/hello/say?name=bosima

Post http://127.0.0.1:8002/hello/set
{"Language":"en"}

## Attention

The request is forwarded directly to http.ServeMux.
No codec, transport, wrapper etc.