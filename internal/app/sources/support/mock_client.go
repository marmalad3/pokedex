package support

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"
)

type MockRoundTripper struct {
	t            *testing.T
	interactions []*MockInteraction
}

type MockInteraction struct {
	ResponseData   []byte
	ResponseStatus int
	ExpectedMethod string
	ExpectedPath   string
}

func (mpt *MockRoundTripper) RoundTrip(r *http.Request) (*http.Response, error) {

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

func GetMockHTTPClient(t *testing.T, interactions ...*MockInteraction) *http.Client {
	t.Helper()

	return &http.Client{
		Transport: &MockRoundTripper{
			t:            t,
			interactions: interactions,
		},
	}
}
