package resource_webhook

import (
	"context"
	"encoding/json"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"gotest.tools/v3/assert"
	"terraform-provider-zendesk/zendesk_webhook_api"
	"testing"
)

func TestWebhookMapper_PutWebhookShowResponseToStateModel(t *testing.T) {
	type args struct {
		ctx                 context.Context
		showWebhookResponse *zendesk_webhook_api.ShowWebhookWrap
		webhookState        *WebhookModel
	}
	webhookValueNull := WebhookModel{Webhook: NewWebhookValueNull()}
	tests := []struct {
		name            string
		args            args
		wantDiagnostics diag.Diagnostics
		assertFunction  AssertModel
	}{
		{name: "all non-sensitive attributes given",
			args: args{context.Background(),
				responseShow200(t),
				&webhookValueNull,
			},
			wantDiagnostics: nil,
			assertFunction:  modelForResponse200AllNonSensitiveAttributes},
		{name: "no headers",
			args: args{context.Background(),
				responseShow200NoHeaders(t),
				&webhookValueNull,
			},
			wantDiagnostics: nil,
			assertFunction:  modelForResponse200NoHeaders},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			we := &WebhookMapper{}
			gotDiagnostics := we.PutWebhookShowResponseToStateModel(tt.args.ctx, tt.args.showWebhookResponse, tt.args.webhookState)
			assert.DeepEqual(t, gotDiagnostics, tt.wantDiagnostics)
			if !gotDiagnostics.HasError() {
				tt.assertFunction(t, tt.args.webhookState)

			}
		})
	}
}

func TestWebhookMapper_PutCreateResponseToStateModel(t *testing.T) {
	type args struct {
		ctx                 context.Context
		showWebhookResponse *zendesk_webhook_api.CreateOrCloneWebhookWrap
		webhookState        *WebhookModel
	}
	tests := []struct {
		name           string
		args           args
		assertFunction AssertModel
	}{
		{name: "all attributes given",
			args: args{context.Background(),
				responseCreate201(t),
				getCreateWebhookModel(),
			},
			assertFunction: modelForCreateResponse200AllAttributes},
		{name: "no headers",
			args: args{context.Background(),
				responseCreate201NoHeaders(t),
				getCreateWebhookModelNoHeaders(),
			},
			assertFunction: modelForCreateResponse200NoHeaders},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			we := &WebhookMapper{}
			diags := we.UpdateAttributesWithCreateResponse(context.Background(), tt.args.showWebhookResponse, tt.args.webhookState)
			assert.Equal(t, false, diags.HasError())
			tt.assertFunction(t, tt.args.webhookState)

		})
	}
}

func responseCreate201(t *testing.T) *zendesk_webhook_api.CreateOrCloneWebhookWrap {
	response := &zendesk_webhook_api.CreateOrCloneWebhookWrap{}
	str := `{"JSON201":{"webhook": {"id": "123456", "name": "test-webhook", 
	"authentication": {"type": "basic_auth", "data": {"username": "test-user", "password": "test-word"},
	"add_position": "header"},
	 "endpoint": "https://example.com", "description": "test webhook",
	"request_format": "json", "http_method": "POST", 
	"custom_headers": {"header1": "value1", "header2": "value2"},
	"status": "active",
	"subscriptions": ["subscription1"],
	"external_source": {"data": {"app_id": "345", "installation_id":"id"}, "type": "app_installation"},
	"signing_secret": { "secret": "secret-value", "algorithm": "SHA256"},
	"created_at": "2024-07-25T09:58:03Z",
    "created_by": "19293454834333"
	}}, "Response":{"StatusCode":201}}`

	err := json.Unmarshal([]byte(str), response)
	if err != nil {
		t.Errorf("Error unmarshalling response: %v", err)
		panic(err)
	}
	return response
}

func responseCreate201NoHeaders(t *testing.T) *zendesk_webhook_api.CreateOrCloneWebhookWrap {
	response := &zendesk_webhook_api.CreateOrCloneWebhookWrap{}
	str := `{"JSON201":{"webhook": {"id": "123456", "name": "test-webhook", 
	"authentication": {"type": "basic_auth", "data": {"username": "test-user", "password": "test-word"},
	"add_position": "header"},
	 "endpoint": "https://example.com", "description": "test webhook",
	"request_format": "json", "http_method": "POST", 
	"status": "active",
	"subscriptions": ["subscription1"],
	"external_source": {"data": {"app_id": "345", "installation_id":"id"}, "type": "app_installation"},
	"signing_secret": { "secret": "secret-value", "algorithm": "SHA256"},
	"created_at": "2024-07-25T09:58:03Z",
    "created_by": "19293454834333"
	}}, "Response":{"StatusCode":201}}`

	err := json.Unmarshal([]byte(str), response)
	if err != nil {
		t.Errorf("Error unmarshalling response: %v", err)
		panic(err)
	}
	return response
}

