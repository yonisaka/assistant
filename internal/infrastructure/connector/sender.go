package connector

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gofiber/fiber/v2/log"
	"net/http"
)

// Send is a function to send HTTP Request
func (c *connector) Send(ctx context.Context, requestOption *RequestOption, result any) error {
	// Set HTTP Request Parameter
	url := fmt.Sprintf("%s%s", c.openai.BaseURL, requestOption.URL)
	log.Infow("Request Created:", "url", url, "method", requestOption.Method)

	request, err := http.NewRequest(
		requestOption.Method,
		url,
		requestOption.Body,
	)
	if err != nil {
		return err
	}

	// Set HTTP Request Header
	for k, v := range c.openai.Header {
		request.Header.Set(k, v)
	}

	for k, v := range requestOption.CustomHeader {
		request.Header.Set(k, v)
	}

	// Do HTTP Request
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	// Decode HTTP Response
	err = json.NewDecoder(response.Body).Decode(&result)
	if err != nil {
		return err
	}

	return nil
}
