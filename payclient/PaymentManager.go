package payclient

import (
	"encoding/json"
	"errors"
)

type PaymentManager struct {
	TillerRestClient    *RestAPIClient
	PaymentIntentSecret string
}

func (paymentProvider *PaymentManager) ConfirmPayment(paymentMethodDetail PaymentMethodDetail) (*PaymentIntent, error) {
	paymentMethod, err := paymentProvider.CreatePaymentMethod(paymentMethodDetail)
	if err != nil {
		return nil, err
	}
	paymentIntent, err := paymentProvider.CreatePaymentIntent(paymentMethod, paymentMethodDetail)
	return paymentIntent, err
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

	response, err := paymentProvider.TillerRestClient.Post("/v1/payment-intents", b, header)
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
	response, err := paymentProvider.TillerRestClient.Post("/v1/payment-methods", b, make(map[string]string))
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
