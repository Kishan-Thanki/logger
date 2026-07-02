# Logger Middlewares

[![Go Reference](https://pkg.go.dev/badge/github.com/kishan-thanki/logger.svg)](https://pkg.go.dev/github.com/kishan-thanki/logger)
[![Go CI](https://github.com/Kishan-Thanki/logger/actions/workflows/go.yml/badge.svg)](https://github.com/Kishan-Thanki/logger/actions/workflows/go.yml)
[![License](https://img.shields.io/badge/License-Apache_2.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)

A collection of composable, high-performance `slog.Handler` middlewares. This toolkit extends Go's standard `log/slog` library with zero-allocation PII redaction, Trace ID propagation, and HTTP telemetry without wrapping or replacing the core logger.

## Philosophy

You provide a standard `slog.JSONHandler` or `slog.TextHandler`, and you snap on the exact behaviors you need. No bloat, no lock-in. 

## Packages

### 1. Context Trace Middleware (`slogctx`)

Propagates Trace IDs from `context.Context` to your log records automatically.

```go
import "github.com/kishan-thanki/logger/slogctx"

baseHandler := slog.NewJSONHandler(os.Stdout, nil)
ctxHandler := slogctx.NewHandler(baseHandler)

slog.SetDefault(slog.New(ctxHandler))

ctx := slogctx.InjectTraceID(context.Background(), "trace-123")
slog.InfoContext(ctx, "Hello") // Output automatically includes {"trace_id": "trace-123"}
```

### 2. PII Redaction (`slogredact`)

Scrub sensitive information from logs automatically before they hit the terminal.

```go
import "github.com/kishan-thanki/logger/slogredact"

baseHandler := slog.NewJSONHandler(os.Stdout, nil)
safeHandler := slogredact.NewHandler(baseHandler, "password", "token")
```
*(Note: A high-performance, zero-allocation `slogredact.ReplaceAttr` is also available for drop-in `HandlerOptions` use!)*

### 3. HTTP Telemetry (`httptelemetry`)

Drop-in `net/http` telemetry middleware that automatically measures latency, captures HTTP status codes, and securely generates and injects Trace IDs into the context.

```go
import "github.com/kishan-thanki/logger/httptelemetry"

mux := http.NewServeMux()
http.ListenAndServe(":8080", httptelemetry.Middleware(mux))
```

## The "Middleware Onion" (Quick Start)

The true power of this package is unlocked by layering the handlers together like an onion. 

```go
// 1. Core Engine
base := slog.NewJSONHandler(os.Stdout, nil)

// 2. Security Layer (Redacts passwords)
safe := slogredact.NewHandler(base, "password", "credit_card")

// 3. Transport Layer (Injects trace_ids)
ctxHandler := slogctx.NewHandler(safe)
log := slog.New(ctxHandler)

// 4. Web Layer (Generates trace_ids & captures HTTP telemetry)
mux := http.NewServeMux()
handler := httptelemetry.Middleware(mux)
```

See the [Examples Directory](examples/README.md) for a complete, runnable application demonstrating this full architecture in action!

## Installation

```sh
go get github.com/kishan-thanki/logger
```

## License

Apache License 2.0. See the [LICENSE](LICENSE) file.
