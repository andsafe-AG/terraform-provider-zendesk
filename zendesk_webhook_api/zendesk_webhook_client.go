package zendesk_webhook_api

/*
Manually edited client, due to anonymous structs issues.
*/
import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/oapi-codegen/runtime"
)

// ListWebhooksParams defines parameters for ListWebhooks.
type ListWebhooksParams struct {
	// Sort Defines the sorting criteria. Only supports name and status.
	Sort *string `form:"sort,omitempty" json:"sort,omitempty"`

	// FilterNameContains Filters the webhooks by a string in the name
	FilterNameContains *string `form:"filter[name_contains],omitempty" json:"filter[name_contains],omitempty"`

	// FilterStatus Filters the webhooks by webhook status
	FilterStatus *string `form:"filter[status],omitempty" json:"filter[status],omitempty"`

	// PageBefore Includes the previous page of results with defined size
	PageBefore *string `form:"page[before],omitempty" json:"page[before],omitempty"`

	// PageAfter Includes the next page of results with defined size
	PageAfter *string `form:"page[after],omitempty" json:"page[after],omitempty"`

	// PageSize Defines a specified number of results per page
	PageSize *string `form:"page[size],omitempty" json:"page[size],omitempty"`
}
type AuthenticationData struct {
	Password *string `json:"password,omitempty"`
	Token    *string `json:"token,omitempty"`
	Username *string `json:"username,omitempty"`
}
type Authentication struct {
	AddPosition string              `json:"add_position,omitempty"`
	Data        *AuthenticationData `json:"data"`
	Type        string              `json:"type"`
}

type ExternalSourceData struct {
	AppId          *string `json:"app_id,omitempty"`
	InstallationId *string `json:"installation_id,omitempty"`
}

type ExternalSource struct {
	Data *ExternalSourceData `json:"data,omitempty"`
	Type *string             `json:"type,omitempty"`
}
type SigningSecret struct {
	Algorithm *string `json:"algorithm,omitempty"`
	Secret    *string `json:"secret,omitempty"`
}

type WebhookWithSensitiveData struct {
	Authentication *Authentication    `json:"authentication,omitempty"`
	CustomHeaders  *map[string]string `json:"custom_headers,omitempty"`
	Description    *string            `json:"description,omitempty"`
	Endpoint       string             `json:"endpoint"`
	ExternalSource *ExternalSource    `json:"external_source,omitempty"`
	HttpMethod     string             `json:"http_method"`
	Name           string             `json:"name"`
	RequestFormat  string             `json:"request_format"`
	SigningSecret  *SigningSecret     `json:"signing_secret,omitempty"`
	Status         string             `json:"status"`
	Subscriptions  *[]string          `json:"subscriptions,omitempty"`
}

// CreateOrCloneWebhookJSONBody defines parameters for CreateOrCloneWebhook.
type CreateOrCloneWebhookJSONBody struct {
	Webhook *WebhookWithSensitiveData `json:"webhook"`
}

// CreateOrCloneWebhookParams defines parameters for CreateOrCloneWebhook.
type CreateOrCloneWebhookParams struct {
	// CloneWebhookId id of the webhook to clone. Only required if cloning a webhook.
	CloneWebhookId *string `form:"clone_webhook_id,omitempty" json:"clone_webhook_id,omitempty"`
}

// TestWebhookJSONBody defines parameters for TestWebhook.
type TestWebhookJSONBody struct {
	Request *struct {
		CustomHeaders   *map[string]string `json:"custom_headers,omitempty"`
		Endpoint        *string            `json:"endpoint,omitempty"`
		HttpMethod      *string            `json:"http_method,omitempty"`
		Payload         *string            `json:"payload,omitempty"`
		QueryParameters *[]struct {
			Key   *string `json:"key,omitempty"`
			Value *string `json:"value,omitempty"`
		} `json:"query_parameters,omitempty"`
		RequestFormat *string `json:"request_format,omitempty"`
	} `json:"request,omitempty"`
}

// TestWebhookParams defines parameters for TestWebhook.
type TestWebhookParams struct {
	// WebhookId The webhook to be tested. Only required for testing an existing webhook.
	WebhookId *string `form:"webhook_id,omitempty" json:"webhook_id,omitempty"`
}

// PatchWebhookJSONBody defines parameters for PatchWebhook.
type PatchWebhookJSONBody struct {
	Webhook *struct {
		Authentication *struct {
			AddPosition *string `json:"add_position,omitempty"`
			Data        *struct {
				Password *string `json:"password,omitempty"`
				Username *string `json:"username,omitempty"`
			} `json:"data,omitempty"`
			Type *string `json:"type,omitempty"`
		} `json:"authentication,omitempty"`
		CreatedAt      *string            `json:"created_at,omitempty"`
		CreatedBy      *string            `json:"created_by,omitempty"`
		CustomHeaders  *map[string]string `json:"custom_headers,omitempty"`
		Description    *string            `json:"description,omitempty"`
		Endpoint       *string            `json:"endpoint,omitempty"`
		ExternalSource *struct {
			Data *struct {
				AppId          *string `json:"app_id,omitempty"`
				InstallationId *string `json:"installation_id,omitempty"`
			} `json:"data,omitempty"`
			Type *string `json:"type,omitempty"`
		} `json:"external_source,omitempty"`
		HttpMethod    *string `json:"http_method,omitempty"`
		Id            *string `json:"id,omitempty"`
		Name          *string `json:"name,omitempty"`
		RequestFormat *string `json:"request_format,omitempty"`
		SigningSecret *struct {
			Algorithm *string `json:"algorithm,omitempty"`
			Secret    *string `json:"secret,omitempty"`
		} `json:"signing_secret,omitempty"`
		Status        *string   `json:"status,omitempty"`
		Subscriptions *[]string `json:"subscriptions,omitempty"`
		UpdatedAt     *string   `json:"updated_at,omitempty"`
		UpdatedBy     *string   `json:"updated_by,omitempty"`
	} `json:"webhook,omitempty"`
}

// UpdateWebhookJSONBody defines parameters for UpdateWebhook.
type UpdateWebhookJSONBody struct {
	Webhook *WebhookWithSensitiveData `json:"webhook,omitempty"`
}

// ListWebhookInvocationsParams defines parameters for ListWebhookInvocations.
type ListWebhookInvocationsParams struct {
	// PageBefore Includes the previous page of invocations with defined size
	PageBefore *string `form:"page[before],omitempty" json:"page[before],omitempty"`

	// PageAfter Includes the next page of invocations with defined size
	PageAfter *string `form:"page[after],omitempty" json:"page[after],omitempty"`

	// PageSize Defines a specific number of invocations per page
	PageSize *string `form:"page[size],omitempty" json:"page[size],omitempty"`

	// Sort Defines a invocation attribute to sort invocations.
	Sort *string `form:"sort,omitempty" json:"sort,omitempty"`

	// FilterStatus Filters invocations by invocation status.
	FilterStatus *string `form:"filter[status],omitempty" json:"filter[status],omitempty"`

	// FilterFromTs Filters invocations by from timestamp. Use ISO 8601 UTC format
	FilterFromTs *string `form:"filter[from_ts],omitempty" json:"filter[from_ts],omitempty"`

	// FilterToTs Filters invocations by timestamp. Use ISO 8601 UTC format
	FilterToTs *string `form:"filter[to_ts],omitempty" json:"filter[to_ts],omitempty"`
}

// CreateOrCloneWebhookJSONRequestBody defines body for CreateOrCloneWebhook for application/json ContentType.
type CreateOrCloneWebhookJSONRequestBody CreateOrCloneWebhookJSONBody

// TestWebhookJSONRequestBody defines body for TestWebhook for application/json ContentType.
type TestWebhookJSONRequestBody TestWebhookJSONBody

// PatchWebhookJSONRequestBody defines body for PatchWebhook for application/json ContentType.
type PatchWebhookJSONRequestBody PatchWebhookJSONBody

// UpdateWebhookJSONRequestBody defines body for UpdateWebhook for application/json ContentType.
type UpdateWebhookJSONRequestBody UpdateWebhookJSONBody

// RequestEditorFn  is the function signature for the RequestEditor callback function
type RequestEditorFn func(ctx context.Context, req *http.Request) error

// Doer performs HTTP requests.
//
// The standard http.Client implements this interface.
type HttpRequestDoer interface {
	Do(req *http.Request) (*http.Response, error)
}

