package httpd

type ErrorRsp struct {
	Code    int    `json:"err_code"`
	Message string `json:"err_message"`
}

type RateRsp struct {
	Rate float64 `json:"rate"`
}