func TestWebhookMapper_MapToCreateRequestBody(t *testing.T) {
	type args struct {
		ctx   context.Context
		model *WebhookModel
	}
	tests := []struct {
		name            string
		args            args
		wantDiagnostics diag.Diagnostics
		wantRequestBody *zendesk_webhook_api.CreateOrCloneWebhookJSONRequestBody
	}{
		{name: "200",
			args: args{context.Background(),
				getCreateWebhookModel()},
			wantDiagnostics: nil,
			wantRequestBody: createRequestBody200()},
		{name: "200",
			args: args{context.Background(),
				getCreateWebhookModelNoHeaders()},
			wantDiagnostics: nil,
			wantRequestBody: createRequestBody200NoHeaders()},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &WebhookMapper{}
			gotRequestBody, gotDiagnostics := m.MapToCreateRequestBody(tt.args.ctx, tt.args.model)
			assert.DeepEqual(t, gotDiagnostics, tt.wantDiagnostics)
			if !gotDiagnostics.HasError() {
				assert.DeepEqual(t, gotRequestBody, tt.wantRequestBody)
			}
		})
	}
}

func TestWebhookMapper_MapToUpdateRequestBody(t *testing.T) {
	type args struct {
		ctx context.Context
		r   *WebhookModel
	}
	tests := []struct {
		name            string
		args            args
		want            zendesk_webhook_api.UpdateWebhookJSONRequestBody
		wantDiagnostics diag.Diagnostics
	}{
		{name: "all attributes filled",
			args: args{ctx: context.Background(),
				r: getUpdateWebhookModel()},
			want:            getUpdateRequestBody(),
			wantDiagnostics: nil},
		{name: "Authentication attributes filled",
			args: args{ctx: context.Background(),
				r: getUpdateWebhookModelOnlyAuthenticationAttributes()},
			want:            getUpdateRequestBodyOnlyAuthenticationAttributes(),
			wantDiagnostics: nil},
		{name: "Endpoint attributes filled",
			args: args{ctx: context.Background(),
				r: getUpdateWebhookModelOnlyEndpointAttributes()},
			want:            getUpdateRequestBodyOnlyEndpointAttributes(),
			wantDiagnostics: nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &WebhookMapper{}
			gotRequestBody, gotDiagnostics := m.MapToUpdateRequestBody(tt.args.ctx, tt.args.r)

			assert.DeepEqual(t, tt.wantDiagnostics, gotDiagnostics)
			if gotDiagnostics == nil || !gotDiagnostics.HasError() {
				assert.DeepEqual(t, tt.want, gotRequestBody)
			}
		})
	}
}

func getUpdateRequestBodyOnlyAuthenticationAttributes() zendesk_webhook_api.UpdateWebhookJSONRequestBody {
	password := "test-word2"
	username := "test-user2"
	return zendesk_webhook_api.UpdateWebhookJSONRequestBody{
		Webhook: &zendesk_webhook_api.WebhookWithSensitiveData{
			Authentication: &zendesk_webhook_api.Authentication{
				AddPosition: "header",
				Type:        "basic_auth",
				Data: &zendesk_webhook_api.AuthenticationData{
					Password: &password,
					Username: &username,
				},
			},
		},
	}
}

func getUpdateRequestBodyOnlyEndpointAttributes() zendesk_webhook_api.UpdateWebhookJSONRequestBody {

	return zendesk_webhook_api.UpdateWebhookJSONRequestBody{
		Webhook: &zendesk_webhook_api.WebhookWithSensitiveData{
			Endpoint:      "https://example.com",
			HttpMethod:    "GET",
			RequestFormat: "json",
		},
	}
}

func getUpdateRequestBody() zendesk_webhook_api.UpdateWebhookJSONRequestBody {
	return zendesk_webhook_api.UpdateWebhookJSONRequestBody{
		Webhook: createRequestBody200().Webhook,
	}
}

