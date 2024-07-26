package resource_webhook

import (
	"context"
	"fmt"
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

func (*WebhookMapper) PutWebhookShowResponseAfterUpdateToStateModel(ctx context.Context, showWebhookResponse *zendesk_webhook_api.ShowWebhookWrap, webhookState *WebhookModel) diag.Diagnostics {

	webhookWithoutSensitive := showWebhookResponse.JSON200.Webhook

	diagnostics := putWebhookCreateResponseBodyToStateModel(ctx, webhookWithoutSensitive, webhookState)

	if diagnostics != nil && diagnostics.HasError() {
		return diagnostics
	}

	return nil
}

func (m *WebhookMapper) UpdateAttributesWithCreateResponse(ctx context.Context, response *zendesk_webhook_api.CreateOrCloneWebhookWrap, model *WebhookModel) diag.Diagnostics {

	webhookWithoutSensitive := response.JSON201.Webhook
	return putWebhookCreateResponseBodyToStateModel(ctx, webhookWithoutSensitive, model)

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
	var diags diag.Diagnostics

	webhookState.Webhook.Id = stringValOrNull(webhookWithoutSensitive.Id)

	diags2 := setAuthenticationFromResponseToState(ctx, webhookWithoutSensitive, webhookState)

	if diags2 != nil && diags2.HasError() {
		tflog.Error(ctx, "Error setting webhook authentication data from the API response to the TF state: ", map[string]interface{}{"error": diags2})
		return diags2

	}
	webhookState.Webhook.CreatedAt = stringValOrNull(webhookWithoutSensitive.CreatedAt)
	webhookState.Webhook.CreatedBy = stringValOrNull(webhookWithoutSensitive.CreatedBy)
	webhookState.Webhook.UpdatedAt = stringValOrNull(webhookWithoutSensitive.UpdatedAt)
	webhookState.Webhook.UpdatedBy = stringValOrNull(webhookWithoutSensitive.UpdatedBy)

	customHeaders, diagsMap := mapAMapOfString(webhookWithoutSensitive.CustomHeaders)
	if diagsMap.HasError() {
		tflog.Error(ctx, "Error reading webhook data from the API: ", map[string]interface{}{"error": diagsMap})
		return diagsMap
	}
	webhookState.Webhook.CustomHeaders = customHeaders
	webhookState.Webhook.Description = stringValOrNull(webhookWithoutSensitive.Description)
	webhookState.Webhook.Endpoint = stringValOrNull(webhookWithoutSensitive.Endpoint)
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
	webhookState.Webhook.HttpMethod = stringValOrNull(webhookWithoutSensitive.HttpMethod)
	webhookState.Webhook.Name = stringValOrNull(webhookWithoutSensitive.Name)
	webhookState.Webhook.RequestFormat = stringValOrNull(webhookWithoutSensitive.RequestFormat)
	secret, diags := mapSecret(ctx, webhookWithoutSensitive.SigningSecret)
	if diags.HasError() {
		tflog.Error(ctx, "Error reading webhook data from the API: ", map[string]interface{}{"error": diags})
		return diags
	}

	webhookState.Webhook.SigningSecret = secret
	webhookState.Webhook.Status = stringValOrNull(webhookWithoutSensitive.Status)

	subscriptionList, diags := mapList(webhookWithoutSensitive.Subscriptions)
	if diags.HasError() {
		tflog.Error(ctx, "Error reading webhook data from the API: ", map[string]interface{}{"error": diags})
		return diags
	}
	webhookState.Webhook.Subscriptions = subscriptionList
	return nil
}

func putWebhookCreateResponseBodyToStateModel(ctx context.Context, webhookWithoutSensitive *zendesk_webhook_api.WebhookWithoutSensitive, webhookPlan *WebhookModel) diag.Diagnostics {
	webhookPlan.WebhookId = types.StringValue(*webhookWithoutSensitive.Id)
	webhookPlan.Webhook.Id = stringValOrNull(webhookWithoutSensitive.Id)

	webhookPlan.Webhook.CreatedAt = stringValOrNull(webhookWithoutSensitive.CreatedAt)
	webhookPlan.Webhook.CreatedBy = stringValOrNull(webhookWithoutSensitive.CreatedBy)
	webhookPlan.Webhook.UpdatedAt = stringValOrNull(webhookWithoutSensitive.UpdatedAt)
	webhookPlan.Webhook.UpdatedBy = stringValOrNull(webhookWithoutSensitive.UpdatedBy)

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

	webhookPlan.Webhook.ExternalSource = externalSourceValue

	secret, diags := mapSecret(ctx, webhookWithoutSensitive.SigningSecret)
	if diags.HasError() {
		tflog.Error(ctx, "Error reading webhook data from the API: ", map[string]interface{}{"error": diags})
		return diags
	}

	webhookPlan.Webhook.SigningSecret = secret

	return nil
}

func setAuthenticationFromResponseToState(ctx context.Context, webhookWithoutSensitive *zendesk_webhook_api.WebhookWithoutSensitive, webhookState *WebhookModel) diag.Diagnostics {

	if webhookWithoutSensitive.Authentication == nil {
		// No Authentication
		nullObject := NewAuthenticationValueNull()
		nullValue, diagnostics := nullObject.ToObjectValue(ctx)
		if diagnostics.HasError() {
			tflog.Error(ctx, "Error reading authentication value from the API: ", map[string]interface{}{"error": diagnostics.Errors()})
			return diagnostics
		}
		webhookState.Webhook.Authentication = nullValue
		return nil
	}
	if webhookState.Webhook.Authentication.IsNull() {
		var diagnostics diag.Diagnostics
		webhookState.Webhook.Authentication, diagnostics = NewAuthenticationValueNull().ToObjectValue(ctx)
		if diagnostics.HasError() {
			tflog.Error(ctx, "Error creating Authentication for the webhook state: ", map[string]interface{}{"error": diagnostics.Errors()})
			return diagnostics
		}
	}

	authenticationValueOld, diags := AuthenticationType.ValueFromObject(AuthenticationType{}, ctx, webhookState.Webhook.Authentication)

	if diags.HasError() {
		tflog.Error(ctx, "Error getting ValueFromObject on Authentication from Webhook State ", map[string]interface{}{"error": diags})
		return diags
	}

	authenticationValueOldCast, ok := authenticationValueOld.(AuthenticationValue)
	if !ok {
		return diag.Diagnostics{diag.NewErrorDiagnostic("Error when casting old state Authentication to AuthenticationValue", "error")}
	}

	authDataMapped, diags := mergeAuthData(ctx, webhookWithoutSensitive, authenticationValueOldCast.Data)
	if diags.HasError() {
		tflog.Error(ctx, "Error reading authentication data from the API: ", map[string]interface{}{"error": diags})
		return diags
	}
	tflog.Debug(ctx, fmt.Sprintf("Mapped Auth Data: %v", authDataMapped))

	authData, diags := authDataMapped.ToObjectValue(ctx)
	if diags.HasError() {
		tflog.Error(ctx, "Error reading authentication data to ObjectValue: ", map[string]interface{}{"error": diags.Errors()})
		return diags
	}

	authenticationValueOldCast.AddPosition = stringValOrNull(webhookWithoutSensitive.Authentication.AddPosition)
	authenticationValueOldCast.Data = authData
	authenticationValueOldCast.AuthenticationType = stringValOrNull(webhookWithoutSensitive.Authentication.Type)

	authentication, diags := authenticationValueOldCast.ToObjectValue(ctx)

	if diags.HasError() {
		tflog.Error(ctx, "Error reading webhook data from the API: ", map[string]interface{}{"error": diags})
		return diags
	}

	webhookState.Webhook.Authentication = authentication

	return diags
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

func mergeAuthData(ctx context.Context, webhookWithoutSensitive *zendesk_webhook_api.WebhookWithoutSensitive, dataOld basetypes.ObjectValue) (*DataValue, diag.Diagnostics) {
	var dataOldObject DataValue
	if dataOld.IsNull() {
		dataOldObject = NewDataValueNull()
	} else {

		dataValueOld, diags := DataType.ValueFromObject(DataType{}, ctx, dataOld)
		if diags.HasError() {
			return nil, diags
		}

		dataValueOldCast, ok := dataValueOld.(DataValue)
		if !ok {
			return nil, diag.Diagnostics{diag.NewErrorDiagnostic("Error when casting old state Data to DataValue", "error")}
		}
		dataOldObject = dataValueOldCast
	}

	oldPassword := dataOldObject.Password
	oldToken := dataOldObject.Token

	if webhookWithoutSensitive.Authentication.Data == nil {
		return &dataOldObject, nil
	}

	newUsername := stringValOrNull(webhookWithoutSensitive.Authentication.Data.Username)
	newPassword := oldPassword
	// if username was removed, then the basic auth is not possible, so password must have been removed as well
	if newUsername.IsNull() {
		newPassword = basetypes.NewStringNull()
		if oldToken.IsNull() {
			newToken := basetypes.NewStringNull()
			return &DataValue{
				Password: newPassword,
				Token:    newToken,
				Username: newUsername,
			}, nil

		}
		return &DataValue{
			Token: oldToken,
		}, nil
	}

	return &DataValue{
		Password: newPassword,
		Token:    types.StringNull(),
		Username: newUsername,
	}, nil
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
	secretAttributes["algorithm"] = stringValOrNull(secret.Algorithm)
	secretAttributes["secret"] = stringValOrNull(secret.Secret)
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
	externalSourceAttributes["type"] = stringValOrNull(webhookWithoutSensitive.ExternalSource.Type)
	externalSourceVal, diags := NewExternalSourceValue(ExternalSourceValue{}.AttributeTypes(ctx), externalSourceAttributes)

	return externalSourceVal, diags

}

func mapExternalSourceData(data *struct {
	AppId          *string `json:"app_id,omitempty"`
	InstallationId *string `json:"installation_id,omitempty"`
}) (ExternalSourceDataValue, diag.Diagnostics) {

	externalSourceDataAttributes := make(map[string]attr.Value)
	externalSourceDataAttributes["app_id"] = stringValOrNull(data.AppId)
	externalSourceDataAttributes["installation_id"] = stringValOrNull(data.InstallationId)
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

func stringValOrNull(value *string) basetypes.StringValue {
	if value == nil {
		return types.StringNull()
	}

	return types.StringValue(*value)
}
