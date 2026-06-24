# Logger

`logger` is a native, highly-optimized wrapper around Go's standard library `log/slog`. It provides advanced capabilities like zero-allocation PII Redaction, automatic Trace ID propagation, and HTTP Telemetry middleware without forcing developers to learn a proprietary logging API.

By strictly utilizing `slog.JSONHandler`, the output is natively structured for seamless ingestion into any log aggregation platform or visualization tool.

## Features

- **Native `slog` Support**: Operates purely as a standard `slog.Handler`.
- **PII Redaction**: Intercepts and masks sensitive keys (e.g. `password`, `credit_card`) with zero-allocation.
- **Trace ID Context Propagation**: Automatically extracts and injects Trace IDs into deeply nested logs.
- **HTTP Middleware**: Drop-in `net/http` telemetry middleware that logs request/response durations and status codes.

## Installation

```sh
go get github.com/kishan-thanki/logger
```

## Quick Start

*For a complete, runnable application demonstrating the logger, see the [Examples Directory](examples/README.md).*

Initialize the handler and set it as your global Go logger:

```go
package main

import (
	"log/slog"
	"net/http"

	"github.com/kishan-thanki/logger"
	"github.com/kishan-thanki/logger/middleware"
)

func main() {
	// 1. Initialize the Logger Toolkit
	handler := logger.New(
		logger.WithLevel("INFO"),
		logger.WithTraceID(true),
		logger.WithRedaction("password", "token", "credit_card"),
	)

	// 2. Set it as the standard Go Logger
	slog.SetDefault(handler)

	// 3. (Optional) Wrap your HTTP Router with the Telemetry Middleware
	mux := http.NewServeMux()
	mux.HandleFunc("/api/login", func(w http.ResponseWriter, r *http.Request) {
		// Log messages natively. "password" is automatically redacted!
		slog.InfoContext(r.Context(), "Processing login request",
			slog.String("username", "admin"),
			slog.String("password", "super_secret"),
		)
	})

	loggedMux := middleware.HTTP(mux)
	http.ListenAndServe(":8080", loggedMux)
}
```

## Advanced Usage

### Context Injection (Trace IDs)

When executing background jobs or deep database calls, you can inject a Trace ID into the context. Any `slog.InfoContext` call will automatically extract and attach it to the JSON output.

```go
ctx := logger.InjectTraceID(context.Background(), "trace-1234")
slog.ErrorContext(ctx, "Database connection failed")
// Output: {"level":"ERROR","msg":"Database connection failed","trace_id":"trace-1234"}
```

## License

This project is licensed under the Apache License 2.0. See the [LICENSE](LICENSE) file for details.
