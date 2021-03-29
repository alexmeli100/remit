package mtn

import (
	"crypto/tls"
	"encoding/json"
	"github.com/pkg/errors"
	"net/http"
)

type Config struct {
	primaryKey string
	apiSecret  string
	userId     string
	baseUrl    string
	targetEnv  string
}

type clientOpts func(c *MomoClient)

type ResponseHandler interface {
	handleResponse(r *http.Response, i interface{}) error
}

type RequestHandler interface {
	Do(r *http.Request) (*http.Response, error)
}

type HttpClient struct {
	client *http.Client
}

// MomoClient helper struct to handle requests to the momo api
type MomoClient struct {
	reqHandler RequestHandler
	resHandler ResponseHandler
}

func (m *HttpClient) Do(r *http.Request) (*http.Response, error) {
	return m.client.Do(r)
}

func (m *HttpClient) handleResponse(res *http.Response, i interface{}) error {
	if i == nil {
		return nil
	}

	if err := decodeBody(res, i); err != nil {
		return errors.Wrap(err, "error decoding response body")
	}

	return nil
}

func createClient(opts ...clientOpts) *MomoClient {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}

	h := &HttpClient{client}
	c := &MomoClient{reqHandler: h, resHandler: h}

	for _, opt := range opts {
		opt(c)
	}

	return c
}

func withAuth(auth *AuthClient) clientOpts {
	return func(c *MomoClient) {
		auth.next = c.reqHandler
		c.reqHandler = auth
	}
}

func withErrorHandler(e *ErrorHandler) clientOpts {
	return func(c *MomoClient) {
		e.next = c.resHandler
		c.resHandler = e
	}
}

func decodeBody(res *http.Response, i interface{}) error {
	decoder := json.NewDecoder(res.Body)

	return decoder.Decode(i)
}
