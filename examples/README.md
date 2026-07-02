# Logger Middlewares: The Unified Example

This directory contains the **Ultimate Example** demonstrating how to integrate all three composable `slog` packages (`slogctx`, `slogredact`, and `httptelemetry`) into a standard Go application simultaneously.

## The Middleware Onion

In `main.go`, we construct a powerful, layered logging engine:
1. **Core Engine**: `slog.NewJSONHandler` formats and prints the final JSON.
2. **Security Layer**: `slogredact.NewHandler` intercepts logs and safely scrubs sensitive keys (like `password`).
3. **Transport Layer**: `slogctx.NewHandler` automatically extracts Trace IDs from the context and injects them into the logs.
4. **Web Layer**: `httptelemetry.Middleware` intercepts HTTP requests, generates Trace IDs, sets the stopwatch, and logs the final HTTP response telemetry.

## Running the Example

1. Start the HTTP server:
```bash
go run ./
```

2. In a separate terminal, trigger the API endpoint with a fake password:
```bash
curl -X POST -d '{"email":"test@test.com", "password":"mysecret"}' http://localhost:8080/login
```

## What to observe in your server terminal:
1. **The Application Log:** You will see a `login attempt` log. Notice that `"password"` has been safely changed to `"***REDACTED***"` and a `"trace_id"` has been automatically appended!
2. **The Telemetry Log:** A microsecond later, you will see a massive `HTTP Request` log. It contains the exact same `"trace_id"`, plus automatically captured data like the HTTP status code (`200`) and the request latency!
