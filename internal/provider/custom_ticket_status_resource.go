package provider

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"net/http"
	"strconv"
	"terraform-provider-zendesk/internal/resource_custom_status"
	"terraform-provider-zendesk/zendesk_api"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/resource"
)

var _ zendesk_api.Client

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &customStatusResource{}
	_ resource.ResourceWithConfigure   = &customStatusResource{}
	_ resource.ResourceWithImportState = &customStatusResource{}
)

func NewCustomStatusResource() resource.Resource {
	return &customStatusResource{}
}

type customStatusResource struct {
	client *zendesk_api.SupportApi
}

// ImportState imports a Custom Status by a given id, when the id value is an integer, or by label=id otherwise /*
func (r *customStatusResource) ImportState(ctx context.Context, request resource.ImportStateRequest, response *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	tflog.Debug(ctx, "Called ImportState custom status")
	id := request.ID
	tflog.Debug(ctx, "request.id: "+id)

	idInt, intParsingError := strconv.ParseInt(id, 10, 64)
	if intParsingError != nil {
		tflog.Debug(ctx, "custom status id could not be parsed as a number: "+intParsingError.Error())
		tflog.Info(ctx, "Will try to find a custom status by the label: "+id)
		customStatusesResponse, errors := r.client.GetClient().ListCustomStatusesWithResponse(ctx, nil, jsonContenttypeHeaderEditor)
		if errors != nil {
			tflog.Error(ctx, "Error listing custom statuses: ", map[string]any{"error": errors.Error()})
			response.Diagnostics.AddError("Error listing custom statuses", errors.Error())
			return
		}
		if customStatusesResponse.HTTPResponse.StatusCode != 200 {
			msg := "API error listing custom statuses: " + customStatusesResponse.HTTPResponse.Status
			tflog.Error(ctx, msg)
			response.Diagnostics.AddError(msg, "List custom statuses failed with status code: "+customStatusesResponse.HTTPResponse.Status)
			return
		}
		customStatuses := customStatusesResponse.JSON200.CustomStatuses
		tflog.Debug(ctx, "Found custom statuses: "+fmt.Sprintf("%v", customStatuses))
		for _, customStatus := range *customStatuses {
			if customStatus.AgentLabel == id {
				idInt = int64(*customStatus.Id)
				break
			}

		}
	}
	if idInt == 0 {
		tflog.Error(ctx, "Could not find custom status with id or label: "+id)
		response.Diagnostics.AddError("Could not find custom status with id or label: "+id, "")
		return
	}

	customStatusIdPath1 := path.Root("custom_status_id")
	customStatusIdPath2 := path.Root("custom_status").AtName("id")
	tflog.Debug(ctx, "ImportState custom status with customStatusIdPath1: "+customStatusIdPath1.String())
	tflog.Debug(ctx, "ImportState custom status with customStatusIdPath2: "+customStatusIdPath2.String())
	response.Diagnostics.Append(response.State.SetAttribute(ctx, customStatusIdPath1, idInt)...)
	response.Diagnostics.Append(response.State.SetAttribute(ctx, customStatusIdPath2, idInt)...)

	if response.Diagnostics.HasError() {

		tflog.Error(ctx, "Error importing state: ", map[string]any{"error": response.Diagnostics.Errors()})
		return
	}
	tflog.Info(ctx, "ImportState custom status completed successfully")
}

// Configure adds the provider configured client to the resource.
func (r *customStatusResource) Configure(ctx context.Context, request resource.ConfigureRequest, response *resource.ConfigureResponse) {
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

	r.client = providerData.supportApi

}

func (r *customStatusResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_custom_status"
}

func (r *customStatusResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = resource_custom_status.CustomStatusResourceSchema(ctx)
}

