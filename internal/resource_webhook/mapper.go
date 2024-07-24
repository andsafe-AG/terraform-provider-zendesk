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

	webhookWithoutSensitive := showWebhookResponse.JSON200.Webhook
	diagnostics := putWebhookResponseBodyToStateModel(ctx, webhookWithoutSensitive, webhookState)
	if diagnostics != nil && diagnostics.HasError() {
		return diagnostics
	}

	return nil
}

func (m *WebhookMapper) PutCreateResponseToStateModel(ctx context.Context, response *zendesk_webhook_api.CreateOrCloneWebhookWrap, model *WebhookModel) diag.Diagnostics {

	webhookWithoutSensitive := response.JSON201.Webhook
	diagnostics := putWebhookResponseBodyToStateModel(ctx, webhookWithoutSensitive, model)
	if diagnostics != nil && diagnostics.HasError() {
		return diagnostics
	}
	return nil
}

func putWebhookResponseBodyToStateModel(ctx context.Context, webhookWithoutSensitive *zendesk_webhook_api.WebhookWithoutSensitive, webhookState *WebhookModel) diag.Diagnostics {
	webhookState.WebhookId = types.StringValue(*webhookWithoutSensitive.Id)
	authentication, diags := mapAuthentication(ctx, webhookWithoutSensitive)

	if diags.HasError() {
		tflog.Error(ctx, "Error reading webhook data from the API: ", map[string]interface{}{"error": diags})
		return diags
	}
	webhookState.Webhook.Authentication = authentication

	webhookState.Webhook.CreatedAt = stringValOrUnknown(webhookWithoutSensitive.CreatedAt)
	webhookState.Webhook.CreatedBy = stringValOrUnknown(webhookWithoutSensitive.CreatedBy)
	webhookState.Webhook.UpdatedAt = stringValOrUnknown(webhookWithoutSensitive.UpdatedAt)
	webhookState.Webhook.UpdatedBy = stringValOrUnknown(webhookWithoutSensitive.UpdatedBy)

	customHeaders, diagsMap := mapAMapOfString(webhookWithoutSensitive.CustomHeaders)
	if diagsMap.HasError() {
		tflog.Error(ctx, "Error reading webhook data from the API: ", map[string]interface{}{"error": diagsMap})
		return diagsMap
	}
	webhookState.Webhook.CustomHeaders = customHeaders
	webhookState.Webhook.Description = stringValOrUnknown(webhookWithoutSensitive.Description)
	webhookState.Webhook.Endpoint = stringValOrUnknown(webhookWithoutSensitive.Endpoint)
	externalSourceMapped, diagsExtSource := mapExternalSourceFromApiResponse(ctx, webhookWithoutSensitive)
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
	webhookState.Webhook.HttpMethod = stringValOrUnknown(webhookWithoutSensitive.HttpMethod)
	webhookState.Webhook.Id = stringValOrUnknown(webhookWithoutSensitive.Id)
	webhookState.Webhook.Name = stringValOrUnknown(webhookWithoutSensitive.Name)
	webhookState.Webhook.RequestFormat = stringValOrUnknown(webhookWithoutSensitive.RequestFormat)
	secret, diags := mapSecret(ctx, webhookWithoutSensitive.SigningSecret)
	if diags.HasError() {
		tflog.Error(ctx, "Error reading webhook data from the API: ", map[string]interface{}{"error": diags})
		return diags
	}

	webhookState.Webhook.SigningSecret = secret
	webhookState.Webhook.Status = stringValOrUnknown(webhookWithoutSensitive.Status)

	subscriptionList, diags := mapList(webhookWithoutSensitive.Subscriptions)
	if diags.HasError() {
		tflog.Error(ctx, "Error reading webhook data from the API: ", map[string]interface{}{"error": diags})
		return diags
	}
	webhookState.Webhook.Subscriptions = subscriptionList
	return nil
}

