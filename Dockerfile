FROM golang:1.24.5-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o proxy-server ./cmd/proxy-server

FROM alpine:3.19
WORKDIR /app
COPY --from=builder /app/proxy-server /app/proxy-server
COPY configs/ /app/configs/
COPY plugins/ /app/plugins/
COPY configs/example.yaml /app/config.yaml
EXPOSE 8080 9000 9001 9100 9090
ENTRYPOINT ["/app/proxy-server", "/app/config.yaml"] 