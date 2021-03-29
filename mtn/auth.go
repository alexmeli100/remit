package mtn

import (
	"encoding/base64"
	"fmt"
	"github.com/pkg/errors"
	"net/http"
	"time"
)

type AccessToken struct {
	Token     string `json:"access_token"`
	TokenType string `json:"token_type"`
	ExpiresIn int64  `json:"expires_in"`
}

type tokenRefresher struct {
	config     *Config
	authorizer Authorizer
	token      *AccessToken
}

// refresh the authorization token if expired otherwise return the old one
func (t *tokenRefresher) refresh() (string, error) {
	if isExpired(t.token) {
		token, err := t.authorizer(t.config)

		if err != nil {
			return "", err
		}

		t.token = token
		return token.Token, nil
	}

	return t.token.Token, nil
}

type Authorizer = func(config *Config) (*AccessToken, error)

// AuthClient this struct ensures a valid token is always available for the next request handler
type AuthClient struct {
	refresher *tokenRefresher
	next      RequestHandler
}

// Do get the token from the refresher and add it to the header for the next request handler
func (a AuthClient) Do(r *http.Request) (*http.Response, error) {
	token, err := a.refresher.refresh()

	if err != nil {
		return nil, errors.Wrap(err, "error getting token")
	}

	r.Header.Add("Authorization", "Bearer "+token)

	return a.next.Do(r)
}

func isExpired(t *AccessToken) bool {
	if t == nil {
		return true
	}

	return time.Now().Unix() > t.ExpiresIn
}

func authRemittance(config *Config) (*AccessToken, error) {
	return authService("/remittance/token/", config)
}

func authService(path string, config *Config) (*AccessToken, error) {
	client := createClient(withErrorHandler(&ErrorHandler{handler: tokenErrHandler}))
	url := config.baseUrl + path
	req, err := http.NewRequest("POST", url, nil)

	if err != nil {
		return nil, errors.Wrap(err, "error creating request")
	}

	req.Header.Set("Ocp-Apim-Subscription-Key", config.primaryKey)
	req.Header.Set("Authorization", "Basic "+getAuthToken(config))

	res, err := client.reqHandler.Do(req)

	if err != nil {
		return nil, errors.Wrap(err, "token request error")
	}

	defer res.Body.Close()
	t := &AccessToken{}

	err = client.resHandler.handleResponse(res, t)

	t.ExpiresIn = time.Now().Unix() + t.ExpiresIn - 6

	return t, err
}

func getAuthToken(config *Config) string {
	key := fmt.Sprintf("%s:%s", config.userId, config.apiSecret)
	return base64.StdEncoding.EncodeToString([]byte(key))
}
