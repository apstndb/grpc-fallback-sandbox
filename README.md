# Example of gRPC Fallback with Cloud Tasks

Currently, it is not safe because security feature is not implemented.

```
$ go run google.golang.org/grpc/examples/helloworld/greeter_server

$ fallback-proxy -address "localhost:50051"

$ ngrok http 1337

$ go run ./ --baseurl=https://12345678.ngrok.io --project=project --location=asia-northeast1 hello
```

## Reference

https://googleapis.github.io/HowToRPC.html#grpc-fallback-experimental

