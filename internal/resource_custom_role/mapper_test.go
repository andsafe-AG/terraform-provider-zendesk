package resource_custom_role

import (
	"context"
	"encoding/json"
	"terraform-provider-zendesk/zendesk_api"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"gotest.tools/v3/assert"
)

const customRoleResponseBody = `{
  "custom_role": {
    "configuration": {
      "chat_access": true,
      "custom_objects": {
        "shipment": { "scopes": ["read", "update"] }
      },
      "end_user_profile_access": "readonly",
      "explore_reports": {
        "groups_ticket_access": "selected-groups",
        "ticket_access_selected_groups": ["1260803907490"]
      },
      "ticket_access": "within-groups",
      "ticket_editing": true
    },
    "created_at": "2012-02-20T22:55:29Z",
    "description": "Can only make private comments on assigned tickets",
    "id": 10127,
    "name": "Partner",
    "role_type": 0,
    "team_member_count": 2,
    "updated_at": "2012-02-21T22:55:29Z"
  }
}`

func customRoleFromResponse(t *testing.T) *zendesk_api.CustomRoleObject {
	var response zendesk_api.CustomRoleResponse
	assert.NilError(t, json.Unmarshal([]byte(customRoleResponseBody), &response))
	return response.CustomRole
}

func TestMapCustomRoleResponseToModel_emptyModel(t *testing.T) {
	ctx := context.Background()
	model := CustomRoleModel{CustomRole: NewCustomRoleValueNull()}

	diags := MapCustomRoleResponseToModel(ctx, customRoleFromResponse(t), &model)

	assert.Assert(t, !diags.HasError(), diags)
	assert.Equal(t, model.CustomRoleId.ValueInt64(), int64(10127))
	assert.Equal(t, model.CustomRole.Id.ValueInt64(), int64(10127))
	assert.Equal(t, model.CustomRole.Name.ValueString(), "Partner")
	assert.Equal(t, model.CustomRole.Description.ValueString(), "Can only make private comments on assigned tickets")
	assert.Equal(t, model.CustomRole.RoleType.ValueInt64(), int64(0))
	assert.Equal(t, model.CustomRole.TeamMemberCount.ValueInt64(), int64(2))
	assert.Equal(t, model.CustomRole.CreatedAt.ValueString(), "2012-02-20T22:55:29Z")
	assert.Equal(t, model.CustomRole.UpdatedAt.ValueString(), "2012-02-21T22:55:29Z")

	configuration := model.CustomRole.Configuration.Attributes()
	assert.Equal(t, configuration["chat_access"], types.BoolValue(true))
	assert.Equal(t, configuration["ticket_access"], types.StringValue("within-groups"))
	assert.Equal(t, configuration["ticket_editing"], types.BoolValue(true))
	assert.Assert(t, configuration["manage_triggers"].IsNull(), "attributes missing in the response must be null")

	customObjects, ok := configuration["custom_objects"].(basetypes.MapValue)
	assert.Assert(t, ok)
	shipment, ok := customObjects.Elements()["shipment"].(CustomObjectsValue)
	assert.Assert(t, ok)
	assert.DeepEqual(t, shipment.Scopes.Elements(), []attr.Value{types.StringValue("read"), types.StringValue("update")})

	exploreReports, ok := configuration["explore_reports"].(basetypes.ObjectValue)
	assert.Assert(t, ok)
	assert.Equal(t, exploreReports.Attributes()["groups_ticket_access"], types.StringValue("selected-groups"))
	assert.Assert(t, exploreReports.Attributes()["brands_ticket_access"].IsNull())
}

func TestMapCustomRoleResponseToModel_keepsConfiguredValuesMissingInResponse(t *testing.T) {
	ctx := context.Background()
	model := CustomRoleModel{CustomRole: NewCustomRoleValueNull()}
	assert.Assert(t, !MapCustomRoleResponseToModel(ctx, customRoleFromResponse(t), &model).HasError())

	// simulate a configured value that is not part of the API response
	attributes := model.CustomRole.Configuration.Attributes()
	attributes["manage_triggers"] = types.BoolValue(true)
	attributes["ticket_access"] = types.StringValue("all")
	configuration, diags := types.ObjectValue(ConfigurationValue{}.AttributeTypes(ctx), attributes)
	assert.Assert(t, !diags.HasError(), diags)
	model.CustomRole.Configuration = configuration

	diags = MapCustomRoleResponseToModel(ctx, customRoleFromResponse(t), &model)

	assert.Assert(t, !diags.HasError(), diags)
	result := model.CustomRole.Configuration.Attributes()
	assert.Equal(t, result["manage_triggers"], types.BoolValue(true), "configured value missing in response must be kept")
	assert.Equal(t, result["ticket_access"], types.StringValue("within-groups"), "value returned by the API must win")
}

func TestMapPlanToCustomRoleRequest(t *testing.T) {
	ctx := context.Background()
	model := CustomRoleModel{CustomRole: NewCustomRoleValueNull()}
	assert.Assert(t, !MapCustomRoleResponseToModel(ctx, customRoleFromResponse(t), &model).HasError())

	// unknown values must not be sent
	attributes := model.CustomRole.Configuration.Attributes()
	attributes["manage_triggers"] = types.BoolUnknown()
	configuration, diags := types.ObjectValue(ConfigurationValue{}.AttributeTypes(ctx), attributes)
	assert.Assert(t, !diags.HasError(), diags)
	model.CustomRole.Configuration = configuration

	request, diags := MapPlanToCustomRoleRequest(ctx, model)

	assert.Assert(t, !diags.HasError(), diags)
	body, err := json.Marshal(request)
	assert.NilError(t, err)
	assert.Equal(t, string(body), `{"custom_role":{"configuration":{"custom_objects":{"shipment":{"scopes":["read","update"]}},`+
		`"end_user_profile_access":"readonly","explore_reports":{"groups_ticket_access":"selected-groups","ticket_access_selected_groups":["1260803907490"]},`+
		`"ticket_access":"within-groups","ticket_editing":true},`+
		`"description":"Can only make private comments on assigned tickets","name":"Partner"}}`)
}

func TestMapPlanToCustomRoleRequest_unknownConfigurationAndDescription(t *testing.T) {
	ctx := context.Background()
	model := CustomRoleModel{
		CustomRole: NewCustomRoleValueMust(CustomRoleValue{}.AttributeTypes(ctx), map[string]attr.Value{
			"configuration":     types.ObjectUnknown(ConfigurationValue{}.AttributeTypes(ctx)),
			"created_at":        types.StringUnknown(),
			"description":       types.StringUnknown(),
			"id":                types.Int64Unknown(),
			"name":              types.StringValue("Agent"),
			"role_type":         types.Int64Unknown(),
			"team_member_count": types.Int64Unknown(),
			"updated_at":        types.StringUnknown(),
		}),
		CustomRoleId: types.Int64Unknown(),
	}

	request, diags := MapPlanToCustomRoleRequest(ctx, model)

	assert.Assert(t, !diags.HasError(), diags)
	body, err := json.Marshal(request)
	assert.NilError(t, err)
	assert.Equal(t, string(body), `{"custom_role":{"name":"Agent"}}`)
}
