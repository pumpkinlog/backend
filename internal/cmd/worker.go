package cmd

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/pumpkinlog/backend/internal/cmdutil"
	"github.com/pumpkinlog/backend/internal/worker"
)

var queues = map[string]worker.NewWorkerFn{
	"presence": worker.NewEvaluationWorker,
}

func WorkerCmd(ctx context.Context) *cobra.Command {
	var queue string
	var concurrency int

	cmd := &cobra.Command{
		Use:   "worker",
		Args:  cobra.ExactArgs(0),
		Short: "Work through the queue.",
		RunE: func(cmd *cobra.Command, args []string) error {

			if queue == "" {
				return fmt.Errorf("queue is required")
			}

			debug, err := cmd.Flags().GetBool("debug")
			if err != nil {
				return err
			}
			logger := cmdutil.NewLogger(debug)

			db, err := cmdutil.NewDatabasePoolWithRetry(ctx, 3)
			if err != nil {
				return err
			}

			_, ch, err := cmdutil.NewRabbitMQClient()
			if err != nil {
				return err
			}

			workerFn, ok := queues[queue]
			if !ok {
				return fmt.Errorf("unknown queue: %s", queue)
			}

			w := workerFn(logger, db, ch, 5)

			if err := w.Start(); err != nil {
				return fmt.Errorf("failed to start worker: %w", err)
			}

			logger.Info("worker started", "queue", queue, "concurrency", concurrency)
			logger.Debug("debug")

			<-ctx.Done()

			if err := w.Stop(); err != nil {
				return fmt.Errorf("failed to stop worker: %w", err)
			}

			logger.Info("worker stopped", "queue", queue)

			return nil
		},
	}

	cmd.Flags().StringVar(&queue, "queue", "", "Queue to consume from")
	cmd.Flags().IntVar(&concurrency, "concurrency", 5, "Number of concurrent workers")

	return cmd
}