func getUpdateWebhookModelOnlyAuthenticationAttributes() *WebhookModel {
	model := WebhookModel{}
	auth := NewAuthenticationValueNull()
	auth.AuthenticationType = types.StringValue("basic_auth")
	dataValue := NewDataValueNull()
	dataValue.Username = types.StringValue("test-user2")
	dataValue.Password = types.StringValue("test-word2")
	auth.Data, _ = dataValue.ToObjectValue(context.Background())
	auth.AddPosition = types.StringValue("header")
	model.Webhook.Authentication, _ = auth.ToObjectValue(context.Background())

	return &model
}
func getUpdateWebhookModelOnlyEndpointAttributes() *WebhookModel {
	model := WebhookModel{}
	model.Webhook.Endpoint = types.StringValue("https://example.com")
	model.Webhook.HttpMethod = types.StringValue("GET")
	model.Webhook.RequestFormat = types.StringValue("json")
	return &model
}

func getCreateWebhookModel() *WebhookModel {
	model := WebhookModel{}
	auth := NewAuthenticationValueNull()
	auth.AuthenticationType = types.StringValue("basic_auth")
	dataValue := NewDataValueNull()
	dataValue.Username = types.StringValue("test-user")
	dataValue.Password = types.StringValue("test-word")
	auth.Data, _ = dataValue.ToObjectValue(context.Background())
	auth.AddPosition = types.StringValue("header")
	model.Webhook.Authentication, _ = auth.ToObjectValue(context.Background())

	headers := make(map[string]attr.Value)
	headers["header1"] = types.StringValue("value1")
	headers["header2"] = types.StringValue("value2")

	model.Webhook.CustomHeaders, _ = types.MapValue(types.StringType, headers)
	model.Webhook.Description = types.StringValue("test webhook")
	model.Webhook.Endpoint = types.StringValue("https://example.com")
	externalSourceValue := NewExternalSourceValueNull()
	externalSourceValue.ExternalSourceType = types.StringValue("app_installation")
	externalSourceData := NewExternalSourceDataValueNull()
	externalSourceData.AppId = types.StringValue("345")
	externalSourceData.InstallationId = types.StringValue("id")
	externalSourceValue.ExternalSourceData, _ = externalSourceData.ToObjectValue(context.Background())
	model.Webhook.ExternalSource, _ = externalSourceValue.ToObjectValue(context.Background())

	model.Webhook.HttpMethod = types.StringValue("POST")
	model.Webhook.Name = types.StringValue("test-webhook")
	model.Webhook.RequestFormat = types.StringValue("json")
	model.Webhook.Status = types.StringValue("active")
	model.Webhook.Subscriptions, _ = types.ListValue(types.StringType, []attr.Value{types.StringValue("subscription1")})
	signingSecretValue := NewSigningSecretValueNull()
	signingSecretValue.Algorithm = types.StringValue("SHA256")
	signingSecretValue.Secret = types.StringValue("secret-value")
	model.Webhook.SigningSecret, _ = signingSecretValue.ToObjectValue(context.Background())
	return &model
}

func getCreateWebhookModelNoHeaders() *WebhookModel {
	model := WebhookModel{}
	auth := NewAuthenticationValueNull()
	auth.AuthenticationType = types.StringValue("basic_auth")
	dataValue := NewDataValueNull()
	dataValue.Username = types.StringValue("test-user")
	dataValue.Password = types.StringValue("test-word")
	auth.Data, _ = dataValue.ToObjectValue(context.Background())
	auth.AddPosition = types.StringValue("header")
	model.Webhook.Authentication, _ = auth.ToObjectValue(context.Background())

	headers := make(map[string]attr.Value)

	model.Webhook.CustomHeaders, _ = types.MapValue(types.StringType, headers)
	model.Webhook.Description = types.StringValue("test webhook")
	model.Webhook.Endpoint = types.StringValue("https://example.com")
	externalSourceValue := NewExternalSourceValueNull()
	externalSourceValue.ExternalSourceType = types.StringValue("app_installation")
	externalSourceData := NewExternalSourceDataValueNull()
	externalSourceData.AppId = types.StringValue("345")
	externalSourceData.InstallationId = types.StringValue("id")
	externalSourceValue.ExternalSourceData, _ = externalSourceData.ToObjectValue(context.Background())
	model.Webhook.ExternalSource, _ = externalSourceValue.ToObjectValue(context.Background())

	model.Webhook.HttpMethod = types.StringValue("POST")
	model.Webhook.Name = types.StringValue("test-webhook")
	model.Webhook.RequestFormat = types.StringValue("json")
	model.Webhook.Status = types.StringValue("active")
	model.Webhook.Subscriptions, _ = types.ListValue(types.StringType, []attr.Value{types.StringValue("subscription1")})
	signingSecretValue := NewSigningSecretValueNull()
	signingSecretValue.Algorithm = types.StringValue("SHA256")
	signingSecretValue.Secret = types.StringValue("secret-value")
	model.Webhook.SigningSecret, _ = signingSecretValue.ToObjectValue(context.Background())
	return &model
}

