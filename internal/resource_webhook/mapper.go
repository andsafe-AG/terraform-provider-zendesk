package resource_webhook

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"terraform-provider-zendesk/zendesk_webhook_api"
)

type WebhookMapper struct {
}

func NewWebhookMapper() *WebhookMapper {
	return &WebhookMapper{}
}

func (*WebhookMapper) PutWebhookShowResponseToStateModel(ctx context.Context, showWebhookResponse *zendesk_webhook_api.ShowWebhookWrap, webhookState *WebhookModel) diag.Diagnostics {

	webhookState.WebhookId = types.StringValue(*showWebhookResponse.JSON200.Webhook.Id)
	authentication, diags := mapAuthentication(ctx, showWebhookResponse)

	if diags.HasError() {
		tflog.Error(ctx, "Error reading webhook data from the API: ", map[string]interface{}{"error": diags})
		return diags
	}
	webhookState.Webhook.Authentication = authentication

	webhookState.Webhook.CreatedAt = stringValOrUnknown(showWebhookResponse.JSON200.Webhook.CreatedAt)
	webhookState.Webhook.CreatedBy = stringValOrUnknown(showWebhookResponse.JSON200.Webhook.CreatedBy)
	webhookState.Webhook.UpdatedAt = stringValOrUnknown(showWebhookResponse.JSON200.Webhook.UpdatedAt)
	webhookState.Webhook.UpdatedBy = stringValOrUnknown(showWebhookResponse.JSON200.Webhook.UpdatedBy)

	customHeaders, diagsMap := mapAMapOfString(showWebhookResponse.JSON200.Webhook.CustomHeaders)
	if diagsMap.HasError() {
		tflog.Error(ctx, "Error reading webhook data from the API: ", map[string]interface{}{"error": diagsMap})
		return diagsMap
	}
	webhookState.Webhook.CustomHeaders = customHeaders
	webhookState.Webhook.Description = stringValOrUnknown(showWebhookResponse.JSON200.Webhook.Description)
	webhookState.Webhook.Endpoint = stringValOrUnknown(showWebhookResponse.JSON200.Webhook.Endpoint)
	externalSourceMapped, diagsExtSource := mapExternalSourceFromApiResponse(ctx, showWebhookResponse)
	if diagsExtSource.HasError() {
		tflog.Error(ctx, "Error reading webhook data from the API: ", map[string]interface{}{"error": diagsExtSource})
		return diagsExtSource
	}

	externalSourceValue, diagnostics := externalSourceMapped.ToObjectValue(ctx)
	if diagnostics.HasError() {
		tflog.Error(ctx, "Error reading webhook data from the API: ", map[string]interface{}{"error": diagnostics})
		return diagnostics
	}

	webhookState.Webhook.ExternalSource = externalSourceValue
	webhookState.Webhook.HttpMethod = stringValOrUnknown(showWebhookResponse.JSON200.Webhook.HttpMethod)
	webhookState.Webhook.Id = stringValOrUnknown(showWebhookResponse.JSON200.Webhook.Id)
	webhookState.Webhook.Name = stringValOrUnknown(showWebhookResponse.JSON200.Webhook.Name)
	webhookState.Webhook.RequestFormat = stringValOrUnknown(showWebhookResponse.JSON200.Webhook.RequestFormat)
	secret, diags := mapSecret(ctx, showWebhookResponse.JSON200.Webhook.SigningSecret)
	if diags.HasError() {
		tflog.Error(ctx, "Error reading webhook data from the API: ", map[string]interface{}{"error": diags})
		return diags
	}

	webhookState.Webhook.SigningSecret = secret
	webhookState.Webhook.Status = stringValOrUnknown(showWebhookResponse.JSON200.Webhook.Status)

	subscriptionList, diags := mapList(showWebhookResponse.JSON200.Webhook.Subscriptions)
	if diags.HasError() {
		tflog.Error(ctx, "Error reading webhook data from the API: ", map[string]interface{}{"error": diags})
		return diags
	}
	webhookState.Webhook.Subscriptions = subscriptionList

	return nil
}

func mapAuthentication(ctx context.Context, showWebhookResponse *zendesk_webhook_api.ShowWebhookWrap) (basetypes.ObjectValue, diag.Diagnostics) {
	if showWebhookResponse.JSON200.Webhook.Authentication == nil {
		return types.ObjectNull(AuthenticationValue{}.AttributeTypes(ctx)), nil
	}

	authMapped, diags := mapAuthenticationFromApiResponse(ctx, *showWebhookResponse)
	if diags.HasError() {
		tflog.Error(ctx, "Error reading webhook data from the API: ", map[string]interface{}{"error": diags})
		return basetypes.ObjectValue{}, diags
	}

	authentication, diags := types.ObjectValueFrom(ctx, authMapped.AttributeTypes(ctx), authMapped)
	return authentication, diags
}

func mapList(listOfString *[]string) (basetypes.ListValue, diag.Diagnostics) {
	if listOfString == nil {
		return basetypes.NewListNull(types.StringType), nil
	}

	list := make([]attr.Value, len(*listOfString))
	for i, v := range *listOfString {
		list[i] = basetypes.NewStringValue(v)
	}

	return basetypes.NewListValue(types.StringType, list)
}

