package config

import (
	"fmt"

	"github.com/buglloc/rateit/internal/provider"
	"github.com/buglloc/rateit/internal/provider/contact"
	"github.com/buglloc/rateit/internal/provider/korona"
)

type Provider struct {
	Kind    ProviderKind `yaml:"kind"`
	Route   string       `yaml:"route"`
	Contact contact.Config
	Korona  korona.Config
}

func (p *Provider) NewProvider() (provider.Provider, error) {
	switch p.Kind {
	case ProviderKindContact:
		return contact.NewProvider(p.Contact)
	case ProviderKindKorona:
		return korona.NewProvider(p.Korona)
	default:
		return nil, fmt.Errorf("unsupported provider: %v", p.Kind)
	}
}
