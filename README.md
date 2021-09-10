```shell
# compile proto
protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative echo/service.proto

# run server
GRPC_GO_LOG_VERBOSITY_LEVEL=99 GRPC_GO_LOG_SEVERITY_LEVEL=info GODEBUG=http2debug=2 go run ./server/main.go

# run client
GRPC_GO_LOG_VERBOSITY_LEVEL=99 GRPC_GO_LOG_SEVERITY_LEVEL=info GODEBUG=http2debug=2 go run ./client/main.go {{ host }}
```