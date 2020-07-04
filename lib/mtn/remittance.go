package mtn

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
	"net/http"
	"time"
)

type TransferStatus string

const millisConversion = int64(time.Millisecond) / int64(time.Nanosecond)

const (
	TransferSuccessFul = "SUCCESSFUL"
	TransferPending    = "PENDING"
	TransferFailed     = "FAILED"
)

const (
	BaseURL   = "https://sandbox.momodeveloper.mtn.com"
	TargetEnv = "sandbox"
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

func NewRemittance(config *Config) (*Remittance, error) {
	refresher := &tokenRefresher{
		config:     config,
		authorizer: authRemittance,
	}
	auth := &AuthClient{
		refresher: refresher,
	}

	client := createClient(withAuth(auth), withErrorHandler(&ErrorHandler{handler: remittanceErrHandler}))

	r := &Remittance{client, config}

	return r, nil
}

func (m *Remittance) Transfer(t *TransferRequest) (string, error) {
	reqBody, err := json.Marshal(t)

	if err != nil {
		return "", errors.Wrap(err, "error creating request body")
	}

	req, err := http.NewRequest("POST", BaseURL+"/remittance/v1_0/transfer", bytes.NewBuffer(reqBody))

	if err != nil {
		return "", errors.Wrap(err, "error creating transfer request")
	}

	refId := uuid.NewV4()
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Reference-Id", refId.String())
	req.Header.Set("X-Target-Environment", TargetEnv)
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
	url := fmt.Sprintf("%s/remittance/v1_0/transfer/%s", BaseURL, refId)
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
// this function keeps pooling the api until it responds with a failed or successful status
// The interval time for pooling is specified in milliseconds
func (m *Remittance) GetFinalStatus(refId string, interval int64, maxTime int64) (*TransferResponse, error) {
	end := time.Now().UnixNano()/millisConversion + maxTime

	for {
		res, err := m.GetTransactionStatus(refId)

		if err != nil {
			return nil, err
		}

		if res.Status == TransferSuccessFul || res.Status == TransferFailed {
			return res, nil
		}

		if end-time.Now().UnixNano()/millisConversion < 0 {
			break
		}

		time.Sleep(time.Duration(interval) * time.Millisecond)
	}

	return nil, ErrorPending{interval, maxTime}
}
