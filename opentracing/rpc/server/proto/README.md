## 编译命令

自 `protobuf-go` 升级到 `second major revision` 之后（仓库改为了[https://github.com/protocolbuffers/protobuf-go](https://github.com/protocolbuffers/protobuf-go)，module 改为了 `google.golang.org/protobuf`），编译 gRPC 的 proto 文件需要一下几个步骤：

1. 需要安装 `protoc`、`protoc-gen-go` 和 `protoc-gen-go-grpc`
2. 升级之后，`service` 代码和 `message` 代码需要单独编译：
    - 工作目录切换到 `server` 下
    - 编译 message：`protoc -I=. --go_out=paths=source_relative:. proto/*.proto`
    - 编译 service：`protoc -I=. --go-grpc_out=paths=source_relative:. proto/*.proto`
3. 或者通过一个命令来完成：`protoc -I=. --go_out=paths=source_relative:. --go-grpc_out=paths=source_relative:. proto/*.proto`