func (r *customStatusResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan resource_custom_status.CustomStatusModel

	tflog.Debug(ctx, "Called Create custom status with plan: "+structToString(plan))

	// Read Terraform plan data into the model
	tflog.Debug(ctx, "Reading Terraform plan data into the model")
	diags := req.Plan.Get(ctx, &plan)

	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		tflog.Error(ctx, "Error reading Terraform plan data into the model ")
		return
	}

	tflog.Debug(ctx, "Mapping custom status to the API Request Payload")
	requestPayload, err := mapPlanToCustomStatusCreateRequestPayload(plan)
	if err != nil {
		tflog.Error(ctx, "Error mapping custom status to the API Request Payload: ", map[string]any{"error": err.Error()})
		resp.Diagnostics.AddError("error mapping custom status to the API Request Payload", err.Error())
		return

	}
	tflog.Debug(ctx, "Mapped custom status to the API Request Payload successfully")

	customStatusCreateResponse, err := r.client.GetClient().CreateCustomStatusWithResponse(ctx, requestPayload)
	if err != nil {
		errorMessage := "API error creating custom status: "
		if customStatusCreateResponse != nil {
			if customStatusCreateResponse.HTTPResponse != nil {
				errorMessage += customStatusCreateResponse.HTTPResponse.Status
			}
		}
		tflog.Error(ctx, errorMessage, map[string]any{"error": err.Error()})
		resp.Diagnostics.AddError(errorMessage, err.Error())
		return
	}
	statusCode := customStatusCreateResponse.HTTPResponse.Status
	tflog.Debug(ctx,
		"API call to create custom status ended with status: "+statusCode)

	if customStatusCreateResponse.HTTPResponse.StatusCode != 201 {
		msg := "API error creating custom status: " + statusCode
		tflog.Error(ctx, msg)
		resp.Diagnostics.AddError(msg, "Create custom status failed with status code: "+statusCode+" and body: "+string(customStatusCreateResponse.Body))
		return
	}

	customStatusResponse := customStatusCreateResponse.JSON201

	tflog.Debug(ctx, "API response body: "+structToString(*customStatusResponse.CustomStatus))

	plan.CustomStatusId = types.Int64Value(int64(*customStatusResponse.CustomStatus.Id))

	mapCustomStatusFromResponseIntoPlan(customStatusResponse.CustomStatus, &plan)

	logDebugCustomStatusModelInPlan(ctx, plan)

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
	if resp.Diagnostics.HasError() {
		tflog.Error(ctx, "Error saving data into Terraform state ", map[string]any{"error": err.Error()})
		return
	}
	tflog.Debug(ctx, "Saved data into Terraform state: "+structToString(resp.State))
	tflog.Debug(ctx, "Create custom status completed successfully.")
}

func logDebugCustomStatusModelInPlan(ctx context.Context, plan resource_custom_status.CustomStatusModel) {
	tflog.Debug(ctx, "Plan is now: "+fmt.Sprintf("%v,", plan.CustomStatusId)+fmt.Sprintf("%v,", plan.CustomStatus.Id)+fmt.Sprintf("%v,", plan.CustomStatus.Active.ValueBool())+fmt.Sprintf("%v,", plan.CustomStatus.CreatedAt)+fmt.Sprintf("%v,", plan.CustomStatus.AgentLabel)+fmt.Sprintf("%v,", plan.CustomStatus.Description)+fmt.Sprintf("%v,", plan.CustomStatus.EndUserDescription)+fmt.Sprintf("%v,", plan.CustomStatus.EndUserLabel)+fmt.Sprintf("%v,", plan.CustomStatus.StatusCategory)+fmt.Sprintf("%v,", plan.CustomStatus.UpdatedAt))
}

func structToString(obj any) string {
	return fmt.Sprintf("%+v", obj)
}

func mapCustomStatusFromResponseIntoPlan(customStatus *zendesk_api.CustomStatusObject, model *resource_custom_status.CustomStatusModel) {
	model.CustomStatus.Id = types.Int64Value(int64(*customStatus.Id))
	model.CustomStatus.Active = types.BoolValue(*customStatus.Active)
	model.CustomStatus.AgentLabel = types.StringValue(customStatus.AgentLabel)
	model.CustomStatus.StatusCategory = types.StringValue(string(customStatus.StatusCategory))
	model.CustomStatus.CreatedAt = types.StringValue(customStatus.CreatedAt.Format(time.RFC3339))
	model.CustomStatus.Description = types.StringValue(*customStatus.Description)
	model.CustomStatus.EndUserLabel = types.StringValue(*customStatus.EndUserLabel)
	model.CustomStatus.EndUserDescription = types.StringValue(*customStatus.EndUserDescription)
	model.CustomStatus.Default = types.BoolValue(*customStatus.Default)
	model.CustomStatus.RawAgentLabel = types.StringValue(*customStatus.RawAgentLabel)
	model.CustomStatus.RawDescription = types.StringValue(*customStatus.RawDescription)
	model.CustomStatus.RawEndUserLabel = types.StringValue(*customStatus.RawEndUserLabel)
	model.CustomStatus.RawEndUserDescription = types.StringValue(*customStatus.RawEndUserDescription)
	model.CustomStatus.UpdatedAt = types.StringValue(customStatus.UpdatedAt.Format(time.RFC3339))

}

