## Go Support for Protocol Buffers

本示例演示了新版本 Go Protobuf 的使用，以及在新版本下使用 `proto` 的 `import` 功能，和编译 gRPC 服务。

自 Go Protobuf 升级到 `second major revision` 之后 (仓库迁移到了 [https://github.com/protocolbuffers/protobuf-go](https://github.com/protocolbuffers/protobuf-go) ，Go Module 地址为 [google.golang.org/protobuf](https://pkg.go.dev/mod/google.golang.org/protobuf) )，编译 gRPC 的 proto 文件需要一下几个步骤：

1. 需要安装 `protoc`、`protoc-gen-go` 和 `protoc-gen-go-grpc`
2. 升级之后，`service` 代码和 `message` 代码需要单独编译：
    - 工作目录切换到项目根目录下 ( `protobuf-go` 目录下 )
    - 编译 message：`protoc --proto_path=. --go_out=paths=source_relative:. **/*.proto` ( `protoc` 使用了 `protoc-gen-go` 进行编译 )
    - 编译 service：`protoc --proto_path=. --go-grpc_out=paths=source_relative:. **/*.proto` ( `protoc` 使用了 `protoc-gen-go-grpc` 进行编译 )
3. 或者通过一个命令来完成：`protoc --proto_path=. --go_out=paths=source_relative:. --go-grpc_out=paths=source_relative:. **/*.proto` (新版本支持多包同时编译了)

参考：

- [新版本 Go Protobuf 官方博客声明](https://blog.golang.org/protobuf-apiv2)
- [protoc 安装](https://github.com/protocolbuffers/protobuf)
- [新版 protoc-gen-go 安装](https://github.com/protocolbuffers/protobuf-go/tree/master/cmd/protoc-gen-go)
- [protoc-gen-go-grpc 安装](https://github.com/grpc/grpc-go/tree/master/cmd/protoc-gen-go-grpc)