func getUpdateWebhookModel() *WebhookModel {
	model := WebhookModel{}
	auth := NewAuthenticationValueNull()
	auth.AuthenticationType = types.StringValue("basic_auth")
	dataValue := NewDataValueNull()
	dataValue.Username = types.StringValue("test-user")
	dataValue.Password = types.StringValue("test-word")
	auth.Data, _ = dataValue.ToObjectValue(context.Background())
	auth.AddPosition = types.StringValue("header")
	model.Webhook.Authentication, _ = auth.ToObjectValue(context.Background())

	headers := make(map[string]attr.Value)
	headers["header1"] = types.StringValue("value1")
	headers["header2"] = types.StringValue("value2")

	model.Webhook.CustomHeaders, _ = types.MapValue(types.StringType, headers)
	model.Webhook.Description = types.StringValue("test webhook")
	model.Webhook.Endpoint = types.StringValue("https://example.com")
	externalSourceValue := NewExternalSourceValueNull()
	externalSourceValue.ExternalSourceType = types.StringValue("app_installation")
	externalSourceData := NewExternalSourceDataValueNull()
	externalSourceData.AppId = types.StringValue("345")
	externalSourceData.InstallationId = types.StringValue("id")
	externalSourceValue.ExternalSourceData, _ = externalSourceData.ToObjectValue(context.Background())
	model.Webhook.ExternalSource, _ = externalSourceValue.ToObjectValue(context.Background())

	model.Webhook.HttpMethod = types.StringValue("POST")
	model.Webhook.Name = types.StringValue("test-webhook")
	model.Webhook.RequestFormat = types.StringValue("json")
	model.Webhook.Status = types.StringValue("active")
	model.Webhook.Subscriptions, _ = types.ListValue(types.StringType, []attr.Value{types.StringValue("subscription1")})
	signingSecretValue := NewSigningSecretValueNull()
	signingSecretValue.Algorithm = types.StringValue("SHA256")
	signingSecretValue.Secret = types.StringValue("secret-value")
	model.Webhook.SigningSecret, _ = signingSecretValue.ToObjectValue(context.Background())
	return &model
}

func createRequestBody200() *zendesk_webhook_api.CreateOrCloneWebhookJSONRequestBody {
	requestBody := zendesk_webhook_api.CreateOrCloneWebhookJSONRequestBody{}

	str := `{"webhook": {"name": "test-webhook", 
	"authentication": {"type": "basic_auth", "data": {"username": "test-user", "password": "test-word"},
	"add_position": "header"},
	 "endpoint": "https://example.com", "description": "test webhook",
	"request_format": "json", "http_method": "POST", 
	"custom_headers": {"header1": "value1", "header2": "value2"},
	"status": "active",
	"subscriptions": ["subscription1"],
	"external_source": {"data": {"app_id": "345", "installation_id":"id"}, "type": "app_installation"},
	"signing_secret": { "secret": "secret-value", "algorithm": "SHA256"}
	}}`

	err := json.Unmarshal([]byte(str), &requestBody)

	if err != nil {
		panic(err)
	}

	return &requestBody
}

func createRequestBody200NoHeaders() *zendesk_webhook_api.CreateOrCloneWebhookJSONRequestBody {
	requestBody := zendesk_webhook_api.CreateOrCloneWebhookJSONRequestBody{}

	str := `{"webhook": {"name": "test-webhook", 
	"authentication": {"type": "basic_auth", "data": {"username": "test-user", "password": "test-word"},
	"add_position": "header"},
	 "endpoint": "https://example.com", "description": "test webhook",
	"request_format": "json", "http_method": "POST", 	
	"status": "active",
	"subscriptions": ["subscription1"],
	"external_source": {"data": {"app_id": "345", "installation_id":"id"}, "type": "app_installation"},
	"signing_secret": { "secret": "secret-value", "algorithm": "SHA256"}
	}}`

	err := json.Unmarshal([]byte(str), &requestBody)

	if err != nil {
		panic(err)
	}

	return &requestBody
}