func mapPlanToCustomStatusCreateRequestPayload(customStatusModel resource_custom_status.CustomStatusModel) (zendesk_api.CustomStatusCreateRequest, error) {

	category, err := mapToStatusCategory(customStatusModel.CustomStatus.StatusCategory.ValueString())
	if err != nil {
		return zendesk_api.CustomStatusCreateRequest{}, err
	}
	active := true
	if !customStatusModel.CustomStatus.Active.IsUnknown() {
		active = customStatusModel.CustomStatus.Active.ValueBool()
	}
	input := zendesk_api.CustomStatusCreateInput{
		Active:             &active,
		AgentLabel:         customStatusModel.CustomStatus.AgentLabel.ValueStringPointer(),
		StatusCategory:     &category,
		Description:        customStatusModel.CustomStatus.Description.ValueStringPointer(),
		EndUserDescription: customStatusModel.CustomStatus.EndUserDescription.ValueStringPointer(),
		EndUserLabel:       customStatusModel.CustomStatus.EndUserLabel.ValueStringPointer(),
	}
	payload := zendesk_api.CustomStatusCreateRequest{
		CustomStatus: &input,
	}
	return payload, nil
}

func mapToStatusCategory(valueString string) (zendesk_api.CustomStatusCreateInputStatusCategory, error) {
	switch valueString {
	case "new":
		return zendesk_api.CustomStatusCreateInputStatusCategoryNew, nil
	case "open":
		return zendesk_api.CustomStatusCreateInputStatusCategoryOpen, nil
	case "pending":
		return zendesk_api.CustomStatusCreateInputStatusCategoryPending, nil
	case "hold":
		return zendesk_api.CustomStatusCreateInputStatusCategoryHold, nil
		// Tickets with a "Closed" status belong to the "StatusCategorySolved" status category.
	case "solved":
		return zendesk_api.CustomStatusCreateInputStatusCategorySolved, nil
	default:
		return zendesk_api.CustomStatusCreateInputStatusCategory(valueString), fmt.Errorf("invalid status category %s", valueString)
	}
}

