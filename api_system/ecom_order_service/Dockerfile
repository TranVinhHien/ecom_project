FROM golang:1.24-alpine AS builder
# thêm timezone 
RUN apk add --no-cache tzdata


WORKDIR /app
COPY . .
RUN go build -o main main.go

FROM alpine:3.20
# Copy timezone data từ builder
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo

WORKDIR /app
COPY --from=builder /app/main .
COPY app.env .
COPY db/migration ./db/migration
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo

EXPOSE 9002
CMD ["/app/main"]
