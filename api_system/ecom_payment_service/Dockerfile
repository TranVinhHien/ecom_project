# --- Stage 1: Builder ---
FROM golang:1.24-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git tzdata ca-certificates

COPY go.mod go.sum ./
RUN go mod download
COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-w -s" \
    -o main main.go

# --- Stage 2: Runner ---
FROM alpine:3.20
WORKDIR /app
# Copy các thứ cần thiết từ Builder
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /app/main .

COPY db/migration ./db/migration

# Setup User (Security)
RUN addgroup -g 1000 appgroup && \
    adduser -D -u 1000 -G appgroup appuser && \
    chown -R appuser:appgroup /app

USER appuser

EXPOSE 9003

CMD ["/app/main"]