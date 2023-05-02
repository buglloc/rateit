package contact

import (
	"context"
	"fmt"

	"github.com/buglloc/rateit/internal/contactsys"
)

type Provider struct {
	cfg Config
	cc  *contactsys.Client
}

func NewProvider(cfg Config) (*Provider, error) {
	cc, err := contactsys.NewClient(
		contactsys.WithVerbose(cfg.Debug),
	)
	if err != nil {
		return nil, err
	}

	return &Provider{
		cfg: cfg,
		cc:  cc,
	}, nil
}

func (p *Provider) CurrentRate(ctx context.Context) (float64, error) {
	sess := p.cc.Session()

	if err := sess.Auth(ctx); err != nil {
		return 0, fmt.Errorf("auth failed: %w", err)
	}

	countries, err := sess.Countries(ctx)
	if err != nil {
		return 0, fmt.Errorf("unable to fetch countries: %w", err)
	}

	var country contactsys.Country
	for _, c := range countries {
		if c.Code == p.cfg.Country || c.Name == p.cfg.Country {
			country = c
			break
		}
	}
	if country.ID == 0 {
		return 0, fmt.Errorf("unable to find country with Code or Name: %s", p.cfg.Country)
	}

	banks, err := sess.Banks(ctx, country)
	if err != nil {
		return 0, fmt.Errorf("unable to fetch banks: %w", err)
	}

	var bank contactsys.Bank
	for _, b := range banks {
		if b.BankCode == p.cfg.Bank {
			bank = b
			break
		}
	}

	if bank.BankData == "" {
		return 0, fmt.Errorf("unable to find bank with Code: %s", p.cfg.Bank)
	}

	tr, err := sess.StartTransaction(ctx, bank)
	if err != nil {
		return 0, fmt.Errorf("unable to start transaction: %w", err)
	}

	err = sess.FillTransferDetails(ctx, tr, contactsys.RandomTransferDetails(p.cfg.Amount))
	if err != nil {
		return 0, fmt.Errorf("unable to fill transfer details: %w", err)
	}

	err = sess.FillTransferParticipants(ctx, tr, &contactsys.TransferParticipants{
		TransferSender:    contactsys.RandomTransferSender(),
		TransferRecipient: contactsys.RandomTransferRecipient(country, bank),
	})
	if err != nil {
		return 0, fmt.Errorf("unable to fill transfer practicans: %w", err)
	}

	rate, err := sess.TransferRate(ctx, tr)
	if err != nil {
		return 0, fmt.Errorf("unable to get transfer rate: %w", err)
	}

	return rate, nil
}

func (p *Provider) Name() string {
	return "contact"
}
