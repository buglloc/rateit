package contactsys

type remoteError struct {
	ErrorCode   string `json:"errorCode"`
	Description string `json:"description"`
}

type Country struct {
	ID   int    `json:"id"`
	Name string `json:"caption"`
	Code string `json:"code"`
}

type Bank struct {
	BankData string `json:"bankData"`
	BankCode string `json:"bankCode"`
	Address  string `json:"address"`
}
