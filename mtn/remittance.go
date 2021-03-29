package mtn

import (
	"bytes"
	"encoding/json"
	"fmt"
	transfer "github.com/alexmeli100/remit/transfer/pkg/service"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"net/http"
	"strconv"
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
		config: config,
		authorizer: &AuthRemittance{
			client: createClient(withErrorHandler(&ErrorHandler{handler: tokenErrHandler})),
			config: config,
		},
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

func (m *Remittance) SendTo(req *transfer.TransferRequest) (*transfer.TransferResponse, error) {
	payee := &Payee{PartyId: req.RecipientNumber, PartyIdType: "MSISDN"}
	tr := &TransferRequest{
		Amount:     strconv.Itoa(int(req.ReceiveAmount)),
		Currency:   req.Currency,
		ExternalID: req.OrderId,
		Payee:      payee,
	}

	refId, err := m.Transfer(tr)

	if err != nil {
		return nil, err
	}

	res, err := m.GetTransactionStatus(refId)

	if err != nil {
		return nil, err
	}

	ret := &transfer.TransferResponse{
		Amount:          req.Amount,
		RecipientId:     req.RecipientId,
		Currency:        req.Currency,
		Service:         req.Service,
		ReceiveCurrency: req.ReceiveCurrency,
		OrderId:         req.OrderId,
		ReceiveAmount:   req.ReceiveAmount,
		SenderId:        req.SenderId,
		Status:          getStatus(res.Status),
		FailReason:      res.Reason,
		RemoteId:        res.FinancialTransactionId,
		Token:           refId,
	}

	return ret, nil
}

func getStatus(status string) string {
	switch status {
	case "SUCCESSFULL":
		return transfer.StatusSuccess
	case "PENDING":
		return transfer.StatusPending
	case "FAILED":
		return transfer.StatusFailed
	default:
		return transfer.StatusUnknown
	}
}