// Client which conforms to the OpenAPI3 specification for this service.
type Client struct {
	// The endpoint of the server conforming to this interface, with scheme,
	// https://api.deepmap.com for example. This can contain a path relative
	// to the server, such as https://api.deepmap.com/dev-test, and all the
	// paths in the swagger spec will be appended to the server.
	Server string

	// Doer for performing requests, typically a *http.Client with any
	// customized settings, such as certificate chains.
	Client HttpRequestDoer

	// A list of callbacks for modifying requests which are generated before sending over
	// the network.
	RequestEditors []RequestEditorFn
}

// ClientOption allows setting custom parameters during construction
type ClientOption func(*Client) error

// Creates a new Client, with reasonable defaults
func NewClient(server string, opts ...ClientOption) (*Client, error) {
	// create a client with sane default values
	client := Client{
		Server: server,
	}
	// mutate client and add all optional params
	for _, o := range opts {
		if err := o(&client); err != nil {
			return nil, err
		}
	}
	// ensure the server URL always has a trailing slash
	if !strings.HasSuffix(client.Server, "/") {
		client.Server += "/"
	}
	// create httpClient, if not already present
	if client.Client == nil {
		client.Client = &http.Client{}
	}
	return &client, nil
}

// WithHTTPClient allows overriding the default Doer, which is
// automatically created using http.Client. This is useful for tests.
func WithHTTPClient(doer HttpRequestDoer) ClientOption {
	return func(c *Client) error {
		c.Client = doer
		return nil
	}
}

// WithRequestEditorFn allows setting up a callback function, which will be
// called right before sending the request. This can be used to mutate the request.
func WithRequestEditorFn(fn RequestEditorFn) ClientOption {
	return func(c *Client) error {
		c.RequestEditors = append(c.RequestEditors, fn)
		return nil
	}
}

// The interface specification for the client above.
type ClientInterface interface {
	// ListWebhooks request
	ListWebhooks(ctx context.Context, params *ListWebhooksParams, reqEditors ...RequestEditorFn) (*http.Response, error)

	// CreateOrCloneWebhookWithBody request with any body
	CreateOrCloneWebhookWithBody(ctx context.Context, params *CreateOrCloneWebhookParams, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error)

	CreateOrCloneWebhook(ctx context.Context, params *CreateOrCloneWebhookParams, body CreateOrCloneWebhookJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error)

	// TestWebhookWithBody request with any body
	TestWebhookWithBody(ctx context.Context, params *TestWebhookParams, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error)

	TestWebhook(ctx context.Context, params *TestWebhookParams, body TestWebhookJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error)

	// DeleteWebhook request
	DeleteWebhook(ctx context.Context, webhookId string, reqEditors ...RequestEditorFn) (*http.Response, error)

	// ShowWebhook request
	ShowWebhook(ctx context.Context, webhookId string, reqEditors ...RequestEditorFn) (*http.Response, error)

	// PatchWebhookWithBody request with any body
	PatchWebhookWithBody(ctx context.Context, webhookId string, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error)

	PatchWebhook(ctx context.Context, webhookId string, body PatchWebhookJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error)

	// UpdateWebhookWithBody request with any body
	UpdateWebhookWithBody(ctx context.Context, webhookId string, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error)

	UpdateWebhook(ctx context.Context, webhookId string, body UpdateWebhookJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error)

	// ListWebhookInvocations request
	ListWebhookInvocations(ctx context.Context, webhookId string, params *ListWebhookInvocationsParams, reqEditors ...RequestEditorFn) (*http.Response, error)

	// ListWebhookInvocationAttempts request
	ListWebhookInvocationAttempts(ctx context.Context, webhookId string, invocationId string, reqEditors ...RequestEditorFn) (*http.Response, error)

	// ShowWebhookSigningSecret request
	ShowWebhookSigningSecret(ctx context.Context, webhookId string, reqEditors ...RequestEditorFn) (*http.Response, error)

	// ResetWebhookSigningSecret request
	ResetWebhookSigningSecret(ctx context.Context, webhookId string, reqEditors ...RequestEditorFn) (*http.Response, error)
}

