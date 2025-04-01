# Builder stage
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o /app/server ./cmd/main

# Final stage
FROM alpine:latest
RUN adduser -D -g '' appuser
WORKDIR /app
COPY --from=builder /app/server .
RUN chown -R appuser:appuser /app
USER appuser
EXPOSE 8080
CMD ["./server"]