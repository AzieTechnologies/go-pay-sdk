package payclient

import (
	"encoding/json"
	"errors"
)

type PaymentManager struct {
	RestAPIClient *RestAPIClient
	Config        *Config
}

func (paymentManager *PaymentManager) PaysafeAPIKey() (string, error) {

	var data = make(map[string]interface{})

	var requestPayload = make(map[string]interface{})
	requestPayload["DevpayId"] = paymentManager.Config.AccountId
	requestPayload["token"] = paymentManager.Config.ShareableKey

	if paymentManager.Config.Sandbox {
		requestPayload["env"] = "sandbox"
	}

	data["RequestDetails"] = requestPayload
	var header = make(map[string]string)

	b, _ := json.Marshal(data)

	response, err := paymentManager.RestAPIClient.Post("/v1/general.svc/paysafe/api-key", b, header)

	if err != nil {
		return "", err
	}

	var objMap map[string]interface{}
	err = json.Unmarshal([]byte(response), &objMap)
	if err != nil {
		return "", err
	}
	providerApiKey, ok := objMap["provider_api_key"]
	if !ok {
		return "", errors.New(string(response))
	}
	return providerApiKey.(string), nil
}

func (paymentManager *PaymentManager) ConfirmPayment(paymentMethodDetail PaymentMethodDetail) (*PaymentIntent, error) {

	paymentMethod, err := paymentManager.CreatePaymentMethod(paymentMethodDetail)
	if err != nil {
		return nil, err
	}
	paymentIntent, err := paymentManager.CreatePaymentIntent(paymentMethod, paymentMethodDetail)
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

func (manager *PaymentManager) CreatePaymentIntent(paymentMethod *PaymentMethod, paymentMethodDetail PaymentMethodDetail) (*PaymentIntent, error) {

	paymentIntentsInfo := make(map[string]interface{})
	paymentIntentsInfo["amount"] = paymentMethodDetail.Amount
	paymentIntentsInfo["currency"] = paymentMethodDetail.Currency
	paymentIntentsInfo["capture_method"] = "automatic"
	paymentIntentsInfo["payment_method_types"] = []string{"card"}
	paymentIntentsInfo["payment_method_id"] = paymentMethod.Id
	paymentIntentsInfo["metadata"] = paymentMethodDetail.MetaData
	paymentIntentsInfo["confirm"] = true

	requestDetails := make(map[string]interface{})

	requestDetails["DevpayId"] = manager.Config.AccountId
	if manager.Config.Sandbox {
		requestDetails["env"] = "sandbox"
	}
	requestDetails["token"] = manager.Config.Secret

	payload := map[string]interface{}{
		"PaymentIntentsInfo": paymentIntentsInfo,
		"RequestDetails":     requestDetails,
	}

	b, _ := json.Marshal(payload)

	var header = make(map[string]string)
	response, err := manager.RestAPIClient.Post("/v1/general/paymentintent", b, header)
	if err != nil {
		return nil, err
	}

	paymentIntentData, err := extractData([]byte(response), "PaymentIntentsResponse")
	if err != nil {
		return nil, err
	}

	var intent PaymentIntent
	err = json.Unmarshal(paymentIntentData, &intent)
	if err != nil {
		return nil, err
	}

	if len(intent.ID) <= 0 {
		return nil, errors.New(string(response))
	}
	return &intent, nil
}

func (manager *PaymentManager) CreatePaymentMethod(paymentMethodDetail PaymentMethodDetail) (*PaymentMethod, error) {

	billingAddress := paymentMethodDetail.BillingDetail.BillingAddress

	paymentMethodInfo := map[string]interface{}{
		"payment_token": paymentMethodDetail.PaymentToken,
		"type":          "card",
		"billing_details": map[string]interface{}{
			"amount":   paymentMethodDetail.Amount,
			"currency": paymentMethodDetail.Currency,
			"address": map[string]interface{}{
				"country": billingAddress.Country,
				"state":   billingAddress.State,
				"zip":     billingAddress.Zip,
				"city":    billingAddress.City,
				"street":  billingAddress.Street,
			},
		},
	}

	requestDetails := make(map[string]interface{})

	requestDetails["DevpayId"] = manager.Config.AccountId
	if manager.Config.Sandbox {
		requestDetails["env"] = "sandbox"
	}
	requestDetails["token"] = manager.Config.Secret

	payload := map[string]interface{}{
		"PaymentMethodInfo": paymentMethodInfo,
		"RequestDetails":    requestDetails,
	}

	b, _ := json.Marshal(payload)

	response, err := manager.RestAPIClient.Post("/v1/paymentmethods/create", b, make(map[string]string))
	if err != nil {
		return nil, err
	}

	paymentMethoData, err := extractData([]byte(response), "PaymentMethodResponse")
	if err != nil {
		return nil, err
	}

	var method PaymentMethod
	err = json.Unmarshal(paymentMethoData, &method)
	if err != nil {
		return nil, err
	}
	if len(method.Id) <= 0 {
		return nil, errors.New(string(response))
	}
	return &method, nil
}

func extractData(data []byte, key string) ([]byte, error) {

	var mappedData map[string]interface{}
	err := json.Unmarshal([]byte(data), &mappedData)
	if err != nil {
		return nil, err
	}

	extractedMap, ok := mappedData[key].(map[string]interface{})
	if !ok {
		return nil, errors.New("failed to process the data")
	}

	return json.Marshal(extractedMap)
}
