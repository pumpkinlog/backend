package worker

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"sync"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rabbitmq/amqp091-go"

	"github.com/pumpkinlog/backend/internal/domain"
	"github.com/pumpkinlog/backend/internal/engine"
	"github.com/pumpkinlog/backend/internal/repository"
	"github.com/pumpkinlog/backend/internal/service"
)

type evaluationWorker struct {
	logger      *slog.Logger
	ch          *amqp091.Channel
	queue       string
	concurrency int

	engine         *engine.Engine
	evaluationSvc  domain.EvaluationService
	presenceRepo   domain.PresenceRepository
	evaluationRepo domain.EvaluationRepository

	msgs    <-chan amqp091.Delivery
	sem     chan struct{}
	wg      sync.WaitGroup
	stopped chan struct{}
}

func NewEvaluationWorker(logger *slog.Logger, conn *pgxpool.Pool, ch *amqp091.Channel, concurrency int) Worker {
	return &evaluationWorker{
		logger:      logger,
		ch:          ch,
		queue:       "presence.worker",
		concurrency: concurrency,

		engine:         engine.NewEngine(),
		evaluationSvc:  service.NewEvaluationService(logger, conn, ch),
		presenceRepo:   repository.NewPostgresPresenceRepository(conn),
		evaluationRepo: repository.NewPostgresEvaluationRepository(conn),

		sem:     make(chan struct{}, concurrency),
		stopped: make(chan struct{}),
	}
}

func (w *evaluationWorker) Start() error {
	if err := w.ch.Qos(w.concurrency, 0, false); err != nil {
		return fmt.Errorf("set QoS: %w", err)
	}

	if err := w.ch.QueueBind(w.queue, "presence.create", "presence", false, nil); err != nil {
		return fmt.Errorf("bind queue: %w", err)
	}

	hostname, err := os.Hostname()
	if err != nil {
		return fmt.Errorf("get hostname: %w", err)
	}

	msgs, err := w.ch.Consume(
		w.queue,
		hostname,
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("consume queue: %w", err)
	}

	w.msgs = msgs

	go w.run()

	return nil
}

func (w *evaluationWorker) run() {
	for msg := range w.msgs {
		w.sem <- struct{}{}
		w.wg.Add(1)

		go func(m amqp091.Delivery) {
			defer w.wg.Done()
			defer func() { <-w.sem }()

			if err := w.handleMessage(m); err != nil {
				w.logger.Error("consume error", "error", err)
				if nackErr := m.Nack(false, true); nackErr != nil {
					w.logger.Error("failed to nack", "error", nackErr)
				}
				return
			}

			if ackErr := m.Ack(false); ackErr != nil {
				w.logger.Error("failed to ack", "error", ackErr)
			}
		}(msg)
	}

	w.wg.Wait()
	close(w.stopped)
}

func (w *evaluationWorker) Stop() error {

	if err := w.ch.Cancel("", false); err != nil {
		w.logger.Warn("cancel consumer failed", "error", err)
	}

	<-w.stopped

	if err := w.ch.Close(); err != nil {
		return fmt.Errorf("close channel: %w", err)
	}

	return nil
}

type presenceMessage struct {
	UserID   int64           `json:"userId"`
	RegionID domain.RegionID `json:"regionId"`
}

func (w *evaluationWorker) handleMessage(msg amqp091.Delivery) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var params presenceMessage
	if err := json.Unmarshal(msg.Body, &params); err != nil {
		return fmt.Errorf("unmarshal message: %w", err)
	}

	w.logger.Info("processing message", "userId", params.UserID, "regionId", params.RegionID)

	opts := &domain.EvaluateOpts{
		Recompute: true,
		Cache:     true,
	}

	if _, err := w.evaluationSvc.EvaluateRegion(ctx, params.UserID, params.RegionID, opts); err != nil {
		return fmt.Errorf("evaluate region: %w", err)
	}

	w.logger.Info("processed message", "userId", params.UserID, "regionId", params.RegionID)

	return nil
}
