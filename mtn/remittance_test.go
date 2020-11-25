package mtn

import (
	"github.com/google/uuid"
	"os"
	"testing"
)

var r *Remittance

func TestMain(m *testing.M) {
	momo := CreateMomoApp(&GlobalConfig{})
	r = momo.NewRemittance(&ProductConfig{
		PrimaryKey: "15294124013343908f01ab8bbb0d95f4",
		ApiSecret:  "87d26b2e19a048ce84c5634e2ceee6b7",
		UserId:     "17465d3b-23ab-454c-82da-67ac6f09f0ce",
	})

	code := m.Run()
	os.Exit(code)
}

func TestTransfer(t *testing.T) {

	testCases := []struct {
		testCase       string
		expectedStatus string
	}{
		{"46733123450", TransferFailed},
		{"46733123451", TransferFailed},
		{"46733123452", TransferFailed},
		{"46733123453", TransferPending},
		{"46733123454", TransferPending},
		{"+256776564739", TransferSuccessFul},
	}

	for i, tc := range testCases {
		tr := createTransferRequest(tc.testCase)
		refId, err := r.Transfer(tr)

		if err != nil {
			t.Errorf("%s\n", err.Error())
		}

		res, err := r.GetTransactionStatus(refId)

		if err != nil {
			t.Errorf("%s\n", err.Error())
		} else {
			if res.Status != tc.expectedStatus {
				t.Errorf("[Test %d]wrong status, expected: %s, got: %s, reason: %s", i, tc.expectedStatus, res.Status, res.Reason)
			}
		}
	}
}

func TestSendTo(t *testing.T) {
	tests := []struct {
		testCase string
		success  bool
	}{
		{"46733123450", false},
		{"46733123451", false},
		{"46733123452", false},
		{"+256776564739", true},
	}

	for i, tc := range tests {
		err := r.SendTo(500, tc.testCase, "EUR")

		if err != nil && tc.success {
			t.Errorf("[Test %d]wrong status, expected error to be nil, got %s", i, err.Error())
		}
	}
}

func createTransferRequest(n string) *TransferRequest {
	id := uuid.New()
	payee := &Payee{
		PartyIdType: "MSISDN",
		PartyId:     n,
	}

	return &TransferRequest{
		Amount:       "50000",
		Currency:     "EUR",
		ExternalID:   id.String(),
		Payee:        payee,
		PayerMessage: "test message",
		PayeeNote:    "test note",
	}
}
