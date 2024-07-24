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

func (m *WebhookMapper) MapToCreateRequestBody(ctx context.Context, model *WebhookModel) (*zendesk_webhook_api.CreateOrCloneWebhookJSONRequestBody, diag.Diagnostics) {

	webhookRequestBody := zendesk_webhook_api.WebhookWithSensitiveData{}
	request := zendesk_webhook_api.CreateOrCloneWebhookJSONRequestBody{Webhook: &webhookRequestBody}

	diagnostics := mapPlanModelToWebhookRequestBody(ctx, model, &webhookRequestBody)
	if diagnostics != nil && diagnostics.HasError() {
		return nil, diagnostics
	}

	return &request, nil
}

func (m *WebhookMapper) MapToUpdateRequestBody(ctx context.Context, r *WebhookModel) (zendesk_webhook_api.UpdateWebhookJSONRequestBody, diag.Diagnostics) {
	webhookRequestBody := zendesk_webhook_api.WebhookWithSensitiveData{}
	request := zendesk_webhook_api.UpdateWebhookJSONRequestBody{Webhook: &webhookRequestBody}
	diags := mapPlanModelToWebhookRequestBody(ctx, r, request.Webhook)
	if diags != nil && diags.HasError() {
		tflog.Error(ctx, "Error reading webhook data from the API: ", map[string]interface{}{"error": diags})
		return request, diags
	}
	return request, nil
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

func mapPlanModelToWebhookRequestBody(ctx context.Context, model *WebhookModel, webhookRequestBody *zendesk_webhook_api.WebhookWithSensitiveData) diag.Diagnostics {
	authentication := model.Webhook.Authentication
	if !authentication.IsNull() && !authentication.IsUnknown() {
		authenticationModel, diags := NewAuthenticationValue(authentication.AttributeTypes(ctx), authentication.Attributes())
		if diags.HasError() {
			tflog.Error(ctx, "Error reading webhook data from the API: ", map[string]interface{}{"error": diags})
			return diags
		}

		webhookRequestBody.Authentication = &zendesk_webhook_api.Authentication{
			Data:        &zendesk_webhook_api.AuthenticationData{},
			AddPosition: authenticationModel.AddPosition.ValueString(),
			Type:        authenticationModel.AuthenticationType.ValueString(),
		}

		if !authenticationModel.Data.IsNull() && !authenticationModel.Data.IsUnknown() {

			authDataModel, diags := NewDataValue(authenticationModel.Data.AttributeTypes(ctx), authenticationModel.Data.Attributes())

			if diags.HasError() {
				tflog.Error(ctx, "Error reading webhook data from the API: ", map[string]interface{}{"error": diags})
				return diags
			}
			webhookRequestBody.Authentication.Data.Username = authDataModel.Username.ValueStringPointer()
			webhookRequestBody.Authentication.Data.Password = authDataModel.Password.ValueStringPointer()
			webhookRequestBody.Authentication.Data.Token = authDataModel.Token.ValueStringPointer()
		}
	}

	customHeaders := getCustomHeaders(ctx, model)
	if customHeaders != nil {
		webhookRequestBody.CustomHeaders = &customHeaders
	}
	webhookRequestBody.Description = model.Webhook.Description.ValueStringPointer()
	webhookRequestBody.Endpoint = model.Webhook.Endpoint.ValueString()

	diagnostics := setExternalSource(ctx, model, webhookRequestBody)
	if diagnostics.HasError() {
		return diagnostics
	}

	webhookRequestBody.HttpMethod = model.Webhook.HttpMethod.ValueString()
	webhookRequestBody.Name = model.Webhook.Name.ValueString()
	webhookRequestBody.RequestFormat = model.Webhook.RequestFormat.ValueString()
	setSigningSecret(model, webhookRequestBody)
	webhookRequestBody.Status = model.Webhook.Status.ValueString()
	subscriptions := getSubscriptions(ctx, model)
	if subscriptions != nil {
		webhookRequestBody.Subscriptions = &subscriptions
	}
	return nil
}

func setSigningSecret(model *WebhookModel, webhookWithSensitiveData *zendesk_webhook_api.WebhookWithSensitiveData) {
	if !model.Webhook.SigningSecret.IsNull() && !model.Webhook.SigningSecret.IsUnknown() {
		webhookWithSensitiveData.SigningSecret = &zendesk_webhook_api.SigningSecret{}
		algorithm := model.Webhook.SigningSecret.Attributes()["algorithm"]
		if !algorithm.IsNull() {
			algorithmString, ok := algorithm.(types.String)
			if !ok {
				panic("unexpected type of algorithm")
			}
			webhookWithSensitiveData.SigningSecret.Algorithm = algorithmString.ValueStringPointer()
		}

		secret := model.Webhook.SigningSecret.Attributes()["secret"]
		if !secret.IsNull() {
			secretString, ok := secret.(types.String)
			if !ok {
				panic("unexpected type of secret")
			}
			webhookWithSensitiveData.SigningSecret.Secret = secretString.ValueStringPointer()
		}
	}
}

func setExternalSource(ctx context.Context, model *WebhookModel, webhookWithSensitiveData *zendesk_webhook_api.WebhookWithSensitiveData) diag.Diagnostics {
	webhookModelBody := model.Webhook
	if !webhookModelBody.ExternalSource.IsNull() && !webhookModelBody.ExternalSource.IsUnknown() {
		externalSourceModel, diags := NewExternalSourceValue(webhookModelBody.ExternalSource.AttributeTypes(ctx), webhookModelBody.ExternalSource.Attributes())
		if diags.HasError() {
			tflog.Error(ctx, "Error reading webhook data from the API: ", map[string]interface{}{"error": diags})
			return diags
		}

		webhookWithSensitiveData.ExternalSource = &zendesk_webhook_api.ExternalSource{
			Data: &zendesk_webhook_api.ExternalSourceData{},
		}

		webhookWithSensitiveData.ExternalSource.Type = externalSourceModel.ExternalSourceType.ValueStringPointer()

		if !externalSourceModel.ExternalSourceData.IsNull() && !externalSourceModel.ExternalSourceData.IsUnknown() {
			externalSourceDataModel, diags := NewExternalSourceDataValue(externalSourceModel.ExternalSourceData.AttributeTypes(ctx), externalSourceModel.ExternalSourceData.Attributes())
			if diags.HasError() {
				tflog.Error(ctx, "Error reading webhook data from the API: ", map[string]interface{}{"error": diags})
				return diags
			}
			webhookWithSensitiveData.ExternalSource.Data.AppId = externalSourceDataModel.AppId.ValueStringPointer()
			webhookWithSensitiveData.ExternalSource.Data.InstallationId = externalSourceDataModel.InstallationId.ValueStringPointer()
		}
	}
	return nil
}

func getSubscriptions(ctx context.Context, model *WebhookModel) []string {
	subscriptions := make([]string, 0)
	if model.Webhook.Subscriptions.IsNull() || model.Webhook.Subscriptions.IsUnknown() {
		return nil
	}
	model.Webhook.Subscriptions.ElementsAs(ctx, &subscriptions, true)
	return subscriptions
}

func getCustomHeaders(ctx context.Context, model *WebhookModel) map[string]string {
	customHeaders := make(map[string]string)
	if model.Webhook.CustomHeaders.IsNull() || model.Webhook.CustomHeaders.IsUnknown() {
		return nil
	}
	model.Webhook.CustomHeaders.ElementsAs(ctx, &customHeaders, true)
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

	authDataMapped := mapAuthData(webhookWithoutSensitive)

	authData, diags := authDataMapped.ToObjectValue(ctx)
	if diags.HasError() {
		tflog.Error(ctx, "Error reading webhook data from the API: ", map[string]interface{}{"error": diags})
		return nil, diags
	}
	auth := AuthenticationValue{
		AddPosition:        stringValOrUnknown(webhookWithoutSensitive.Authentication.AddPosition),
		Data:               authData,
		AuthenticationType: stringValOrUnknown(webhookWithoutSensitive.Authentication.Type),
	}

	return &auth, nil
}

func mapAuthData(webhookWithoutSensitive *zendesk_webhook_api.WebhookWithoutSensitive) *DataValue {

	if webhookWithoutSensitive.Authentication.Data == nil {
		nullObjectData := NewDataValueNull()
		return &nullObjectData
	}

	username := webhookWithoutSensitive.Authentication.Data.Username

	authData := DataValue{
		Password: types.StringUnknown(),
		Token:    types.StringUnknown(),
		Username: stringValOrUnknown(username),
	}

	return &authData
}

func stringValOrUnknown(value *string) basetypes.StringValue {
	if value == nil {
		return types.StringUnknown()
	}

	return types.StringValue(*value)
}