func mapSecret(ctx context.Context, secret *struct {
	Algorithm *string `json:"algorithm,omitempty"`
	Secret    *string `json:"secret,omitempty"`
}) (basetypes.ObjectValue, diag.Diagnostics) {
	if secret == nil {
		return basetypes.NewObjectNull(SigningSecretValue{}.AttributeTypes(context.Background())), nil

	}
	secretAttributes := make(map[string]attr.Value)
	secretAttributes["algorithm"] = stringValOrUnknown(secret.Algorithm)
	secretAttributes["secret"] = stringValOrUnknown(secret.Secret)
	secretValue, diags := NewSigningSecretValue(SigningSecretValue{}.AttributeTypes(context.Background()), secretAttributes)
	if diags.HasError() {
		return basetypes.ObjectValue{}, diags
	}

	value, diagnostics := secretValue.ToObjectValue(ctx)
	return value, diagnostics
}

func mapExternalSourceFromApiResponse(ctx context.Context, response *zendesk_webhook_api.ShowWebhookWrap) (ExternalSourceValue, diag.Diagnostics) {

	if response.JSON200.Webhook.ExternalSource == nil {
		return NewExternalSourceValueNull(), nil
	}
	externalSourceAttributes := make(map[string]attr.Value)
	externalSourceData, diags := mapExternalSourceData(response.JSON200.Webhook.ExternalSource.Data)
	if diags.HasError() {
		tflog.Error(ctx, "Error reading webhook data from the API: ", map[string]interface{}{"error": diags})
		return ExternalSourceValue{}, diags
	}
	externalSourceDataObject, diagnostics := externalSourceData.ToObjectValue(ctx)
	if diagnostics.HasError() {
		tflog.Error(ctx, "Error reading webhook data from the API: ", map[string]interface{}{"error": diagnostics})
		return ExternalSourceValue{}, diagnostics
	}
	externalSourceAttributes["external_source_data"] = externalSourceDataObject
	externalSourceAttributes["type"] = stringValOrUnknown(response.JSON200.Webhook.ExternalSource.Type)
	externalSourceVal, diags := NewExternalSourceValue(ExternalSourceValue{}.AttributeTypes(ctx), externalSourceAttributes)

	return externalSourceVal, diags

}

func mapExternalSourceData(data *struct {
	AppId          *string `json:"app_id,omitempty"`
	InstallationId *string `json:"installation_id,omitempty"`
}) (ExternalSourceDataValue, diag.Diagnostics) {

	externalSourceDataAttributes := make(map[string]attr.Value)
	externalSourceDataAttributes["app_id"] = stringValOrUnknown(data.AppId)
	externalSourceDataAttributes["installation_id"] = stringValOrUnknown(data.InstallationId)
	externalSourceData, diags := NewExternalSourceDataValue(ExternalSourceDataValue{}.AttributeTypes(context.Background()), externalSourceDataAttributes)
	return externalSourceData, diags
}

func mapAMapOfString(headers *map[string]string) (basetypes.MapValue, diag.Diagnostics) {
	if headers == nil {
		return basetypes.NewMapNull(types.StringType), nil
	}

	customHeaders := make(map[string]attr.Value)
	for key, value := range *headers {
		customHeaders[key] = basetypes.NewStringValue(value)
	}

	return basetypes.NewMapValue(types.StringType, customHeaders)
}

func mapAuthenticationFromApiResponse(ctx context.Context, wrap zendesk_webhook_api.ShowWebhookWrap) (*AuthenticationValue, diag.Diagnostics) {
	if wrap.JSON200.Webhook.Authentication == nil {
		nullObject := NewAuthenticationValueNull()
		return &nullObject, nil
	}

	authDataMapped, diags := mapAuthData(wrap)
	if diags.HasError() {
		tflog.Error(ctx, "Error reading webhook data from the API: ", map[string]interface{}{"error": diags})
		return nil, diags
	}
	authData, diags := authDataMapped.ToObjectValue(ctx)
	auth := AuthenticationValue{
		AddPosition:        stringValOrUnknown(wrap.JSON200.Webhook.Authentication.AddPosition),
		Data:               authData,
		AuthenticationType: stringValOrUnknown(wrap.JSON200.Webhook.Authentication.Type),
	}

	return &auth, nil
}

func mapAuthData(wrap zendesk_webhook_api.ShowWebhookWrap) (*DataValue, diag.Diagnostics) {

	if wrap.JSON200.Webhook.Authentication.Data == nil {
		nullObjectData := NewDataValueNull()
		return &nullObjectData, nil
	}

	username := wrap.JSON200.Webhook.Authentication.Data.Username

	authData := DataValue{
		Password: types.StringUnknown(),
		Token:    types.StringUnknown(),
		Username: stringValOrUnknown(username),
	}

	return &authData, nil
}

func stringValOrUnknown(value *string) basetypes.StringValue {
	if value == nil {
		return types.StringUnknown()
	}

	return types.StringValue(*value)
}
