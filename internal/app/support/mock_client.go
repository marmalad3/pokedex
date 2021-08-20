package support

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"
)

// MockRoundTripper implements http.RoundTripper to return
// mocked responses based on matching interactions
type MockRoundTripper struct {
	t            *testing.T
	interactions []*MockInteraction
}

// MockIntegration defines a HTTP response to be returned
// based on an expected reqeust
type MockInteraction struct {
	ResponseData   []byte
	ResponseStatus int
	ExpectedMethod string
	ExpectedPath   string
}

// RoundTrip returns a the http.Response object based on mocked interactions
func (mpt *MockRoundTripper) RoundTrip(r *http.Request) (*http.Response, error) {
	mpt.t.Helper()
	for _, interaction := range mpt.interactions {
		if interaction.ExpectedMethod == r.Method && interaction.ExpectedPath == r.URL.Path {
			return &http.Response{
				StatusCode: interaction.ResponseStatus,
				Body:       ioutil.NopCloser(bytes.NewBuffer(interaction.ResponseData)),
			}, nil
		}
	}
	return nil, fmt.Errorf("mock not found for method '%s' and path '%s'", r.Method, r.URL.Path)
}

// GetMockHTTPClient returns a http.Client which will return mocked interactions
// based on the request data
func GetMockHTTPClient(t *testing.T, interactions ...*MockInteraction) *http.Client {
	t.Helper()

	return &http.Client{
		Transport: &MockRoundTripper{
			t:            t,
			interactions: interactions,
		},
	}
}
