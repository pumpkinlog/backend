package cmd

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/pumpkinlog/backend/internal/api"
	"github.com/pumpkinlog/backend/internal/cmdutil"
)

func APICmd(ctx context.Context) *cobra.Command {
	var port int

	cmd := &cobra.Command{
		Use:   "api",
		Args:  cobra.ExactArgs(0),
		Short: "Runs the RESTful API.",
		RunE: func(cmd *cobra.Command, args []string) error {

			debug, err := cmd.Flags().GetBool("debug")
			if err != nil {
				return err
			}
			logger := cmdutil.NewLogger(debug)

			db, err := cmdutil.NewDatabasePoolWithRetry(ctx, 3)
			if err != nil {
				return err
			}
			defer db.Close()

			_, ch, err := cmdutil.NewRabbitMQClient()
			if err != nil {
				return err
			}

			api := api.NewAPI(logger, db, ch)
			srv := api.Server(port)

			go func() { _ = srv.ListenAndServe() }()

			logger.Info("started api", "port", port)

			<-ctx.Done()

			_ = srv.Shutdown(ctx)

			logger.Info("shutdown api")

			return nil
		},
	}

	cmd.Flags().IntVar(&port, "port", 4000, "Port to run the API on")

	return cmd
}