func (m *WebhookMapper) MapToRequestBody(ctx context.Context, model *WebhookModel) (*zendesk_webhook_api.CreateOrCloneWebhookJSONRequestBody, diag.Diagnostics) {
	request := zendesk_webhook_api.CreateOrCloneWebhookJSONRequestBody{}

	authentication := model.Webhook.Authentication
	if !authentication.IsNull() && !authentication.IsUnknown() {
		authenticationModel, diags := NewAuthenticationValue(authentication.AttributeTypes(ctx), authentication.Attributes())
		if diags.HasError() {
			tflog.Error(ctx, "Error reading webhook data from the API: ", map[string]interface{}{"error": diags})
			return nil, diags
		}

		request.Webhook.Authentication = &zendesk_webhook_api.Authentication{
			Data:        zendesk_webhook_api.AuthenticationData{},
			AddPosition: authenticationModel.AddPosition.ValueString(),
			Type:        authenticationModel.AuthenticationType.ValueString(),
		}

		if !authenticationModel.Data.IsNull() && !authenticationModel.Data.IsUnknown() {

			authDataModel, diags := NewDataValue(authenticationModel.Data.AttributeTypes(ctx), authenticationModel.Data.Attributes())

			if diags.HasError() {
				tflog.Error(ctx, "Error reading webhook data from the API: ", map[string]interface{}{"error": diags})
				return &request, diags
			}
			request.Webhook.Authentication.Data.Username = authDataModel.Username.ValueStringPointer()
			request.Webhook.Authentication.Data.Password = authDataModel.Password.ValueStringPointer()
			request.Webhook.Authentication.Data.Token = authDataModel.Token.ValueStringPointer()
		}
	}

	customHeaders := getCustomHeaders(model)
	request.Webhook.CustomHeaders = &customHeaders
	request.Webhook.Description = model.Webhook.Description.ValueStringPointer()
	request.Webhook.Endpoint = model.Webhook.Endpoint.ValueString()

	diagnostics := setExternalSource(ctx, model, &request)
	if diagnostics.HasError() {
		return &request, diagnostics
	}

	request.Webhook.HttpMethod = model.Webhook.HttpMethod.ValueString()
	request.Webhook.Name = model.Webhook.Name.ValueString()
	request.Webhook.RequestFormat = model.Webhook.RequestFormat.ValueString()
	setSigningSecret(model, &request)
	request.Webhook.Status = model.Webhook.Status.ValueString()
	subscriptions := getSubscriptions(model)
	request.Webhook.Subscriptions = &subscriptions

	return &request, nil
}

func setSigningSecret(model *WebhookModel, request *zendesk_webhook_api.CreateOrCloneWebhookJSONRequestBody) {
	if !model.Webhook.SigningSecret.IsNull() && !model.Webhook.SigningSecret.IsUnknown() {
		request.Webhook.SigningSecret = &zendesk_webhook_api.SigningSecret{}
		algorithm := model.Webhook.SigningSecret.Attributes()["algorithm"]
		if !algorithm.IsNull() {
			request.Webhook.SigningSecret.Algorithm = algorithm.(types.String).ValueStringPointer()
		}

		secret := model.Webhook.SigningSecret.Attributes()["secret"]
		if !secret.IsNull() {
			request.Webhook.SigningSecret.Secret = secret.(types.String).ValueStringPointer()
		}
	}
}

