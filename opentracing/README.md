# opentracing

OpenTracing 是一套 Tracing 规范，它将分布式链路中的一些事物和行为概念化、模型化，并抽象成与之对应的 `interface`，来进行编程约束，统一业务代码。

开发者只需将实现了 OpenTracing 的 Tracing 系统（如 Zipkin, Jaeger 等）的客户端 `tracer` 实例注入到 `opentracing` 的全局 `tracer` 接口，之后只需要调用 `opentracing` 的接口方法就可以进行链路跟踪的相关行为，大大降低了开发者的心智负担。

## 基本概念

- `tracer`
- `trace`
- `span`(`tag`, `log`)
- `spanContext`
- `carrier`
- ...

这些概念请参考如下链接，这里不做重复阐述：

- [OpenTracing 官网](https://opentracing.io/)
- [OpenTracing 中文文档](https://wu-sheng.gitbooks.io/opentracing-io/content/)
- [OpenTracing API for Go](https://pkg.go.dev/github.com/opentracing/opentracing-go?tab=overview)
- [OpenTracing Go 代码仓库](https://github.com/opentracing/opentracing-go)

## 项目结构

- `http` 目录是一个 HTTP 服务，在同一个进程中演示了跨进程的网络调用和链路跟踪。
- `rpc` 目录演示了在 gRPC 调用中使用 OpenTracing 做链路跟踪，其中，`client` 目录是 RPC 客户端，`server` 是 RPC 服务端，`interceptors` 目录中含有一些客户端和服务端的 gRPC 拦截器。
- `utils` 目录是共用的一些函数。

## 链路跟踪思路

在请求的调用链路中，`span` 代表链路中的一个环节。在服务调用过程中，客户端向服务端发起请求时，需将客户端 `span` 信息 (`spanContext`) 附着于 `carrier`(可以是 `HTTP Headers` 或者 gRPC 的 `metadata` 等) 之上，传送给服务端。服务端接收到请求，首先判断请求的 `carrier` 中是否携带客户端的 `spanContext`，有则继承它 (`childOf`) 并创建 `span`，无则创建一个 `root span`。 

每一个 `span` 都会向 Tracing 系统上报信息(`span.Finish()`)，由于 `span` 之间存在继承关系，故 Tracing 系统可以聚合统计出调用链路。