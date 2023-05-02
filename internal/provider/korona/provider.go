package korona

import (
	"context"

	"github.com/buglloc/rateit/internal/koronapay"
	"github.com/buglloc/rateit/internal/provider"
)

var _ provider.Provider = (*Provider)(nil)

type Provider struct {
	req *koronapay.TariffReq
	kc  *koronapay.Client
}

func NewProvider(cfg Config) (*Provider, error) {
	kc, err := koronapay.NewClient(
		koronapay.WithVerbose(cfg.Debug),
	)
	if err != nil {
		return nil, err
	}

	return &Provider{
		kc: kc,
		req: &koronapay.TariffReq{
			Amount: cfg.Amount,
			Sender: koronapay.Participant{
				Country:  cfg.Sender.Country,
				Currency: cfg.Sender.Currency,
			},
			Receiver: koronapay.Participant{
				Country:  cfg.Receiver.Country,
				Currency: cfg.Receiver.Currency,
			},
			PaymentMethod:   cfg.PaymentMethod,
			ReceivingMethod: cfg.ReceivingMethod,
		},
	}, nil
}

func (p *Provider) CurrentRate(ctx context.Context) (float64, error) {
	return p.kc.Tariff(ctx, p.req)
}

func (p *Provider) Name() string {
	return "korona"
}
