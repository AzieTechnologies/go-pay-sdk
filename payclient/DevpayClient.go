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
	TilledEndPoint       = "https://api.tilled.com"
	TilledSanboxEndPoint = "https://sandbox-api.tilled.com"

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

	// Get Paysafe API key
	provider := &PaymentProvider{RestAPIClient: restClient}
	paysafeAPIKey, err := provider.PaysafeAPIKey()

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

	// Create payment manager
	paymentManager := &PaymentManager{RestAPIClient: restClient,
		Config:              devpayClient.Config,
		PaymentIntentSecret: devpayClient.Config.Secret}

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
	header["Authorization"] = "Bearer " + devpayClient.Config.ShareableKey
	header["Tilled-Account"] = devpayClient.Config.AccountId
	header["Content-Type"] = "application/json"

	// Select end point
	var endPoint = TilledEndPoint
	if devpayClient.Config.Sandbox {
		endPoint = TilledSanboxEndPoint
	}

	// Create Rest API
	restClient := &RestAPIClient{BaseUrl: endPoint, Headers: header}
	return restClient
}
