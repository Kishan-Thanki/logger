# Logger Examples

This directory contains a complete, runnable example demonstrating how to integrate the `logger` toolkit into a standard Go application.

## Running the Example

1. Start the HTTP server:
```bash
go run main.go
```

2. In a separate terminal, trigger the API endpoint to see the logger in action:
```bash
curl -s http://localhost:8080/api/login
```

## What to observe in the output:
- **Trace IDs:** You will see the exact same `trace_id` attached to both the internal application log and the final HTTP request log.
- **PII Redaction:** Notice how the `password` field in the JSON output is automatically replaced with `***REDACTED***`.
- **Structured Telemetry:** Notice how the HTTP middleware automatically nests the request and response data into beautiful `request: {}` and `response: {}` JSON blocks.
