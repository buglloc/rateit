package koronapay

import "strconv"

const (
	CountryIDRussia   CountryID = "RUS"
	CountryIDThailand CountryID = "THA"

	CurrencyIDRUB CurrencyID = 810
	CurrencyIDUSD CurrencyID = 840

	PaymentMethodDebitCard          PaymentMethod = "debitCard"
	PaymentMethodAccountViaDeeMoney PaymentMethod = "accountViaDeeMoney"
)

type CountryID string

func (c CountryID) String() string {
	return string(c)
}

type PaymentMethod string

func (p PaymentMethod) String() string {
	return string(p)
}

type CurrencyID int

func (c CurrencyID) String() string {
	return strconv.Itoa(int(c))
}

type TariffReq struct {
	Amount          float64
	Sender          Participant
	Receiver        Participant
	PaymentMethod   PaymentMethod
	ReceivingMethod PaymentMethod
}

type Participant struct {
	Country  CountryID
	Currency CurrencyID
}

type remoteError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Type    string `json:"error"`
}

type tariffsRsp struct {
	SendingCurrency                currency `json:"sendingCurrency"`
	SendingAmount                  int      `json:"sendingAmount"`
	SendingAmountDiscount          int      `json:"sendingAmountDiscount"`
	SendingAmountWithoutCommission int      `json:"sendingAmountWithoutCommission"`
	SendingCommission              int      `json:"sendingCommission"`
	SendingCommissionDiscount      int      `json:"sendingCommissionDiscount"`
	SendingTransferCommission      int      `json:"sendingTransferCommission"`
	PaidNotificationCommission     int      `json:"paidNotificationCommission"`
	ReceivingCurrency              currency `json:"receivingCurrency"`
	ReceivingAmount                int      `json:"receivingAmount"`
	ReceivingAmountComment         string   `json:"receivingAmountComment"`
	ExchangeRate                   float64  `json:"exchangeRate"`
	ExchangeRateType               string   `json:"exchangeRateType"`
	ExchangeRateDiscount           int      `json:"exchangeRateDiscount"`
	Profit                         int      `json:"profit"`
}

type currency struct {
	ID   string `json:"id"`
	Code string `json:"code"`
	Name string `json:"name"`
}
