package commands

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"

	"github.com/buglloc/rateit/internal/httpd"
)

var startCmd = &cobra.Command{
	Use:           "start",
	SilenceUsage:  true,
	SilenceErrors: true,
	Short:         "Starts server",
	RunE: func(_ *cobra.Command, _ []string) error {
		app, err := httpd.NewServer(cfg)
		if err != nil {
			return fmt.Errorf("unable to create httpd service: %w", err)
		}

		errChan := make(chan error, 1)
		doneChan := make(chan struct{})
		go func() {
			defer close(doneChan)

			err := app.ListenAndServe()
			if err != nil {
				errChan <- err
				return
			}
		}()

		stopChan := make(chan os.Signal, 1)
		signal.Notify(stopChan, syscall.SIGINT, syscall.SIGTERM)
		select {
		case <-stopChan:
			log.Info().Msg("shutting down gracefully by signal")

			ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
			defer cancel()

			if err := app.Shutdown(ctx); err != nil {
				log.Error().Err(err).Msg("shutdown failed")
			}
		case err := <-errChan:
			log.Error().Err(err).Msg("start failed")
			return err
		case <-doneChan:
		}

		return nil
	},
}
