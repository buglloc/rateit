package contact

type Config struct {
	Debug   bool    `yaml:"debug" json:"debug"`
	Amount  float64 `yaml:"amount" json:"amount"`
	Country string  `yaml:"country" json:"country"`
	Bank    string  `yaml:"bank" json:"bank"`
}
