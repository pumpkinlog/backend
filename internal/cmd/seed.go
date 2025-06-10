package cmd

import (
	"context"
	"errors"

	"github.com/spf13/cobra"

	"github.com/pumpkinlog/backend/internal/cmdutil"
	"github.com/pumpkinlog/backend/internal/seed"
)

func SeedCmd(ctx context.Context) *cobra.Command {
	var file string

	cmd := &cobra.Command{
		Use:   "seed",
		Args:  cobra.ExactArgs(0),
		Short: "Seed the database with initial data.",
		RunE: func(cmd *cobra.Command, args []string) error {
			if file == "" {
				return errors.New("file is required")
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

			return seed.NewSeeder(logger, db).Seed(ctx, file)
		},
	}

	cmd.Flags().StringVar(&file, "file", "/seed.json", "The JSON data file to seed with")

	return cmd
}
