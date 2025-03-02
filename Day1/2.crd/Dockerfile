FROM golang:1.23-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git

COPY go.mod go.sum ./
RUN go mod tidy

COPY main.go ./

RUN go get -d ./...
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o secret-operator main.go

FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/secret-operator .

RUN chmod +x secret-operator

CMD ["./secret-operator"]