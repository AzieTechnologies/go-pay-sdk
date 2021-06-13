package payclient

import "time"

type PaymentDetail struct {
	Amount         int            `json:"amount"`
	Currency       Currency       `json:"currency"`
	Card           Card           `json:"card"`
	BillingAddress BillingAddress `json:"billingAddress"`
	MetaData       map[string]string
	Name           string `json:"name"`
	Email          string `json:"email"`
	Phone          string `json:"phone"`
}

type Card struct {
	CardNum    string            `json:"cardNum"`
	CardExpiry map[string]string `json:"cardExpiry"`
	Cvv        string            `json:"cvv"`
}

type BillingAddress struct {
	Country string `json:"country"`
	Zip     string `json:"zip"`
	State   string `json:"state"`
	Street  string `json:"street"`
	City    string `json:"city"`
}

type PaymentMethodDetail struct {
	PaymentToken  string            `json:"payment_token"`
	Type          string            `json:"type"`
	BillingDetail BillingDetail     `json:"billing_details"`
	Amount        int               `json:"-"`
	Currency      Currency          `json:"-"`
	MetaData      map[string]string `json:"-"`
}

type BillingDetail struct {
	Name           string         `json:"name"`
	Email          string         `json:"email,omitempty"`
	Phone          string         `json:"phone,omitempty"`
	BillingAddress BillingAddress `json:"address"`
}

type PaymentMethod struct {
	Type string `json:"type"`
	Id   string `json:"id"`
}

type PaymentIntent struct {
	AmountCapturable          int         `json:"amount_capturable"`
	AmountReceived            int         `json:"amount_received"`
	LastPaymentError          interface{} `json:"last_payment_error"`
	ID                        string      `json:"id"`
	AccountID                 string      `json:"account_id"`
	Amount                    int         `json:"amount"`
	Currency                  string      `json:"currency"`
	PaymentMethodTypes        []string    `json:"payment_method_types"`
	Status                    string      `json:"status"`
	StatementDescriptorSuffix interface{} `json:"statement_descriptor_suffix"`
	CaptureMethod             string      `json:"capture_method"`
	ClientSecret              string      `json:"client_secret"`
	PlatformFeeAmount         interface{} `json:"platform_fee_amount"`
	Metadata                  interface{} `json:"metadata"`
	CanceledAt                interface{} `json:"canceled_at"`
	CancellationReason        interface{} `json:"cancellation_reason"`
	OccurrenceType            interface{} `json:"occurrence_type"`
	CreatedAt                 time.Time   `json:"created_at"`
	UpdatedAt                 time.Time   `json:"updated_at"`
	PaymentMethod             struct {
		Card struct {
			Brand    string `json:"brand"`
			Last4    string `json:"last4"`
			ExpYear  int    `json:"exp_year"`
			ExpMonth int    `json:"exp_month"`
		} `json:"card"`
		AchDebit   interface{} `json:"ach_debit"`
		Chargeable bool        `json:"chargeable"`
		ID         string      `json:"id"`
		Type       string      `json:"type"`
		CustomerID interface{} `json:"customer_id"`
		NickName   interface{} `json:"nick_name"`
		Details    struct {
			Brand    string `json:"brand"`
			Last4    string `json:"last4"`
			ExpYear  int    `json:"exp_year"`
			ExpMonth int    `json:"exp_month"`
		} `json:"details"`
		ExpiresAt      interface{} `json:"expires_at"`
		ApplePay       bool        `json:"apple_pay"`
		BillingDetails struct {
			Name    string `json:"name"`
			Address struct {
				Zip     string `json:"zip"`
				City    string `json:"city"`
				State   string `json:"state"`
				Street  string `json:"street"`
				Country string `json:"country"`
			} `json:"address"`
		} `json:"billing_details"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
	} `json:"payment_method"`
	Charges []struct {
		AmountRefunded     int           `json:"amount_refunded"`
		Captured           bool          `json:"captured"`
		Refunded           bool          `json:"refunded"`
		ID                 string        `json:"id"`
		PaymentIntentID    string        `json:"payment_intent_id"`
		Status             string        `json:"status"`
		AmountCaptured     int           `json:"amount_captured"`
		CapturedAt         time.Time     `json:"captured_at"`
		PaymentMethodID    string        `json:"payment_method_id"`
		FailureMessage     interface{}   `json:"failure_message"`
		CreatedAt          time.Time     `json:"created_at"`
		UpdatedAt          time.Time     `json:"updated_at"`
		Refunds            []interface{} `json:"refunds"`
		PlatformFee        interface{}   `json:"platform_fee"`
		BalanceTransaction struct {
			Status      string      `json:"status"`
			ID          string      `json:"id"`
			AccountID   string      `json:"account_id"`
			Amount      int         `json:"amount"`
			Currency    string      `json:"currency"`
			Description interface{} `json:"description"`
			Fee         int         `json:"fee"`
			FeeDetails  []struct {
				Type        string `json:"type"`
				Amount      int    `json:"amount"`
				Currency    string `json:"currency"`
				Description string `json:"description"`
			} `json:"fee_details"`
			Net         int       `json:"net"`
			SourceID    string    `json:"source_id"`
			Type        string    `json:"type"`
			AvailableOn time.Time `json:"available_on"`
			CreatedAt   time.Time `json:"created_at"`
			UpdatedAt   time.Time `json:"updated_at"`
		} `json:"balance_transaction"`
	} `json:"charges"`
}
