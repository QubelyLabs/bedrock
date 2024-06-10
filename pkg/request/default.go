package request

import (
	"bytes"
	"encoding/json"
	"net/http"
)

var (
	Default = NewRequest()
)

// Post sends a POST request with a JSON body
func Post(url string, body interface{}, queries map[string]string, headers map[string]string, retryCount int) (*HttpResponse, error) {
	bodyBytes, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	return Default.Request(http.MethodPost, url, bytes.NewReader(bodyBytes), queries, headers, retryCount)
}

// Get sends a GET request
func Get(url string, queries map[string]string, headers map[string]string, retryCount int) (*HttpResponse, error) {
	return Default.Request(http.MethodGet, url, nil, queries, headers, retryCount)
}

// Patch sends a PATCH request with a JSON body
func Patch(url string, body interface{}, queries map[string]string, headers map[string]string, retryCount int) (*HttpResponse, error) {
	bodyBytes, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	return Default.Request(http.MethodPatch, url, bytes.NewReader(bodyBytes), queries, headers, retryCount)
}

// Put sends a PUT request with a JSON body
func Put(url string, body interface{}, queries map[string]string, headers map[string]string, retryCount int) (*HttpResponse, error) {
	bodyBytes, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	return Default.Request(http.MethodPut, url, bytes.NewReader(bodyBytes), queries, headers, retryCount)
}

// Delete sends a DELETE request
func Delete(url string, queries map[string]string, headers map[string]string, retryCount int) (*HttpResponse, error) {
	return Default.Request(http.MethodDelete, url, nil, queries, headers, retryCount)
}
