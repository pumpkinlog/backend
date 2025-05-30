package worker

import (
	"log/slog"

	"github.com/rabbitmq/amqp091-go"

	"github.com/pumpkinlog/backend/internal/repository"
)

type NewWorkerFn func(logger *slog.Logger, conn repository.Connection, ch *amqp091.Channel, concurrency int) Worker

type Worker interface {
	Start() error
	Stop() error
}
