package zendesk_api

import (
	"encoding/base64"
	"log"

	"github.com/oapi-codegen/oapi-codegen/v2/pkg/securityprovider"
)

type SupportApi struct {
	supportApiClient *ClientWithResponses
}

// serverUrl := "https://...."
// email := "jdoe@example.com"

func NewSupportApi(serverUrl string, email string, apiKey string) *SupportApi {
	// https://developer.zendesk.com/api-reference/introduction/security-and-auth/#api-token

	headerToken := email + "/token:" + apiKey
	headerTokenBase64 := base64.StdEncoding.EncodeToString([]byte(headerToken))
	headerValue := "Basic " + headerTokenBase64

	securityProviderApiKey, err := securityprovider.NewSecurityProviderApiKey(
		"header",
		"Authorization",
		headerValue)

	if err != nil {
		log.Fatal(err)
	}

	client, err := NewClientWithResponses(serverUrl, WithRequestEditorFn(securityProviderApiKey.Intercept))
	if err != nil {
		log.Fatal(err)
	}

	return &SupportApi{
		supportApiClient: client,
	}
}

func (s *SupportApi) GetClient() *ClientWithResponses {
	return s.supportApiClient
}