func setExternalSource(ctx context.Context, model *WebhookModel, request *zendesk_webhook_api.CreateOrCloneWebhookJSONRequestBody) diag.Diagnostics {
	if !model.Webhook.ExternalSource.IsNull() && !model.Webhook.ExternalSource.IsUnknown() {
		externalSourceModel, diags := NewExternalSourceValue(model.Webhook.ExternalSource.AttributeTypes(ctx), model.Webhook.ExternalSource.Attributes())
		if diags.HasError() {
			tflog.Error(ctx, "Error reading webhook data from the API: ", map[string]interface{}{"error": diags})
			return diags
		}

		request.Webhook.ExternalSource = &zendesk_webhook_api.ExternalSource{
			Data: &zendesk_webhook_api.ExternalSourceData{},
		}

		request.Webhook.ExternalSource.Type = externalSourceModel.ExternalSourceType.ValueStringPointer()

		if !externalSourceModel.ExternalSourceData.IsNull() && !externalSourceModel.ExternalSourceData.IsUnknown() {
			externalSourceDataModel, diags := NewExternalSourceDataValue(externalSourceModel.ExternalSourceData.AttributeTypes(ctx), externalSourceModel.ExternalSourceData.Attributes())
			if diags.HasError() {
				tflog.Error(ctx, "Error reading webhook data from the API: ", map[string]interface{}{"error": diags})
				return diags
			}
			request.Webhook.ExternalSource.Data.AppId = externalSourceDataModel.AppId.ValueStringPointer()
			request.Webhook.ExternalSource.Data.InstallationId = externalSourceDataModel.InstallationId.ValueStringPointer()
		}
	}
	return nil
}

func getSubscriptions(model *WebhookModel) []string {
	subscriptions := make([]string, 0)
	for _, value := range model.Webhook.Subscriptions.Elements() {
		subscriptions = append(subscriptions, value.(types.String).ValueString())
	}
	return subscriptions
}

func getCustomHeaders(model *WebhookModel) map[string]string {
	customHeaders := make(map[string]string)
	for key, value := range model.Webhook.CustomHeaders.Elements() {
		customHeaders[key] = value.(types.String).ValueString()
	}
	return customHeaders
}

func mapAuthentication(ctx context.Context, webhookWithoutSensitive *zendesk_webhook_api.WebhookWithoutSensitive) (basetypes.ObjectValue, diag.Diagnostics) {

	if webhookWithoutSensitive.Authentication == nil {
		return types.ObjectNull(AuthenticationValue{}.AttributeTypes(ctx)), nil
	}

	authMapped, diags := mapAuthenticationFromApiResponse(ctx, webhookWithoutSensitive)
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

func mapExternalSourceFromApiResponse(ctx context.Context, webhookWithoutSensitive *zendesk_webhook_api.WebhookWithoutSensitive) (ExternalSourceValue, diag.Diagnostics) {

	if webhookWithoutSensitive.ExternalSource == nil {
		return NewExternalSourceValueNull(), nil
	}
	externalSourceAttributes := make(map[string]attr.Value)
	externalSourceData, diags := mapExternalSourceData(webhookWithoutSensitive.ExternalSource.Data)
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
	externalSourceAttributes["type"] = stringValOrUnknown(webhookWithoutSensitive.ExternalSource.Type)
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

func mapAuthenticationFromApiResponse(ctx context.Context, webhookWithoutSensitive *zendesk_webhook_api.WebhookWithoutSensitive) (*AuthenticationValue, diag.Diagnostics) {

	if webhookWithoutSensitive.Authentication == nil {
		nullObject := NewAuthenticationValueNull()
		return &nullObject, nil
	}

	authDataMapped, diags := mapAuthData(webhookWithoutSensitive)
	if diags.HasError() {
		tflog.Error(ctx, "Error reading webhook data from the API: ", map[string]interface{}{"error": diags})
		return nil, diags
	}
	authData, diags := authDataMapped.ToObjectValue(ctx)
	auth := AuthenticationValue{
		AddPosition:        stringValOrUnknown(webhookWithoutSensitive.Authentication.AddPosition),
		Data:               authData,
		AuthenticationType: stringValOrUnknown(webhookWithoutSensitive.Authentication.Type),
	}

	return &auth, nil
}

func mapAuthData(webhookWithoutSensitive *zendesk_webhook_api.WebhookWithoutSensitive) (*DataValue, diag.Diagnostics) {

	if webhookWithoutSensitive.Authentication.Data == nil {
		nullObjectData := NewDataValueNull()
		return &nullObjectData, nil
	}

	username := webhookWithoutSensitive.Authentication.Data.Username

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
