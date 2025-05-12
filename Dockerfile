FROM golang:1.23.8-alpine3.21 AS builder

WORKDIR /app

RUN go install github.com/pressly/goose/v3/cmd/goose@latest
COPY go.mod go.sum ./
RUN go mod download

COPY . .
COPY .env .env

RUN CGO_ENABLED=0 go build -o app cmd/ikakbolit/main.go

FROM alpine:3.21

WORKDIR /app
RUN apk add --no-cache ca-certificates

COPY --from=builder /app/app .

EXPOSE 8080 50051
CMD ["./app"]
