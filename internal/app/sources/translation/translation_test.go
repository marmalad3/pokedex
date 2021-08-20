package translation

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/marmalad3/pokemon/internal/app/sources/baseapi"
	"github.com/marmalad3/pokemon/internal/app/support"
	"github.com/stretchr/testify/assert"
)

func TestApiClientYodaTranslate(t *testing.T) {

	mockResponse := ApiResponse{
		Contents: ApiResponseContents{
			Translated: "A very nice pokemon slash houseplant, oddish is",
		},
	}

	respPayload, err := json.Marshal(mockResponse)
	assert.Nil(t, err)

	httpClient := support.GetMockHTTPClient(t, &support.MockInteraction{
		ResponseData:   respPayload,
		ResponseStatus: http.StatusOK,
		ExpectedMethod: http.MethodPost,
		ExpectedPath:   "/translate/yoda.json",
	})

	translationClient, err := NewTranslationAPIClient(baseapi.WithHTTPClient(httpClient))
	assert.Nil(t, err)

	translation, err := translationClient.TranslateYoda(context.Background(), "Oddish is a very nice pokemon slash houseplant")
	assert.Nil(t, err)

	assert.Equal(t, mockResponse.Contents.Translated, translation)
}

func TestApiClientShakespeareTranslate(t *testing.T) {

	mockResponse := ApiResponse{
		Contents: ApiResponseContents{
			Translated: "Don't water oddish too much or he might kicketh the bucket",
		},
	}

	respPayload, err := json.Marshal(mockResponse)
	assert.Nil(t, err)

	httpClient := support.GetMockHTTPClient(t, &support.MockInteraction{
		ResponseData:   respPayload,
		ResponseStatus: http.StatusOK,
		ExpectedMethod: http.MethodPost,
		ExpectedPath:   "/translate/shakespeare.json",
	})

	translationClient, err := NewTranslationAPIClient(baseapi.WithHTTPClient(httpClient))
	assert.Nil(t, err)

	translation, err := translationClient.TranslateShakespeare(context.Background(), "Don't water Oddish too much or he might die")
	assert.Nil(t, err)

	assert.Equal(t, mockResponse.Contents.Translated, translation)
}