func (c *Client) ListWebhooks(ctx context.Context, params *ListWebhooksParams, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewListWebhooksRequest(c.Server, params)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) CreateOrCloneWebhookWithBody(ctx context.Context, params *CreateOrCloneWebhookParams, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewCreateOrCloneWebhookRequestWithBody(c.Server, params, contentType, body)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) CreateOrCloneWebhook(ctx context.Context, params *CreateOrCloneWebhookParams, body CreateOrCloneWebhookJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewCreateOrCloneWebhookRequest(c.Server, params, body)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) TestWebhookWithBody(ctx context.Context, params *TestWebhookParams, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewTestWebhookRequestWithBody(c.Server, params, contentType, body)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) TestWebhook(ctx context.Context, params *TestWebhookParams, body TestWebhookJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewTestWebhookRequest(c.Server, params, body)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) DeleteWebhook(ctx context.Context, webhookId string, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewDeleteWebhookRequest(c.Server, webhookId)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) ShowWebhook(ctx context.Context, webhookId string, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewShowWebhookRequest(c.Server, webhookId)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) PatchWebhookWithBody(ctx context.Context, webhookId string, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewPatchWebhookRequestWithBody(c.Server, webhookId, contentType, body)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) PatchWebhook(ctx context.Context, webhookId string, body PatchWebhookJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewPatchWebhookRequest(c.Server, webhookId, body)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) UpdateWebhookWithBody(ctx context.Context, webhookId string, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewUpdateWebhookRequestWithBody(c.Server, webhookId, contentType, body)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) UpdateWebhook(ctx context.Context, webhookId string, body UpdateWebhookJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewUpdateWebhookRequest(c.Server, webhookId, body)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) ListWebhookInvocations(ctx context.Context, webhookId string, params *ListWebhookInvocationsParams, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewListWebhookInvocationsRequest(c.Server, webhookId, params)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) ListWebhookInvocationAttempts(ctx context.Context, webhookId string, invocationId string, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewListWebhookInvocationAttemptsRequest(c.Server, webhookId, invocationId)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) ShowWebhookSigningSecret(ctx context.Context, webhookId string, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewShowWebhookSigningSecretRequest(c.Server, webhookId)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) ResetWebhookSigningSecret(ctx context.Context, webhookId string, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewResetWebhookSigningSecretRequest(c.Server, webhookId)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

// NewListWebhooksRequest generates requests for ListWebhooks
func NewListWebhooksRequest(server string, params *ListWebhooksParams) (*http.Request, error) {
	var err error

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := "/api/v2/webhooks"
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	if params != nil {
		queryValues := queryURL.Query()

		if params.Sort != nil {

			if queryFrag, err := runtime.StyleParamWithLocation("form", true, "sort", runtime.ParamLocationQuery, *params.Sort); err != nil {
				return nil, err
			} else if parsed, err := url.ParseQuery(queryFrag); err != nil {
				return nil, err
			} else {
				for k, v := range parsed {
					for _, v2 := range v {
						queryValues.Add(k, v2)
					}
				}
			}

		}

		if params.FilterNameContains != nil {

			if queryFrag, err := runtime.StyleParamWithLocation("form", true, "filter[name_contains]", runtime.ParamLocationQuery, *params.FilterNameContains); err != nil {
				return nil, err
			} else if parsed, err := url.ParseQuery(queryFrag); err != nil {
				return nil, err
			} else {
				for k, v := range parsed {
					for _, v2 := range v {
						queryValues.Add(k, v2)
					}
				}
			}

		}

		if params.FilterStatus != nil {

			if queryFrag, err := runtime.StyleParamWithLocation("form", true, "filter[status]", runtime.ParamLocationQuery, *params.FilterStatus); err != nil {
				return nil, err
			} else if parsed, err := url.ParseQuery(queryFrag); err != nil {
				return nil, err
			} else {
				for k, v := range parsed {
					for _, v2 := range v {
						queryValues.Add(k, v2)
					}
				}
			}

		}

		if params.PageBefore != nil {

			if queryFrag, err := runtime.StyleParamWithLocation("form", true, "page[before]", runtime.ParamLocationQuery, *params.PageBefore); err != nil {
				return nil, err
			} else if parsed, err := url.ParseQuery(queryFrag); err != nil {
				return nil, err
			} else {
				for k, v := range parsed {
					for _, v2 := range v {
						queryValues.Add(k, v2)
					}
				}
			}

		}

		if params.PageAfter != nil {

			if queryFrag, err := runtime.StyleParamWithLocation("form", true, "page[after]", runtime.ParamLocationQuery, *params.PageAfter); err != nil {
				return nil, err
			} else if parsed, err := url.ParseQuery(queryFrag); err != nil {
				return nil, err
			} else {
				for k, v := range parsed {
					for _, v2 := range v {
						queryValues.Add(k, v2)
					}
				}
			}

		}

		if params.PageSize != nil {

			if queryFrag, err := runtime.StyleParamWithLocation("form", true, "page[size]", runtime.ParamLocationQuery, *params.PageSize); err != nil {
				return nil, err
			} else if parsed, err := url.ParseQuery(queryFrag); err != nil {
				return nil, err
			} else {
				for k, v := range parsed {
					for _, v2 := range v {
						queryValues.Add(k, v2)
					}
				}
			}

		}

		queryURL.RawQuery = queryValues.Encode()
	}

	req, err := http.NewRequest("GET", queryURL.String(), nil)
	if err != nil {
		return nil, err
	}

	return req, nil
}

// NewCreateOrCloneWebhookRequest calls the generic CreateOrCloneWebhook builder with application/json body
func NewCreateOrCloneWebhookRequest(server string, params *CreateOrCloneWebhookParams, body CreateOrCloneWebhookJSONRequestBody) (*http.Request, error) {
	var bodyReader io.Reader
	buf, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	bodyReader = bytes.NewReader(buf)
	return NewCreateOrCloneWebhookRequestWithBody(server, params, "application/json", bodyReader)
}

// NewCreateOrCloneWebhookRequestWithBody generates requests for CreateOrCloneWebhook with any type of body
func NewCreateOrCloneWebhookRequestWithBody(server string, params *CreateOrCloneWebhookParams, contentType string, body io.Reader) (*http.Request, error) {
	var err error

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := "/api/v2/webhooks"
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	if params != nil {
		queryValues := queryURL.Query()

		if params.CloneWebhookId != nil {

			if queryFrag, err := runtime.StyleParamWithLocation("form", true, "clone_webhook_id", runtime.ParamLocationQuery, *params.CloneWebhookId); err != nil {
				return nil, err
			} else if parsed, err := url.ParseQuery(queryFrag); err != nil {
				return nil, err
			} else {
				for k, v := range parsed {
					for _, v2 := range v {
						queryValues.Add(k, v2)
					}
				}
			}

		}

		queryURL.RawQuery = queryValues.Encode()
	}

	req, err := http.NewRequest("POST", queryURL.String(), body)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", contentType)

	return req, nil
}

// NewTestWebhookRequest calls the generic TestWebhook builder with application/json body
func NewTestWebhookRequest(server string, params *TestWebhookParams, body TestWebhookJSONRequestBody) (*http.Request, error) {
	var bodyReader io.Reader
	buf, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	bodyReader = bytes.NewReader(buf)
	return NewTestWebhookRequestWithBody(server, params, "application/json", bodyReader)
}

// NewTestWebhookRequestWithBody generates requests for TestWebhook with any type of body
func NewTestWebhookRequestWithBody(server string, params *TestWebhookParams, contentType string, body io.Reader) (*http.Request, error) {
	var err error

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := "/api/v2/webhooks/test"
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	if params != nil {
		queryValues := queryURL.Query()

		if params.WebhookId != nil {

			if queryFrag, err := runtime.StyleParamWithLocation("form", true, "webhook_id", runtime.ParamLocationQuery, *params.WebhookId); err != nil {
				return nil, err
			} else if parsed, err := url.ParseQuery(queryFrag); err != nil {
				return nil, err
			} else {
				for k, v := range parsed {
					for _, v2 := range v {
						queryValues.Add(k, v2)
					}
				}
			}

		}

		queryURL.RawQuery = queryValues.Encode()
	}

	req, err := http.NewRequest("POST", queryURL.String(), body)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", contentType)

	return req, nil
}

// NewDeleteWebhookRequest generates requests for DeleteWebhook
func NewDeleteWebhookRequest(server string, webhookId string) (*http.Request, error) {
	var err error

	var pathParam0 string

	pathParam0, err = runtime.StyleParamWithLocation("simple", false, "webhook_id", runtime.ParamLocationPath, webhookId)
	if err != nil {
		return nil, err
	}

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/api/v2/webhooks/%s", pathParam0)
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("DELETE", queryURL.String(), nil)
	if err != nil {
		return nil, err
	}

	return req, nil
}

// NewShowWebhookRequest generates requests for ShowWebhook
func NewShowWebhookRequest(server string, webhookId string) (*http.Request, error) {
	var err error

	var pathParam0 string

	pathParam0, err = runtime.StyleParamWithLocation("simple", false, "webhook_id", runtime.ParamLocationPath, webhookId)
	if err != nil {
		return nil, err
	}

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/api/v2/webhooks/%s", pathParam0)
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("GET", queryURL.String(), nil)
	if err != nil {
		return nil, err
	}

	return req, nil
}

// NewPatchWebhookRequest calls the generic PatchWebhook builder with application/json body
func NewPatchWebhookRequest(server string, webhookId string, body PatchWebhookJSONRequestBody) (*http.Request, error) {
	var bodyReader io.Reader
	buf, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	bodyReader = bytes.NewReader(buf)
	return NewPatchWebhookRequestWithBody(server, webhookId, "application/json", bodyReader)
}

// NewPatchWebhookRequestWithBody generates requests for PatchWebhook with any type of body
func NewPatchWebhookRequestWithBody(server string, webhookId string, contentType string, body io.Reader) (*http.Request, error) {
	var err error

	var pathParam0 string

	pathParam0, err = runtime.StyleParamWithLocation("simple", false, "webhook_id", runtime.ParamLocationPath, webhookId)
	if err != nil {
		return nil, err
	}

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/api/v2/webhooks/%s", pathParam0)
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("PATCH", queryURL.String(), body)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", contentType)

	return req, nil
}

// NewUpdateWebhookRequest calls the generic UpdateWebhook builder with application/json body
func NewUpdateWebhookRequest(server string, webhookId string, body UpdateWebhookJSONRequestBody) (*http.Request, error) {
	var bodyReader io.Reader
	buf, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	bodyReader = bytes.NewReader(buf)
	return NewUpdateWebhookRequestWithBody(server, webhookId, "application/json", bodyReader)
}

// NewUpdateWebhookRequestWithBody generates requests for UpdateWebhook with any type of body
func NewUpdateWebhookRequestWithBody(server string, webhookId string, contentType string, body io.Reader) (*http.Request, error) {
	var err error

	var pathParam0 string

	pathParam0, err = runtime.StyleParamWithLocation("simple", false, "webhook_id", runtime.ParamLocationPath, webhookId)
	if err != nil {
		return nil, err
	}

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/api/v2/webhooks/%s", pathParam0)
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("PUT", queryURL.String(), body)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", contentType)

	return req, nil
}

// NewListWebhookInvocationsRequest generates requests for ListWebhookInvocations
func NewListWebhookInvocationsRequest(server string, webhookId string, params *ListWebhookInvocationsParams) (*http.Request, error) {
	var err error

	var pathParam0 string

	pathParam0, err = runtime.StyleParamWithLocation("simple", false, "webhook_id", runtime.ParamLocationPath, webhookId)
	if err != nil {
		return nil, err
	}

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/api/v2/webhooks/%s/invocations", pathParam0)
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	if params != nil {
		queryValues := queryURL.Query()

		if params.PageBefore != nil {

			if queryFrag, err := runtime.StyleParamWithLocation("form", true, "page[before]", runtime.ParamLocationQuery, *params.PageBefore); err != nil {
				return nil, err
			} else if parsed, err := url.ParseQuery(queryFrag); err != nil {
				return nil, err
			} else {
				for k, v := range parsed {
					for _, v2 := range v {
						queryValues.Add(k, v2)
					}
				}
			}

		}

		if params.PageAfter != nil {

			if queryFrag, err := runtime.StyleParamWithLocation("form", true, "page[after]", runtime.ParamLocationQuery, *params.PageAfter); err != nil {
				return nil, err
			} else if parsed, err := url.ParseQuery(queryFrag); err != nil {
				return nil, err
			} else {
				for k, v := range parsed {
					for _, v2 := range v {
						queryValues.Add(k, v2)
					}
				}
			}

		}

		if params.PageSize != nil {

			if queryFrag, err := runtime.StyleParamWithLocation("form", true, "page[size]", runtime.ParamLocationQuery, *params.PageSize); err != nil {
				return nil, err
			} else if parsed, err := url.ParseQuery(queryFrag); err != nil {
				return nil, err
			} else {
				for k, v := range parsed {
					for _, v2 := range v {
						queryValues.Add(k, v2)
					}
				}
			}

		}

		if params.Sort != nil {

			if queryFrag, err := runtime.StyleParamWithLocation("form", true, "sort", runtime.ParamLocationQuery, *params.Sort); err != nil {
				return nil, err
			} else if parsed, err := url.ParseQuery(queryFrag); err != nil {
				return nil, err
			} else {
				for k, v := range parsed {
					for _, v2 := range v {
						queryValues.Add(k, v2)
					}
				}
			}

		}

		if params.FilterStatus != nil {

			if queryFrag, err := runtime.StyleParamWithLocation("form", true, "filter[status]", runtime.ParamLocationQuery, *params.FilterStatus); err != nil {
				return nil, err
			} else if parsed, err := url.ParseQuery(queryFrag); err != nil {
				return nil, err
			} else {
				for k, v := range parsed {
					for _, v2 := range v {
						queryValues.Add(k, v2)
					}
				}
			}

		}

		if params.FilterFromTs != nil {

			if queryFrag, err := runtime.StyleParamWithLocation("form", true, "filter[from_ts]", runtime.ParamLocationQuery, *params.FilterFromTs); err != nil {
				return nil, err
			} else if parsed, err := url.ParseQuery(queryFrag); err != nil {
				return nil, err
			} else {
				for k, v := range parsed {
					for _, v2 := range v {
						queryValues.Add(k, v2)
					}
				}
			}

		}

		if params.FilterToTs != nil {

			if queryFrag, err := runtime.StyleParamWithLocation("form", true, "filter[to_ts]", runtime.ParamLocationQuery, *params.FilterToTs); err != nil {
				return nil, err
			} else if parsed, err := url.ParseQuery(queryFrag); err != nil {
				return nil, err
			} else {
				for k, v := range parsed {
					for _, v2 := range v {
						queryValues.Add(k, v2)
					}
				}
			}

		}

		queryURL.RawQuery = queryValues.Encode()
	}

	req, err := http.NewRequest("GET", queryURL.String(), nil)
	if err != nil {
		return nil, err
	}

	return req, nil
}

// NewListWebhookInvocationAttemptsRequest generates requests for ListWebhookInvocationAttempts
func NewListWebhookInvocationAttemptsRequest(server string, webhookId string, invocationId string) (*http.Request, error) {
	var err error

	var pathParam0 string

	pathParam0, err = runtime.StyleParamWithLocation("simple", false, "webhook_id", runtime.ParamLocationPath, webhookId)
	if err != nil {
		return nil, err
	}

	var pathParam1 string

	pathParam1, err = runtime.StyleParamWithLocation("simple", false, "invocation_id", runtime.ParamLocationPath, invocationId)
	if err != nil {
		return nil, err
	}

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/api/v2/webhooks/%s/invocations/%s/attempts", pathParam0, pathParam1)
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("GET", queryURL.String(), nil)
	if err != nil {
		return nil, err
	}

	return req, nil
}

// NewShowWebhookSigningSecretRequest generates requests for ShowWebhookSigningSecret
func NewShowWebhookSigningSecretRequest(server string, webhookId string) (*http.Request, error) {
	var err error

	var pathParam0 string

	pathParam0, err = runtime.StyleParamWithLocation("simple", false, "webhook_id", runtime.ParamLocationPath, webhookId)
	if err != nil {
		return nil, err
	}

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/api/v2/webhooks/%s/signing_secret", pathParam0)
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("GET", queryURL.String(), nil)
	if err != nil {
		return nil, err
	}

	return req, nil
}

// NewResetWebhookSigningSecretRequest generates requests for ResetWebhookSigningSecret
func NewResetWebhookSigningSecretRequest(server string, webhookId string) (*http.Request, error) {
	var err error

	var pathParam0 string

	pathParam0, err = runtime.StyleParamWithLocation("simple", false, "webhook_id", runtime.ParamLocationPath, webhookId)
	if err != nil {
		return nil, err
	}

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/api/v2/webhooks/%s/signing_secret", pathParam0)
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", queryURL.String(), nil)
	if err != nil {
		return nil, err
	}

	return req, nil
}

func (c *Client) applyEditors(ctx context.Context, req *http.Request, additionalEditors []RequestEditorFn) error {
	for _, r := range c.RequestEditors {
		if err := r(ctx, req); err != nil {
			return err
		}
	}
	for _, r := range additionalEditors {
		if err := r(ctx, req); err != nil {
			return err
		}
	}
	return nil
}

// ClientWithResponses builds on ClientInterface to offer response payloads
type ClientWithResponses struct {
	ClientInterface
}

// NewClientWithResponses creates a new ClientWithResponses, which wraps
// Client with return type handling
func NewClientWithResponses(server string, opts ...ClientOption) (*ClientWithResponses, error) {
	client, err := NewClient(server, opts...)
	if err != nil {
		return nil, err
	}
	return &ClientWithResponses{client}, nil
}

// WithBaseURL overrides the baseURL.
func WithBaseURL(baseURL string) ClientOption {
	return func(c *Client) error {
		newBaseURL, err := url.Parse(baseURL)
		if err != nil {
			return err
		}
		c.Server = newBaseURL.String()
		return nil
	}
}

// ClientWithResponsesInterface is the interface specification for the client with responses above.
type ClientWithResponsesInterface interface {
	// ListWebhooksWithResponse request
	ListWebhooksWithResponse(ctx context.Context, params *ListWebhooksParams, reqEditors ...RequestEditorFn) (*ListWebhooksWrap, error)

	// CreateOrCloneWebhookWithBodyWithResponse request with any body
	CreateOrCloneWebhookWithBodyWithResponse(ctx context.Context, params *CreateOrCloneWebhookParams, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*CreateOrCloneWebhookWrap, error)

	CreateOrCloneWebhookWithResponse(ctx context.Context, params *CreateOrCloneWebhookParams, body CreateOrCloneWebhookJSONRequestBody, reqEditors ...RequestEditorFn) (*CreateOrCloneWebhookWrap, error)

	// TestWebhookWithBodyWithResponse request with any body
	TestWebhookWithBodyWithResponse(ctx context.Context, params *TestWebhookParams, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*TestWebhookWrap, error)

	TestWebhookWithResponse(ctx context.Context, params *TestWebhookParams, body TestWebhookJSONRequestBody, reqEditors ...RequestEditorFn) (*TestWebhookWrap, error)

	// DeleteWebhookWithResponse request
	DeleteWebhookWithResponse(ctx context.Context, webhookId string, reqEditors ...RequestEditorFn) (*DeleteWebhookWrap, error)

	// ShowWebhookWithResponse request
	ShowWebhookWithResponse(ctx context.Context, webhookId string, reqEditors ...RequestEditorFn) (*ShowWebhookWrap, error)

	// PatchWebhookWithBodyWithResponse request with any body
	PatchWebhookWithBodyWithResponse(ctx context.Context, webhookId string, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*PatchWebhookWrap, error)

	PatchWebhookWithResponse(ctx context.Context, webhookId string, body PatchWebhookJSONRequestBody, reqEditors ...RequestEditorFn) (*PatchWebhookWrap, error)

	// UpdateWebhookWithBodyWithResponse request with any body
	UpdateWebhookWithBodyWithResponse(ctx context.Context, webhookId string, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*UpdateWebhookWrap, error)

	UpdateWebhookWithResponse(ctx context.Context, webhookId string, body UpdateWebhookJSONRequestBody, reqEditors ...RequestEditorFn) (*UpdateWebhookWrap, error)

	// ListWebhookInvocationsWithResponse request
	ListWebhookInvocationsWithResponse(ctx context.Context, webhookId string, params *ListWebhookInvocationsParams, reqEditors ...RequestEditorFn) (*ListWebhookInvocationsWrap, error)

	// ListWebhookInvocationAttemptsWithResponse request
	ListWebhookInvocationAttemptsWithResponse(ctx context.Context, webhookId string, invocationId string, reqEditors ...RequestEditorFn) (*ListWebhookInvocationAttemptsWrap, error)

	// ShowWebhookSigningSecretWithResponse request
	ShowWebhookSigningSecretWithResponse(ctx context.Context, webhookId string, reqEditors ...RequestEditorFn) (*ShowWebhookSigningSecretWrap, error)

	// ResetWebhookSigningSecretWithResponse request
	ResetWebhookSigningSecretWithResponse(ctx context.Context, webhookId string, reqEditors ...RequestEditorFn) (*ResetWebhookSigningSecretWrap, error)
}

type ListWebhooksWrap struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON200      *struct {
		Links *struct {
			Next *string `json:"next,omitempty"`
			Prev *string `json:"prev,omitempty"`
		} `json:"links,omitempty"`
		Meta *struct {
			AfterCursor  *string `json:"after_cursor,omitempty"`
			BeforeCursor *string `json:"before_cursor,omitempty"`
			HasMore      *string `json:"has_more,omitempty"`
		} `json:"meta,omitempty"`
		Webhooks *[]WebhookWithoutSensitive `json:"webhooks,omitempty"`
	}
}

