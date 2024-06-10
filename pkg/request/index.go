package request

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/QubelyLabs/bedrock/pkg/util"
)

// HttpResponse represents the response structure
type HttpResponse struct {
	Status     bool        `json:"status"`
	StatusCode int         `json:"status_code"`
	StatusText string      `json:"status_text"`
	Title      string      `json:"title,omitempty"`
	Message    string      `json:"message,omitempty"`
	Data       util.Object `json:"data,omitempty"`
}

const (
	timeout = 10 * time.Second
)

type request struct {
	client *http.Client
}

func NewRequest() *request {
	return &request{
		client: &http.Client{Timeout: timeout},
	}
}

func (s *request) Request(method, url string, body *bytes.Reader, queries map[string]string, headers map[string]string, retryCount int) (*HttpResponse, error) {
	var req *http.Request
	var err error

	if method == "GET" || body == nil {
		req, err = http.NewRequest(method, url, nil)
	} else {
		req, err = http.NewRequest(method, url, body)
	}

	if err != nil {
		return nil, err
	}

	q := req.URL.Query()
	for key, value := range queries {
		q.Add(key, value)
	}
	req.URL.RawQuery = q.Encode()

	for key, value := range headers {
		req.Header.Set(key, value)
	}

	var response *http.Response
	response, err = s.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	data := util.Object{}

	if response.Body != nil {
		if response.Header.Get("Content-Encoding") == "gzip" {
			reader, err := gzip.NewReader(response.Body)
			if err != nil {
				return nil, err
			}
			defer reader.Close()

			err = json.NewDecoder(reader).Decode(&data)
			if err != nil {
				return nil, err
			}
		} else {
			err := json.NewDecoder(response.Body).Decode(&data)
			if err != nil {
				return nil, err
			}
		}
	}

	if response.StatusCode < 200 || response.StatusCode > 299 {
		return &HttpResponse{false, response.StatusCode, response.Status, "Error", "Something went wrong", data}, errors.New(response.Status)
	}

	return &HttpResponse{true, response.StatusCode, response.Status, "", "", data}, err
}
