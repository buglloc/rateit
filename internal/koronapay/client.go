package koronapay

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/buglloc/rateit/internal/ja3"
)

const (
	DefaultUpstream = "https://koronapay.com"
	DefaultRetries  = 3
	DefaultTimeout  = 5 * time.Minute
)

type Client struct {
	httpc *resty.Client
	log   zerolog.Logger
}

func NewClient(opts ...Option) (*Client, error) {
	client := &Client{
		log: log.With().Str("source", "koronapay").Logger(),
		httpc: ja3.NewResty().
			SetRetryCount(DefaultRetries).
			SetTimeout(DefaultTimeout).
			SetHeader("X-Application", "Qpay-Web/3.0").
			SetHeader("Accept", "application/vnd.cft-data.v2.112+json"),
	}

	defaultOpts := []Option{
		WithUpstream(DefaultUpstream),
	}

	for _, opt := range defaultOpts {
		opt(client)
	}

	for _, opt := range opts {
		opt(client)
	}

	return client, nil
}

func (c *Client) Tariff(ctx context.Context, req *TariffReq) (float64, error) {
	c.log.Info().
		Any("req", req).
		Msg("fetch tariffs")

	var rsp []tariffsRsp
	var remoteErr remoteError
	httpRsp, err := c.httpc.R().
		SetContext(ctx).
		SetError(&remoteErr).
		SetResult(&rsp).
		ForceContentType("application/json").
		SetQueryParams(map[string]string{
			"sendingCountryId":        req.Sender.Country.String(),
			"sendingCurrencyId":       req.Sender.Currency.String(),
			"receivingCountryId":      req.Receiver.Country.String(),
			"receivingCurrencyId":     req.Receiver.Currency.String(),
			"paymentMethod":           req.PaymentMethod.String(),
			"receivingMethod":         req.ReceivingMethod.String(),
			"receivingAmount":         strconv.Itoa(int(req.Amount * 100)),
			"paidNotificationEnabled": "false",
		}).
		Get("/transfers/tariffs")
	if err != nil {
		return 0, fmt.Errorf("request failed: %w", err)
	}

	if remoteErr.Code != 0 {
		return 0, fmt.Errorf("remote error %q: %s", remoteErr.Code, remoteErr.Message)
	}

	if !httpRsp.IsSuccess() {
		return 0, fmt.Errorf("non-200 status code: %s", httpRsp.Status())
	}

	if len(rsp) != 1 || rsp[0].SendingAmount == 0 {
		return 0, fmt.Errorf("unexpected response: %s", string(httpRsp.Body()))
	}

	receivingAmount, ok := parseReceivingComment(rsp[0].ReceivingAmountComment)
	if !ok {
		return 0, fmt.Errorf("unexpected receiving amount comment: %s", rsp[0].ReceivingAmountComment)
	}

	return float64(rsp[0].SendingAmount) / float64(receivingAmount*100), nil
}

func parseReceivingComment(in string) (int, bool) {
	if !strings.HasPrefix(in, "~ ") {
		return 0, false
	}

	in = in[2:]
	idx := strings.IndexByte(in, ' ')
	if idx == -1 {
		return 0, false
	}

	out, err := strconv.Atoi(in[:idx])
	return out, err == nil
}