// Status returns HTTPResponse.Status
func (r ListWebhooksWrap) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r ListWebhooksWrap) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type WebhookWithoutSensitive struct {
	Authentication *struct {
		AddPosition *string `json:"add_position,omitempty"`
		Data        *struct {
			Username *string `json:"username,omitempty"`
		} `json:"data,omitempty"`
		Type *string `json:"type,omitempty"`
	} `json:"authentication,omitempty"`
	CreatedAt      *string            `json:"created_at,omitempty"`
	CreatedBy      *string            `json:"created_by,omitempty"`
	CustomHeaders  *map[string]string `json:"custom_headers,omitempty"`
	Description    *string            `json:"description,omitempty"`
	Endpoint       *string            `json:"endpoint,omitempty"`
	ExternalSource *struct {
		Data *struct {
			AppId          *string `json:"app_id,omitempty"`
			InstallationId *string `json:"installation_id,omitempty"`
		} `json:"data,omitempty"`
		Type *string `json:"type,omitempty"`
	} `json:"external_source,omitempty"`
	HttpMethod    *string `json:"http_method,omitempty"`
	Id            *string `json:"id,omitempty"`
	Name          *string `json:"name,omitempty"`
	RequestFormat *string `json:"request_format,omitempty"`
	SigningSecret *struct {
		Algorithm *string `json:"algorithm,omitempty"`
		Secret    *string `json:"secret,omitempty"`
	} `json:"signing_secret,omitempty"`
	Status        *string   `json:"status,omitempty"`
	Subscriptions *[]string `json:"subscriptions,omitempty"`
	UpdatedAt     *string   `json:"updated_at,omitempty"`
	UpdatedBy     *string   `json:"updated_by,omitempty"`
}

