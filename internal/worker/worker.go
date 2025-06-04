package worker

import (
	"log/slog"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rabbitmq/amqp091-go"
)

type NewWorkerFn func(logger *slog.Logger, conn *pgxpool.Pool, ch *amqp091.Channel, concurrency int) Worker

type Worker interface {
	Start() error
	Stop() error
}
