package commands

import (
	"fmt"
	"log"
	"os"

	"github.com/rs/zerolog"
	"github.com/spf13/cobra"

	"github.com/buglloc/rateit/internal/config"
)

var (
	configs []string
	cfg     *config.Config
)

var rootCmd = &cobra.Command{
	Use:           "rateit",
	SilenceUsage:  true,
	SilenceErrors: true,
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	cobra.OnInitialize(
		initConfig,
		initLogger,
	)

	flags := rootCmd.PersistentFlags()
	flags.StringSliceVar(&configs, "config", nil, "config file")

	rootCmd.AddCommand(
		startCmd,
	)
}

func initConfig() {
	var err error
	cfg, err = config.LoadConfig(configs...)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "unable to load config: %v\n", err)
		os.Exit(1)
	}
}

func initLogger() {
	log.SetOutput(os.Stderr)

	if cfg.Debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	} else {
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}
}