type CreateOrCloneWebhookWrap struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON201      *struct {
		Webhook *WebhookWithoutSensitive `json:"webhook,omitempty"`
	}
	JSON400 *struct {
		Errors *[]struct {
			Code   *string `json:"code,omitempty"`
			Detail *string `json:"detail,omitempty"`
			Id     *string `json:"id,omitempty"`
			Source *struct {
				Parameter *string `json:"parameter,omitempty"`
				Pointer   *string `json:"pointer,omitempty"`
			} `json:"source,omitempty"`
			Title *string `json:"title,omitempty"`
		} `json:"errors,omitempty"`
	}
	JSON403 *struct {
		Errors *[]struct {
			Code   *string `json:"code,omitempty"`
			Detail *string `json:"detail,omitempty"`
			Id     *string `json:"id,omitempty"`
			Source *struct {
				Parameter *string `json:"parameter,omitempty"`
				Pointer   *string `json:"pointer,omitempty"`
			} `json:"source,omitempty"`
			Title *string `json:"title,omitempty"`
		} `json:"errors,omitempty"`
	}
}

// Status returns HTTPResponse.Status
func (r CreateOrCloneWebhookWrap) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r CreateOrCloneWebhookWrap) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type TestWebhookWrap struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON200      *struct {
		Response *struct {
			Headers *[]struct {
				Key   *string `json:"key,omitempty"`
				Value *string `json:"value,omitempty"`
			} `json:"headers,omitempty"`
			Payload *string `json:"payload,omitempty"`
			Status  *string `json:"status,omitempty"`
		} `json:"response,omitempty"`
	}
	JSON400 *struct {
		Errors *[]struct {
			Code   *string `json:"code,omitempty"`
			Detail *string `json:"detail,omitempty"`
			Id     *string `json:"id,omitempty"`
			Source *struct {
				Parameter *string `json:"parameter,omitempty"`
				Pointer   *string `json:"pointer,omitempty"`
			} `json:"source,omitempty"`
			Title *string `json:"title,omitempty"`
		} `json:"errors,omitempty"`
	}
}

// Status returns HTTPResponse.Status
func (r TestWebhookWrap) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r TestWebhookWrap) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type DeleteWebhookWrap struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON400      *struct {
		Errors *[]struct {
			Code   *string `json:"code,omitempty"`
			Detail *string `json:"detail,omitempty"`
			Id     *string `json:"id,omitempty"`
			Source *struct {
				Parameter *string `json:"parameter,omitempty"`
				Pointer   *string `json:"pointer,omitempty"`
			} `json:"source,omitempty"`
			Title *string `json:"title,omitempty"`
		} `json:"errors,omitempty"`
	}
	JSON404 *struct {
		Errors *[]struct {
			Code   *string `json:"code,omitempty"`
			Detail *string `json:"detail,omitempty"`
			Id     *string `json:"id,omitempty"`
			Source *struct {
				Parameter *string `json:"parameter,omitempty"`
				Pointer   *string `json:"pointer,omitempty"`
			} `json:"source,omitempty"`
			Title *string `json:"title,omitempty"`
		} `json:"errors,omitempty"`
	}
}

// Status returns HTTPResponse.Status
func (r DeleteWebhookWrap) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r DeleteWebhookWrap) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type ShowWebhookWrap struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON200      *struct {
		Webhook *WebhookWithoutSensitive `json:"webhook,omitempty"`
	}
	JSON400 *struct {
		Errors *[]struct {
			Code   *string `json:"code,omitempty"`
			Detail *string `json:"detail,omitempty"`
			Id     *string `json:"id,omitempty"`
			Source *struct {
				Parameter *string `json:"parameter,omitempty"`
				Pointer   *string `json:"pointer,omitempty"`
			} `json:"source,omitempty"`
			Title *string `json:"title,omitempty"`
		} `json:"errors,omitempty"`
	}
	JSON404 *struct {
		Errors *[]struct {
			Code   *string `json:"code,omitempty"`
			Detail *string `json:"detail,omitempty"`
			Id     *string `json:"id,omitempty"`
			Source *struct {
				Parameter *string `json:"parameter,omitempty"`
				Pointer   *string `json:"pointer,omitempty"`
			} `json:"source,omitempty"`
			Title *string `json:"title,omitempty"`
		} `json:"errors,omitempty"`
	}
}

// Status returns HTTPResponse.Status
func (r ShowWebhookWrap) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r ShowWebhookWrap) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type PatchWebhookWrap struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON400      *struct {
		Errors *[]struct {
			Code   *string `json:"code,omitempty"`
			Detail *string `json:"detail,omitempty"`
			Id     *string `json:"id,omitempty"`
			Source *struct {
				Parameter *string `json:"parameter,omitempty"`
				Pointer   *string `json:"pointer,omitempty"`
			} `json:"source,omitempty"`
			Title *string `json:"title,omitempty"`
		} `json:"errors,omitempty"`
	}
	JSON404 *struct {
		Errors *[]struct {
			Code   *string `json:"code,omitempty"`
			Detail *string `json:"detail,omitempty"`
			Id     *string `json:"id,omitempty"`
			Source *struct {
				Parameter *string `json:"parameter,omitempty"`
				Pointer   *string `json:"pointer,omitempty"`
			} `json:"source,omitempty"`
			Title *string `json:"title,omitempty"`
		} `json:"errors,omitempty"`
	}
}

// Status returns HTTPResponse.Status
func (r PatchWebhookWrap) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r PatchWebhookWrap) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type UpdateWebhookWrap struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON400      *struct {
		Errors *[]struct {
			Code   *string `json:"code,omitempty"`
			Detail *string `json:"detail,omitempty"`
			Id     *string `json:"id,omitempty"`
			Source *struct {
				Parameter *string `json:"parameter,omitempty"`
				Pointer   *string `json:"pointer,omitempty"`
			} `json:"source,omitempty"`
			Title *string `json:"title,omitempty"`
		} `json:"errors,omitempty"`
	}
	JSON404 *struct {
		Errors *[]struct {
			Code   *string `json:"code,omitempty"`
			Detail *string `json:"detail,omitempty"`
			Id     *string `json:"id,omitempty"`
			Source *struct {
				Parameter *string `json:"parameter,omitempty"`
				Pointer   *string `json:"pointer,omitempty"`
			} `json:"source,omitempty"`
			Title *string `json:"title,omitempty"`
		} `json:"errors,omitempty"`
	}
}

