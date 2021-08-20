package translation

import (
	"context"
	"net/url"
	"strings"

	"github.com/marmalad3/pokemon/internal/app/sources/baseapi"
)

// TranslationAPIClient is a client to integrate with
// api.funtranslations.com
type TranslationAPIClient struct {
	*baseapi.ApiClient
}

const (
	baseUrlEnvVar  = "TRANSLATION_API_URL"
	defaultBaseUrl = "https://api.funtranslations.com/translate"
)

// NewTranslationAPIClient accepts baseapi.ApiClientOpt parameters to
// configure the base URL and http.Client on a new TranslationAPIClient
// instance.
// An error will be returned if the base URL fails to be parsed
func NewTranslationAPIClient(opts ...baseapi.ApiClientOpt) (*TranslationAPIClient, error) {
	baseClient, err := baseapi.GenerateNewAPIClient(defaultBaseUrl, baseUrlEnvVar).NewApiClient(opts...)
	if err != nil {
		return nil, err
	}
	return &TranslationAPIClient{
		baseClient,
	}, nil
}

// ApiResponse defines how data is returned from the translation REST API
type ApiResponse struct {
	Contents ApiResponseContents `json:"contents"`
}

// ApiResponseContents contains translated data
type ApiResponseContents struct {
	Translated string `json:"translated"`
}

func (ar *ApiResponse) getTranslation() string {
	return ar.Contents.Translated
}

func (p *TranslationAPIClient) requestTranslation(ctx context.Context, path, text string) (string, error) {
	var responseObj *ApiResponse

	form := url.Values{}
	form.Add("text", text)
	requestBody := strings.NewReader(form.Encode())

	headers := url.Values{}
	headers.Add("Content-Type", "application/x-www-form-urlencoded")

	err := p.DoRequest(ctx, "POST", path, requestBody, headers, &responseObj)
	if err != nil {
		return "", err
	}
	return responseObj.getTranslation(), nil
}

// TranslateShakespeare calls the translation API and parses out requested data
func (p *TranslationAPIClient) TranslateShakespeare(ctx context.Context, text string) (string, error) {
	return p.requestTranslation(ctx, "shakespeare.json", text)
}

// TranslateYoda calls the translation API and parses out requested data
func (p *TranslationAPIClient) TranslateYoda(ctx context.Context, text string) (string, error) {
	return p.requestTranslation(ctx, "yoda.json", text)
}
