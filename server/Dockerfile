# ── Stage 1: Build ─────────────────────────────────────────────────────────
FROM golang:1.21-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o heartchat-server .

# ── Stage 2: Minimal runtime image ─────────────────────────────────────────
FROM alpine:3.19

# ca-certificates required for TLS calls to DeepSeek API and WeChat APIs
RUN apk --no-cache add ca-certificates

WORKDIR /app

COPY --from=builder /app/heartchat-server .

# Persistent data directory for per-user state (mount a volume here in prod)
RUN mkdir -p /app/data
ENV DATA_DIR=/app/data

EXPOSE 8080

CMD ["./heartchat-server"]
