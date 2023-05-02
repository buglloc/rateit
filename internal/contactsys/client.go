package contactsys

import (
	"net/http"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/buglloc/rateit/internal/ja3"
)

const (
	DefaultUpstream  = "https://online.contact-sys.com"
	DefaultPartnerID = "D5267BED-18CC-4661-B03A-65934CAE1CA4"
	DefaultRetries   = 3
	DefaultTimeout   = 5 * time.Minute
)

type Client struct {
	tr   http.RoundTripper
	log  zerolog.Logger
	opts []Option
}

func NewClient(opts ...Option) (*Client, error) {
	return &Client{
		log: log.With().Str("source", "contactsys").Logger(),
		tr:  ja3.NewTransport(),
		opts: append(
			[]Option{
				WithUpstream(DefaultUpstream),
				WithPartnerID(DefaultPartnerID),
			},
			opts...,
		),
	}, nil
}

func (c *Client) Session() *Session {
	sess := &Session{
		log: c.log,
		httpc: ja3.NewRestyWithTransport(c.tr).
			SetRetryCount(DefaultRetries).
			SetTimeout(DefaultTimeout),
	}

	for _, opt := range c.opts {
		opt(sess)
	}

	return sess
}