func (r *customStatusResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state resource_custom_status.CustomStatusModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		tflog.Error(ctx, "Error reading Terraform prior state data into the model: ", map[string]any{"error": resp.Diagnostics.Errors()})
		return
	}

	tflog.Debug(ctx, "Called Read custom status with state: "+structToString(state))

	// Get refreshed CustomStatus value from Zendesk
	statusId, diagsStatusId := mapCustomStatusId(ctx, state)
	if diagsStatusId != nil {
		resp.Diagnostics.Append(diagsStatusId...)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	tflog.Debug(ctx, "Mapped custom status Id from state: <"+structToString(statusId)+">")

	customStatusReadResponse, err := r.client.GetClient().ShowCustomStatusWithResponse(ctx, statusId, jsonContenttypeHeaderEditor)

	if err != nil {
		msg := "Could not read Zendesk CustomStatus with id= " + state.CustomStatusId.String() + ": " + err.Error()
		tflog.Error(ctx, msg)
		resp.Diagnostics.AddError(
			"Error Reading Zendesk CustomStatus",
			msg,
		)
		return
	}

	if customStatusReadResponse.HTTPResponse.StatusCode != 200 {
		msg := "Error Reading Zendesk CustomStatus with id= " + state.CustomStatusId.String() + " and status: " + customStatusReadResponse.HTTPResponse.Status + " and body: <" + string(customStatusReadResponse.Body) + ">"
		tflog.Error(ctx, msg)

		resp.Diagnostics.AddError(
			"Failure Reading Zendesk CustomStatus",
			msg,
		)
		return

	}

	customStatusValue := customStatusReadResponse.JSON200.CustomStatus
	tflog.Debug(ctx, "Mapping Response Body to State: <"+structToString(customStatusValue)+">")

	mapCustomStatusFromResponseIntoPlan(customStatusValue, &state)

	logDebugCustomStatusModelInPlan(ctx, state)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func jsonContenttypeHeaderEditor(ctx context.Context, req *http.Request) error {
	req.Header.Set("Content-Type", "application/json")
	return nil
}

func mapCustomStatusId(ctx context.Context, state resource_custom_status.CustomStatusModel) (zendesk_api.CustomStatusId, diag.Diagnostics) {
	intVal, diags := state.CustomStatusId.ToInt64Value(ctx)
	if diags.HasError() {
		tflog.Error(ctx, "Error converting CustomStatusId to int64", map[string]any{"error": diags.Errors()})
		return 0, diags
	}
	var statusId = int(intVal.ValueInt64())
	return statusId, nil
}

func (r *customStatusResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan resource_custom_status.CustomStatusModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		tflog.Error(ctx, "Error reading Terraform plan data into the model: ", map[string]any{"error": resp.Diagnostics.Errors()})
		return
	}
	tflog.Debug(ctx, "Called Update custom status with plan: "+structToString(plan))

	var state resource_custom_status.CustomStatusModel

	// Read Terraform State data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		tflog.Error(ctx, "Error reading Terraform State data into the model: ", map[string]any{"error": resp.Diagnostics.Errors()})
		return
	}

	tflog.Debug(ctx, "Update custom status State is: "+structToString(state))

	// Get refreshed CustomStatus value from Zendesk
	statusId, diagsStatusId := mapCustomStatusId(ctx, state)

	if diagsStatusId != nil {
		resp.Diagnostics.Append(diagsStatusId...)
		if resp.Diagnostics.HasError() {
			tflog.Error(ctx, "Error mapping CustomStatusId: ", map[string]any{"error": resp.Diagnostics.Errors()})
			return
		}
	}

	// Read Terraform prior state data into the model
	var priorState resource_custom_status.CustomStatusModel
	//Check that Status Category has not changed
	req.State.Get(ctx, &priorState)
	if !priorState.CustomStatus.StatusCategory.Equal(plan.CustomStatus.StatusCategory) {
		tflog.Debug(ctx, "Status Category has changed from: "+priorState.CustomStatus.StatusCategory.ValueString()+" to: "+plan.CustomStatus.StatusCategory.ValueString())
		resp.Diagnostics.AddError("Status Category cannot be updated", "Status Category cannot be updated, replace Status instead")
		return
	}

	tflog.Debug(ctx, "Mapped custom status Id from state: "+structToString(statusId))

	// Map plan to API request payload

	updateRequestBody := mapPlanToCustomStatusUpdateRequestPayload(plan)

	// Update Custom Status call logic
	customStatusUpdateResponse, err := r.client.GetClient().UpdateCustomStatusWithResponse(ctx, statusId, updateRequestBody)

	if err != nil {
		errorMessage := "API error updating custom status: "
		if customStatusUpdateResponse != nil {
			if customStatusUpdateResponse.HTTPResponse != nil {
				errorMessage += customStatusUpdateResponse.HTTPResponse.Status
			}
		}
		tflog.Error(ctx, errorMessage, map[string]any{"error": err.Error()})
		resp.Diagnostics.AddError(errorMessage, err.Error())
		return

	}
	statusCode := customStatusUpdateResponse.HTTPResponse.Status
	tflog.Debug(ctx, "API call to update custom status ended with status: "+statusCode)

	if customStatusUpdateResponse.HTTPResponse.StatusCode != 200 {
		msg := "API error updating custom status: " + statusCode
		tflog.Error(ctx, msg)
		resp.Diagnostics.AddError(msg, "Update custom status failed with status code: "+statusCode+" and body: "+string(customStatusUpdateResponse.Body))
		return

	}

	customStatusResponse := customStatusUpdateResponse.JSON200
	tflog.Debug(ctx, "API response body: "+structToString(*customStatusResponse.CustomStatus))

	plan.CustomStatusId = types.Int64Value(int64(*customStatusResponse.CustomStatus.Id))

	mapCustomStatusFromResponseIntoPlan(customStatusResponse.CustomStatus, &plan)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)

	logDebugCustomStatusModelInPlan(ctx, plan)

	if resp.Diagnostics.HasError() {
		tflog.Error(ctx, "Error saving data into Terraform state, but resource was successfully updated ", map[string]any{"error": err.Error()})
		return

	}

	tflog.Debug(ctx, "Saved data into Terraform state: "+structToString(resp.State))
	tflog.Debug(ctx, "Update custom status completed successfully.")
}

func mapPlanToCustomStatusUpdateRequestPayload(plan resource_custom_status.CustomStatusModel) zendesk_api.CustomStatusUpdateRequest {
	active := true
	if !plan.CustomStatus.Active.IsUnknown() {
		active = plan.CustomStatus.Active.ValueBool()
	}
	input := zendesk_api.CustomStatusUpdateInput{
		Active:             &active,
		AgentLabel:         plan.CustomStatus.AgentLabel.ValueStringPointer(),
		Description:        plan.CustomStatus.Description.ValueStringPointer(),
		EndUserDescription: plan.CustomStatus.EndUserDescription.ValueStringPointer(),
		EndUserLabel:       plan.CustomStatus.EndUserLabel.ValueStringPointer(),
	}
	payload := zendesk_api.CustomStatusUpdateRequest{
		CustomStatus: &input,
	}
	return payload

}