// Status returns HTTPResponse.Status
func (r UpdateWebhookWrap) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r UpdateWebhookWrap) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type ListWebhookInvocationsWrap struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON200      *struct {
		Invocations *[]struct {
			Id                *string `json:"id,omitempty"`
			LatestCompletedAt *string `json:"latest_completed_at,omitempty"`
			Status            *string `json:"status,omitempty"`
			StatusCode        *string `json:"status_code,omitempty"`
		} `json:"invocations,omitempty"`
		Links *struct {
			Next *string `json:"next,omitempty"`
			Prev *string `json:"prev,omitempty"`
		} `json:"links,omitempty"`
		Meta *struct {
			AfterCursor  *string `json:"after_cursor,omitempty"`
			BeforeCursor *string `json:"before_cursor,omitempty"`
			HasMore      *string `json:"has_more,omitempty"`
		} `json:"meta,omitempty"`
	}
}

// Status returns HTTPResponse.Status
func (r ListWebhookInvocationsWrap) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r ListWebhookInvocationsWrap) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type ListWebhookInvocationAttemptsWrap struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON200      *struct {
		Attempts *[]struct {
			CompletedAt  *string `json:"completed_at,omitempty"`
			Id           *string `json:"id,omitempty"`
			InvocationId *string `json:"invocation_id,omitempty"`
			Request      *struct {
				Headers *[]struct {
					Key   *string `json:"key,omitempty"`
					Value *string `json:"value,omitempty"`
				} `json:"headers,omitempty"`
				Payload *struct {
					Dolorcf    *string `json:"dolorcf,omitempty"`
					Euf        *string `json:"euf,omitempty"`
					Ex5        *string `json:"ex_5,omitempty"`
					In9f2      *string `json:"in9f2,omitempty"`
					Ullamcoce8 *string `json:"ullamcoce8,omitempty"`
					Ut1f       *string `json:"ut_1f,omitempty"`
				} `json:"payload,omitempty"`
			} `json:"request,omitempty"`
			Response *struct {
				Headers *[]struct {
					Key   *string `json:"key,omitempty"`
					Value *string `json:"value,omitempty"`
				} `json:"headers,omitempty"`
				Payload *struct {
					Aliqua46   *string `json:"aliqua_46,omitempty"`
					Cupidatat2 *string `json:"cupidatat_2,omitempty"`
					DoloreB2   *string `json:"dolore_b2,omitempty"`
					Tempor0b5  *string `json:"tempor_0b5,omitempty"`
				} `json:"payload,omitempty"`
			} `json:"response,omitempty"`
			Status     *string `json:"status,omitempty"`
			StatusCode *string `json:"status_code,omitempty"`
		} `json:"attempts,omitempty"`
	}
}

// Status returns HTTPResponse.Status
func (r ListWebhookInvocationAttemptsWrap) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r ListWebhookInvocationAttemptsWrap) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type ShowWebhookSigningSecretWrap struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON200      *struct {
		SigningSecret *struct {
			Algorithm *string `json:"algorithm,omitempty"`
			Secret    *string `json:"secret,omitempty"`
		} `json:"signing_secret,omitempty"`
	}
	JSON403 *struct {
		Errors *[]struct {
			Code   *string `json:"code,omitempty"`
			Detail *string `json:"detail,omitempty"`
			Id     *string `json:"id,omitempty"`
			Source *struct {
				Parameter *string `json:"parameter,omitempty"`
				Pointer   *string `json:"pointer,omitempty"`
			} `json:"source,omitempty"`
			Title *string `json:"title,omitempty"`
		} `json:"errors,omitempty"`
	}
	JSON404 *struct {
		Errors *[]struct {
			Code   *string `json:"code,omitempty"`
			Detail *string `json:"detail,omitempty"`
			Id     *string `json:"id,omitempty"`
			Source *struct {
				Parameter *string `json:"parameter,omitempty"`
				Pointer   *string `json:"pointer,omitempty"`
			} `json:"source,omitempty"`
			Title *string `json:"title,omitempty"`
		} `json:"errors,omitempty"`
	}
}

// Status returns HTTPResponse.Status
func (r ShowWebhookSigningSecretWrap) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r ShowWebhookSigningSecretWrap) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type ResetWebhookSigningSecretWrap struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON201      *struct {
		SigningSecret *struct {
			Algorithm *string `json:"algorithm,omitempty"`
			Secret    *string `json:"secret,omitempty"`
		} `json:"signing_secret,omitempty"`
	}
	JSON403 *struct {
		Errors *[]struct {
			Code   *string `json:"code,omitempty"`
			Detail *string `json:"detail,omitempty"`
			Id     *string `json:"id,omitempty"`
			Source *struct {
				Parameter *string `json:"parameter,omitempty"`
				Pointer   *string `json:"pointer,omitempty"`
			} `json:"source,omitempty"`
			Title *string `json:"title,omitempty"`
		} `json:"errors,omitempty"`
	}
	JSON404 *struct {
		Errors *[]struct {
			Code   *string `json:"code,omitempty"`
			Detail *string `json:"detail,omitempty"`
			Id     *string `json:"id,omitempty"`
			Source *struct {
				Parameter *string `json:"parameter,omitempty"`
				Pointer   *string `json:"pointer,omitempty"`
			} `json:"source,omitempty"`
			Title *string `json:"title,omitempty"`
		} `json:"errors,omitempty"`
	}
}

// Status returns HTTPResponse.Status
func (r ResetWebhookSigningSecretWrap) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r ResetWebhookSigningSecretWrap) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

