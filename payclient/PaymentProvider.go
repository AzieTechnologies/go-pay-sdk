package payclient

import (
	"encoding/json"
	"errors"
)

type PaymentProvider struct {
	RestAPIClient *RestAPIClient
}

func (paymentProvider *PaymentProvider) PaysafeAPIKey() (string, error) {
	response, err := paymentProvider.RestAPIClient.Get("/v1/payment-providers/paysafe/api-key")

	if err != nil {
		return "", err
	}

	var objMap map[string]string
	err = json.Unmarshal([]byte(response), &objMap)
	if err != nil {
		return "", err
	}
	providerApiKey, ok := objMap["provider_api_key"]
	if !ok {
		return "", errors.New(string(response))
	}
	return providerApiKey, nil
}
