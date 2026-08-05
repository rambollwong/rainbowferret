# RainbowFerret

RainbowFerret 是一个轻量级的 Go HTTP 框架，直接构建在 Go 1.22+ 增强版
[`net/http.ServeMux`](https://pkg.go.dev/net/http#ServeMux) 之上。它提供了
中间件链、子分组嵌套、自动 404/405 区分、泛型业务处理器，以及一系列开箱即用的
内置中间件——无需引入任何第三方路由库。

[![Go Version](https://img.shields.io/badge/Go-1.22%2B-00ADD8?logo=go)](https://go.dev/)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)

## 特性

- **零路由依赖** — 基于标准库 `net/http.ServeMux`，直接使用原生的方法前缀模式
  （`GET /path`、`POST /path`）和路径通配符（`{id}`）。
- **泛型逻辑处理器** — `LogicFunc[T]` 函数接收一个类型化的请求结构体，返回响应
  和错误；框架自动完成 JSON 解码与编码。
- **中间件链** — `Middleware` 即经典的 `func(http.Handler) http.Handler` 模式。
  通过 `Chain` 组合多个中间件，或按分组/路由单独挂载。
- **前缀子分组** — `Group("/prefix")` 创建子分组，继承父分组的中间件并共享底层
  mux，中间件执行顺序在分组树中可预测。
- **404 与 405 区分** — 框架分别记录已注册的路径和方法映射。对已知路径使用未注册
  方法发起的请求返回 `405 Method Not Allowed`（含 `Allow` 头）；真正未知的路径
  返回 `404 Not Found`。
- **参数绑定** — `Bind` 根据 `Content-Type` 自动将 JSON、XML 和表单数据
  反序列化为 Go 结构体（类似 Gin 的 `ShouldBind`）；`PathParam` / `QueryParam`
  系列辅助函数便捷读取 URL 参数。
- **响应渲染** — 提供 JSON、XML、纯文本以及流式输出的便捷函数。

## 安装

```bash
go get github.com/rambollwong/rainbowferret
```

需要 **Go 1.22 及以上**版本。

## 快速开始

```go
package main

import (
    "net/http"

    "github.com/rambollwong/rainbowferret/ferret"
)

func main() {
    sm := http.NewServeMux()
    root := ferret.NewRootGroup(sm)

    root.Get("/", func(w http.ResponseWriter, r *http.Request) {
        w.Write([]byte("Hello, World!"))
    })

    root.Get("/hello/{name}", func(w http.ResponseWriter, r *http.Request) {
        name := r.PathValue("name")
        w.Write([]byte("Hello, " + name + "!"))
    })

    http.ListenAndServe(":8080", root.Handler())
}
```

### 子分组

```go
root := ferret.NewRootGroup(sm, middleware.Logger())
api  := root.Group("/api", middleware.Timeout(5*time.Second))
v1   := api.Group("/v1")

v1.Get("/users", func(w http.ResponseWriter, r *http.Request) {
    w.Write([]byte("[user1, user2]"))
})
// 注册路径为: GET /api/v1/users
```

### 泛型逻辑处理器

```go
type CreateReq struct {
    Name  string `json:"name"`
    Email string `json:"email"`
}

ferret.PostLogic(v1, "/users", func(ctx context.Context, req CreateReq) (any, error) {
    // 请求体自动解码为 CreateReq，响应自动 JSON 编码
    return map[string]string{"id": "42", "name": req.Name}, nil
})
```

## 核心概念

### Router 接口

`ferret.Router` 是框架的核心抽象，声明了按 HTTP 动词注册处理器（`Get`、`Post`、
`Put`、`Delete` 等）、挂载中间件（`Use`）、创建子分组（`Group`）以及自定义
404/405 处理器（`NotFound`、`MethodNotAllowed`）等方法。唯一的实现类型是 `*Group`。

### Group（分组）

`Group` 封装了 `http.ServeMux`，管理以下内容：

- **中间件栈** — 分组级中间件会追加到该分组下每条路由的前面。
- **路径前缀** — 通过 `Group("/v1")` 创建的子分组会自动为其所有路由添加 `"/v1"` 前缀。
- **共享 mux 与路由树** — 同根派生的所有分组共享同一个底层 mux 和纯路径树，
  用于精确的 405 检测。

使用 `ferret.NewRootGroup(sm, middleware...)` 创建根分组，再通过
`root.Group("/prefix", middleware...)` 创建子分组。

### 中间件

```go
type Middleware func(next http.Handler) http.Handler
```

中间件按 FIFO 顺序执行（先添加的最外层，请求进入时最先执行、响应返回时最后执行）。
使用 `ferret.Chain(ms...)` 将多个中间件组合为一个。中间件可在三个层级挂载：

1. **根分组** — 作用于应用中的所有路由。
2. **子分组** — 作用于该前缀下的所有路由。
3. **单条路由** — 作为 `Get`/`Post` 等的最后一个参数，在分组中间件之后执行。

### 泛型逻辑处理器

```go
type LogicFunc[T any] func(ctx context.Context, req T) (res any, err error)
```

`LogicFunc[T]` 封装了 **解码 → 执行 → 编码** 三步流程：框架将请求体解码为类型
`T`，调用你的业务函数，最后将结果以 JSON 写入响应。返回 `nil, nil` 时自动产生
`204 No Content`。

辅助注册函数 — `ferret.GetLogic`、`ferret.PostLogic`、`ferret.DeleteLogic` 等 —
可直接将 `LogicFunc[T]` 注册到分组上。

## 内置中间件

| 中间件 | 签名 | 说明 |
| ------ | ---- | ---- |
| **Recoverer** | `middleware.Recoverer` | 从 panic 中恢复，记录堆栈，返回 `500`（若 panic 值为 `*types.HTTPError` 则使用其状态码）。 |
| **Logger** | `middleware.Logger()` | 记录每个请求的日志：方法、路径、状态码、耗时、响应体大小。可通过 `LoggerWithConfig(cfg)` 配置。 |
| **RequestID** | `middleware.RequestID()` | 为每个请求确保唯一 ID。从 `X-Request-ID` 头读取（支持分布式追踪透传），或生成带主机名前缀的 ID。使用 `crypto/rand`。 |
| **Timeout** | `middleware.Timeout(d)` | 在指定时长后取消请求 context。若 handler 尚未响应则返回 `504 Gateway Timeout`。 |
| **ClientIP** | `middleware.ClientIP()` | 从 `RemoteAddr` 解析客户端 IP，可开启 `TrustProxy` 支持代理头。 |
| **ContentType** | `middleware.ContentType("application/json", ...)` | 拒绝 `Content-Type` 不在允许列表中的请求，返回 `415`。跳过无请求体的方法（GET、HEAD 等）。 |
| **ContentCharset** | `middleware.ContentCharset("utf-8")` | 拒绝 `Content-Type` 中 charset 参数不合规的请求。 |
| **ContentEncoding** | `middleware.ContentEncoding("identity")` | 拒绝 `Content-Encoding` 头不合规的请求。 |
| **Compress** | `middleware.Compress(level ...)` | 通过 `Accept-Encoding` 协商，使用 gzip 或 deflate 压缩响应体。 |

所有中间件均遵循标准的 `func(http.Handler) http.Handler` 签名，可通过
`ferret.Chain` 自由组合。

## 工具包

### `util` — 处理器辅助函数

| 函数 | 说明 |
| ---- | ---- |
| `HandleT(logicFn)` | 将 `LogicFunc[T]` 包装为标准 `http.HandlerFunc`。 |
| `Bind(r, &v)` | 根据 `Content-Type` 自动将请求体绑定到结构体（JSON、XML、表单）。 |
| `DecodeJSON / DecodeXML / DecodeForm` | 为 JSON、XML 或表单数据解码请求体。 |
| `WriteJSON / WriteXML / WriteText / WriteNoContent` | 写入常见格式的响应。 |
| `WriteStream(w, code, contentType, reader)` | 流式输出响应（适用于文件下载或 SSE）。 |
| `PathParam(r, name)` / `QueryParam(r, name)` | 读取路径通配符和查询参数，提供 `Int`、`Float`、`Bool`、`Int64` 等类型变体。 |
| `Validator` 接口 | 若绑定的结构体实现了 `Validate() error`，绑定后自动调用校验。 |

### `types` — 共享类型

| 类型 | 说明 |
| ---- | ---- |
| `HTTPError` | 携带 HTTP 状态码的标准错误类型。 |
| `Logic[T]` / `LogicFunc[T]` | 泛型处理器的接口与函数类型。 |
| `BadRequest / NotFound / Internal / …` | 常用 `HTTPError` 的工厂函数。 |

## 示例

| 示例 | 运行命令 |
| ---- | -------- |
| Hello World | `go run _examples/hello-world/main.go` |
| 子分组与中间件 | `go run _examples/sub-group/main.go` |
| 中间件组合 | `go run _examples/middleware/main.go` |
| REST API（使用 Logic 处理器） | `go run _examples/rest-api/main.go` |

**hello-world** 示例展示了最简用法：一个根分组配两条 GET 路由，使用标准
`http.HandlerFunc` 处理器。

**sub-group** 示例展示了子分组嵌套和中间件继承（Logger → Timeout → RequestID）。

**middleware** 示例将所有内置中间件组合为一个生产就绪的服务。

**rest-api** 示例展示了中间件（`Logger`、`RequestID`、`ContentType`）、
子分组（`/api/v1`）以及泛型 `LogicFunc` 处理器 — 请求体自动解码为 JSON，
响应自动 JSON 编码。

## 开源协议

[MIT](LICENSE)
