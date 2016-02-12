package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/alexyer/orderup/orderup"
)

const (
	API_PREFIX     = "api/v1"
	ERROR_RESPONSE = "error"
)

type Response struct {
	Response string
	Errors   []string
	Orders   []orderup.Order
}

func APICall(endpoint, method string, payload interface{}) (*Response, error) {
	cred, err := ReadCredentials()
	if err != nil {
		return nil, err
	}

	buf, err := encodePayload(payload)
	if err != nil {
		return nil, err
	}

	return doCall(endpoint, method, cred, buf)
}

func encodePayload(payload interface{}) ([]byte, error) {
	return json.Marshal(payload)
}

func doCall(endpoint, method string, cred *Credentials, buf []byte) (*Response, error) {
	client := http.Client{}

	req, err := http.NewRequest(method,
		fmt.Sprintf("http://%s:%d/%s%s", cred.Host, cred.Port, API_PREFIX, endpoint),
		bytes.NewReader(buf))
	if err != nil {
		return nil, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	buf, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(string(buf))
	}

	// Decode response
	content := Response{}
	if err := json.Unmarshal(buf, &content); err != nil {
		return nil, err
	}

	if content.Response == ERROR_RESPONSE {
		return nil, errors.New(content.Errors[0])
	}

	return &content, nil
}
