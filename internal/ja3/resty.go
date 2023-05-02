package ja3

import (
	"net"
	"net/http"
	"time"

	"github.com/go-resty/resty/v2"
)

const (
	dialTimeout = 10 * time.Second
	keepAlive   = 60 * time.Second
)

func NewResty() *resty.Client {
	return NewRestyWithTransport(NewTransport())
}

func NewRestyWithTransport(tr http.RoundTripper) *resty.Client {
	return resty.New().
		SetTransport(tr).
		SetHeader("User-Agent", UserAgent).
		SetHeader("Accept-Language", "en")
}

func NewTransport() http.RoundTripper {
	dialer := &net.Dialer{
		Timeout:   dialTimeout,
		KeepAlive: keepAlive,
	}

	return newRoundTripper(dialer)
}
