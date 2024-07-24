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

func response200(t *testing.T) *zendesk_webhook_api.ShowWebhookWrap {
	var response zendesk_webhook_api.ShowWebhookWrap
	str := `{"JSON200": {"webhook": {"name": "test-webhook", "id": "123456", "created_at": "2024-07-22T15:02:03Z",
	"updated_at": "2024-08-23T09:01:03.12345Z", "endpoint": "https://example.com", "description": "test webhook",
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

func TestWebhookMapper_PutWebhookShowResponseToStateModel(t *testing.T) {
	type args struct {
		ctx                 context.Context
		showWebhookResponse *zendesk_webhook_api.ShowWebhookWrap
		webhookState        *WebhookModel
	}
	tests := []struct {
		name            string
		args            args
		wantDiagnostics diag.Diagnostics
		assertFunction  AssertModel
	}{
		{name: "all attributes given",
			args: args{context.Background(),
				response200(t),
				&WebhookModel{},
			},
			wantDiagnostics: nil,
			assertFunction:  modelForResponse200AllAttributes},
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

func TestWebhookMapper_MapToRequestBody(t *testing.T) {
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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &WebhookMapper{}
			gotRequestBody, gotDiagnostics := m.MapToRequestBody(tt.args.ctx, tt.args.model)
			assert.DeepEqual(t, gotDiagnostics, tt.wantDiagnostics)
			if !gotDiagnostics.HasError() {
				assert.DeepEqual(t, gotRequestBody, tt.wantRequestBody)
			}
		})
	}
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

func createRequestBody200() *zendesk_webhook_api.CreateOrCloneWebhookJSONRequestBody {
	requestBody := zendesk_webhook_api.CreateOrCloneWebhookJSONRequestBody{}

	str := `{"webhook": {"name": "test-webhook", 
	"authentication": {"type": "basic_auth", "data": {"username": "test-user", "password": "test-word"},
	"add_position": "header"},
	 "endpoint": "https://example.com", "description": "test webhook",
	"request_format": "json", "http_method": "POST", 
	"custom_headers": {"header1": "value1", "header2": "value2"},
	"external_source": {"data": {"app_id": "345", "installation_id":"id"}, "type": "app_installation"},
	"status": "active",
	"subscriptions": ["subscription1"],
	"signing_secret": { "secret": "secret-value", "algorithm": "SHA256"}
	}}`

	err := json.Unmarshal([]byte(str), &requestBody)

	if err != nil {
		panic(err)
	}

	return &requestBody
}

func modelForResponse200AllAttributes(t *testing.T, model *WebhookModel) {

	assert.Equal(t, types.StringValue("123456"), model.WebhookId)
	assert.Equal(t, "123456", model.Webhook.Id.ValueString())
	assert.Equal(t, "test-webhook", model.Webhook.Name.ValueString())
	assert.Equal(t, "2024-07-22T15:02:03Z", model.Webhook.CreatedAt.ValueString())
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

}

type AssertModel func(t *testing.T, model *WebhookModel)
