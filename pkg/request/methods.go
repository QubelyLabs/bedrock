package request

import (
	"bytes"
	"encoding/json"
	"net/http"
)

// Post sends a POST request with a JSON body
func (s *request) Post(url string, body interface{}, queries map[string]string, headers map[string]string, retryCount int) (*HttpResponse, error) {
	bodyBytes, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	return s.Request(http.MethodPost, url, bytes.NewReader(bodyBytes), queries, headers, retryCount)
}

// Get sends a GET request
func (s *request) Get(url string, queries map[string]string, headers map[string]string, retryCount int) (*HttpResponse, error) {
	return s.Request(http.MethodGet, url, nil, queries, headers, retryCount)
}

// Patch sends a PATCH request with a JSON body
func (s *request) Patch(url string, body interface{}, queries map[string]string, headers map[string]string, retryCount int) (*HttpResponse, error) {
	bodyBytes, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	return s.Request(http.MethodPatch, url, bytes.NewReader(bodyBytes), queries, headers, retryCount)
}

// Put sends a PUT request with a JSON body
func (s *request) Put(url string, body interface{}, queries map[string]string, headers map[string]string, retryCount int) (*HttpResponse, error) {
	bodyBytes, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	return s.Request(http.MethodPut, url, bytes.NewReader(bodyBytes), queries, headers, retryCount)
}

// Delete sends a DELETE request
func (s *request) Delete(url string, queries map[string]string, headers map[string]string, retryCount int) (*HttpResponse, error) {
	return s.Request(http.MethodDelete, url, nil, queries, headers, retryCount)
}
