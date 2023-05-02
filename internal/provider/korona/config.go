package korona

import "github.com/buglloc/rateit/internal/koronapay"

type PaymentMethod = koronapay.PaymentMethod
type CountryID = koronapay.CountryID
type CurrencyID = koronapay.CurrencyID

type Config struct {
	Debug           bool          `yaml:"debug" json:"debug"`
	Amount          float64       `yaml:"amount" json:"amount"`
	Sender          Participant   `yaml:"sender" json:"sender"`
	Receiver        Participant   `yaml:"receiver" json:"receiver"`
	PaymentMethod   PaymentMethod `yaml:"paymentMethod" json:"paymentMethod"`
	ReceivingMethod PaymentMethod `yaml:"receivingMethod" json:"receivingMethod"`
}

type Participant struct {
	Country  CountryID  `yaml:"country" json:"country"`
	Currency CurrencyID `yaml:"currency" json:"currency"`
}