func responseShow200(t *testing.T) *zendesk_webhook_api.ShowWebhookWrap {
	var response zendesk_webhook_api.ShowWebhookWrap
	str := `{"JSON200": {"webhook": {"name": "test-webhook", "id": "123456", "created_at": "2024-07-25T09:58:03Z",
	"updated_at": "2024-08-23T09:01:03.12345Z", "endpoint": "https://example.com", "description": "test webhook",
	"authentication": {"type": "basic_auth", "data": {"username": "test-user"}, "add_position": "header"},
	"request_format": "json", "http_method": "POST", "created_by": "tester", "updated_by": "tester2",
	"custom_headers": {"header1": "value1", "header2": "value2"},
	"external_source": {"data": {"app_id": "345", "installation_id":"id"}, "type": "app_installation"},
	"status": "active",
	"subscriptions": ["subscription1"],
	"signing_secret": { "secret": "secret-value", "algorithm": "SHA256"}
	}}}`

	err := json.Unmarshal([]byte(str), &response)
	if err != nil {
		t.Errorf("Error unmarshalling response: %v", err)
	}
	return &response
}

func responseShow200NoHeaders(t *testing.T) *zendesk_webhook_api.ShowWebhookWrap {
	var response zendesk_webhook_api.ShowWebhookWrap
	str := `{"JSON200": {"webhook": {"name": "test-webhook", "id": "123456", "created_at": "2024-07-25T09:58:03Z",
	"updated_at": "2024-08-23T09:01:03.12345Z", "endpoint": "https://example.com", "description": "test webhook",
	"authentication": {"type": "basic_auth", "data": {"username": "test-user"}, "add_position": "header"},
	"request_format": "json", "http_method": "POST", "created_by": "tester", "updated_by": "tester2",
	"custom_headers": null,
	"external_source": {"data": {"app_id": "345", "installation_id":"id"}, "type": "app_installation"},
	"status": "active",
	"subscriptions": ["subscription1"],
	"signing_secret": { "secret": "secret-value", "algorithm": "SHA256"}
	}}}`

	err := json.Unmarshal([]byte(str), &response)
	if err != nil {
		t.Errorf("Error unmarshalling response: %v", err)
	}
	return &response
}

func modelForCreateResponse200AllAttributes(t *testing.T, model *WebhookModel) {

	assert.Equal(t, types.StringValue("123456"), model.WebhookId)
	assert.Equal(t, "123456", model.Webhook.Id.ValueString())
	assert.Equal(t, "test-webhook", model.Webhook.Name.ValueString())
	assert.Equal(t, "2024-07-25T09:58:03Z", model.Webhook.CreatedAt.ValueString())
	assert.Equal(t, true, model.Webhook.UpdatedAt.IsNull())
	assert.Equal(t, "19293454834333", model.Webhook.CreatedBy.ValueString())
	assert.Equal(t, true, model.Webhook.UpdatedBy.IsNull())
	assert.Equal(t, "https://example.com", model.Webhook.Endpoint.ValueString())
	assert.Equal(t, "test webhook", model.Webhook.Description.ValueString())
	assert.Equal(t, "json", model.Webhook.RequestFormat.ValueString())
	assert.Equal(t, "POST", model.Webhook.HttpMethod.ValueString())
	assert.Equal(t, "\"app_installation\"", model.Webhook.ExternalSource.Attributes()["type"].String())

	externalSourceData, ok := model.Webhook.ExternalSource.Attributes()["external_source_data"].(types.Object)
	assert.Equal(t, true, ok)
	assert.Equal(t, "\"345\"", externalSourceData.Attributes()["app_id"].String())
	assert.Equal(t, "\"id\"", externalSourceData.Attributes()["installation_id"].String())
	assert.Equal(t, "\"value1\"", model.Webhook.CustomHeaders.Elements()["header1"].String())
	assert.Equal(t, "\"value2\"", model.Webhook.CustomHeaders.Elements()["header2"].String())
	assert.Equal(t, "active", model.Webhook.Status.ValueString())
	assert.Equal(t, "\"subscription1\"", model.Webhook.Subscriptions.Elements()[0].String())
	assert.Equal(t, "\"secret-value\"", model.Webhook.SigningSecret.Attributes()["secret"].String())
	assert.Equal(t, "\"SHA256\"", model.Webhook.SigningSecret.Attributes()["algorithm"].String())

	assert.Equal(t, "\"header\"", model.Webhook.Authentication.Attributes()["add_position"].String())
	assert.Equal(t, "\"basic_auth\"", model.Webhook.Authentication.Attributes()["type"].String())
	dataObj, ok := model.Webhook.Authentication.Attributes()["data"].(types.Object)
	assert.Equal(t, true, ok)
	assert.Equal(t, "\"test-user\"", dataObj.Attributes()["username"].String())
	assert.Equal(t, "\"test-word\"", dataObj.Attributes()["password"].String())

}

