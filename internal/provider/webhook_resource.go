package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"strconv"
	"terraform-provider-zendesk/internal/resource_webhook"
	"terraform-provider-zendesk/zendesk_webhook_api"

	"github.com/hashicorp/terraform-plugin-framework/resource"
)

var _ resource.Resource = (*webhookResource)(nil)

func NewWebhookResource() resource.Resource {
	return &webhookResource{}
}

type webhookResource struct {
	client *zendesk_webhook_api.WebhookApi
}

func (r *webhookResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_webhook"
}

// Configure adds the provider configured client to the resource.
func (r *webhookResource) Configure(ctx context.Context, request resource.ConfigureRequest, response *resource.ConfigureResponse) {
	// Add a nil check when handling ProviderData because Terraform
	// sets that data after it calls the ConfigureProvider RPC.
	if request.ProviderData == nil {
		return
	}

	providerData, ok := request.ProviderData.(zendeskProviderData)

	if !ok {
		response.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *zendeskProviderData, got: %T. Please report this issue to the provider developers.",
				request.ProviderData),
		)
		return
	}

	r.client = providerData.webhookApi

}

func (r *webhookResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = resource_webhook.WebhookResourceSchema(ctx)
}

// ImportState imports a Webhook by a given id, when the id value is an integer, or by name=id otherwise /*
func (r *webhookResource) ImportState(ctx context.Context, request resource.ImportStateRequest, response *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	tflog.Debug(ctx, "Called ImportState of a webhook resource")
	id := request.ID
	tflog.Debug(ctx, "request.id: "+id)

	webhookIdPath1 := path.Root("webhook_id")
	webhookIdPath2 := path.Root("webhook").AtName("id")
	var webhookIdOfMatching string

	reqEditors := make([]zendesk_webhook_api.RequestEditorFn, 0)
	webhookShowResponse, err := r.client.GetClient().ShowWebhookWithResponse(ctx, id, reqEditors...)
	if err != nil {
		tflog.Error(ctx, "Error reading webhook data from the API: ", map[string]interface{}{"error": err})
		response.Diagnostics.AddError("Error reading webhook data from the API", err.Error())
		return
	}

	if webhookShowResponse.StatusCode() == 404 {

		reqEditors := make([]zendesk_webhook_api.RequestEditorFn, 0)
		webhooks, err := r.client.GetClient().ListWebhooksWithResponse(ctx, nil, reqEditors...)
		if err != nil {
			tflog.Error(ctx, "Error reading webhooks list from the API: ", map[string]interface{}{"error": err})
			response.Diagnostics.AddError("Error reading webhooks list from the API", err.Error())
			return
		}

		if webhooks.StatusCode() != 200 {
			tflog.Error(ctx, "Error reading webhooks list from the API: ", map[string]interface{}{"error": string(webhooks.Body)})
			response.Diagnostics.AddError("Error reading webhooks list from the API", fmt.Sprintf("%+v", string(webhooks.Body)))
			return
		}

		list := webhooks.JSON200.Webhooks

		for _, webhook := range *list {
			if webhook.Name == &id {
				webhookIdOfMatching = *webhook.Id
				break
			}
		}

	} else if webhookShowResponse.StatusCode() != 200 {
		tflog.Error(ctx, "Error reading webhook data from the API: ", map[string]interface{}{"error": string(webhookShowResponse.Body)})
		response.Diagnostics.AddError("Error reading webhook data from the API", fmt.Sprintf("%+v", string(webhookShowResponse.Body)))
		return
	} else {
		webhookIdOfMatching = *webhookShowResponse.JSON200.Webhook.Id
	}

	response.Diagnostics.Append(response.State.SetAttribute(ctx, webhookIdPath1, webhookIdOfMatching)...)
	response.Diagnostics.Append(response.State.SetAttribute(ctx, webhookIdPath2, webhookIdOfMatching)...)

	if response.Diagnostics.HasError() {

		tflog.Error(ctx, "Error importing state: ", map[string]any{"error": response.Diagnostics.Errors()})
		return
	}

	tflog.Info(ctx, "ImportState custom status completed successfully")
}

