FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o main main.go

FROM alpine:3.20
WORKDIR /app
COPY --from=builder /app/main .
COPY app.env .
COPY db/migration ./db/migration

EXPOSE 9002
CMD ["/app/main"]