func modelForCreateResponse200NoHeaders(t *testing.T, model *WebhookModel) {

	assert.Equal(t, types.StringValue("123456"), model.WebhookId)
	assert.Equal(t, "123456", model.Webhook.Id.ValueString())
	assert.Equal(t, "test-webhook", model.Webhook.Name.ValueString())
	assert.Equal(t, "2024-07-25T09:58:03Z", model.Webhook.CreatedAt.ValueString())
	assert.Equal(t, true, model.Webhook.UpdatedAt.IsNull())
	assert.Equal(t, "19293454834333", model.Webhook.CreatedBy.ValueString())
	assert.Equal(t, true, model.Webhook.UpdatedBy.IsNull())
	assert.Equal(t, "https://example.com", model.Webhook.Endpoint.ValueString())
	assert.Equal(t, "test webhook", model.Webhook.Description.ValueString())
	assert.Equal(t, "json", model.Webhook.RequestFormat.ValueString())
	assert.Equal(t, "POST", model.Webhook.HttpMethod.ValueString())
	assert.Equal(t, "\"app_installation\"", model.Webhook.ExternalSource.Attributes()["type"].String())

	externalSourceData, ok := model.Webhook.ExternalSource.Attributes()["external_source_data"].(types.Object)
	assert.Equal(t, true, ok)
	assert.Equal(t, "\"345\"", externalSourceData.Attributes()["app_id"].String())
	assert.Equal(t, "\"id\"", externalSourceData.Attributes()["installation_id"].String())
	assert.Equal(t, true, model.Webhook.CustomHeaders.IsNull())
	assert.Equal(t, "active", model.Webhook.Status.ValueString())
	assert.Equal(t, "\"subscription1\"", model.Webhook.Subscriptions.Elements()[0].String())
	assert.Equal(t, "\"secret-value\"", model.Webhook.SigningSecret.Attributes()["secret"].String())
	assert.Equal(t, "\"SHA256\"", model.Webhook.SigningSecret.Attributes()["algorithm"].String())

	assert.Equal(t, "\"header\"", model.Webhook.Authentication.Attributes()["add_position"].String())
	assert.Equal(t, "\"basic_auth\"", model.Webhook.Authentication.Attributes()["type"].String())
	dataObj, ok := model.Webhook.Authentication.Attributes()["data"].(types.Object)
	assert.Equal(t, true, ok)
	assert.Equal(t, "\"test-user\"", dataObj.Attributes()["username"].String())
	assert.Equal(t, "\"test-word\"", dataObj.Attributes()["password"].String())

}