func (r *webhookResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var planModel resource_webhook.WebhookModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &planModel)...)

	if resp.Diagnostics.HasError() {
		tflog.Error(ctx, "Error reading Terraform plan data into the model: ", map[string]any{"error": resp.Diagnostics.Errors()})
		return
	}

	tflog.Debug(ctx, "Create webhook resource with model: "+fmt.Sprintf("%+v", &planModel))

	webhookMapper := resource_webhook.NewWebhookMapper()
	requestBody, diags := webhookMapper.MapToCreateRequestBody(ctx, &planModel)
	if diags != nil && diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}
	tflog.Debug(ctx, "Create webhook resource with request body: "+fmt.Sprintf("%+v", requestBody))
	reqEditors := make([]zendesk_webhook_api.RequestEditorFn, 0)
	createResponse, err := r.client.GetClient().CreateOrCloneWebhookWithResponse(ctx, nil, *requestBody, reqEditors...)
	if err != nil {
		tflog.Error(ctx, "Error creating webhook data from the API: ", map[string]interface{}{"error": err})
		resp.Diagnostics.AddError("Error creating webhook data from the API", err.Error())
		return
	}
	tflog.Debug(ctx, fmt.Sprintf("Received API response: Status %v, Body: %+v", createResponse.StatusCode(), string(createResponse.Body)))
	if createResponse.StatusCode() != 201 {
		tflog.Error(ctx, "Error creating webhook data from the API: ", map[string]interface{}{"error": string(createResponse.Body)})
		body, _ := json.Marshal(requestBody)
		resp.Diagnostics.AddError("Error creating webhook data from the API", fmt.Sprintf("Code: %v Body: %+v Request: %v", createResponse.StatusCode(), string(createResponse.Body), string(body)))
		return
	}
	if createResponse.JSON201.Webhook.SigningSecret == nil {
		createResponse.JSON201.Webhook.SigningSecret = &struct {
			Algorithm *string `json:"algorithm,omitempty"`
			Secret    *string `json:"secret,omitempty"`
		}{
			Algorithm: nil,
			Secret:    nil,
		}
	}
	signingSecret := createResponse.JSON201.Webhook.SigningSecret

	id := *createResponse.JSON201.Webhook.Id
	fetchSigningSecret(ctx, r.client, id, signingSecret)

	diags2 := webhookMapper.UpdateAttributesWithCreateResponse(ctx, createResponse, &planModel)
	if diags2 != nil && diags2.HasError() {
		tflog.Error(ctx, "Error updating webhook model with create response: ", map[string]any{"error": diags2.Errors()})
		resp.Diagnostics.Append(diags2...)
		return
	}

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &planModel)...)
}

func fetchSigningSecret(ctx context.Context, client *zendesk_webhook_api.WebhookApi, id string, signingSecret *struct {
	Algorithm *string `json:"algorithm,omitempty"`
	Secret    *string `json:"secret,omitempty"`
}) {
	reqEditors := make([]zendesk_webhook_api.RequestEditorFn, 0)
	responseSecret, errSecret := client.GetClient().ShowWebhookSigningSecretWithResponse(ctx, id, reqEditors...)
	if errSecret != nil {
		tflog.Error(ctx, "Error reading webhook signing secret from the API: ", map[string]interface{}{"error": errSecret})
	} else if responseSecret != nil && responseSecret.StatusCode() == 200 {
		signingSecret.Secret = responseSecret.JSON200.SigningSecret.Secret
		signingSecret.Algorithm = responseSecret.JSON200.SigningSecret.Algorithm
	}
}

func (r *webhookResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state resource_webhook.WebhookModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		logErrors(ctx, resp)
		return
	}

	tflog.Debug(ctx, "Read webhook resource with id: "+fmt.Sprintf("%+v", state.WebhookId))
	webhookId := state.WebhookId.ValueString()

	// Get Webhook data from the API
	reqEditors := make([]zendesk_webhook_api.RequestEditorFn, 0)
	webhookResponse, err := r.client.GetClient().ShowWebhookWithResponse(ctx, webhookId, reqEditors...)

	if err != nil {
		tflog.Error(ctx, "Error reading webhook data from the API: ", map[string]interface{}{"error": err})
		resp.Diagnostics.AddError("Error reading webhook data from the API", err.Error())
		return
	}

	if webhookResponse.StatusCode() != 200 {
		message := fmt.Sprintf("Error reading webhook data from the API, StatusCode: %v , Request Param: %v", webhookResponse.StatusCode(), webhookId)
		tflog.Error(ctx, message,
			map[string]interface{}{"error": string(webhookResponse.Body)})
		resp.Diagnostics.AddError("Error reading webhook data from the API", fmt.Sprintf("%+v", string(webhookResponse.Body)))
		return
	}

	if webhookResponse.JSON200.Webhook.SigningSecret == nil {
		webhookResponse.JSON200.Webhook.SigningSecret = &struct {
			Algorithm *string `json:"algorithm,omitempty"`
			Secret    *string `json:"secret,omitempty"`
		}{
			Algorithm: nil,
			Secret:    nil,
		}
	}
	fetchSigningSecret(ctx, r.client, webhookId, webhookResponse.JSON200.Webhook.SigningSecret)
	webhookMapper := resource_webhook.NewWebhookMapper()

	webhookMapper.PutWebhookShowResponseToStateModel(ctx, webhookResponse, &state)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func logErrors(ctx context.Context, resp *resource.ReadResponse) {
	tflog.Error(ctx, "Error reading Terraform prior state data into the model: ", map[string]any{"error": resp.Diagnostics.Errors()})
}

