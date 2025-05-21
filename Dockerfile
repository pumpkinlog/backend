FROM golang:1.24.3-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o /app/pumpkinlog cmd/pumpkinlog/main.go

FROM alpine:3.21.3

WORKDIR /app

COPY --from=builder /app/pumpkinlog /app/pumpkinlog
COPY ./docs /app/docs

ENTRYPOINT ["/app/pumpkinlog"]