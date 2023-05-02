package config

import (
	"encoding"
	"encoding/json"
	"fmt"

	"gopkg.in/yaml.v3"
)

var _ yaml.Unmarshaler = (*ProviderKind)(nil)
var _ yaml.Marshaler = (*ProviderKind)(nil)
var _ json.Unmarshaler = (*ProviderKind)(nil)
var _ json.Marshaler = (*ProviderKind)(nil)
var _ encoding.TextUnmarshaler = (*ProviderKind)(nil)

type ProviderKind uint8

const (
	ProviderKindNone ProviderKind = iota
	ProviderKindContact
	ProviderKindKorona
)

func (k ProviderKind) String() string {
	switch k {
	case ProviderKindContact:
		return "contact"
	case ProviderKindKorona:
		return "korona"
	case ProviderKindNone:
		return "none"
	default:
		return fmt.Sprintf("provider#%d", k)
	}
}

func (k *ProviderKind) fromString(v string) error {
	switch v {
	case "contact":
		*k = ProviderKindContact
	case "korona":
		*k = ProviderKindKorona
	case "none", "":
		*k = ProviderKindNone
	default:
		return fmt.Errorf("unknown provider: %s", v)
	}

	return nil
}

func (k ProviderKind) MarshalYAML() (interface{}, error) {
	return k.String(), nil
}

func (k *ProviderKind) UnmarshalYAML(value *yaml.Node) error {
	var s string
	if err := value.Decode(&s); err != nil {
		return err
	}

	return k.fromString(s)
}

func (k ProviderKind) MarshalJSON() ([]byte, error) {
	return json.Marshal(k.String())
}

func (k *ProviderKind) UnmarshalJSON(in []byte) error {
	var s string
	if err := json.Unmarshal(in, &s); err != nil {
		return err
	}

	return k.fromString(s)
}

func (k *ProviderKind) UnmarshalText(v []byte) error {
	return k.fromString(string(v))
}
