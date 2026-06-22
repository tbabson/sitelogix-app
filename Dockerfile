# Build stage
FROM golang:1.22 AS builder

WORKDIR /app

# golang:1.22 (Debian bookworm) already includes git and ca-certificates
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -ldflags="-s -w" -o /app/bin/sitelogix ./cmd/server

# Run stage — Debian slim avoids Alpine CDN; keeps wget for healthcheck
FROM debian:bookworm-slim

RUN apt-get update && \
    apt-get install -y --no-install-recommends ca-certificates tzdata wget && \
    rm -rf /var/lib/apt/lists/* && \
    groupadd -r appgroup && useradd -r -g appgroup appuser

WORKDIR /app

COPY --from=builder /app/bin/sitelogix .
COPY --from=builder /app/migrations ./migrations

USER appuser

EXPOSE 5005

HEALTHCHECK --interval=30s --timeout=5s --start-period=10s --retries=3 \
    CMD wget -qO- http://localhost:5005/health || exit 1

ENTRYPOINT ["./sitelogix"]
