package commands

import (
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"

	"github.com/buglloc/rateit/internal/httpd"
)

var watchArgs struct {
	Upstream string
	Period   time.Duration
}

var watchCmd = &cobra.Command{
	Use:           "watch",
	SilenceUsage:  true,
	SilenceErrors: true,
	Short:         "Watch rates",
	RunE: func(_ *cobra.Command, _ []string) error {
		type providerInfo struct {
			provider string
			route    string
		}
		providers := make([]providerInfo, len(cfg.Providers))
		for i, p := range cfg.Providers {
			pp, err := p.NewProvider()
			if err != nil {
				return fmt.Errorf("unable to create provider: %w", err)
			}

			providers[i] = providerInfo{
				provider: pp.Name(),
				route:    p.Route,
			}
		}

		httpc := resty.New().
			SetBaseURL(watchArgs.Upstream).
			SetTimeout(5 * time.Minute)

		stopChan := make(chan os.Signal, 1)
		signal.Notify(stopChan, syscall.SIGINT, syscall.SIGTERM)

		type rateInfo struct {
			provider string
			rate     float64
		}

		var out strings.Builder
		for {
			nextDumpAt := time.Now().Add(watchArgs.Period).Truncate(watchArgs.Period)
			log.Info().
				Time("planned_at", nextDumpAt).
				Msg("wait next sync")

			select {
			case <-stopChan:
				return nil
			case <-time.After(time.Until(nextDumpAt)):
				log.Info().Msg("start sync")
				rates := make([]rateInfo, 0, len(providers))
				for _, p := range providers {
					var rsp httpd.RateRsp
					_, err := httpc.R().
						SetResult(&rsp).
						Get("/api/v1/rate/" + p.route)
					if err != nil {
						log.Error().Err(err).Str("provider", p.provider).Msg("sync failed")
					}

					rates = append(rates, rateInfo{
						provider: p.provider,
						rate:     rsp.Rate,
					})
				}

				out.Reset()
				out.WriteString(time.Now().Format(time.RFC822))
				out.WriteString("\t")
				for i, r := range rates {
					if i != 0 {
						out.WriteByte('\t')
					}
					out.WriteString(r.provider)
					out.WriteByte('=')
					out.WriteString(fmt.Sprintf("%.4f", r.rate))
				}
				fmt.Println(out.String())
			}
		}
	},
}

func init() {
	flags := watchCmd.PersistentFlags()
	flags.StringVar(&watchArgs.Upstream, "upstream", "http://localhost:3000", "rateit upstream")
	flags.DurationVar(&watchArgs.Period, "period", time.Hour, "period to fetch")
}
