package cmd

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/pumpkinlog/backend/internal/cmdutil"
	"github.com/pumpkinlog/backend/internal/seed"
)

func SeedCmd(ctx context.Context) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "seed",
		Args:  cobra.ExactArgs(0),
		Short: "Seed the database with initial data.",
		RunE: func(cmd *cobra.Command, args []string) error {

			logger := cmdutil.NewLogger("seeder")

			db, err := cmdutil.NewDatabasePoolWithRetry(ctx, 3)
			if err != nil {
				return err
			}

			seed.NewSeeder(logger, db).Seed(ctx)

			return nil
		},
	}

	return cmd
}
