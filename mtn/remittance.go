package mtn

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"net/http"
	"strconv"
	"time"
)

type TransferStatus string

const (
	TransferSuccessFul = "SUCCESSFUL"
	TransferPending    = "PENDING"
	TransferFailed     = "FAILED"
)

type Payee struct {
	PartyIdType string `json:"partyIdType"`
	PartyId     string `json:"partyId"`
}

type FailReason struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type TransferRequest struct {
	Amount       string `json:"amount"`
	Currency     string `json:"currency"`
	ExternalID   string `json:"externalId"`
	Payee        *Payee `json:"payee"`
	PayerMessage string `json:"payerMessage"`
	PayeeNote    string `json:"payeeNote"`
}

// The momo api documentation specifies the reason field as a FailedReason struct(see above)
// but the response has a string instead which caused the unmarshalling to fail.
// So I changed it to string and hopefully it stays this way
type TransferResponse struct {
	Amount                 string `json:"amount"`
	Currency               string `json:"currency"`
	FinancialTransactionId string `json:"financialTransactionId"`
	ExternalId             string `json:"externalId"`
	Payee                  *Payee `json:"payee"`
	PayerMessage           string `json:"payerMessage"`
	PayeeNote              string `json:"payeeNote"`
	Status                 string `json:"status"`
	Reason                 string `json:"reason"`
}

type Remittance struct {
	client *MomoClient
	config *Config
}

func NewRemittance(config *Config) *Remittance {
	refresher := &tokenRefresher{
		config:     config,
		authorizer: authRemittance,
	}
	auth := &AuthClient{
		refresher: refresher,
	}

	client := createClient(withAuth(auth), withErrorHandler(&ErrorHandler{handler: remittanceErrHandler}))

	r := &Remittance{client, config}

	return r
}

func (m *Remittance) Transfer(t *TransferRequest) (string, error) {
	reqBody, err := json.Marshal(t)

	if err != nil {
		return "", errors.Wrap(err, "error creating request body")
	}

	req, err := http.NewRequest("POST", m.config.baseUrl+"/remittance/v1_0/transfer", bytes.NewBuffer(reqBody))

	if err != nil {
		return "", errors.Wrap(err, "error creating transfer request")
	}

	refId := uuid.New()

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Reference-Id", refId.String())
	req.Header.Set("X-Target-Environment", m.config.targetEnv)
	req.Header.Set("Ocp-Apim-Subscription-Key", m.config.primaryKey)

	res, err := m.client.reqHandler.Do(req)

	if err != nil {
		return "", err
	}

	defer res.Body.Close()
	err = m.client.resHandler.handleResponse(res, nil)

	if err != nil {
		return "", err
	}

	return refId.String(), nil
}

func (m *Remittance) GetTransactionStatus(refId string) (*TransferResponse, error) {
	url := fmt.Sprintf("%s/remittance/v1_0/transfer/%s", m.config.baseUrl, refId)
	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		return nil, errors.Wrap(err, "error creating request")
	}

	req.Header.Set("X-Target-Environment", TargetEnv)
	req.Header.Set("Ocp-Apim-Subscription-Key", m.config.primaryKey)

	res, err := m.client.reqHandler.Do(req)

	if err != nil {
		return nil, err
	}

	defer res.Body.Close()
	var trRes TransferResponse
	err = m.client.resHandler.handleResponse(res, &trRes)

	if err != nil {
		return nil, err
	}

	return &trRes, nil
}

// get the final status of the transaction
// this function keeps polling the api for a max duration of timeout until it responds with a failed or successful status
// The interval time for polling is specified in milliseconds
func (m *Remittance) GetFinalStatus(refId string, interval time.Duration, timeout time.Duration) (*TransferResponse, error) {
	ticker := time.NewTicker(interval)

	for {
		select {
		case _ = <-ticker.C:
			res, err := m.GetTransactionStatus(refId)

			if err != nil {
				return nil, err
			}

			if res.Status == TransferSuccessFul || res.Status == TransferFailed {
				return res, nil
			}
		case <-time.After(timeout):
			ticker.Stop()
			return nil, ErrorPending
		}
	}
}

func (m *Remittance) SendTo(amount float64, recipient, currency string) error {
	payee := &Payee{PartyId: recipient, PartyIdType: "MSISDN"}

	tr := &TransferRequest{
		Amount:   strconv.FormatFloat(amount, 'f', 2, 64),
		Currency: currency,
		Payee:    payee,
	}

	refId, err := m.Transfer(tr)

	if err != nil {
		return err
	}

	res, err := m.GetFinalStatus(refId, time.Millisecond*100, time.Second)

	if err != nil {
		return err
	}

	if res.Status == TransferFailed {
		return ErrTransferFailed{
			message: res.Reason,
		}
	}

	return nil
}
