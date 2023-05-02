package provider

import "context"

type Provider interface {
	Name() string
	CurrentRate(ctx context.Context) (float64, error)
}
