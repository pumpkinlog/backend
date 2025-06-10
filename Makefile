DATABASE_DSN ?= "postgres://postgres:password@localhost:5432/pumpkinlog"
RABBITMQ_URL ?= "amqp://user:password@localhost:5672"

lint:
	@golangci-lint run

test:
	@go test -v ./...

run:
	@DATABASE_DSN=$(DATABASE_DSN) RABBITMQ_URL=$(RABBITMQ_URL) go run cmd/pumpkinlog/main.go api --port 4000 --debug

migrate_up:
	@migrate -path migrations/ -database $(DATABASE_DSN)?sslmode=disable up

migrate_down:
	@migrate -path migrations/ -database $(DATABASE_DSN)?sslmode=disable down

seed:
	@DATABASE_DSN=$(DATABASE_DSN) go run cmd/pumpkinlog/main.go seed

.PHONY: lint test run migrate_up migrate_down seed all