# gRPC Examples

Goland and Nodejs gRPC examples

## Components

- /go-client - contains gRPC golang client
- /go-server - contains gRPC golang server
- /node-client - contains gRPC nodejs client (NestJS)
- /node-server - contains gRPC nodejs server (Loopback)

Both nodejs and golang client communicate with golang and nodejs servers.

## generate golang proto files

1. go to proto folder
2. run:

```bash
protoc \
    --go_out=Mgrpc/service_config/service_config.proto=/internal/proto/grpc_service_config:. \
    --go-grpc_out=Mgrpc/service_config/service_config.proto=/internal/proto/grpc_service_config:. \
    --go_opt=paths=source_relative \
    --go-grpc_opt=paths=source_relative \
    random.proto
```
