package orange

import (
	"crypto/tls"
	"encoding/json"
	"github.com/pkg/errors"
	"net/http"
)

type ErrorHandlerFunc func(res *http.Response) error

type ErrorHandler struct {
	handler ErrorHandlerFunc
	next    ResponseHandler
}

func (e *ErrorHandler) handleResponse(res *http.Response, i interface{}) error {
	if res.StatusCode >= 200 && res.StatusCode < 300 {
		return e.next.handleResponse(res, i)
	} else if res.StatusCode == http.StatusUnauthorized {
		return errors.Errorf("status: 401\nmessage: status unauthorized")
	}

	return e.handler(res)
}

func momoErrorHandler(res *http.Response) error {
	return errors.Errorf("Orange money API failed with response: %s", res.Status)
}

type Config struct {
	apiKey       string
	apiSecret    string
	clientKey    string
	clientSecret string
	baseUrl      string
	msisdn       string
	pin          string
	targetEnv    string
}

var accessTokenBody = map[string]string{"grant_type": "client_credentials"}

type clientOpts func(c *OrangeMomoClient)

type ResponseHandler interface {
	handleResponse(r *http.Response, i interface{}) error
}

type RequestHandler interface {
	Do(r *http.Request) (*http.Response, error)
}

type HttpClient struct {
	client *http.Client
}

// OrangeMomoClient MomoClient helper struct to handle requests to the momo api
type OrangeMomoClient struct {
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

func createClient(opts ...clientOpts) *OrangeMomoClient {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}

	h := &HttpClient{client}
	c := &OrangeMomoClient{reqHandler: h, resHandler: h}

	for _, opt := range opts {
		opt(c)
	}

	return c
}

func withAuth(auth *AuthClient) clientOpts {
	return func(c *OrangeMomoClient) {
		auth.next = c.reqHandler
		c.reqHandler = auth
	}
}

func withErrorHandler(e *ErrorHandler) clientOpts {
	return func(c *OrangeMomoClient) {
		e.next = c.resHandler
		c.resHandler = e
	}
}

func decodeBody(res *http.Response, i interface{}) error {
	decoder := json.NewDecoder(res.Body)

	return decoder.Decode(i)
}
