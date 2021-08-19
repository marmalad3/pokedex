package baseapi

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
)

// ApiClientGenerator is a baseapi for APIs which allows default information
// to be set. If baseUrlEnv is set, the base URL is created from the env variable
type ApiClientGenerator struct {
	defaultBaseURL string
	baseURLEnv     string
}

// GenerateNewAPIClient
func GenerateNewAPIClient(defaultBaseUrl string, baseUrlEnv string) *ApiClientGenerator {
	return &ApiClientGenerator{
		defaultBaseURL: defaultBaseUrl,
		baseURLEnv:     baseUrlEnv,
	}
}

// ApiClient interfaces with the PokeAPI.co REST API
type ApiClient struct {
	client  *http.Client
	baseUrl *url.URL
}

// ApiClientOpt defines optional ApiClient functional params
type ApiClientOpt func(input *ApiClient) error

// NewApiClient creates a new ApiClient object using params
func (acg *ApiClientGenerator) NewApiClient(opts ...ApiClientOpt) (*ApiClient, error) {
	defaultBaseUrl := acg.defaultBaseURL
	if envBaseUrl := os.Getenv(acg.baseURLEnv); envBaseUrl != "" {
		defaultBaseUrl = envBaseUrl
	}
	parsedBaseUrl, err := url.Parse(defaultBaseUrl)
	if err != nil {
		return nil, err
	}

	client := &ApiClient{
		client:  http.DefaultClient,
		baseUrl: parsedBaseUrl,
	}

	for _, optFn := range opts {
		optFn(client)
	}
	return client, nil
}

// WithHTTPClient is an functional param for setting a http.Client
func WithHTTPClient(c *http.Client) ApiClientOpt {
	return func(as *ApiClient) error {
		as.client = c
		return nil
	}
}

// WithBaseURL is a functional param for setting a base URL for the poke API
func WithBaseURL(u string) ApiClientOpt {
	return func(as *ApiClient) error {

		parsedBaseUrl, err := url.Parse(u)
		if err != nil {
			return err
		}

		as.baseUrl = parsedBaseUrl
		return nil
	}
}

// ApiError is returned when a 4xx+ status code is returned
type ApiError struct {
	StatusCode int
}

// Error returns an error string
func (ae *ApiError) Error() string {
	return fmt.Sprintf("Status '%d' returned from server", ae.StatusCode)
}

// IsRetryable returns true if a request can be retried
func (ae *ApiError) IsRetryable() bool {
	return ae.StatusCode >= 500
}

// doRequest sends a request to the source API and parses the response
func (ac *ApiClient) DoRequest(ctx context.Context, method string, path string, requestBody io.Reader, additionalHeaders url.Values, responseObj interface{}) error {

	req, err := http.NewRequestWithContext(ctx, method, fmt.Sprintf("%s/%s", ac.baseUrl, path), requestBody)
	if err != nil {
		return err
	}

	for k, values := range additionalHeaders {
		for _, val := range values {
			req.Header.Add(k, val)
		}
	}

	resp, err := ac.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return &ApiError{
			StatusCode: resp.StatusCode,
		}
	}

	payload, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if err = json.Unmarshal(payload, &responseObj); err != nil {
		return err
	}

	return nil
}
