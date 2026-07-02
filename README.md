# go-http

An HTTP server library written in Go from scratch, using raw TCP connections.

## Features

- HTTP/1.0 request parsing with support for custom header and body size limits
- Radix tree-based router with static, parameterized (`:param`), and catch-all (`*path`) route patterns
- Per-method routing (GET, POST, PUT, DELETE, HEAD, CONNECT, OPTIONS, TRACE, PATCH)
- Customizable not-found handler
- Terminal colorization utility for log output

## Usage

```go
package main

import (
    "codeberg.org/JollyFrolics/go-http/http"
)

func main() {
    server := http.Default()

    server.GET("/hello", func(conn *http.ConnectionInfo) {
        conn.Write([]byte("HTTP/1.0 200 OK\r\n\r\nHello, world!"))
    })

    server.GET("/users/:id", func(conn *http.ConnectionInfo) {
        conn.Write([]byte("HTTP/1.0 200 OK\r\n\r\nUser ID: " + conn.URL))
    })

    server.Start()
}
```