// ListWebhooksWithResponse request returning *ListWebhooksWrap
func (c *ClientWithResponses) ListWebhooksWithResponse(ctx context.Context, params *ListWebhooksParams, reqEditors ...RequestEditorFn) (*ListWebhooksWrap, error) {
	rsp, err := c.ListWebhooks(ctx, params, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseListWebhooksWrap(rsp)
}

// CreateOrCloneWebhookWithBodyWithResponse request with arbitrary body returning *CreateOrCloneWebhookWrap
func (c *ClientWithResponses) CreateOrCloneWebhookWithBodyWithResponse(ctx context.Context, params *CreateOrCloneWebhookParams, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*CreateOrCloneWebhookWrap, error) {
	rsp, err := c.CreateOrCloneWebhookWithBody(ctx, params, contentType, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseCreateOrCloneWebhookWrap(rsp)
}

func (c *ClientWithResponses) CreateOrCloneWebhookWithResponse(ctx context.Context, params *CreateOrCloneWebhookParams, body CreateOrCloneWebhookJSONRequestBody, reqEditors ...RequestEditorFn) (*CreateOrCloneWebhookWrap, error) {
	rsp, err := c.CreateOrCloneWebhook(ctx, params, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseCreateOrCloneWebhookWrap(rsp)
}

// TestWebhookWithBodyWithResponse request with arbitrary body returning *TestWebhookWrap
func (c *ClientWithResponses) TestWebhookWithBodyWithResponse(ctx context.Context, params *TestWebhookParams, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*TestWebhookWrap, error) {
	rsp, err := c.TestWebhookWithBody(ctx, params, contentType, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseTestWebhookWrap(rsp)
}

func (c *ClientWithResponses) TestWebhookWithResponse(ctx context.Context, params *TestWebhookParams, body TestWebhookJSONRequestBody, reqEditors ...RequestEditorFn) (*TestWebhookWrap, error) {
	rsp, err := c.TestWebhook(ctx, params, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseTestWebhookWrap(rsp)
}

// DeleteWebhookWithResponse request returning *DeleteWebhookWrap
func (c *ClientWithResponses) DeleteWebhookWithResponse(ctx context.Context, webhookId string, reqEditors ...RequestEditorFn) (*DeleteWebhookWrap, error) {
	rsp, err := c.DeleteWebhook(ctx, webhookId, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseDeleteWebhookWrap(rsp)
}

// ShowWebhookWithResponse request returning *ShowWebhookWrap
func (c *ClientWithResponses) ShowWebhookWithResponse(ctx context.Context, webhookId string, reqEditors ...RequestEditorFn) (*ShowWebhookWrap, error) {
	rsp, err := c.ShowWebhook(ctx, webhookId, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseShowWebhookWrap(rsp)
}

// PatchWebhookWithBodyWithResponse request with arbitrary body returning *PatchWebhookWrap
func (c *ClientWithResponses) PatchWebhookWithBodyWithResponse(ctx context.Context, webhookId string, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*PatchWebhookWrap, error) {
	rsp, err := c.PatchWebhookWithBody(ctx, webhookId, contentType, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParsePatchWebhookWrap(rsp)
}

func (c *ClientWithResponses) PatchWebhookWithResponse(ctx context.Context, webhookId string, body PatchWebhookJSONRequestBody, reqEditors ...RequestEditorFn) (*PatchWebhookWrap, error) {
	rsp, err := c.PatchWebhook(ctx, webhookId, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParsePatchWebhookWrap(rsp)
}

// UpdateWebhookWithBodyWithResponse request with arbitrary body returning *UpdateWebhookWrap
func (c *ClientWithResponses) UpdateWebhookWithBodyWithResponse(ctx context.Context, webhookId string, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*UpdateWebhookWrap, error) {
	rsp, err := c.UpdateWebhookWithBody(ctx, webhookId, contentType, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseUpdateWebhookWrap(rsp)
}

func (c *ClientWithResponses) UpdateWebhookWithResponse(ctx context.Context, webhookId string, body UpdateWebhookJSONRequestBody, reqEditors ...RequestEditorFn) (*UpdateWebhookWrap, error) {
	rsp, err := c.UpdateWebhook(ctx, webhookId, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseUpdateWebhookWrap(rsp)
}

// ListWebhookInvocationsWithResponse request returning *ListWebhookInvocationsWrap
func (c *ClientWithResponses) ListWebhookInvocationsWithResponse(ctx context.Context, webhookId string, params *ListWebhookInvocationsParams, reqEditors ...RequestEditorFn) (*ListWebhookInvocationsWrap, error) {
	rsp, err := c.ListWebhookInvocations(ctx, webhookId, params, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseListWebhookInvocationsWrap(rsp)
}

// ListWebhookInvocationAttemptsWithResponse request returning *ListWebhookInvocationAttemptsWrap
func (c *ClientWithResponses) ListWebhookInvocationAttemptsWithResponse(ctx context.Context, webhookId string, invocationId string, reqEditors ...RequestEditorFn) (*ListWebhookInvocationAttemptsWrap, error) {
	rsp, err := c.ListWebhookInvocationAttempts(ctx, webhookId, invocationId, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseListWebhookInvocationAttemptsWrap(rsp)
}

// ShowWebhookSigningSecretWithResponse request returning *ShowWebhookSigningSecretWrap
func (c *ClientWithResponses) ShowWebhookSigningSecretWithResponse(ctx context.Context, webhookId string, reqEditors ...RequestEditorFn) (*ShowWebhookSigningSecretWrap, error) {
	rsp, err := c.ShowWebhookSigningSecret(ctx, webhookId, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseShowWebhookSigningSecretWrap(rsp)
}

// ResetWebhookSigningSecretWithResponse request returning *ResetWebhookSigningSecretWrap
func (c *ClientWithResponses) ResetWebhookSigningSecretWithResponse(ctx context.Context, webhookId string, reqEditors ...RequestEditorFn) (*ResetWebhookSigningSecretWrap, error) {
	rsp, err := c.ResetWebhookSigningSecret(ctx, webhookId, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseResetWebhookSigningSecretWrap(rsp)
}

// ParseListWebhooksWrap parses an HTTP response from a ListWebhooksWithResponse call
func ParseListWebhooksWrap(rsp *http.Response) (*ListWebhooksWrap, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &ListWebhooksWrap{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 200:
		var dest struct {
			Links *struct {
				Next *string `json:"next,omitempty"`
				Prev *string `json:"prev,omitempty"`
			} `json:"links,omitempty"`
			Meta *struct {
				AfterCursor  *string `json:"after_cursor,omitempty"`
				BeforeCursor *string `json:"before_cursor,omitempty"`
				HasMore      *string `json:"has_more,omitempty"`
			} `json:"meta,omitempty"`
			Webhooks *[]WebhookWithoutSensitive `json:"webhooks,omitempty"`
		}
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON200 = &dest

	}

	return response, nil
}

// ParseCreateOrCloneWebhookWrap parses an HTTP response from a CreateOrCloneWebhookWithResponse call
func ParseCreateOrCloneWebhookWrap(rsp *http.Response) (*CreateOrCloneWebhookWrap, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &CreateOrCloneWebhookWrap{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 201:
		var dest struct {
			Webhook *WebhookWithoutSensitive `json:"webhook,omitempty"`
		}
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON201 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 400:
		var dest struct {
			Errors *[]struct {
				Code   *string `json:"code,omitempty"`
				Detail *string `json:"detail,omitempty"`
				Id     *string `json:"id,omitempty"`
				Source *struct {
					Parameter *string `json:"parameter,omitempty"`
					Pointer   *string `json:"pointer,omitempty"`
				} `json:"source,omitempty"`
				Title *string `json:"title,omitempty"`
			} `json:"errors,omitempty"`
		}
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON400 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 403:
		var dest struct {
			Errors *[]struct {
				Code   *string `json:"code,omitempty"`
				Detail *string `json:"detail,omitempty"`
				Id     *string `json:"id,omitempty"`
				Source *struct {
					Parameter *string `json:"parameter,omitempty"`
					Pointer   *string `json:"pointer,omitempty"`
				} `json:"source,omitempty"`
				Title *string `json:"title,omitempty"`
			} `json:"errors,omitempty"`
		}
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON403 = &dest

	}

	return response, nil
}

// ParseTestWebhookWrap parses an HTTP response from a TestWebhookWithResponse call
func ParseTestWebhookWrap(rsp *http.Response) (*TestWebhookWrap, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &TestWebhookWrap{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 200:
		var dest struct {
			Response *struct {
				Headers *[]struct {
					Key   *string `json:"key,omitempty"`
					Value *string `json:"value,omitempty"`
				} `json:"headers,omitempty"`
				Payload *string `json:"payload,omitempty"`
				Status  *string `json:"status,omitempty"`
			} `json:"response,omitempty"`
		}
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON200 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 400:
		var dest struct {
			Errors *[]struct {
				Code   *string `json:"code,omitempty"`
				Detail *string `json:"detail,omitempty"`
				Id     *string `json:"id,omitempty"`
				Source *struct {
					Parameter *string `json:"parameter,omitempty"`
					Pointer   *string `json:"pointer,omitempty"`
				} `json:"source,omitempty"`
				Title *string `json:"title,omitempty"`
			} `json:"errors,omitempty"`
		}
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON400 = &dest

	}

	return response, nil
}

// ParseDeleteWebhookWrap parses an HTTP response from a DeleteWebhookWithResponse call
func ParseDeleteWebhookWrap(rsp *http.Response) (*DeleteWebhookWrap, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &DeleteWebhookWrap{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 400:
		var dest struct {
			Errors *[]struct {
				Code   *string `json:"code,omitempty"`
				Detail *string `json:"detail,omitempty"`
				Id     *string `json:"id,omitempty"`
				Source *struct {
					Parameter *string `json:"parameter,omitempty"`
					Pointer   *string `json:"pointer,omitempty"`
				} `json:"source,omitempty"`
				Title *string `json:"title,omitempty"`
			} `json:"errors,omitempty"`
		}
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON400 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 404:
		var dest struct {
			Errors *[]struct {
				Code   *string `json:"code,omitempty"`
				Detail *string `json:"detail,omitempty"`
				Id     *string `json:"id,omitempty"`
				Source *struct {
					Parameter *string `json:"parameter,omitempty"`
					Pointer   *string `json:"pointer,omitempty"`
				} `json:"source,omitempty"`
				Title *string `json:"title,omitempty"`
			} `json:"errors,omitempty"`
		}
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON404 = &dest

	}

	return response, nil
}

// ParseShowWebhookWrap parses an HTTP response from a ShowWebhookWithResponse call
func ParseShowWebhookWrap(rsp *http.Response) (*ShowWebhookWrap, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &ShowWebhookWrap{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 200:
		var dest struct {
			Webhook *WebhookWithoutSensitive `json:"webhook,omitempty"`
		}
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON200 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 400:
		var dest struct {
			Errors *[]struct {
				Code   *string `json:"code,omitempty"`
				Detail *string `json:"detail,omitempty"`
				Id     *string `json:"id,omitempty"`
				Source *struct {
					Parameter *string `json:"parameter,omitempty"`
					Pointer   *string `json:"pointer,omitempty"`
				} `json:"source,omitempty"`
				Title *string `json:"title,omitempty"`
			} `json:"errors,omitempty"`
		}
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON400 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 404:
		var dest struct {
			Errors *[]struct {
				Code   *string `json:"code,omitempty"`
				Detail *string `json:"detail,omitempty"`
				Id     *string `json:"id,omitempty"`
				Source *struct {
					Parameter *string `json:"parameter,omitempty"`
					Pointer   *string `json:"pointer,omitempty"`
				} `json:"source,omitempty"`
				Title *string `json:"title,omitempty"`
			} `json:"errors,omitempty"`
		}
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON404 = &dest

	}

	return response, nil
}

// ParsePatchWebhookWrap parses an HTTP response from a PatchWebhookWithResponse call
func ParsePatchWebhookWrap(rsp *http.Response) (*PatchWebhookWrap, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &PatchWebhookWrap{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 400:
		var dest struct {
			Errors *[]struct {
				Code   *string `json:"code,omitempty"`
				Detail *string `json:"detail,omitempty"`
				Id     *string `json:"id,omitempty"`
				Source *struct {
					Parameter *string `json:"parameter,omitempty"`
					Pointer   *string `json:"pointer,omitempty"`
				} `json:"source,omitempty"`
				Title *string `json:"title,omitempty"`
			} `json:"errors,omitempty"`
		}
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON400 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 404:
		var dest struct {
			Errors *[]struct {
				Code   *string `json:"code,omitempty"`
				Detail *string `json:"detail,omitempty"`
				Id     *string `json:"id,omitempty"`
				Source *struct {
					Parameter *string `json:"parameter,omitempty"`
					Pointer   *string `json:"pointer,omitempty"`
				} `json:"source,omitempty"`
				Title *string `json:"title,omitempty"`
			} `json:"errors,omitempty"`
		}
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON404 = &dest

	}

	return response, nil
}

// ParseUpdateWebhookWrap parses an HTTP response from a UpdateWebhookWithResponse call
func ParseUpdateWebhookWrap(rsp *http.Response) (*UpdateWebhookWrap, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &UpdateWebhookWrap{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 400:
		var dest struct {
			Errors *[]struct {
				Code   *string `json:"code,omitempty"`
				Detail *string `json:"detail,omitempty"`
				Id     *string `json:"id,omitempty"`
				Source *struct {
					Parameter *string `json:"parameter,omitempty"`
					Pointer   *string `json:"pointer,omitempty"`
				} `json:"source,omitempty"`
				Title *string `json:"title,omitempty"`
			} `json:"errors,omitempty"`
		}
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON400 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 404:
		var dest struct {
			Errors *[]struct {
				Code   *string `json:"code,omitempty"`
				Detail *string `json:"detail,omitempty"`
				Id     *string `json:"id,omitempty"`
				Source *struct {
					Parameter *string `json:"parameter,omitempty"`
					Pointer   *string `json:"pointer,omitempty"`
				} `json:"source,omitempty"`
				Title *string `json:"title,omitempty"`
			} `json:"errors,omitempty"`
		}
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON404 = &dest

	}

	return response, nil
}

// ParseListWebhookInvocationsWrap parses an HTTP response from a ListWebhookInvocationsWithResponse call
func ParseListWebhookInvocationsWrap(rsp *http.Response) (*ListWebhookInvocationsWrap, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &ListWebhookInvocationsWrap{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 200:
		var dest struct {
			Invocations *[]struct {
				Id                *string `json:"id,omitempty"`
				LatestCompletedAt *string `json:"latest_completed_at,omitempty"`
				Status            *string `json:"status,omitempty"`
				StatusCode        *string `json:"status_code,omitempty"`
			} `json:"invocations,omitempty"`
			Links *struct {
				Next *string `json:"next,omitempty"`
				Prev *string `json:"prev,omitempty"`
			} `json:"links,omitempty"`
			Meta *struct {
				AfterCursor  *string `json:"after_cursor,omitempty"`
				BeforeCursor *string `json:"before_cursor,omitempty"`
				HasMore      *string `json:"has_more,omitempty"`
			} `json:"meta,omitempty"`
		}
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON200 = &dest

	}

	return response, nil
}

// ParseListWebhookInvocationAttemptsWrap parses an HTTP response from a ListWebhookInvocationAttemptsWithResponse call
func ParseListWebhookInvocationAttemptsWrap(rsp *http.Response) (*ListWebhookInvocationAttemptsWrap, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &ListWebhookInvocationAttemptsWrap{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 200:
		var dest struct {
			Attempts *[]struct {
				CompletedAt  *string `json:"completed_at,omitempty"`
				Id           *string `json:"id,omitempty"`
				InvocationId *string `json:"invocation_id,omitempty"`
				Request      *struct {
					Headers *[]struct {
						Key   *string `json:"key,omitempty"`
						Value *string `json:"value,omitempty"`
					} `json:"headers,omitempty"`
					Payload *struct {
						Dolorcf    *string `json:"dolorcf,omitempty"`
						Euf        *string `json:"euf,omitempty"`
						Ex5        *string `json:"ex_5,omitempty"`
						In9f2      *string `json:"in9f2,omitempty"`
						Ullamcoce8 *string `json:"ullamcoce8,omitempty"`
						Ut1f       *string `json:"ut_1f,omitempty"`
					} `json:"payload,omitempty"`
				} `json:"request,omitempty"`
				Response *struct {
					Headers *[]struct {
						Key   *string `json:"key,omitempty"`
						Value *string `json:"value,omitempty"`
					} `json:"headers,omitempty"`
					Payload *struct {
						Aliqua46   *string `json:"aliqua_46,omitempty"`
						Cupidatat2 *string `json:"cupidatat_2,omitempty"`
						DoloreB2   *string `json:"dolore_b2,omitempty"`
						Tempor0b5  *string `json:"tempor_0b5,omitempty"`
					} `json:"payload,omitempty"`
				} `json:"response,omitempty"`
				Status     *string `json:"status,omitempty"`
				StatusCode *string `json:"status_code,omitempty"`
			} `json:"attempts,omitempty"`
		}
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON200 = &dest

	}

	return response, nil
}

// ParseShowWebhookSigningSecretWrap parses an HTTP response from a ShowWebhookSigningSecretWithResponse call
func ParseShowWebhookSigningSecretWrap(rsp *http.Response) (*ShowWebhookSigningSecretWrap, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &ShowWebhookSigningSecretWrap{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 200:
		var dest struct {
			SigningSecret *struct {
				Algorithm *string `json:"algorithm,omitempty"`
				Secret    *string `json:"secret,omitempty"`
			} `json:"signing_secret,omitempty"`
		}
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON200 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 403:
		var dest struct {
			Errors *[]struct {
				Code   *string `json:"code,omitempty"`
				Detail *string `json:"detail,omitempty"`
				Id     *string `json:"id,omitempty"`
				Source *struct {
					Parameter *string `json:"parameter,omitempty"`
					Pointer   *string `json:"pointer,omitempty"`
				} `json:"source,omitempty"`
				Title *string `json:"title,omitempty"`
			} `json:"errors,omitempty"`
		}
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON403 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 404:
		var dest struct {
			Errors *[]struct {
				Code   *string `json:"code,omitempty"`
				Detail *string `json:"detail,omitempty"`
				Id     *string `json:"id,omitempty"`
				Source *struct {
					Parameter *string `json:"parameter,omitempty"`
					Pointer   *string `json:"pointer,omitempty"`
				} `json:"source,omitempty"`
				Title *string `json:"title,omitempty"`
			} `json:"errors,omitempty"`
		}
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON404 = &dest

	}

	return response, nil
}

// ParseResetWebhookSigningSecretWrap parses an HTTP response from a ResetWebhookSigningSecretWithResponse call
func ParseResetWebhookSigningSecretWrap(rsp *http.Response) (*ResetWebhookSigningSecretWrap, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &ResetWebhookSigningSecretWrap{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 201:
		var dest struct {
			SigningSecret *struct {
				Algorithm *string `json:"algorithm,omitempty"`
				Secret    *string `json:"secret,omitempty"`
			} `json:"signing_secret,omitempty"`
		}
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON201 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 403:
		var dest struct {
			Errors *[]struct {
				Code   *string `json:"code,omitempty"`
				Detail *string `json:"detail,omitempty"`
				Id     *string `json:"id,omitempty"`
				Source *struct {
					Parameter *string `json:"parameter,omitempty"`
					Pointer   *string `json:"pointer,omitempty"`
				} `json:"source,omitempty"`
				Title *string `json:"title,omitempty"`
			} `json:"errors,omitempty"`
		}
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON403 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 404:
		var dest struct {
			Errors *[]struct {
				Code   *string `json:"code,omitempty"`
				Detail *string `json:"detail,omitempty"`
				Id     *string `json:"id,omitempty"`
				Source *struct {
					Parameter *string `json:"parameter,omitempty"`
					Pointer   *string `json:"pointer,omitempty"`
				} `json:"source,omitempty"`
				Title *string `json:"title,omitempty"`
			} `json:"errors,omitempty"`
		}
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON404 = &dest

	}

	return response, nil
}
