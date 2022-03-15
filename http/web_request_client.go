package http

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	go_http "net/http"
	"net/url"
)

type webRequestClient struct {
	client        *go_http.Client
	marshalFunc   func(v interface{}) ([]byte, error)
	unmarshalFunc func(data []byte, v interface{}) error
}

// NewWebRequestClient creates a wrapper utility which handles http communication.
func NewWebRequestClient(client *go_http.Client, marshalFunc func(v interface{}) ([]byte, error), unmarshalFunc func(data []byte, v interface{}) error) *webRequestClient {
	return &webRequestClient{
		client:        client,
		marshalFunc:   marshalFunc,
		unmarshalFunc: unmarshalFunc,
	}
}

// Get sends a GET http request.
func (w *webRequestClient) Get(ctx context.Context, uri string, headers map[string]string, queryParams map[string]interface{}, responseParser interface{}) (resHeaders go_http.Header, resBody interface{}, statusCode int, err error) {
	if queryParams != nil {
		params := url.Values{}
		for k, v := range queryParams {
			params.Add(k, fmt.Sprint(v))
		}
		uri = uri + "?" + params.Encode()
	}

	httpReq, err := go_http.NewRequestWithContext(ctx, "GET", uri, nil)
	if err != nil {
		errStr := fmt.Errorf("could not create new request: %w", err)
		return nil, nil, 0, errStr
	}

	for k, v := range headers {
		httpReq.Header.Set(k, v)
	}

	httpRes, err := w.client.Do(httpReq)
	if err != nil {
		errStr := fmt.Errorf("error executing request: %w", err)
		return nil, nil, 0, errStr
	}

	defer httpRes.Body.Close()

	bodyBytes, err := ioutil.ReadAll(httpRes.Body)
	if err != nil {
		errStr := fmt.Errorf("could not read response body: %w", err)
		return nil, nil, 0, errStr
	}

	err = w.unmarshalFunc(bodyBytes, responseParser)
	if err != nil {
		errStr := fmt.Errorf("could not unmarshal response into input interface: %w", err)
		return nil, nil, 0, errStr
	}
	return httpRes.Header, responseParser, httpRes.StatusCode, nil
}

// Post sends a POST http request using a struct as payload.
//
// Use PostSerializedBody method if your payload input is string.
func (w *webRequestClient) Post(ctx context.Context, uri string, headers map[string]string, queryParams map[string]interface{}, request, responseParser interface{}) (resHeaders go_http.Header, resBody interface{}, statusCode int, err error) {
	reqAsBytes, err := w.marshalFunc(request)
	if err != nil {
		errStr := fmt.Errorf("could not convert request to byte array: %w", err)
		return nil, nil, 0, errStr
	}

	if queryParams != nil {
		params := url.Values{}
		for k, v := range queryParams {
			params.Add(k, fmt.Sprint(v))
		}
		uri = uri + "?" + params.Encode()
	}

	httpReq, err := go_http.NewRequestWithContext(ctx, "POST", uri, bytes.NewBuffer(reqAsBytes))
	if err != nil {
		errStr := fmt.Errorf("could not create new request: %w", err)
		return nil, nil, 0, errStr
	}

	for k, v := range headers {
		httpReq.Header.Set(k, v)
	}

	httpRes, err := w.client.Do(httpReq)
	if err != nil {
		errStr := fmt.Errorf("error executing request: %w", err)
		return nil, nil, 0, errStr
	}

	defer httpRes.Body.Close()

	bodyBytes, err := ioutil.ReadAll(httpRes.Body)
	if err != nil {
		errStr := fmt.Errorf("could not read response body: %w", err)
		return nil, nil, 0, errStr
	}

	err = w.unmarshalFunc(bodyBytes, responseParser)
	if err != nil {
		errStr := fmt.Errorf("could not unmarshal response into input interface: %w", err)
		return nil, nil, 0, errStr
	}
	return httpRes.Header, responseParser, httpRes.StatusCode, nil
}

// PostSerializedBody sends a POST http request with a string payload.
func (w *webRequestClient) PostSerializedBody(ctx context.Context, uri string, headers map[string]string, queryParams map[string]interface{}, request string, responseParser interface{}) (resHeaders go_http.Header, resBody interface{}, statusCode int, err error) {
	if queryParams != nil {
		params := url.Values{}
		for k, v := range queryParams {
			params.Add(k, fmt.Sprint(v))
		}
		uri = uri + "?" + params.Encode()
	}

	httpReq, err := go_http.NewRequestWithContext(ctx, "POST", uri, bytes.NewBuffer([]byte(request)))
	if err != nil {
		errStr := fmt.Errorf("could not create new request: %w", err)
		return nil, nil, 0, errStr
	}

	for k, v := range headers {
		httpReq.Header.Set(k, v)
	}

	httpRes, err := w.client.Do(httpReq)
	if err != nil {
		errStr := fmt.Errorf("error executing request: %w", err)
		return nil, nil, 0, errStr
	}

	defer httpRes.Body.Close()

	bodyBytes, err := ioutil.ReadAll(httpRes.Body)
	if err != nil {
		errStr := fmt.Errorf("could not read response body: %w", err)
		return nil, nil, 0, errStr
	}

	err = w.unmarshalFunc(bodyBytes, responseParser)
	if err != nil {
		errStr := fmt.Errorf("could not unmarshal response into input interface: %w", err)
		return nil, nil, 0, errStr
	}
	return httpRes.Header, responseParser, httpRes.StatusCode, nil
}

func (w *webRequestClient) CreateBasicAuthHeaderValue(username, password string) string {
	auth := username + ":" + password
	encoded := base64.StdEncoding.EncodeToString([]byte(auth))
	return "Basic " + encoded
}
