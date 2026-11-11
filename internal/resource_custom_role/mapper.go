package resource_custom_role

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"math/big"
	"terraform-provider-zendesk/zendesk_api"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

// readOnlyConfigurationAttributes are documented as read-only by the Zendesk API and are never sent in requests.
var readOnlyConfigurationAttributes = []string{
	"chat_access",
	"group_access",
	"light_agent",
	"moderate_forums",
	"organization_notes_editing",
}

// MapPlanToCustomRoleRequest maps the Terraform plan to the request body used for creating and updating a custom role.
// Unknown and null values are not sent, so Zendesk keeps its defaults (on create) or the current values (on update).
func MapPlanToCustomRoleRequest(ctx context.Context, plan CustomRoleModel) (zendesk_api.CustomRoleRequest, diag.Diagnostics) {
	var diags diag.Diagnostics

	input := zendesk_api.CustomRoleInput{
		Name:        plan.CustomRole.Name.ValueStringPointer(),
		Description: valueStringPointerIfKnown(plan.CustomRole.Description),
	}

	configuration, err := attrValueToJSON(ctx, plan.CustomRole.Configuration)
	if err != nil {
		diags.AddError("Error mapping custom role configuration", err.Error())
		return zendesk_api.CustomRoleRequest{}, diags
	}
	if configurationMap, ok := configuration.(map[string]any); ok {
		for _, readOnly := range readOnlyConfigurationAttributes {
			delete(configurationMap, readOnly)
		}
		var apiConfiguration zendesk_api.CustomRoleConfigurationObject
		if err := convertViaJSON(configurationMap, &apiConfiguration); err != nil {
			diags.AddError("Error mapping custom role configuration", err.Error())
			return zendesk_api.CustomRoleRequest{}, diags
		}
		input.Configuration = &apiConfiguration
	}

	return zendesk_api.CustomRoleRequest{CustomRole: &input}, diags
}

// MapCustomRoleResponseToModel puts the custom role returned by the Zendesk API into the given model.
// Configuration attributes missing in the response keep their known value from the model, so that an
// attribute configured by the user but not returned by Zendesk does not produce an inconsistent result.
func MapCustomRoleResponseToModel(ctx context.Context, customRole *zendesk_api.CustomRoleObject, model *CustomRoleModel) diag.Diagnostics {
	var diags diag.Diagnostics

	if customRole == nil {
		diags.AddError("Empty custom role response", "The Zendesk API did not return a custom role")
		return diags
	}

	if model.CustomRole.IsNull() || model.CustomRole.IsUnknown() {
		model.CustomRole = NewCustomRoleValueMust(CustomRoleValue{}.AttributeTypes(ctx), map[string]attr.Value{
			"configuration":     types.ObjectNull(ConfigurationValue{}.AttributeTypes(ctx)),
			"created_at":        types.StringNull(),
			"description":       types.StringNull(),
			"id":                types.Int64Null(),
			"name":              types.StringNull(),
			"role_type":         types.Int64Null(),
			"team_member_count": types.Int64Null(),
			"updated_at":        types.StringNull(),
		})
	}

	model.CustomRole.Id = types.Int64PointerValue(customRole.Id)
	model.CustomRole.Name = types.StringValue(customRole.Name)
	model.CustomRole.Description = types.StringPointerValue(customRole.Description)
	model.CustomRole.RoleType = intPointerToInt64Value(customRole.RoleType)
	model.CustomRole.TeamMemberCount = intPointerToInt64Value(customRole.TeamMemberCount)
	model.CustomRole.CreatedAt = timePointerToStringValue(customRole.CreatedAt)
	model.CustomRole.UpdatedAt = timePointerToStringValue(customRole.UpdatedAt)
	if customRole.Id != nil {
		model.CustomRoleId = types.Int64Value(*customRole.Id)
	}

	configuration, configurationDiags := mapConfigurationResponse(ctx, customRole.Configuration, model.CustomRole.Configuration)
	diags.Append(configurationDiags...)
	if diags.HasError() {
		return diags
	}
	model.CustomRole.Configuration = configuration

	return diags
}

