package cmd

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/pumpkinlog/backend/internal/cmdutil"
)

func SchedulerCmd(ctx context.Context) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "scheduler",
		Args:  cobra.ExactArgs(0),
		Short: "Schedules jobs and runs several maintenance tasks periodically.",
		RunE: func(cmd *cobra.Command, args []string) error {

			logger := cmdutil.NewLogger("scheduler")

			_, err := cmdutil.NewDatabasePoolWithRetry(ctx, 3)
			if err != nil {
				return err
			}

			logger.Info("Scheduler started")

			return nil
		},
	}

	return cmd
}
