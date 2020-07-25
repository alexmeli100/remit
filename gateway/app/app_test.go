package app

import (
	"bytes"
	"testing"
)

func TestDecodeBody(t *testing.T) {
	req := `
		{
		    "user": {
		        "firstName": "Alex",
		        "lastName": "Meli",
		        "email": "alexmeli600@gmail.com",
		        "uuid": "89ypfh98eo348q-opdg09",
		        "country": "canada"
		    },
		    "password": "winniealex"
		}
		`

	r, err := decodeBody(bytes.NewBuffer([]byte(req)))

	if err != nil {
		t.Errorf("decode error: %v\n", err)
	} else {
		t.Log(r.FirstName)
	}

}
