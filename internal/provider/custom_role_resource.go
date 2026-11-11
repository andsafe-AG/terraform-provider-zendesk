package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"terraform-provider-zendesk/internal/resource_custom_role"
	"terraform-provider-zendesk/zendesk_api"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &customRoleResource{}
	_ resource.ResourceWithConfigure   = &customRoleResource{}
	_ resource.ResourceWithImportState = &customRoleResource{}
)

func NewCustomRoleResource() resource.Resource {
	return &customRoleResource{}
}

type customRoleResource struct {
	client *zendesk_api.SupportApi
}

func (r *customRoleResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_custom_role"
}

func (r *customRoleResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = resource_custom_role.CustomRoleResourceSchema(ctx)
}

// Configure adds the provider configured client to the resource.
func (r *customRoleResource) Configure(ctx context.Context, request resource.ConfigureRequest, response *resource.ConfigureResponse) {
	// Add a nil check when handling ProviderData because Terraform
	// sets that data after it calls the ConfigureProvider RPC.
	if request.ProviderData == nil {
		return
	}

	providerData, ok := request.ProviderData.(zendeskProviderData)
	if !ok {
		response.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected zendeskProviderData, got: %T. Please report this issue to the provider developers.",
				request.ProviderData),
		)
		return
	}

	r.client = providerData.supportApi
}

// ImportState imports a Custom Role by a given id, when the id value is an integer, or by its name otherwise /*
func (r *customRoleResource) ImportState(ctx context.Context, request resource.ImportStateRequest, response *resource.ImportStateResponse) {
	id := request.ID
	tflog.Debug(ctx, "Called ImportState custom role with id: "+id)

	roleId, intParsingError := strconv.ParseInt(id, 10, 64)
	if intParsingError != nil {
		tflog.Info(ctx, "Custom role id could not be parsed as a number, will try to find a custom role by the name: "+id)
		listResponse, err := r.client.GetClient().ListCustomRolesWithResponse(ctx, jsonContenttypeHeaderEditor)
		if err != nil {
			tflog.Error(ctx, "Error listing custom roles: ", map[string]any{"error": err.Error()})
			response.Diagnostics.AddError("Error listing custom roles", err.Error())
			return
		}
		if listResponse.HTTPResponse.StatusCode != http.StatusOK || listResponse.JSON200 == nil || listResponse.JSON200.CustomRoles == nil {
			msg := "API error listing custom roles: " + listResponse.HTTPResponse.Status
			tflog.Error(ctx, msg)
			response.Diagnostics.AddError(msg, "List custom roles failed with body: "+string(listResponse.Body))
			return
		}
		for _, customRole := range *listResponse.JSON200.CustomRoles {
			if customRole.Name == id && customRole.Id != nil {
				roleId = *customRole.Id
				break
			}
		}
	}
	if roleId == 0 {
		tflog.Error(ctx, "Could not find custom role with id or name: "+id)
		response.Diagnostics.AddError("Could not find custom role with id or name: "+id, "")
		return
	}

	response.Diagnostics.Append(response.State.SetAttribute(ctx, path.Root("custom_role_id"), roleId)...)
	response.Diagnostics.Append(response.State.SetAttribute(ctx, path.Root("custom_role").AtName("id"), roleId)...)
	if response.Diagnostics.HasError() {
		tflog.Error(ctx, "Error importing state: ", map[string]any{"error": response.Diagnostics.Errors()})
		return
	}
	tflog.Info(ctx, "ImportState custom role completed successfully")
}

