package resource_webhook

import (
	"context"
	"encoding/json"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"gotest.tools/v3/assert"
	"reflect"
	"terraform-provider-zendesk/zendesk_webhook_api"
	"testing"
)

func TestNewWebhookMapper(t *testing.T) {
	tests := []struct {
		name string
		want *WebhookMapper
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewWebhookMapper(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewWebhookMapper() = %v, want %v", got, tt.want)
			}
		})
	}
}
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
	assert.Equal(t, "app_installation", model.Webhook.ExternalSource.Attributes()["type"].(types.String).ValueString())
	assert.Equal(t, "345", model.Webhook.ExternalSource.Attributes()["external_source_data"].(types.Object).Attributes()["app_id"].(types.String).ValueString())
	assert.Equal(t, "id", model.Webhook.ExternalSource.Attributes()["external_source_data"].(types.Object).Attributes()["installation_id"].(types.String).ValueString())
	assert.Equal(t, "value1", model.Webhook.CustomHeaders.Elements()["header1"].(types.String).ValueString())
	assert.Equal(t, "value2", model.Webhook.CustomHeaders.Elements()["header2"].(types.String).ValueString())
	assert.Equal(t, "active", model.Webhook.Status.ValueString())
	assert.Equal(t, "subscription1", model.Webhook.Subscriptions.Elements()[0].(types.String).ValueString())
	assert.Equal(t, "secret-value", model.Webhook.SigningSecret.Attributes()["secret"].(types.String).ValueString())
	assert.Equal(t, "SHA256", model.Webhook.SigningSecret.Attributes()["algorithm"].(types.String).ValueString())

}

type AssertModel func(t *testing.T, model *WebhookModel)
