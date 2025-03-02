FROM golang:1.20 AS builder

WORKDIR /app

COPY main.go .

RUN go build -o secret-operator main.go

FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/secret-operator .

CMD ["./secret-operator"]