func (r *customStatusResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state resource_custom_status.CustomStatusModel

	// Read Terraform prior state state into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	tflog.Debug(ctx, "Called Delete custom status with state: "+structToString(state))

	statusId, diagStatusId := mapCustomStatusId(ctx, state)
	if diagStatusId != nil {
		resp.Diagnostics.Append(diagStatusId...)
		if resp.Diagnostics.HasError() {
			tflog.Error(ctx, "Error mapping CustomStatusId: ", map[string]any{"error": resp.Diagnostics.Errors()})
			return
		}
	}

	customStatusReadResponse, err := r.client.GetClient().ShowCustomStatusWithResponse(ctx, statusId, jsonContenttypeHeaderEditor)

	if err != nil {
		msg := "Could not read Zendesk CustomStatus with id= " + state.CustomStatusId.String() + ": " + err.Error()
		tflog.Error(ctx, msg)
		resp.Diagnostics.AddError(
			"Error Reading Zendesk CustomStatus",
			msg,
		)
		return
	}

	if customStatusReadResponse.HTTPResponse.StatusCode != 200 {
		msg := "Error Reading Zendesk CustomStatus with id= " + state.CustomStatusId.String() + " and status: " + customStatusReadResponse.HTTPResponse.Status + " and body: <" + string(customStatusReadResponse.Body) + ">"
		tflog.Error(ctx, msg)

		resp.Diagnostics.AddError(
			"Failure Reading Zendesk CustomStatus",
			msg,
		)
		return

	}

	customStatusValue := customStatusReadResponse.JSON200.CustomStatus
	tflog.Debug(ctx, "Mapping Response Body to State: <"+structToString(customStatusValue)+">")

	if !state.CustomStatus.Active.ValueBool() || !*customStatusValue.Active {
		resp.Diagnostics.AddWarning("Custom Status is already deactivated, will be now removed from the Terraform State.",
			"Custom Status is already deactivated. Deletion can be implemented only by deactivation. It will be now removed from the Terraform State.")
		tflog.Warn(ctx, "Custom Status is already deactivated! It will be now removed from the Terraform State.")

	} else {

		updateRequestBody := mapCustomStatusToDeactivatedInput(false, customStatusValue)

		status, err := r.client.GetClient().UpdateCustomStatusWithResponse(ctx, statusId, updateRequestBody)
		if err != nil {
			resp.Diagnostics.AddError("Error deactivating custom status", err.Error())
			tflog.Error(ctx, "Error deactivating custom status: ", map[string]any{"error": err.Error()})
			return

		}

		tflog.Debug(ctx, "Deactivated custom status with status: "+structToString(status))
		if status.HTTPResponse.StatusCode != 200 {
			resp.Diagnostics.AddError("Error deactivating custom status", "Deactivation failed with status code: "+status.HTTPResponse.Status)
			tflog.Error(ctx, "Error deactivating custom status: ", map[string]any{"status": status.HTTPResponse.Status})
			return
		}

		tflog.Debug(ctx, "Updated custom status with active=false successfully")
	}

	state.CustomStatus.Active = types.BoolValue(false)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)

	logDebugCustomStatusModelInPlan(ctx, state)

	if resp.Diagnostics.HasError() {
		tflog.Error(ctx, "Error saving data into Terraform state, but resource was successfully updated ", map[string]any{"error": err.Error()})
		return

	}

	tflog.Debug(ctx, "Saved data into Terraform state: "+structToString(resp.State))
	tflog.Debug(ctx, "Deactivate custom status completed successfully.")

}

func mapCustomStatusToDeactivatedInput(deactivated bool, customStatusValue *zendesk_api.CustomStatusObject) zendesk_api.CustomStatusUpdateRequest {
	return zendesk_api.CustomStatusUpdateRequest{
		CustomStatus: &zendesk_api.CustomStatusUpdateInput{
			Active:             &deactivated,
			AgentLabel:         &customStatusValue.AgentLabel,
			Description:        customStatusValue.Description,
			EndUserDescription: customStatusValue.EndUserDescription,
			EndUserLabel:       customStatusValue.EndUserLabel,
		},
	}
}
