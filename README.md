# RainbowFerret

A lightweight, idiomatic HTTP framework for Go, built directly on the enhanced
[`net/http.ServeMux`](https://pkg.go.dev/net/http#ServeMux) introduced in Go 1.22+.
It brings middleware chaining, sub-group nesting, automatic 404/405 distinction,
generic handlers, and a set of production-ready built-in middleware — all
without pulling in a third-party router.

[![Go Version](https://img.shields.io/badge/Go-1.22%2B-00ADD8?logo=go)](https://go.dev/)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)

[[中文文档](README_zh_cn.md)]

## Features

- **Zero router dependency** — sits on top of `net/http.ServeMux` and uses its
  native method-prefixed patterns (`GET /path`, `POST /path`) and path wildcards
  (`{id}`).
- **Generic handlers** — `HandlerFunc[T]` functions receive a typed request
  struct, return a typed response and an error; the framework handles JSON
  decoding/encoding automatically.
- **Middleware chain** — `Middleware` is the familiar `func(http.Handler)
  http.Handler`. Use `Chain` to compose multiple middleware or attach them
  per-group / per-route.
- **Sub-groups with prefix inheritance** — `Group("/prefix")` creates a child
  group that inherits the parent's middleware and shares the underlying mux, so
  cross-group middleware ordering is predictable.
- **404 vs 405 distinction** — the framework tracks registered paths separately
  from method mappings. A request to a known path with an unrecognised method
  receives a `405 Method Not Allowed` with a populated `Allow` header; a truly
  unknown path gets `404 Not Found`.
- **Param binding** — `Bind` automatically unmarshals JSON, XML, and form
  data into Go structs (similar to Gin's `ShouldBind`), while `PathParam` /
  `QueryParam` helpers provide easy access to URL parameters.
- **Response rendering** — convenience functions for JSON, XML, plain text, and
  streaming output.

## Installation

```bash
go get github.com/rambollwong/rainbowferret
```

Requires **Go 1.22 or later**.

## Quick Start

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

### Sub-groups

```go
root := ferret.NewRootGroup(sm, middleware.Logger())
api  := root.Group("/api", middleware.Timeout(5*time.Second))
v1   := api.Group("/v1")

v1.Get("/users", func(w http.ResponseWriter, r *http.Request) {
    w.Write([]byte("[user1, user2]"))
})
// Registered as: GET /api/v1/users
```

### Generic handlers

```go
type CreateReq struct {
    Name  string `json:"name"`
    Email string `json:"email"`
}

ferret.PostHandler(v1, "/users", func(ctx context.Context, req CreateReq) (any, error) {
    // req is auto-decoded from JSON body, response auto-encoded
    return map[string]string{"id": "42", "name": req.Name}, nil
})
```

### Static files

```go
// Local directory
root.Static("/assets", "./public")

// Embed with //go:embed
//go:embed dist
var distFS embed.FS
sub, _ := fs.Sub(distFS, "dist")
root.StaticFS("/", sub)
```

## Core Concepts

### Router interface

`ferret.Router` is the central abstraction. It declares methods for registering
handlers by HTTP verb (`Get`, `Post`, `Put`, `Delete`, …), attaching middleware
(`Use`), creating sub-groups (`Group`), and customising 404/405 handlers
(`NotFound`, `MethodNotAllowed`). The single concrete implementation is `*Group`.

### Group

A `Group` wraps a `http.ServeMux` and manages:

- **Middleware stack** — group-level middleware is prepended to every route
  registered under that group.
- **Prefix** — child groups created via `Group("/v1")` automatically prepend
  `"/v1"` to all their routes.
- **Shared mux & route tree** — all groups derived from the same root share one
  underlying mux and a path-only tree used for accurate 405 detection.

Create a root group with `ferret.NewRootGroup(sm, middleware...)`, then call
`root.Group("/prefix", middleware...)` to create a child.

### Middleware

```go
type Middleware func(next http.Handler) http.Handler
```

Middleware is applied in FIFO order (first added = outermost, runs first on the
way in and last on the way out). Use `ferret.Chain(ms...)` to compose multiple
middleware into a single `Middleware` value. Middleware can be attached at three
levels:

1. **Root group** — applies to every route in the application.
2. **Sub-group** — applies to all routes under that prefix.
3. **Per-route** — the final argument to `Get`/`Post`/etc., runs after group
   middleware.

### Generic handlers

```go
type HandlerFunc[T any] func(ctx context.Context, req T) (res any, err error)
```

`HandlerFunc[T]` encapsulates a three-step flow: **decode → execute → encode**.
The framework decodes the request body into type `T`, calls your function, and
writes the result as JSON. Returning `nil, nil` produces `204 No Content`.

Helper registration functions — `ferret.GetHandler`, `ferret.PostHandler`,
`ferret.DeleteHandler`, etc. — wire a `HandlerFunc[T]` directly onto a group.

## Built-in Middleware

| Middleware | Signature | Description |
| ---------- | --------- | ----------- |
| **Recoverer** | `middleware.Recoverer` | Recovers from panics, logs the stack trace, and responds with `500` (or the code carried by `*types.HTTPError`). |
| **Logger** | `middleware.Logger()` | Logs every request: method, path, status, duration, body size. Configurable via `LoggerWithConfig(cfg)`. |
| **RequestID** | `middleware.RequestID()` | Ensures a unique request ID on every request. Reads from `X-Request-ID` header (for trace propagation) or generates one with a hostname prefix. Uses `crypto/rand`. |
| **Timeout** | `middleware.Timeout(d)` | Cancels the request context after the given duration. Sends `504 Gateway Timeout` if the handler hasn't responded. |
| **ClientIP** | `middleware.ClientIP()` | Resolves the client IP from `RemoteAddr`, with optional proxy-header support (`TrustProxy`). |
| **ContentType** | `middleware.ContentType("application/json", ...)` | Rejects requests whose `Content-Type` is not in the allowed list with `415`. Skips body-less methods (GET, HEAD, …). |
| **ContentCharset** | `middleware.ContentCharset("utf-8")` | Rejects requests with an unsupported charset parameter in `Content-Type`. |
| **ContentEncoding** | `middleware.ContentEncoding("identity")` | Rejects requests with an unsupported `Content-Encoding` header. |
| **Compress** | `middleware.Compress(level ...)` | Compresses response bodies with gzip or deflate, negotiated via `Accept-Encoding`. |

All middleware follow the standard `func(http.Handler) http.Handler` signature
and can be combined freely with `ferret.Chain`.

## Utility Packages

### `util` — helpers for handlers

| Function | Description |
| -------- | ----------- |
| `HandleT(handlerFn)` | Wraps a `HandlerFunc[T]` into a standard `http.HandlerFunc`. |
| `Bind(r, &v)` | Auto-binds the request body to a struct based on `Content-Type` (JSON, XML, form). |
| `DecodeJSON / DecodeXML / DecodeForm` | Decode the request body for JSON, XML, or form data. |
| `WriteJSON / WriteXML / WriteText / WriteNoContent` | Write common response formats. |
| `WriteStream(w, code, contentType, reader)` | Stream arbitrary data to the response (useful for file downloads or SSE). |
| `PathParam(r, name)` / `QueryParam(r, name)` | Read path wildcards and query parameters with typed variants (`Int`, `Float`, `Bool`, `Int64`). |
| `Validator` interface | If the bound struct implements `Validate() error`, it is called automatically after binding. |

### `types` — shared types

| Type | Description |
| ---- | ----------- |
| `HTTPError` | A standard error carrying an HTTP status code. |
| `Handler[T]` / `HandlerFunc[T]` | The generic handler interface and function type. |
| `BadRequest / NotFound / Internal / …` | Factory functions for common `HTTPError` status codes. |

## Examples

| Example | Command |
| ------- | ------- |
| Hello World | `go run _examples/hello-world/main.go` |
| Static files | `go run _examples/static-files/main.go` |
| Sub-groups & middleware | `go run _examples/sub-group/main.go` |
| Middleware combo | `go run _examples/middleware/main.go` |
| REST API with generic handlers | `go run _examples/rest-api/main.go` |

The **hello-world** example shows minimal usage: a root group with two GET
routes and standard `http.HandlerFunc` handlers.

The **sub-group** example shows sub-group nesting with inherited
middleware (Logger → Timeout → RequestID).

The **middleware** example combines all built-in middleware into a single
production-ready service.

The **rest-api** example demonstrates middleware (`Logger`, `RequestID`,
`ContentType`), sub-groups (`/api/v1`), and generic `HandlerFunc` handlers that
automatically decode JSON requests and encode JSON responses.

## License

[MIT](LICENSE)
