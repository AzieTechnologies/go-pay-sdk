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

func (tillerRestAPI *RestAPIClient) Get(path string) ([]byte, error) {

	URL, _ := url.Parse(tillerRestAPI.BaseUrl + path)
	req, _ := http.NewRequest("GET", URL.String(), nil)
	for key, value := range tillerRestAPI.Headers {
		req.Header.Add(key, value)
	}

	tillerRestAPI.logRequestDetails("Get", URL.String(), tillerRestAPI.Headers, make([]byte, 0))
	resp, err := http.DefaultClient.Do(req)

	if err == nil {

		defer resp.Body.Close()
		contents, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		tillerRestAPI.logResponseDetails(err, contents)
		if resp.StatusCode == http.StatusOK {
			return contents, nil
		} else {
			return nil, errors.New(string(contents))
		}

	}
	return nil, err
}

func (tillerRestAPI *RestAPIClient) Post(path string, data []byte, headers map[string]string) ([]byte, error) {

	URL, _ := url.Parse(tillerRestAPI.BaseUrl + path)
	req, _ := http.NewRequest("POST", URL.String(), bytes.NewReader(data))

	for k, v := range headers {
		tillerRestAPI.Headers[k] = v
	}
	for key, value := range tillerRestAPI.Headers {
		req.Header.Add(key, value)
	}

	tillerRestAPI.logRequestDetails("Post", URL.String(), tillerRestAPI.Headers, data)
	resp, err := http.DefaultClient.Do(req)

	if err == nil {
		defer resp.Body.Close()
		contents, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		tillerRestAPI.logResponseDetails(err, contents)

		if resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusCreated {
			return contents, nil
		} else {
			return nil, errors.New(string(contents))
		}
	}
	return nil, err
}

func (tillerRestAPI *RestAPIClient) logResponseDetails(err error, resp []byte) {
	if tillerRestAPI.InfoLogger != nil {
		tillerRestAPI.InfoLogger.Printf("\n Response - %s \n", string(resp))
	}
}

func (tillerRestAPI *RestAPIClient) logRequestDetails(method string, url string, headers map[string]string, body []byte) {
	if tillerRestAPI.InfoLogger != nil {
		tillerRestAPI.InfoLogger.Printf("\n ------------- \n Method - %s\n URL - %s\n Headers - %s \n Data - %s \n", method, url, headers, string(body))
	}
}
