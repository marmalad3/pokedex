package source

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"testing"
)

type mockRoundTripper struct {
	payload    []byte
	statusCode int
}

func (mpt *mockRoundTripper) RoundTrip(r *http.Request) (*http.Response, error) {
	return &http.Response{
		StatusCode: mpt.statusCode,
		Body:       ioutil.NopCloser(bytes.NewBuffer(mpt.payload)),
	}, nil

}

func getMockRoundTripper(t *testing.T, inputData []byte, statusCode int) http.RoundTripper {
	t.Helper()
	return &mockRoundTripper{
		payload:    inputData,
		statusCode: statusCode,
	}
}
