package baseapi

import (
	"net/http"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestApiClientGeneratorWithDefaultBaseUrl(t *testing.T) {
	baseUrl := "http://example.com"

	baseapi := GenerateNewAPIClient(baseUrl, "SOME_ENV_THAT_DOESNT_EXIT")
	apiClient, err := baseapi.NewApiClient()
	assert.Nil(t, err)

	assert.Equal(t, baseUrl, apiClient.baseUrl.String())
}

func TestApiClientGeneratorEnvBaseUrlTakesPresedenceOverDefault(t *testing.T) {
	baseUrl := "http://example.com"
	envUrlKey := "TEST_ENV"
	envUrlVal := "http://google.com"
	assert.Nil(t, os.Setenv(envUrlKey, envUrlVal))

	baseapi := GenerateNewAPIClient(baseUrl, envUrlKey)
	apiClient, err := baseapi.NewApiClient()
	assert.Nil(t, err)

	assert.Equal(t, envUrlVal, apiClient.baseUrl.String())
}

func TestApiClientDefaultHttpClient(t *testing.T) {
	baseapi := GenerateNewAPIClient("default-url", "SOME_ENV")
	apiClient, err := baseapi.NewApiClient()
	assert.Nil(t, err)

	assert.Equal(t, http.DefaultClient, apiClient.client)
}

func TestApiClientCustomHttpClient(t *testing.T) {
	baseapi := GenerateNewAPIClient("default-url", "SOME_ENV")

	customClient := &http.Client{}

	apiClient, err := baseapi.NewApiClient(WithHTTPClient(customClient))
	assert.Nil(t, err)

	assert.Equal(t, customClient, apiClient.client)
}
