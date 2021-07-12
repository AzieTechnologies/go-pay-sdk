package payclient

import (
	"encoding/json"
	"errors"
)

type PaymentManager struct {
	RestAPIClient       *RestAPIClient
	Config              *Config
	PaymentIntentSecret string
}

func (paymentProvider *PaymentManager) ConfirmPayment(paymentMethodDetail PaymentMethodDetail) (*PaymentIntent, error) {

	_, err := paymentProvider.CreateDevpayPaymentIntent(createDevpayRestClient(), paymentMethodDetail)
	if err != nil {
		return nil, err
	}

	paymentMethod, err := paymentProvider.CreatePaymentMethod(paymentMethodDetail)
	if err != nil {
		return nil, err
	}
	paymentIntent, err := paymentProvider.CreatePaymentIntent(paymentMethod, paymentMethodDetail)
	return paymentIntent, err
}

func createDevpayRestClient() *RestAPIClient {
	var header = make(map[string]string)
	header["Content-Type"] = "application/json"

	// Select end point
	var endPoint = "https://api.devpay.io"

	// Create Rest API
	restClient := &RestAPIClient{BaseUrl: endPoint, Headers: header}
	return restClient
}

func (manager *PaymentManager) CreateDevpayPaymentIntent(client *RestAPIClient, paymentMethodDetail PaymentMethodDetail) (bool, error) {
	dataMap := make(map[string]interface{})
	paymentIntentsInfo := make(map[string]interface{})

	paymentIntentsInfo["amount"] = paymentMethodDetail.Amount
	paymentIntentsInfo["currency"] = paymentMethodDetail.Currency
	paymentIntentsInfo["capture_method"] = "automatic"
	paymentIntentsInfo["payment_method_types"] = []string{"card"}

	requestDetails := make(map[string]interface{})

	requestDetails["DevpayId"] = manager.Config.AccountId
	if manager.Config.Sandbox {
		requestDetails["env"] = "sandbox"
	}
	requestDetails["token"] = manager.Config.Secret

	dataMap["PaymentIntentsInfo"] = paymentIntentsInfo
	dataMap["RequestDetails"] = requestDetails

	b, _ := json.Marshal(dataMap)

	var header = make(map[string]string)

	response, err := manager.RestAPIClient.Post("/v1/general/paymentintent", b, header)
	if err != nil {
		return false, err
	}

	var intentData map[string]interface{}
	err = json.Unmarshal([]byte(response), &intentData)
	if err != nil {
		return false, err
	}

	responseMap, ok := intentData["Response"].(map[string]interface{})
	if !ok {
		return false, errors.New("failed to process the reponse")
	}

	success, isOk := responseMap["status"].(bool)
	if !isOk {
		return false, errors.New("failed to process the reponse")
	}

	if success {
		return success, nil
	}
	return false, errors.New("failed to create dev-pay payment intent")
}

func (paymentProvider *PaymentManager) CreatePaymentIntent(paymentMethod *PaymentMethod, paymentMethodDetail PaymentMethodDetail) (*PaymentIntent, error) {

	dataMap := make(map[string]interface{})
	dataMap["amount"] = paymentMethodDetail.Amount
	dataMap["currency"] = paymentMethodDetail.Currency
	dataMap["capture_method"] = "automatic"
	dataMap["payment_method_types"] = []string{"card"}
	dataMap["payment_method_id"] = paymentMethod.Id
	dataMap["metadata"] = paymentMethodDetail.MetaData
	dataMap["confirm"] = true

	b, _ := json.Marshal(dataMap)

	var header = make(map[string]string)
	header["Authorization"] = "Bearer " + paymentProvider.PaymentIntentSecret

	response, err := paymentProvider.RestAPIClient.Post("/v1/payment-intents", b, header)
	if err != nil {
		return nil, err
	}

	var intent PaymentIntent
	err = json.Unmarshal([]byte(response), &intent)
	if err != nil {
		return nil, err
	}

	if len(intent.ID) <= 0 {
		return nil, errors.New(string(response))
	}
	return &intent, nil
}

func (paymentProvider *PaymentManager) CreatePaymentMethod(paymentMethodDetail PaymentMethodDetail) (*PaymentMethod, error) {

	b, _ := json.Marshal(paymentMethodDetail)
	response, err := paymentProvider.RestAPIClient.Post("/v1/payment-methods", b, make(map[string]string))
	if err != nil {
		return nil, err
	}

	var method PaymentMethod
	err = json.Unmarshal([]byte(response), &method)
	if err != nil {
		return nil, err
	}
	if len(method.Id) <= 0 {
		return nil, errors.New(string(response))
	}
	return &method, nil
}
