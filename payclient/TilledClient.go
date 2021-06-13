package payclient

import (
	"log"
	"os"
)

type Currency string

const (
	USD  = "usd"
	AUD  = "aud"
	CAD  = "cad"
	DKK  = "dkk"
	EUR  = "eur"
	HKD  = "hkd"
	JPY  = "jpy"
	NZD  = "nzd"
	NOK  = "nok"
	GBD  = "gbp"
	ZAR  = "zar"
	SEK  = "sek"
	XCHF = "chf"
)

type TilledClient struct {
	Config *Config
}

type Config struct {
	ShareableKey         string
	Secret               string
	AccountId            string
	Sandbox              bool
	EnableVerboseLogging bool
}

var (
	TillerEndPoint       = "https://api.tilled.com"
	TillerSanboxEndPoint = "https://sandbox-api.tilled.com"
)

var (
	InfoLogger *log.Logger
)

func New(config *Config) *TilledClient {
	if config == nil {
		config = &Config{}
		config.Sandbox = false
		config.EnableVerboseLogging = false
	}

	client := &TilledClient{
		Config: config,
	}
	return client
}

func (tilledClient *TilledClient) ConfirmPayment(paymntDetail PaymentDetail) (paymentResponse *PaymentIntent, err error) {

	var InfoLogger *log.Logger = nil
	// Create InfoLogger if verbose logging is enabled
	if tilledClient.Config.EnableVerboseLogging {
		InfoLogger = log.New(os.Stdout, "INFO ", log.Ldate|log.Ltime|log.Lshortfile)
	}

	restClient := tilledClient.createRestClient()
	restClient.InfoLogger = InfoLogger

	// Get Paysafe API key
	provider := &PaymentProvider{TillerRestClient: restClient}
	paysafeAPIKey, err := provider.PaysafeAPIKey()

	if err != nil {
		return nil, err
	}

	// Get payment token
	paysafe := &Paysafe{Sandbox: tilledClient.Config.Sandbox,
		APIKey:     paysafeAPIKey,
		InfoLogger: InfoLogger}

	paymentToken, err := paysafe.Tokenize(paymntDetail)
	if err != nil {
		return nil, err
	}

	// Create payment manager
	paymentManager := &PaymentManager{TillerRestClient: restClient,
		PaymentIntentSecret: tilledClient.Config.Secret}

	bd := BillingDetail{BillingAddress: paymntDetail.BillingAddress,
		Name:  paymntDetail.Name,
		Email: paymntDetail.Email,
		Phone: paymntDetail.Phone}

	pmd := PaymentMethodDetail{Amount: paymntDetail.Amount,
		Currency:      paymntDetail.Currency,
		Type:          "card",
		PaymentToken:  paymentToken,
		BillingDetail: bd,
		MetaData:      paymntDetail.MetaData}

	// Confirm payment
	paymentIntent, err := paymentManager.ConfirmPayment(pmd)

	if err != nil {
		return nil, err
	}

	return paymentIntent, nil
}

func (tilledClient *TilledClient) createRestClient() *RestAPIClient {
	var header = make(map[string]string)
	header["Authorization"] = "Bearer " + tilledClient.Config.ShareableKey
	header["Tilled-Account"] = tilledClient.Config.AccountId
	header["Content-Type"] = "application/json"

	// Select end point
	var endPoint = TillerEndPoint
	if tilledClient.Config.Sandbox {
		endPoint = TillerSanboxEndPoint
	}

	// Create Rest API
	restClient := &RestAPIClient{BaseUrl: endPoint, Headers: header}
	return restClient
}
