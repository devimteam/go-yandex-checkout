package go_yandex_checkout

import "time"

type ProcessingResponse struct {
	Type        string `json:"type"`
	Description string `json:"description"`
	RetryAfter  *uint  `json:"retry_after"` // in milliseconds

	Id        *string `json:"id"`
	Code      *string `json:"code"`
	Parameter *string `json:"parameter"`
}

func (pr ProcessingResponse) GetSleepTime() time.Duration {
	if nil == pr.RetryAfter {
		return 0
	}
	return time.Millisecond * time.Duration(*pr.RetryAfter)
}

const (
	CARD_TYPE__MASTER_CARD = "MasterCard"
)

type Amount struct {
	Value    StrFloat64 `json:"value"`
	Currency string     `json:"currency"`
}

type Card struct {
	Last4       string   `json:"last4"`
	ExpiryMonth StrInt64 `json:"expiry_month"`
	ExpiryYear  StrInt64 `json:"expiry_year"`
	CardType    string   `json:"card_type"`
}

type Receipt struct {
	Items []struct {
		Description string `json:"description"`
		Quantity    string `json:"quantity"`
		Amount      Amount `json:"amount"`
		VatCode     int    `json:"vat_code"` // https://kassa.yandex.ru/docs/guides/#kody-stawok-nds
	} `json:"items"`
	TaxSystemCode *int   `json:"tax_system_code"`
	Phone         string `json:"phone"` // https://ru.wikipedia.org/wiki/E.164
	Email         string
}

const (
	CONFIRMATION_TYPE__REDIRECT = "redirect"
	CONFIRMATION_TYPE__EXTERNAL = "external"
)

type Confirmation struct {
	Type            string `json:"type"`
	Enforce         bool   `json:"enforce"`
	ReturnUrl       string `json:"return_url"`
	ConfirmationUrl string `json:"confirmation_url"`
}

type Metadata map[string]interface{}

const (
	PAYMENT_METHOD_TYPE__SBERBANK       = "sberbank"
	PAYMENT_METHOD_TYPE__BANK_CARD      = "bank_card"
	PAYMENT_METHOD_TYPE__CASH           = "cash"
	PAYMENT_METHOD_TYPE__YANDEX_MONEY   = "yandex_money"
	PAYMENT_METHOD_TYPE__QIWI           = "qiwi"
	PAYMENT_METHOD_TYPE__ALFABANK       = "alfabank"
	PAYMENT_METHOD_TYPE__WEBMONEY       = "webmoney"
	PAYMENT_METHOD_TYPE__APPLE_PAY      = "apple_pay"
	PAYMENT_METHOD_TYPE__MOBILE_BALANCE = "mobile_balance"
)

type PaymentMethod struct {
	Type  string `json:"type"`
	ID    string `json:"id"`
	Saved bool   `json:"saved"`
	Card  Card   `json:"card"`
	Title string `json:"title"`
}

type CreatePaymentRequest struct {
	Amount            Amount `json:"amount"`
	PaymentMethodData struct {
		Type string `json:"type"`
	} `json:"payment_method_data"`
	Confirmation Confirmation `json:"confirmation"`
}

type CreatePaymentResponse struct {
	ID            string        `json:"id"`
	Status        string        `json:"status"`
	Paid          bool          `json:"paid"`
	Amount        Amount        `json:"amount"`
	CreatedAt     time.Time     `json:"created_at"`
	Metadata      Metadata      `json:"metadata"`
	PaymentMethod PaymentMethod `json:"payment_method"`
}

const (
	PAYMENT_STATUS__WAIT      = "waiting_for_capture"
	PAYMENT_STATUS__PENDING   = "pending"
	PAYMENT_STATUS__SUCCEEDED = "succeeded"
	PAYMENT_STATUS__CANCELED  = "canceled"
)

type Recipient struct {
	GatewayID string `json:"gateway_id"`
}

type PaymentResponse struct {
	ID                  string        `json:"id"`
	Status              string        `json:"status"`
	Amount              Amount        `json:"amount"`
	Description         *string       `json:"description"`
	Recipient           *Recipient    `json:"recipient"`
	PaymentMethod       PaymentMethod `json:"payment_method"`
	CapturedAt          *time.Time    `json:"captured_at"`
	CreatedAt           time.Time     `json:"created_at"`
	ExpiresAt           *time.Time    `json:"expires_at"`
	Confirmation        *Confirmation `json:"confirmation"`
	Test                *bool         `json:"test"`
	RefundedAmount      *Amount       `json:"refunded_amount"`
	Paid                bool          `json:"paid"`
	ReceiptRegistration *string       `json:"receipt_registration"`
	Metadata            Metadata      `json:"metadata"`
}

type GetPaymentResponse PaymentResponse

type CapturePaymentRequest struct {
	Amount  Amount  `json:"amount"`
	Receipt Receipt `json:"receipt"`
}

type CapturePaymentResponse PaymentResponse

type CancelPaymentRequest struct{}

type CancelPaymentResponse PaymentResponse

type CreateRefundRequest struct {
	PaymentID   string  `json:"PaymentID"`
	Amount      Amount  `json:"amount"`
	Description *string `json:"description"`
	Receipt     Receipt `json:"receipt"`
}

const (
	REFUND_STATUS__CANCELED  = "canceled"
	REFUND_STATUS__SUCCEEDED = "succeeded"
)

const (
	REFUND_RECEIPT_REGISTRATION_PENDING   = "pending"
	REFUND_RECEIPT_REGISTRATION_SUCCEEDED = "succeeded"
	REFUND_RECEIPT_REGISTRATION_CANCELED  = "canceled"
)

type RefundResponse struct {
	ID                  string    `json:"id"`
	PaymentID           string    `json:"PaymentID"`
	Status              string    `json:"status"`
	CreatedAt           time.Time `json:"created_at"`
	Amount              Amount    `json:"amount"`
	ReceiptRegistration string    `json:"receipt_registration"`
	Description         *string   `json:"description"`
}

type CreateRefundResponse RefundResponse

type GetRefundResponse RefundResponse
