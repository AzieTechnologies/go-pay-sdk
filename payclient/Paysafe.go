package payclient

import (
	"encoding/json"
	"errors"
	"log"
)

type Paysafe struct {
	Sandbox    bool
	APIKey     string
	InfoLogger *log.Logger
}

func (paysafe *Paysafe) Tokenize(paymntDetail PaymentDetail) (string, error) {

	var path = "/js/api/v1/tokenize"
	restClient := paysafe.createRestClient()
	restClient.InfoLogger = paysafe.InfoLogger
	b, err := json.Marshal(paymntDetail)

	if err != nil {
		return "", err
	}

	response, err := restClient.Post(path, b, make(map[string]string))

	if err != nil {
		return "", err
	}

	var objMap map[string]string
	err = json.Unmarshal([]byte(response), &objMap)
	if err != nil {
		return "", err
	}

	paymentToken, ok := objMap["paymentToken"]
	if !ok {
		return "", errors.New(string(response))
	}
	return paymentToken, nil
}

func (paysafe *Paysafe) createRestClient() *RestAPIClient {
	var BaseUrl = "https://hosted.paysafe.com"
	if paysafe.Sandbox {
		BaseUrl = "https://hosted.test.paysafe.com"
	}
	var header = make(map[string]string)
	header["X-Paysafe-Credentials"] = "Basic " + paysafe.APIKey
	header["Content-Type"] = "application/json"
	restClient := &RestAPIClient{BaseUrl: BaseUrl, Headers: header}
	return restClient
}
