FROM golang:1.23.4-alpine3.21 AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -ldflags="-s -w" -o main cmd/main.go

FROM alpine:3.21

WORKDIR /app

COPY --from=builder /app/main .

CMD ["./main"]