func (r *customRoleResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan resource_custom_role.CustomRoleModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		tflog.Error(ctx, "Error reading Terraform plan data into the model")
		return
	}
	tflog.Debug(ctx, "Called Create custom role with plan: "+structToString(plan))

	requestBody, diags := resource_custom_role.MapPlanToCustomRoleRequest(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	createResponse, err := r.client.GetClient().CreateCustomRoleWithResponse(ctx, requestBody)
	if err != nil {
		tflog.Error(ctx, "API error creating custom role", map[string]any{"error": err.Error()})
		resp.Diagnostics.AddError("API error creating custom role", err.Error())
		return
	}
	statusCode := createResponse.HTTPResponse.StatusCode
	tflog.Debug(ctx, "API call to create custom role ended with status: "+createResponse.HTTPResponse.Status)

	// The API reference documents 200 for this call, 201 is accepted as well
	if statusCode != http.StatusOK && statusCode != http.StatusCreated {
		msg := "API error creating custom role: " + createResponse.HTTPResponse.Status
		tflog.Error(ctx, msg)
		resp.Diagnostics.AddError(msg, "Create custom role failed with body: "+string(createResponse.Body))
		return
	}

	customRole, diags := decodeCustomRoleResponse(createResponse.Body)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resource_custom_role.MapCustomRoleResponseToModel(ctx, customRole, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
	tflog.Debug(ctx, "Create custom role completed.")
}

func (r *customRoleResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state resource_custom_role.CustomRoleModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		tflog.Error(ctx, "Error reading Terraform prior state data into the model: ", map[string]any{"error": resp.Diagnostics.Errors()})
		return
	}
	tflog.Debug(ctx, "Called Read custom role with state: "+structToString(state))

	roleId := state.CustomRoleId.ValueInt64()
	readResponse, err := r.client.GetClient().ShowCustomRoleByIdWithResponse(ctx, roleId, jsonContenttypeHeaderEditor)
	if err != nil {
		msg := "Could not read Zendesk custom role with id= " + state.CustomRoleId.String() + ": " + err.Error()
		tflog.Error(ctx, msg)
		resp.Diagnostics.AddError("Error Reading Zendesk custom role", msg)
		return
	}

	if readResponse.HTTPResponse.StatusCode == http.StatusNotFound {
		tflog.Warn(ctx, "Custom role with id= "+state.CustomRoleId.String()+" not found, removing it from the Terraform state")
		resp.State.RemoveResource(ctx)
		return
	}

	if readResponse.HTTPResponse.StatusCode != http.StatusOK {
		msg := "Error Reading Zendesk custom role with id= " + state.CustomRoleId.String() + " and status: " + readResponse.HTTPResponse.Status + " and body: <" + string(readResponse.Body) + ">"
		tflog.Error(ctx, msg)
		resp.Diagnostics.AddError("Failure Reading Zendesk custom role", msg)
		return
	}

	customRole, diags := decodeCustomRoleResponse(readResponse.Body)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resource_custom_role.MapCustomRoleResponseToModel(ctx, customRole, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *customRoleResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan resource_custom_role.CustomRoleModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)

	var state resource_custom_role.CustomRoleModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		tflog.Error(ctx, "Error reading Terraform plan or state data into the model: ", map[string]any{"error": resp.Diagnostics.Errors()})
		return
	}
	tflog.Debug(ctx, "Called Update custom role with plan: "+structToString(plan))

	requestBody, diags := resource_custom_role.MapPlanToCustomRoleRequest(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	roleId := state.CustomRoleId.ValueInt64()
	updateResponse, err := r.client.GetClient().UpdateCustomRoleByIdWithResponse(ctx, roleId, requestBody)
	if err != nil {
		tflog.Error(ctx, "API error updating custom role", map[string]any{"error": err.Error()})
		resp.Diagnostics.AddError("API error updating custom role", err.Error())
		return
	}
	tflog.Debug(ctx, "API call to update custom role ended with status: "+updateResponse.HTTPResponse.Status)

	if updateResponse.HTTPResponse.StatusCode != http.StatusOK {
		msg := "API error updating custom role: " + updateResponse.HTTPResponse.Status
		tflog.Error(ctx, msg)
		resp.Diagnostics.AddError(msg, "Update custom role failed with body: "+string(updateResponse.Body))
		return
	}

	customRole, diags := decodeCustomRoleResponse(updateResponse.Body)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resource_custom_role.MapCustomRoleResponseToModel(ctx, customRole, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
	tflog.Debug(ctx, "Update custom role completed.")
}

func (r *customRoleResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state resource_custom_role.CustomRoleModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Debug(ctx, "Called Delete custom role with state: "+structToString(state))

	roleId := state.CustomRoleId.ValueInt64()
	deleteResponse, err := r.client.GetClient().DeleteCustomRoleByIdWithResponse(ctx, roleId, jsonContenttypeHeaderEditor)
	if err != nil {
		tflog.Error(ctx, "API error deleting custom role", map[string]any{"error": err.Error()})
		resp.Diagnostics.AddError("API error deleting custom role", err.Error())
		return
	}

	switch deleteResponse.HTTPResponse.StatusCode {
	case http.StatusNoContent, http.StatusOK:
		tflog.Debug(ctx, "Delete custom role completed.")
	case http.StatusNotFound:
		tflog.Warn(ctx, "Custom role with id= "+state.CustomRoleId.String()+" was already deleted")
	default:
		msg := "API error deleting custom role: " + deleteResponse.HTTPResponse.Status
		tflog.Error(ctx, msg)
		resp.Diagnostics.AddError(msg, "Delete custom role failed with body: "+string(deleteResponse.Body))
	}
}

// decodeCustomRoleResponse decodes the response body directly, because the generated client only parses
// the documented status code (200), while creating a custom role may also answer with 201.
func decodeCustomRoleResponse(body []byte) (*zendesk_api.CustomRoleObject, diag.Diagnostics) {
	var diags diag.Diagnostics
	var response zendesk_api.CustomRoleResponse
	if err := json.Unmarshal(body, &response); err != nil {
		diags.AddError("Error decoding custom role response", err.Error()+", body: "+string(body))
		return nil, diags
	}
	if response.CustomRole == nil {
		diags.AddError("Error decoding custom role response", "custom_role is missing in body: "+string(body))
		return nil, diags
	}
	return response.CustomRole, diags
}
