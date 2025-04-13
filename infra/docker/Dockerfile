FROM golang:1.20.3-alpine as builder

COPY ../.. /github.com/bovinxx/auth-service/source
WORKDIR /github.com/bovinxx/auth-service/source

RUN go mod download
RUN go build -o ./bin/auth_service cmd/grpc_server/main.go

FROM alpine:latest

WORKDIR /root/
COPY --from=builder /github.com/bovinxx/auth-service/source/bin/auth_service .

CMD ["./auth_service"]