package provider

import (
	"context"
	"fmt"
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

/*type webhookResourceModel struct {
	Id types.String `tfsdk:"id"`
}*/

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

	providerData, ok := request.ProviderData.(*zendeskProviderData)

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

	idInt, intParsingError := strconv.ParseInt(id, 10, 64)
	if intParsingError != nil {
		// TODO - Add logic to handle the case where the id is not an integer
	}
	if idInt == 0 {
		tflog.Error(ctx, "Could not find webhook with id or name: "+id)
		response.Diagnostics.AddError("Could not find webhook with id or name: "+id, "")
		return
	}

	//TODO - Add logic to retrieve the webhook by id

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

	tflog.Debug(ctx, "Create webhook resource with id: "+fmt.Sprintf("%+v", planModel.WebhookId))
	webhookMapper := resource_webhook.NewWebhookMapper()
	requestBody, diags := webhookMapper.MapToRequestBody(ctx, &planModel)
	if diags != nil && diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	createResponse, err := r.client.GetClient().CreateOrCloneWebhookWithResponse(ctx, nil, *requestBody, nil)
	if err != nil {
		tflog.Error(ctx, "Error creating webhook data from the API: ", map[string]interface{}{"error": err})
		resp.Diagnostics.AddError("Error creating webhook data from the API", err.Error())
		return
	}

	diags2 := webhookMapper.PutCreateResponseToStateModel(ctx, createResponse, &planModel)
	if diags2 != nil && diags2.HasError() {
		resp.Diagnostics.Append(diags2...)
		return
	}
	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &planModel)...)
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
	webhookId := state.WebhookId.String()

	// Get Webhook data from the API
	webhookResponse, err := r.client.GetClient().ShowWebhookWithResponse(ctx, webhookId, nil)

	if err != nil {
		tflog.Error(ctx, "Error reading webhook data from the API: ", map[string]interface{}{"error": err})
		resp.Diagnostics.AddError("Error reading webhook data from the API", err.Error())
		return
	}

	if webhookResponse.StatusCode() != 200 {
		tflog.Error(ctx, "Error reading webhook data from the API: ", map[string]interface{}{"error": webhookResponse.Body})
		resp.Diagnostics.AddError("Error reading webhook data from the API", fmt.Sprintf("%+v", webhookResponse.Body))
		return
	}

	webhookMapper := resource_webhook.NewWebhookMapper()

	webhookMapper.PutWebhookShowResponseToStateModel(ctx, webhookResponse, &state)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func logErrors(ctx context.Context, resp *resource.ReadResponse) {
	tflog.Error(ctx, "Error reading Terraform prior state data into the model: ", map[string]any{"error": resp.Diagnostics.Errors()})
}

func (r *webhookResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data resource_webhook.WebhookModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Update API call logic

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *webhookResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data resource_webhook.WebhookModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Delete API call logic
}