func mapConfigurationResponse(ctx context.Context, apiConfiguration *zendesk_api.CustomRoleConfigurationObject, prior basetypes.ObjectValue) (basetypes.ObjectValue, diag.Diagnostics) {
	var diags diag.Diagnostics
	configurationType := basetypes.ObjectType{AttrTypes: ConfigurationValue{}.AttributeTypes(ctx)}

	responseValues := map[string]any{}
	if apiConfiguration != nil {
		if err := convertViaJSON(apiConfiguration, &responseValues); err != nil {
			diags.AddError("Error mapping custom role configuration response", err.Error())
			return prior, diags
		}
	}

	if !prior.IsNull() && !prior.IsUnknown() {
		priorValues, err := attrValueToJSON(ctx, prior)
		if err != nil {
			diags.AddError("Error mapping custom role configuration response", err.Error())
			return prior, diags
		}
		if priorMap, ok := priorValues.(map[string]any); ok {
			for key, value := range priorMap {
				if _, returned := responseValues[key]; !returned {
					responseValues[key] = value
				}
			}
		}
	}

	terraformValue, err := jsonToTerraformValue(configurationType.TerraformType(ctx), responseValues)
	if err != nil {
		diags.AddError("Error mapping custom role configuration response", err.Error())
		return prior, diags
	}
	value, err := configurationType.ValueFromTerraform(ctx, terraformValue)
	if err != nil {
		diags.AddError("Error mapping custom role configuration response", err.Error())
		return prior, diags
	}
	objectValue, ok := value.(basetypes.ObjectValue)
	if !ok {
		diags.AddError("Error mapping custom role configuration response", fmt.Sprintf("unexpected value type %T", value))
		return prior, diags
	}
	return objectValue, diags
}

// attrValueToJSON converts a Terraform value into its JSON representation (bool, string, json.Number, []any, map[string]any).
// Null and unknown values return nil and are omitted from objects and maps.
func attrValueToJSON(ctx context.Context, value attr.Value) (any, error) {
	terraformValue, err := value.ToTerraformValue(ctx)
	if err != nil {
		return nil, err
	}
	return terraformValueToJSON(terraformValue)
}

func terraformValueToJSON(value tftypes.Value) (any, error) {
	if value.IsNull() || !value.IsKnown() {
		return nil, nil
	}

	valueType := value.Type()
	switch {
	case valueType.Is(tftypes.Bool):
		var b bool
		err := value.As(&b)
		return b, err
	case valueType.Is(tftypes.String):
		var s string
		err := value.As(&s)
		return s, err
	case valueType.Is(tftypes.Number):
		var n big.Float
		if err := value.As(&n); err != nil {
			return nil, err
		}
		return json.Number(n.Text('f', -1)), nil
	case valueType.Is(tftypes.List{}), valueType.Is(tftypes.Set{}), valueType.Is(tftypes.Tuple{}):
		var elements []tftypes.Value
		if err := value.As(&elements); err != nil {
			return nil, err
		}
		result := make([]any, 0, len(elements))
		for _, element := range elements {
			converted, err := terraformValueToJSON(element)
			if err != nil {
				return nil, err
			}
			if converted != nil {
				result = append(result, converted)
			}
		}
		return result, nil
	case valueType.Is(tftypes.Map{}), valueType.Is(tftypes.Object{}):
		var elements map[string]tftypes.Value
		if err := value.As(&elements); err != nil {
			return nil, err
		}
		result := make(map[string]any, len(elements))
		for key, element := range elements {
			converted, err := terraformValueToJSON(element)
			if err != nil {
				return nil, err
			}
			if converted != nil {
				result[key] = converted
			}
		}
		return result, nil
	default:
		return nil, fmt.Errorf("unsupported terraform type %s", valueType)
	}
}

