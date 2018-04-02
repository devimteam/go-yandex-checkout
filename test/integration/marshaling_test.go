package integration

import (
	"testing"
	"github.com/devimteam/go-yandex-checkout"
	"encoding/json"
)

func TestGetPaymentResponse(t *testing.T) {
	sample := `{
  "id": "215d8da0-000f-50be-b000-0003308c89be",
  "status": "waiting_for_capture",
  "paid": true,
  "amount": {
    "value": "4.00",
    "currency": "RUB"
  },
  "created_at": "2017-09-27T10:13:52.588Z",
  "metadata": {},
  "payment_method": {
    "type": "bank_card",
    "id": "012b57fe-1904-493a-a2de-48cd42c4f356",
    "saved": false,
    "card": {
      "last4": "4448",
      "expiry_month": "04",
      "expiry_year": "2020",
      "card_type": "MasterCard"
    },
    "title": "Bank card *4448"
  }
}`
	r := &go_yandex_checkout.GetPaymentResponse{}
	json.Unmarshal([]byte(sample), r)
	expected := "215d8da0-000f-50be-b000-0003308c89be"
	if r.ID != expected {
		t.Error("Expected", expected, ",got", r.ID)
	}
	if r.Amount.Value != 4.0 {
		t.Error("Expected 4.0, got", r.Amount.Value)
	}
	if r.PaymentMethod.Saved != false {
		t.Error("Expected false, got true")
	}
}
