package http

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"
)

type HttpRequest struct {
	Method  string
	URL     string
	Body    interface{}
	Headers map[string]interface{}
}

type HttpClient struct {
	client *http.Client
}

var (
	once     sync.Once
	instance *HttpClient
)

func GetHttpClientInstance() *HttpClient {
	once.Do(func() {
		instance = &HttpClient{
			client: &http.Client{},
		}
	})
	return instance
}
func (c *HttpClient) SendHttpRequest(request HttpRequest) ([]byte, error) {
	jsonBody, err := json.Marshal(request.Body)
	if err != nil {
		return nil, err
	}

	httpRequest, err := http.NewRequest(request.Method, request.URL, bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, err
	}

	for name, value := range request.Headers {
		strValue, ok := value.(string)
		if ok {
			httpRequest.Header.Add(name, strValue)
		} else {
			return nil, fmt.Errorf("header value for %v is not a string", name)
		}
	}

	response, err := c.client.Do(httpRequest)
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Printf("Error closing response body: %v", err)
		}
	}(response.Body)

	if response.StatusCode != http.StatusOK {
		bodyBytes, err := io.ReadAll(response.Body)
		if err != nil {
			return nil, fmt.Errorf("unexpected status code: %v", response.StatusCode)
		}
		return nil, fmt.Errorf("unexpected status code: %v body: %v", response.StatusCode, string(bodyBytes))
	}

	return io.ReadAll(response.Body)
}