// jsonToTerraformValue converts a decoded JSON value into a Terraform value of the given type.
// Object attributes missing in the JSON value become null.
func jsonToTerraformValue(valueType tftypes.Type, value any) (tftypes.Value, error) {
	if value == nil {
		return tftypes.NewValue(valueType, nil), nil
	}

	switch {
	case valueType.Is(tftypes.Bool):
		b, ok := value.(bool)
		if !ok {
			return tftypes.Value{}, fmt.Errorf("expected bool, got %T", value)
		}
		return tftypes.NewValue(valueType, b), nil
	case valueType.Is(tftypes.String):
		s, ok := value.(string)
		if !ok {
			return tftypes.Value{}, fmt.Errorf("expected string, got %T", value)
		}
		return tftypes.NewValue(valueType, s), nil
	case valueType.Is(tftypes.Number):
		var text string
		switch n := value.(type) {
		case json.Number:
			text = n.String()
		case float64:
			text = fmt.Sprintf("%v", n)
		default:
			return tftypes.Value{}, fmt.Errorf("expected number, got %T", value)
		}
		f, _, err := big.ParseFloat(text, 10, 512, big.ToNearestEven)
		if err != nil {
			return tftypes.Value{}, err
		}
		return tftypes.NewValue(valueType, f), nil
	}

	switch typed := valueType.(type) {
	case tftypes.List:
		return jsonArrayToTerraformValue(valueType, typed.ElementType, value)
	case tftypes.Set:
		return jsonArrayToTerraformValue(valueType, typed.ElementType, value)
	case tftypes.Map:
		jsonMap, ok := value.(map[string]any)
		if !ok {
			return tftypes.Value{}, fmt.Errorf("expected object, got %T", value)
		}
		elements := make(map[string]tftypes.Value, len(jsonMap))
		for key, element := range jsonMap {
			converted, err := jsonToTerraformValue(typed.ElementType, element)
			if err != nil {
				return tftypes.Value{}, fmt.Errorf("%s: %w", key, err)
			}
			elements[key] = converted
		}
		return tftypes.NewValue(valueType, elements), nil
	case tftypes.Object:
		jsonMap, ok := value.(map[string]any)
		if !ok {
			return tftypes.Value{}, fmt.Errorf("expected object, got %T", value)
		}
		attributes := make(map[string]tftypes.Value, len(typed.AttributeTypes))
		for key, attributeType := range typed.AttributeTypes {
			converted, err := jsonToTerraformValue(attributeType, jsonMap[key])
			if err != nil {
				return tftypes.Value{}, fmt.Errorf("%s: %w", key, err)
			}
			attributes[key] = converted
		}
		return tftypes.NewValue(valueType, attributes), nil
	default:
		return tftypes.Value{}, fmt.Errorf("unsupported terraform type %s", valueType)
	}
}

func jsonArrayToTerraformValue(valueType tftypes.Type, elementType tftypes.Type, value any) (tftypes.Value, error) {
	jsonArray, ok := value.([]any)
	if !ok {
		return tftypes.Value{}, fmt.Errorf("expected array, got %T", value)
	}
	elements := make([]tftypes.Value, 0, len(jsonArray))
	for i, element := range jsonArray {
		converted, err := jsonToTerraformValue(elementType, element)
		if err != nil {
			return tftypes.Value{}, fmt.Errorf("[%d]: %w", i, err)
		}
		elements = append(elements, converted)
	}
	return tftypes.NewValue(valueType, elements), nil
}

// convertViaJSON converts between Go representations by marshalling the source and unmarshalling into the target.
// Numbers are decoded as json.Number to keep int64 ids exact.
func convertViaJSON(source any, target any) error {
	data, err := json.Marshal(source)
	if err != nil {
		return err
	}
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.UseNumber()
	return decoder.Decode(target)
}

func valueStringPointerIfKnown(value basetypes.StringValue) *string {
	if value.IsUnknown() {
		return nil
	}
	return value.ValueStringPointer()
}

func intPointerToInt64Value(value *int) basetypes.Int64Value {
	if value == nil {
		return types.Int64Null()
	}
	return types.Int64Value(int64(*value))
}

func timePointerToStringValue(value *time.Time) basetypes.StringValue {
	if value == nil {
		return types.StringNull()
	}
	return types.StringValue(value.Format(time.RFC3339))
}