func (r *webhookResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var planModel resource_webhook.WebhookModel

	tflog.Debug(ctx, "Update webhook resource")
	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &planModel)...)

	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Update webhook resource with model: %+v", &planModel))

	// Read Terraform state data into the model
	var stateModel resource_webhook.WebhookModel
	resp.Diagnostics.Append(req.State.Get(ctx, &stateModel)...)

	tflog.Debug(ctx, fmt.Sprintf("Update webhook resource with previous state: %+v", &stateModel))

	if resp.Diagnostics.HasError() {
		return
	}

	webhookMapper := resource_webhook.NewWebhookMapper()

	mappedRequestBody, diags := webhookMapper.MapToUpdateRequestBody(ctx, &planModel)

	if diags != nil && diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}
	tflog.Debug(ctx, fmt.Sprintf("Update webhook resource with mapped request body: %+v", &mappedRequestBody))

	reqEditors := make([]zendesk_webhook_api.RequestEditorFn, 0)
	tflog.Debug(ctx, fmt.Sprintf("Will update webhook with id %v with request body: %+v and reqEditors %v", stateModel.WebhookId.ValueString(), &mappedRequestBody, &reqEditors))
	response, err := r.client.GetClient().UpdateWebhookWithResponse(ctx, stateModel.WebhookId.ValueString(), mappedRequestBody, reqEditors...)

	if err != nil {
		resp.Diagnostics.AddError("Error updating webhook data from the API", err.Error())
		return
	}

	if response.StatusCode() == 404 {

		resp.Diagnostics.AddError("Error updating webhook data", "Webhook not found: "+string(response.Body))
		return

	}

	if response.StatusCode() == 400 {
		responseErrors, err := json.Marshal(response.JSON400.Errors)
		if err != nil {
			resp.Diagnostics.AddError("Error updating webhook data from the API", "Bad Request: "+err.Error())
			return
		}
		resp.Diagnostics.AddError("Error updating webhook data from the API", "Bad Request: "+string(responseErrors))
		return

	}
	if response.StatusCode() != 204 {
		detail := "Unexpected response status code: " + strconv.Itoa(response.StatusCode()) + ", Response Body: " + fmt.Sprintf("%+v", string(response.Body))
		resp.Diagnostics.AddError("Error updating webhook data from the API", detail)
		return
	}
	reqEditors2 := make([]zendesk_webhook_api.RequestEditorFn, 0)
	tflog.Debug(ctx, fmt.Sprintf("Will get webhook with id %v with reqEditors %v", stateModel.WebhookId.ValueString(), reqEditors2))
	showResponse, err := r.client.GetClient().ShowWebhookWithResponse(ctx, stateModel.WebhookId.ValueString(), reqEditors2...)

	if err != nil {
		resp.Diagnostics.AddError("Error reading webhook data from the API after successful update", err.Error())
		return
	}

	if showResponse.StatusCode() != 200 {
		resp.Diagnostics.AddError(fmt.Sprintf("Error reading webhook data from the API after successful update with status %v", showResponse.StatusCode()), fmt.Sprintf("%+v", string(showResponse.Body)))
		return
	}

	if showResponse.JSON200.Webhook.SigningSecret == nil {
		showResponse.JSON200.Webhook.SigningSecret = &struct {
			Algorithm *string `json:"algorithm,omitempty"`
			Secret    *string `json:"secret,omitempty"`
		}{
			Algorithm: nil,
			Secret:    nil,
		}
	}
	fetchSigningSecret(ctx, r.client, stateModel.WebhookId.ValueString(), showResponse.JSON200.Webhook.SigningSecret)

	diagnostics := webhookMapper.PutWebhookShowResponseAfterUpdateToStateModel(ctx, showResponse, &planModel)
	if diagnostics != nil && diagnostics.HasError() {
		tflog.Error(ctx, "Error updating webhook model with show response after successful update: ", map[string]any{"error": diagnostics.Errors()})
		resp.Diagnostics.Append(diagnostics...)
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &planModel)...)
}

func (r *webhookResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data resource_webhook.WebhookModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	reqEditors := make([]zendesk_webhook_api.RequestEditorFn, 0)
	response, err := r.client.GetClient().DeleteWebhookWithResponse(ctx, data.WebhookId.ValueString(), reqEditors...)
	if err != nil {
		resp.Diagnostics.AddError("Error deleting webhook data from the API", err.Error())
		return
	}
	if response.StatusCode() == 404 {
		resp.Diagnostics.AddError("Error deleting webhook data from the API", "Webhook not found")
		return
	}

	if response.StatusCode() == 400 {
		responseErrors, err := json.Marshal(response.JSON400.Errors)
		if err != nil {
			resp.Diagnostics.AddError("Error deleting webhook data from the API", "Bad Request: "+err.Error())
			return
		}
		resp.Diagnostics.AddError("Error deleting webhook data from the API", "Bad Request: "+string(responseErrors))
		return

	}

	if response.StatusCode() != 204 {
		detail := "Unexpected response status code: " + strconv.Itoa(response.StatusCode()) + ", Response Body: " + fmt.Sprintf("%+v", string(response.Body))
		resp.Diagnostics.AddError("Error deleting webhook data from the API", detail)
		return
	}

	tflog.Info(ctx, fmt.Sprintf("Delete webhook with id %v completed successfully", data.WebhookId.ValueString()))

}
