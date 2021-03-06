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

type DevpayClient struct {
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

	// API end-points
	DevpayAPI = "https://api.devpay.io"

	// Info logger, used when EnableVerboseLogging set to true
	InfoLogger *log.Logger
)

func New(config *Config) *DevpayClient {
	if config == nil {
		config = &Config{}
		config.Sandbox = false
		config.EnableVerboseLogging = false
	}

	client := &DevpayClient{
		Config: config,
	}
	return client
}

func (devpayClient *DevpayClient) ConfirmPayment(paymntDetail PaymentDetail) (paymentResponse *PaymentIntent, err error) {

	var InfoLogger *log.Logger = nil
	// Create InfoLogger if verbose logging is enabled
	if devpayClient.Config.EnableVerboseLogging {
		InfoLogger = log.New(os.Stdout, "INFO ", log.Ldate|log.Ltime|log.Lshortfile)
	}

	restClient := devpayClient.createRestClient()
	restClient.InfoLogger = InfoLogger

	// Create payment manager
	paymentManager := &PaymentManager{RestAPIClient: restClient,
		Config: devpayClient.Config}

	paysafeAPIKey, err := paymentManager.PaysafeAPIKey()
	if err != nil {
		return nil, err
	}

	// Get payment token
	paysafe := &Paysafe{Sandbox: devpayClient.Config.Sandbox,
		APIKey:     paysafeAPIKey,
		InfoLogger: InfoLogger}

	paymentToken, err := paysafe.Tokenize(paymntDetail)
	if err != nil {
		return nil, err
	}

	bd := BillingDetail{BillingAddress: paymntDetail.BillingAddress,
		Name:  paymntDetail.Name,
		Email: paymntDetail.Email,
		Phone: paymntDetail.Phone}

	payMethodDetail := PaymentMethodDetail{Amount: paymntDetail.Amount,
		Currency:      paymntDetail.Currency,
		Type:          "card",
		PaymentToken:  paymentToken,
		BillingDetail: bd,
		MetaData:      paymntDetail.MetaData}

	// Confirm payment
	return paymentManager.ConfirmPayment(payMethodDetail)
}

func (devpayClient *DevpayClient) createRestClient() *RestAPIClient {
	var header = make(map[string]string)
	header["Content-Type"] = "application/json"

	// Create Rest API
	restClient := &RestAPIClient{BaseUrl: DevpayAPI, Headers: header}
	return restClient
}
