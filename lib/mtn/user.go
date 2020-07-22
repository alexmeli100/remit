package mtn

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
	"net/http"
)

type User struct {
	pk     string
	UserId string
	client *MomoClient
}

func NewUser(host string, pk string) (*User, error) {
	userId, err := uuid.NewV4()

	if err != nil {
		return nil, err
	}

	client := createClient(withErrorHandler(&ErrorHandler{handler: remittanceErrHandler}))
	body, err := json.Marshal(map[string]string{"providerCallbackHost": host})

	if err != nil {
		return nil, errors.Wrap(err, "error creating request body")
	}

	req, err := http.NewRequest("POST", BaseURL+"/v1_0/apiuser", bytes.NewBuffer(body))

	if err != nil {
		return nil, errors.Wrap(err, "error creating request")
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Reference-Id", userId.String())
	req.Header.Set("Ocp-Apim-Subscription-Key", pk)

	res, err := client.reqHandler.Do(req)

	if err != nil {
		return nil, err
	}

	defer res.Body.Close()
	err = client.resHandler.handleResponse(res, nil)

	if err != nil {
		return nil, err
	}

	user := &User{pk: pk, client: client, UserId: userId.String()}

	return user, nil
}

func (u *User) Login() (string, error) {
	url := fmt.Sprintf("%s/v1_0/apiuser/%s/apikey", BaseURL, u.UserId)
	req, err := http.NewRequest("POST", url, nil)

	if err != nil {
		return "", errors.Wrap(err, "error creating request")
	}

	req.Header.Set("Ocp-Apim-Subscription-Key", u.pk)
	res, err := u.client.reqHandler.Do(req)

	if err != nil {
		return "", err
	}

	var credentials struct {
		ApiKey string `json:"apiKey"`
	}

	err = u.client.resHandler.handleResponse(res, &credentials)

	if err != nil {
		return "", nil
	}

	return credentials.ApiKey, nil
}
