package interfaces

import (
	"github.com/go-kratos/kratos/v2/transport"
)

// Server 是我们框架内所有服务类型的顶层抽象。
// 它通过内嵌 transport.Server，确保了任何实现了我们 Server 接口的类型，
// 同时也自动满足了 Kratos App 所需的 transport.Server 接口。
type Server interface {
	transport.Server // <-- 核心：确保与 Kratos App 的完全兼容
}

// Client 是一个标记接口，代表一个客户端连接实例，例如 *grpc.ClientConn。
// 由于不同协议的客户端（如 gRPC, HTTP）没有统一的接口，
// 我们使用一个空接口来提供灵活性，调用方需要进行类型断言。
type Client interface{}
