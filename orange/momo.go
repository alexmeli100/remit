package orange

import (
	"bytes"
	"encoding/json"
	"fmt"
	transfer "github.com/alexmeli100/remit/transfer/pkg/service"
	"github.com/pkg/errors"
	"net/http"
)

type Cashin struct {
	client *OrangeMomoClient
	config *Config
}

type CashinRequest struct {
	ChannelUserMsisdn string `json:"channelUserMsisdn"`
	Amount            int    `json:"amount"`
	SubscriberMsisdn  string `json:"subscriberMsisdn"`
	Pin               string `json:"pin"`
	OrderId           string `json:"orderId"`
	PayToken          string `json:"payToken"`
	NotifUrl          string `json:"notifUrl,omitempty"`
}

type cashinResponseData struct {
	CreateTime        string `json:"createtime"`
	ChannelUserMsisdn string `json:"channelUserMsisdn"`
	Amount            string `json:"amount"`
	PayToken          string `json:"payToken"`
	Txnid             string `json:"txnid"`
	Passcode          string `json:"passcode"`
	Txnmode           string `json:"txnmode"`
	Txnmessage        string `json:"txnmessage"`
	Txnstatus         string `json:"txnstatus"`
	OrderId           string `json:"orderid"`
	Status            string `json:"status"`
}

type CashinResponse struct {
	Message string             `json:"message"`
	Data    cashinResponseData `json:"data"`
}

func NewCashIn(config *Config) *Cashin {
	refresher := &tokenRefresher{
		config: config,
		authorizer: &AuthOrangeMoney{
			client: createClient(withErrorHandler(&ErrorHandler{handler: momoErrorHandler})),
			config: config,
		},
	}

	auth := &AuthClient{
		refresher: refresher,
	}

	client := createClient(withAuth(auth), withErrorHandler(&ErrorHandler{handler: momoErrorHandler}))

	r := &Cashin{client, config}

	return r
}

type initData struct {
	PayToken string `json:"payToken"`
}

type initResponse struct {
	Message string   `json:"message"`
	Data    initData `json:"data"`
}

func (c *Cashin) init(xAuthToken string) (string, error) {
	req, err := http.NewRequest("POST", c.config.baseUrl+"/cashin/init", nil)

	if err != nil {
		return "", errors.Wrap(err, "error creating init request")
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Auth-Token", xAuthToken)

	res, err := c.client.reqHandler.Do(req)

	if err != nil {
		return "", err
	}

	var initRes initResponse
	err = c.client.resHandler.handleResponse(res, &initRes)

	if err != nil {
		return "", err
	}

	return initRes.Data.PayToken, nil
}

func (c *Cashin) pay(tReq *CashinRequest) (*CashinResponse, error) {
	xAuthToken := GetAuthToken(c.config.apiKey, c.config.apiSecret)
	token, err := c.init(xAuthToken)
	tReq.PayToken = token

	if err != nil {
		return nil, errors.Wrap(err, "error getting pay token")
	}

	reqBody, err := json.Marshal(tReq)

	if err != nil {
		return nil, errors.Wrap(err, "error creating request body")
	}

	req, err := http.NewRequest("POST", c.config.baseUrl+"/cashin/pay", bytes.NewBuffer(reqBody))

	if err != nil {
		return nil, errors.Wrap(err, "error creating transfer request")
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Auth-Token", xAuthToken)
	res, err := c.client.reqHandler.Do(req)

	if err != nil {
		return nil, errors.Wrap(err, "error sending pay request")
	}

	var cRes CashinResponse
	err = c.client.resHandler.handleResponse(res, &cRes)

	if err != nil {
		return nil, errors.Wrap(err, "error getting pay response")
	}

	return &cRes, nil
}

func (c *Cashin) paymentStatus(token string) (*CashinResponse, error) {
	xAuthToken := GetAuthToken(c.config.apiKey, c.config.apiSecret)
	url := fmt.Sprintf("%s/cashin/pay/%s", c.config.baseUrl, token)
	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		return nil, errors.Wrap(err, "error creating transfer request")
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Auth-Token", xAuthToken)
	res, err := c.client.reqHandler.Do(req)

	if err != nil {
		return nil, errors.Wrap(err, "error sending pay request")
	}

	var cRes CashinResponse
	err = c.client.resHandler.handleResponse(res, &cRes)

	if err != nil {
		return nil, errors.Wrap(err, "error getting pay response")
	}

	return &cRes, nil
}

func (c *Cashin) SendTo(req *transfer.TransferRequest) (*transfer.TransferResponse, error) {
	cReq := &CashinRequest{
		ChannelUserMsisdn: req.RecipientNumber,
		Amount:            int(req.Amount),
		SubscriberMsisdn:  c.config.msisdn,
		Pin:               c.config.pin,
		OrderId:           req.OrderId,
	}

	res, err := c.pay(cReq)

	if err != nil {
		return nil, errors.Wrap(err, "Failed to perform orange momo transfer")
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
		Status:          getStatus(res.Data.Status),
		FailReason:      res.Message,
		RemoteId:        res.Data.Txnid,
		Token:           res.Data.PayToken,
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
	case "INITIATED":
		return transfer.StatusInitiated
	case "EXPIRED":
		return transfer.StatusExpired
	default:
		return transfer.StatusUnknown
	}
}
