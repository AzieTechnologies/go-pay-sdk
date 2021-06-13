package payclient

import (
	"bytes"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)

type RestAPIClient struct {
	BaseUrl    string
	Headers    map[string]string
	InfoLogger *log.Logger
}

func (restAPIClient *RestAPIClient) Get(path string) ([]byte, error) {

	URL, err := url.Parse(restAPIClient.BaseUrl + path)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("GET", URL.String(), nil)
	if err != nil {
		return nil, err
	}

	for key, value := range restAPIClient.Headers {
		req.Header.Add(key, value)
	}

	restAPIClient.logRequestDetails("Get", URL.String(), restAPIClient.Headers, make([]byte, 0))
	resp, err := http.DefaultClient.Do(req)

	if err == nil {

		defer resp.Body.Close()
		contents, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		restAPIClient.logResponseDetails(err, contents)
		if resp.StatusCode == http.StatusOK {
			return contents, nil
		} else {
			return nil, errors.New(string(contents))
		}

	}
	return nil, err
}

func (restAPIClient *RestAPIClient) Post(path string, data []byte, headers map[string]string) ([]byte, error) {

	URL, err := url.Parse(restAPIClient.BaseUrl + path)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", URL.String(), bytes.NewReader(data))
	if err != nil {
		return nil, err
	}

	// Merge keys
	for k, v := range headers {
		restAPIClient.Headers[k] = v
	}
	for key, value := range restAPIClient.Headers {
		req.Header.Add(key, value)
	}

	restAPIClient.logRequestDetails("Post", URL.String(), restAPIClient.Headers, data)
	resp, err := http.DefaultClient.Do(req)

	if err == nil {
		defer resp.Body.Close()
		contents, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		restAPIClient.logResponseDetails(err, contents)

		if resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusCreated {
			return contents, nil
		} else {
			return nil, errors.New(string(contents))
		}
	}
	return nil, err
}

func (restAPIClient *RestAPIClient) logResponseDetails(err error, resp []byte) {
	if restAPIClient.InfoLogger != nil {
		restAPIClient.InfoLogger.Printf("\n Response - %s \n", string(resp))
	}
}

func (restAPIClient *RestAPIClient) logRequestDetails(method string, url string, headers map[string]string, body []byte) {
	if restAPIClient.InfoLogger != nil {
		restAPIClient.InfoLogger.Printf("\n ------------- \n Method - %s\n URL - %s\n Headers - %s \n Data - %s \n", method, url, headers, string(body))
	}
}
