package go_yandex_money

import "time"

const (
	CARD_TYPE__MASTER_CARD = "MasterCard"
)

type Amount struct {
	Value    float64 `json:"value"`
	Currency string  `json:"currency"`
}

type Card struct {
	Last4       string `json:"last4"`
	ExpiryMonth int    `json:"expiry_month"`
	ExpiryYear  int    `json:"expiry_month"`
	CardType    string `json:"card_type"`
}

type PaymentMethod struct {
	Type  string `json:"type"`
	Id    string `json:"id"`
	Saved bool   `json:"saved"`
	Card  Card   `json:"card"`
	Title string `json:"title"`
}

type CreatePaymentRequest struct {
	Amount            Amount `json:"amount"`
	PaymentMethodData struct {
		Type string `json:"type"`
	} `json:"payment_method_data"`
	Confirmation struct {
		Type      string `json:"type"`
		ReturnUrl string `json:"return_url"`
	} `json:"confirmation"`
}

type CreatePaymentResponse struct {
	Id            string        `json:"id"`
	Status        string        `json:"status"`
	Paid          bool          `json:"paid"`
	Amount        Amount        `json:"amount"`
	CreatedAt     time.Time     `json:"created_at"`
	Metadata      interface{}   `json:"metadata"`
	PaymentMethod PaymentMethod `json:"payment_method"`
}

const (
	PAYMENT_STATUS__WAIT      = "waiting_for_capture"
	PAYMENT_STATUS__PENDING   = "pending"
	PAYMENT_STATUS__SUCCEEDED = "succeeded"
	PAYMENT_STATUS__CANCELED  = "canceled"
)

type PaymentResponse struct {
	Id            string        `json:"id"`
	Status        string        `json:"status"`
	Paid          bool          `json:"paid"`
	Amount        Amount        `json:"amount"`
	Description   *string       `json:"description"`
	CreatedAt     time.Time     `json:"created_at"`
	Metadata      interface{}   `json:"metadata"`
	PaymentMethod PaymentMethod `json:"payment_method"`
}

type GetPaymentResponse PaymentResponse

type CapturePaymentRequest struct {
	Amount Amount `json:"amount"`
}

type CapturePaymentResponse PaymentResponse

type CancelPaymentRequest struct{}

type CancelPaymentResponse PaymentResponse
