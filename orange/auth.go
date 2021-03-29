package orange

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"net/http"
	"time"
)

type AccessToken struct {
	Token     string `json:"access_token"`
	Scope     string `json:"scope"`
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
		token, err := t.authorizer.authorize()

		if err != nil {
			return "", err
		}

		t.token = token
		return token.Token, nil
	}

	return t.token.Token, nil
}

type Authorizer interface {
	authorize() (*AccessToken, error)
}

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

type AuthOrangeMoney struct {
	client *OrangeMomoClient
	config *Config
}

func (a *AuthOrangeMoney) authorize() (*AccessToken, error) {
	reqBody, err := json.Marshal(accessTokenBody)

	if err != nil {
		return nil, errors.Wrap(err, "error creating access token request body")
	}

	req, err := http.NewRequest("POST", "https://apiw.orange.cm/token", bytes.NewBuffer(reqBody))

	if err != nil {
		return nil, errors.Wrap(err, "error creating request")
	}

	req.Header.Set("Content-type", "application/x-www-form-urlencoded")
	req.Header.Set("Authorization", "Basic "+GetAuthToken(a.config.clientKey, a.config.clientSecret))

	res, err := a.client.reqHandler.Do(req)

	if err != nil {
		return nil, errors.Wrap(err, "token request error")
	}

	defer res.Body.Close()
	t := &AccessToken{}

	err = a.client.resHandler.handleResponse(res, t)

	t.ExpiresIn = time.Now().Unix() + t.ExpiresIn - 6

	return t, err
}

func GetAuthToken(key, secret string) string {
	k := fmt.Sprintf("%s:%s", key, secret)
	return base64.StdEncoding.EncodeToString([]byte(k))
}
