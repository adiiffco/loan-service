package clients

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"loanapp/adapters/logger"

	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/sirupsen/logrus"
)

type ApiClient struct {
	BaseURL    string
	HttpClient *http.Client
	Log        *logger.Log
	Ctx        context.Context
}

func NewClient() *ApiClient {
	t := http.DefaultTransport.(*http.Transport).Clone()
	return &ApiClient{
		HttpClient: &http.Client{
			Transport: newrelic.NewRoundTripper(t),
		},
		Log: &logger.Log{
			Tag: "ApiClient",
		},
	}
}

func buildURLWithParams(requestURL string, data map[string]interface{}) (string, error) {
	URL, err := url.Parse(requestURL)
	if err != nil {
		return "", err
	}

	parameters := url.Values{}
	for k, v := range data {
		switch val := v.(type) {
		case []int64:
			for _, value := range val {
				parameters.Add(k, fmt.Sprintf("%v", value))
			}
		case []string:
			for _, value := range val {
				parameters.Add(k, fmt.Sprintf("%v", value))
			}
		default:
			parameters.Add(k, fmt.Sprintf("%v", v))
		}
	}

	URL.RawQuery = parameters.Encode()

	return URL.String(), nil
}

func (c *ApiClient) addRequestHeaders(req *http.Request, headers map[string]string) {
	req.Header.Set("Content-Type", "application/json")

	for key, value := range headers {
		req.Header.Set(key, value)
	}
}

func (c *ApiClient) doRequest(req *http.Request, out interface{}) error {
	response, err := c.HttpClient.Do(req)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode >= http.StatusOK &&
		response.StatusCode < http.StatusMultipleChoices {
		body, err := io.ReadAll(response.Body)
		if err != nil {
			return err
		}
		if err := json.Unmarshal(body, out); err != nil {
			return err
		}
		return nil
	} else {
		responseData, err1 := io.ReadAll(response.Body)
		if err1 != nil {
			c.Log.LogData(c.Ctx, logrus.ErrorLevel, "", logrus.Fields{"request": req, "error": err1.Error()}, "HTTP_RESPONSE_ERROR")
			return errors.New("HTTP request failed with status code " + response.Status)
		}
		c.Log.LogData(c.Ctx, logrus.ErrorLevel, "", logrus.Fields{"request": req, "response": string(responseData)}, "HTTP_RESPONSE_ERROR")
		return errors.New("HTTP request failed with status code " + response.Status)
	}
}

func (c *ApiClient) Get(
	ctx context.Context,
	out interface{},
	path string,
	queryParams map[string]interface{},
	headers map[string]string,
) error {
	var err error
	URL := fmt.Sprintf("%s%s", c.BaseURL, path)
	if URL, err = buildURLWithParams(URL, queryParams); err != nil {
		return err
	}
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, URL, nil)

	c.addRequestHeaders(req, headers)
	return c.doRequest(req, out)
}

func (c *ApiClient) Post(
	ctx context.Context,
	out interface{},
	path string,
	data interface{},
	headers map[string]string,
) error {

	URL := fmt.Sprintf("%s%s", c.BaseURL, path)
	jsonValue, _ := json.Marshal(data)
	req, _ := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		URL,
		bytes.NewBuffer(jsonValue),
	)

	c.addRequestHeaders(req, headers)
	return c.doRequest(req, out)
}
