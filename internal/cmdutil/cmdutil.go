package cmdutil

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rabbitmq/amqp091-go"
)

func NewLogger(debug bool) *slog.Logger {
	opts := &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}

	if debug {
		opts.Level = slog.LevelDebug
	}

	return slog.New(slog.NewTextHandler(os.Stdout, opts))
}

func NewDatabasePool(ctx context.Context) (*pgxpool.Pool, error) {

	url := os.Getenv("DATABASE_DSN")
	if url == "" {
		return nil, fmt.Errorf("DATABASE_DSN is not set")
	}

	if os.Getenv("ENV") != "production" {
		// We disable SSL for local development
		url += "?sslmode=disable"
	}

	connStr := fmt.Sprintf(
		"%s&pool_max_conns=%d&pool_min_conns=%d",
		url,
		10,
		1,
	)

	config, err := pgxpool.ParseConfig(connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to parse connection string: %w", err)
	}

	config.MaxConnLifetime = 1 * time.Hour
	config.MaxConnIdleTime = 30 * time.Second

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, fmt.Errorf("failed to create connection pool: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return pool, nil
}

func NewDatabasePoolWithRetry(ctx context.Context, retries int) (*pgxpool.Pool, error) {

	var db *pgxpool.Pool
	var err error

	for range retries {

		db, err = NewDatabasePool(ctx)
		if err == nil {
			return db, nil
		}

		time.Sleep(5 * time.Second)
	}

	return nil, fmt.Errorf("could not connect to database after %d attempts: %w", retries, err)
}

func NewRabbitMQClient() (*amqp091.Connection, *amqp091.Channel, error) {

	url := os.Getenv("RABBITMQ_URL")
	if url == "" {
		return nil, nil, fmt.Errorf("RABBITMQ_URL is not set")
	}

	conn, err := amqp091.Dial(url)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to open channel: %w", err)
	}

	return conn, ch, nil
}
