package contactsys

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-resty/resty/v2"
	"github.com/rs/zerolog"
)

type Session struct {
	httpc     *resty.Client
	partnerID string
	log       zerolog.Logger
}

func (s *Session) Auth(ctx context.Context) error {
	s.log.Info().Msg("authenticate")

	var rsp struct {
		AccessToken string `json:"accessToken"`
		Type        string `json:"type"`
	}

	var remoteErr remoteError
	httpRsp, err := s.httpc.R().
		SetContext(ctx).
		SetError(&remoteErr).
		SetResult(&rsp).
		SetBody(map[string]string{
			"tokenType": "SplitTokenV2",
			"grantType": "anonymous",
			"ticket":    s.partnerID,
		}).
		Post("/auth/token")
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}

	if remoteErr.ErrorCode != "" {
		return fmt.Errorf("remote error %q: %s", remoteErr.ErrorCode, remoteErr.Description)
	}

	if !httpRsp.IsSuccess() {
		return fmt.Errorf("non-200 status code: %s", httpRsp.Status())
	}

	if rsp.AccessToken == "" || rsp.Type == "" {
		return fmt.Errorf("unexpected response: %s", string(httpRsp.Body()))
	}

	s.httpc.SetHeader("Authorization", fmt.Sprintf("%s %s", rsp.Type, rsp.AccessToken))
	s.log.Info().Msg("authenticated")
	return nil
}

func (s *Session) Countries(ctx context.Context) ([]Country, error) {
	s.log.Info().Msg("list countries")

	var rsp []Country
	return rsp, s.execute(
		http.MethodGet,
		"/countries",
		s.httpc.R().
			SetContext(ctx).
			SetResult(&rsp),
	)
}

func (s *Session) Banks(ctx context.Context, country Country) ([]Bank, error) {
	s.log.Info().
		Str("country_code", country.Code).
		Msg("list banks")

	var rsp []Bank
	return rsp, s.execute(
		http.MethodGet,
		"/banks",
		s.httpc.R().
			SetContext(ctx).
			SetQueryParams(map[string]string{
				"countryId":     strconv.Itoa(country.ID),
				"deliveryType":  "ACCOUNT",
				"recipientType": "INDIVIDUAL",
			}).
			SetResult(&rsp),
	)
}

func (s *Session) StartTransaction(ctx context.Context, bank Bank) (*Transaction, error) {
	s.log.Info().
		Str("bank_code", bank.BankCode).
		Msg("start transaction")

	var rsp Transaction
	return &rsp, s.execute(
		http.MethodPost,
		"/trns/cash",
		s.httpc.R().
			SetContext(ctx).
			SetBody(map[string]string{
				"bankData": bank.BankData,
			}).
			SetResult(&rsp),
	)
}

func (s *Session) FillTransferDetails(ctx context.Context, transaction *Transaction, details *TransferDetails) error {
	s.log.Info().
		Str("transaction_id", transaction.ID).
		Any("details", details).
		Msg("fill transfer details")

	return s.execute(
		http.MethodPut,
		"/trns/{id}/fields",
		s.httpc.R().
			SetContext(ctx).
			SetPathParam("id", transaction.ID).
			SetBody(details),
	)
}

func (s *Session) FillTransferParticipants(ctx context.Context, transaction *Transaction, participants *TransferParticipants) error {
	s.log.Info().
		Str("transaction_id", transaction.ID).
		Any("participants", participants).
		Msg("fill transfer participants")

	return s.execute(
		http.MethodPut,
		"/trns/{id}/fields",
		s.httpc.R().
			SetContext(ctx).
			SetPathParam("id", transaction.ID).
			SetBody(participants),
	)
}

func (s *Session) TransferRate(ctx context.Context, transaction *Transaction) (float64, error) {
	s.log.Info().
		Str("transaction_id", transaction.ID).
		Msg("request transfer fees")

	/*
		example:
			{
			  "type": "EXACT",
			  "value": [
			    "403.85"
			  ],
			  "rate": "80.77",
			  "totalAmount": "113480.45",
			  "amount": "113480.45",
			  "currency": "RUB",
			  "payoutAmount": "47502.00",
			  "payoutCurrency": "THB",
			  "payoutRate": "33.93",
			  "transactionAmount": "1400.00",
			  "transactionCurrency": "USD",
			  "payoutToEnterRate": "2.39",
			  "enterToPayoutRate": null
			}
	*/
	var rsp struct {
		TotalAmount  float64 `json:"totalAmount,string"`
		PayoutAmount float64 `json:"payoutAmount,string"`
	}

	err := s.execute(
		http.MethodPost,
		"/trns/{id}/fees",
		s.httpc.R().
			SetContext(ctx).
			SetPathParam("id", transaction.ID).
			SetResult(&rsp),
	)
	if err != nil {
		return 0, err
	}

	if rsp.TotalAmount == 0 || rsp.PayoutAmount == 0 {
		return 0, errors.New("invalid fees: amount are zero")
	}
	return rsp.TotalAmount / rsp.PayoutAmount, nil
}

func (s *Session) execute(method string, url string, req *resty.Request) error {
	var remoteErr remoteError
	rsp, err := req.SetError(&remoteErr).Execute(method, url)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}

	if remoteErr.ErrorCode != "" {
		return fmt.Errorf("remote error %q: %s", remoteErr.ErrorCode, remoteErr.Description)
	}

	if !rsp.IsSuccess() {
		return fmt.Errorf("non-200 status code: %s", rsp.Status())
	}

	return nil
}