func modelForResponse200AllNonSensitiveAttributes(t *testing.T, model *WebhookModel) {

	assert.Equal(t, types.StringValue("123456"), model.WebhookId)
	assert.Equal(t, "123456", model.Webhook.Id.ValueString())
	assert.Equal(t, "test-webhook", model.Webhook.Name.ValueString())
	assert.Equal(t, "2024-07-25T09:58:03Z", model.Webhook.CreatedAt.ValueString())
	assert.Equal(t, "2024-08-23T09:01:03.12345Z", model.Webhook.UpdatedAt.ValueString())
	assert.Equal(t, "tester", model.Webhook.CreatedBy.ValueString())
	assert.Equal(t, "tester2", model.Webhook.UpdatedBy.ValueString())
	assert.Equal(t, "https://example.com", model.Webhook.Endpoint.ValueString())
	assert.Equal(t, "test webhook", model.Webhook.Description.ValueString())
	assert.Equal(t, "json", model.Webhook.RequestFormat.ValueString())
	assert.Equal(t, "POST", model.Webhook.HttpMethod.ValueString())
	assert.Equal(t, "\"app_installation\"", model.Webhook.ExternalSource.Attributes()["type"].String())

	externalSourceData, ok := model.Webhook.ExternalSource.Attributes()["external_source_data"].(types.Object)
	assert.Equal(t, true, ok)
	assert.Equal(t, "\"345\"", externalSourceData.Attributes()["app_id"].String())
	assert.Equal(t, "\"id\"", externalSourceData.Attributes()["installation_id"].String())
	assert.Equal(t, "\"value1\"", model.Webhook.CustomHeaders.Elements()["header1"].String())
	assert.Equal(t, "\"value2\"", model.Webhook.CustomHeaders.Elements()["header2"].String())
	assert.Equal(t, "active", model.Webhook.Status.ValueString())
	assert.Equal(t, "\"subscription1\"", model.Webhook.Subscriptions.Elements()[0].String())
	assert.Equal(t, "\"secret-value\"", model.Webhook.SigningSecret.Attributes()["secret"].String())
	assert.Equal(t, "\"SHA256\"", model.Webhook.SigningSecret.Attributes()["algorithm"].String())

	assert.Equal(t, "\"header\"", model.Webhook.Authentication.Attributes()["add_position"].String())
	assert.Equal(t, "\"basic_auth\"", model.Webhook.Authentication.Attributes()["type"].String())

	object, ok := model.Webhook.Authentication.Attributes()["data"].(types.Object)
	assert.Equal(t, true, ok)
	assert.Equal(t, "\"test-user\"", object.Attributes()["username"].String())
	assert.Equal(t, true, object.Attributes()["password"].IsNull())

}

func modelForResponse200NoHeaders(t *testing.T, model *WebhookModel) {

	assert.Equal(t, types.StringValue("123456"), model.WebhookId)
	assert.Equal(t, "123456", model.Webhook.Id.ValueString())
	assert.Equal(t, "test-webhook", model.Webhook.Name.ValueString())
	assert.Equal(t, "2024-07-25T09:58:03Z", model.Webhook.CreatedAt.ValueString())
	assert.Equal(t, "2024-08-23T09:01:03.12345Z", model.Webhook.UpdatedAt.ValueString())
	assert.Equal(t, "tester", model.Webhook.CreatedBy.ValueString())
	assert.Equal(t, "tester2", model.Webhook.UpdatedBy.ValueString())
	assert.Equal(t, "https://example.com", model.Webhook.Endpoint.ValueString())
	assert.Equal(t, "test webhook", model.Webhook.Description.ValueString())
	assert.Equal(t, "json", model.Webhook.RequestFormat.ValueString())
	assert.Equal(t, "POST", model.Webhook.HttpMethod.ValueString())
	assert.Equal(t, "\"app_installation\"", model.Webhook.ExternalSource.Attributes()["type"].String())

	externalSourceData, ok := model.Webhook.ExternalSource.Attributes()["external_source_data"].(types.Object)
	assert.Equal(t, true, ok)
	assert.Equal(t, "\"345\"", externalSourceData.Attributes()["app_id"].String())
	assert.Equal(t, "\"id\"", externalSourceData.Attributes()["installation_id"].String())
	assert.Equal(t, true, model.Webhook.CustomHeaders.IsNull())
	assert.Equal(t, "active", model.Webhook.Status.ValueString())
	assert.Equal(t, "\"subscription1\"", model.Webhook.Subscriptions.Elements()[0].String())
	assert.Equal(t, "\"secret-value\"", model.Webhook.SigningSecret.Attributes()["secret"].String())
	assert.Equal(t, "\"SHA256\"", model.Webhook.SigningSecret.Attributes()["algorithm"].String())

	assert.Equal(t, "\"header\"", model.Webhook.Authentication.Attributes()["add_position"].String())
	assert.Equal(t, "\"basic_auth\"", model.Webhook.Authentication.Attributes()["type"].String())

	object, ok := model.Webhook.Authentication.Attributes()["data"].(types.Object)
	assert.Equal(t, true, ok)
	assert.Equal(t, "\"test-user\"", object.Attributes()["username"].String())
	assert.Equal(t, true, object.Attributes()["password"].IsNull())

}

type AssertModel func(t *testing.T, model *WebhookModel)
