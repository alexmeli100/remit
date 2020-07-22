package mtn

import (
	uuid "github.com/satori/go.uuid"
	"testing"
)

func TestTransfer(t *testing.T) {
	r := createRemittance()

	tests := []struct {
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

	for i, tc := range tests {
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

func createRemittance() *Remittance {
	r := NewRemittance(&Config{
		primaryKey: "15294124013343908f01ab8bbb0d95f4",
		apiSecret:  "06bb32fe271d467d8d26ec31c6c2d6a3",
		userId:     "8cb50e22-47ee-4047-b7fc-a0b530f54568",
	})

	return r
}

func createTransferRequest(n string) *TransferRequest {
	id := uuid.NewV4()
	payee := &Payee{
		PartyIdType: "MSISDN",
		PartyId:     n,
	}

	return &TransferRequest{
		Amount:       "500",
		Currency:     "EUR",
		ExternalID:   id.String(),
		Payee:        payee,
		PayerMessage: "test message",
		PayeeNote:    "test note",
	}
}
