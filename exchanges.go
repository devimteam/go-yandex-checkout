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
	Id            string      `json:"id"`
	Status        string      `json:"status"`
	Paid          bool        `json:"paid"`
	Amount        Amount      `json:"amount"`
	CreatedAt     time.Time   `json:"created_at"`
	Metadata      interface{} `json:"metadata"`
	PaymentMethod struct {
		Type  string `json:"type"`
		Id    string `json:"id"`
		Saved bool   `json:"saved"`
		Card  Card   `json:"card"`
		Title string `json:"title"`
	} `json:"payment_method"`
}
