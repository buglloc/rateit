package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Debug     bool       `yaml:"debug"`
	Addr      string     `yaml:"addr"`
	Providers []Provider `yaml:"providers"`
}

func LoadConfig(configs ...string) (*Config, error) {
	out := &Config{
		Debug: false,
		Addr:  "127.0.0.1:3000",
	}

	if len(configs) == 0 {
		return out, nil
	}

	for _, cfgPath := range configs {
		err := func() error {
			f, err := os.Open(cfgPath)
			if err != nil {
				return fmt.Errorf("unable to open config file: %w", err)
			}
			defer func() { _ = f.Close() }()

			if err := yaml.NewDecoder(f).Decode(&out); err != nil {
				return fmt.Errorf("invalid config: %w", err)
			}

			return nil
		}()
		if err != nil {
			return nil, fmt.Errorf("unable to load config %q: %w", cfgPath, err)
		}
	}

	return out, nil
}
