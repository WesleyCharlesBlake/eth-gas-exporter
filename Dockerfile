FROM golang:1.18-alpine AS builder

WORKDIR /app

COPY . .

RUN go build -o exporter .

FROM alpine

WORKDIR /app

COPY --from=builder /app/exporter .

CMD ["/app/exporter"]