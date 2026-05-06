# Stage 1: Build
FROM golang:1.25 AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o /app/fred-mcp ./cmd/fred-mcp

# Stage 2: Runtime
FROM scratch
COPY --from=builder /app/fred-mcp /usr/local/bin/fred-mcp
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
ENV PORT=4000
EXPOSE 4000
ENTRYPOINT ["/usr/local/bin/fred-mcp", "serve"]
