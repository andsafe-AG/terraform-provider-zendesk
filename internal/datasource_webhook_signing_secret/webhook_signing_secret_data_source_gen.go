// Code generated by terraform-plugin-framework-generator DO NOT EDIT.

package datasource_webhook_signing_secret

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
)

func WebhookSigningSecretDataSourceSchema(ctx context.Context) schema.Schema {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"signing_secret": schema.SingleNestedAttribute{
				Attributes: map[string]schema.Attribute{
					"algorithm": schema.StringAttribute{
						Computed: true,
					},
					"secret": schema.StringAttribute{
						Computed:  true,
						Sensitive: true,
					},
				},
				CustomType: SigningSecretType{
					ObjectType: types.ObjectType{
						AttrTypes: SigningSecretValue{}.AttributeTypes(ctx),
					},
				},
				Computed: true,
			},
			"webhook_id": schema.StringAttribute{
				Required:            true,
				Description:         "(Required) Webhook id",
				MarkdownDescription: "(Required) Webhook id",
			},
		},
	}
}

type WebhookSigningSecretModel struct {
	SigningSecret SigningSecretValue `tfsdk:"signing_secret"`
	WebhookId     types.String       `tfsdk:"webhook_id"`
}

var _ basetypes.ObjectTypable = SigningSecretType{}

type SigningSecretType struct {
	basetypes.ObjectType
}

func (t SigningSecretType) Equal(o attr.Type) bool {
	other, ok := o.(SigningSecretType)

	if !ok {
		return false
	}

	return t.ObjectType.Equal(other.ObjectType)
}

func (t SigningSecretType) String() string {
	return "SigningSecretType"
}

func (t SigningSecretType) ValueFromObject(ctx context.Context, in basetypes.ObjectValue) (basetypes.ObjectValuable, diag.Diagnostics) {
	var diags diag.Diagnostics

	attributes := in.Attributes()

	algorithmAttribute, ok := attributes["algorithm"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`algorithm is missing from object`)

		return nil, diags
	}

	algorithmVal, ok := algorithmAttribute.(basetypes.StringValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`algorithm expected to be basetypes.StringValue, was: %T`, algorithmAttribute))
	}

	secretAttribute, ok := attributes["secret"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`secret is missing from object`)

		return nil, diags
	}

	secretVal, ok := secretAttribute.(basetypes.StringValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`secret expected to be basetypes.StringValue, was: %T`, secretAttribute))
	}

	if diags.HasError() {
		return nil, diags
	}

	return SigningSecretValue{
		Algorithm: algorithmVal,
		Secret:    secretVal,
		state:     attr.ValueStateKnown,
	}, diags
}

func NewSigningSecretValueNull() SigningSecretValue {
	return SigningSecretValue{
		state: attr.ValueStateNull,
	}
}

func NewSigningSecretValueUnknown() SigningSecretValue {
	return SigningSecretValue{
		state: attr.ValueStateUnknown,
	}
}

func NewSigningSecretValue(attributeTypes map[string]attr.Type, attributes map[string]attr.Value) (SigningSecretValue, diag.Diagnostics) {
	var diags diag.Diagnostics

	// Reference: https://github.com/hashicorp/terraform-plugin-framework/issues/521
	ctx := context.Background()

	for name, attributeType := range attributeTypes {
		attribute, ok := attributes[name]

		if !ok {
			diags.AddError(
				"Missing SigningSecretValue Attribute Value",
				"While creating a SigningSecretValue value, a missing attribute value was detected. "+
					"A SigningSecretValue must contain values for all attributes, even if null or unknown. "+
					"This is always an issue with the provider and should be reported to the provider developers.\n\n"+
					fmt.Sprintf("SigningSecretValue Attribute Name (%s) Expected Type: %s", name, attributeType.String()),
			)

			continue
		}

		if !attributeType.Equal(attribute.Type(ctx)) {
			diags.AddError(
				"Invalid SigningSecretValue Attribute Type",
				"While creating a SigningSecretValue value, an invalid attribute value was detected. "+
					"A SigningSecretValue must use a matching attribute type for the value. "+
					"This is always an issue with the provider and should be reported to the provider developers.\n\n"+
					fmt.Sprintf("SigningSecretValue Attribute Name (%s) Expected Type: %s\n", name, attributeType.String())+
					fmt.Sprintf("SigningSecretValue Attribute Name (%s) Given Type: %s", name, attribute.Type(ctx)),
			)
		}
	}

	for name := range attributes {
		_, ok := attributeTypes[name]

		if !ok {
			diags.AddError(
				"Extra SigningSecretValue Attribute Value",
				"While creating a SigningSecretValue value, an extra attribute value was detected. "+
					"A SigningSecretValue must not contain values beyond the expected attribute types. "+
					"This is always an issue with the provider and should be reported to the provider developers.\n\n"+
					fmt.Sprintf("Extra SigningSecretValue Attribute Name: %s", name),
			)
		}
	}

	if diags.HasError() {
		return NewSigningSecretValueUnknown(), diags
	}

	algorithmAttribute, ok := attributes["algorithm"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`algorithm is missing from object`)

		return NewSigningSecretValueUnknown(), diags
	}

	algorithmVal, ok := algorithmAttribute.(basetypes.StringValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`algorithm expected to be basetypes.StringValue, was: %T`, algorithmAttribute))
	}

	secretAttribute, ok := attributes["secret"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`secret is missing from object`)

		return NewSigningSecretValueUnknown(), diags
	}

	secretVal, ok := secretAttribute.(basetypes.StringValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`secret expected to be basetypes.StringValue, was: %T`, secretAttribute))
	}

	if diags.HasError() {
		return NewSigningSecretValueUnknown(), diags
	}

	return SigningSecretValue{
		Algorithm: algorithmVal,
		Secret:    secretVal,
		state:     attr.ValueStateKnown,
	}, diags
}

func NewSigningSecretValueMust(attributeTypes map[string]attr.Type, attributes map[string]attr.Value) SigningSecretValue {
	object, diags := NewSigningSecretValue(attributeTypes, attributes)

	if diags.HasError() {
		// This could potentially be added to the diag package.
		diagsStrings := make([]string, 0, len(diags))

		for _, diagnostic := range diags {
			diagsStrings = append(diagsStrings, fmt.Sprintf(
				"%s | %s | %s",
				diagnostic.Severity(),
				diagnostic.Summary(),
				diagnostic.Detail()))
		}

		panic("NewSigningSecretValueMust received error(s): " + strings.Join(diagsStrings, "\n"))
	}

	return object
}

func (t SigningSecretType) ValueFromTerraform(ctx context.Context, in tftypes.Value) (attr.Value, error) {
	if in.Type() == nil {
		return NewSigningSecretValueNull(), nil
	}

	if !in.Type().Equal(t.TerraformType(ctx)) {
		return nil, fmt.Errorf("expected %s, got %s", t.TerraformType(ctx), in.Type())
	}

	if !in.IsKnown() {
		return NewSigningSecretValueUnknown(), nil
	}

	if in.IsNull() {
		return NewSigningSecretValueNull(), nil
	}

	attributes := map[string]attr.Value{}

	val := map[string]tftypes.Value{}

	err := in.As(&val)

	if err != nil {
		return nil, err
	}

	for k, v := range val {
		a, err := t.AttrTypes[k].ValueFromTerraform(ctx, v)

		if err != nil {
			return nil, err
		}

		attributes[k] = a
	}

	return NewSigningSecretValueMust(SigningSecretValue{}.AttributeTypes(ctx), attributes), nil
}

func (t SigningSecretType) ValueType(ctx context.Context) attr.Value {
	return SigningSecretValue{}
}

var _ basetypes.ObjectValuable = SigningSecretValue{}

type SigningSecretValue struct {
	Algorithm basetypes.StringValue `tfsdk:"algorithm"`
	Secret    basetypes.StringValue `tfsdk:"secret"`
	state     attr.ValueState
}

func (v SigningSecretValue) ToTerraformValue(ctx context.Context) (tftypes.Value, error) {
	attrTypes := make(map[string]tftypes.Type, 2)

	var val tftypes.Value
	var err error

	attrTypes["algorithm"] = basetypes.StringType{}.TerraformType(ctx)
	attrTypes["secret"] = basetypes.StringType{}.TerraformType(ctx)

	objectType := tftypes.Object{AttributeTypes: attrTypes}

	switch v.state {
	case attr.ValueStateKnown:
		vals := make(map[string]tftypes.Value, 2)

		val, err = v.Algorithm.ToTerraformValue(ctx)

		if err != nil {
			return tftypes.NewValue(objectType, tftypes.UnknownValue), err
		}

		vals["algorithm"] = val

		val, err = v.Secret.ToTerraformValue(ctx)

		if err != nil {
			return tftypes.NewValue(objectType, tftypes.UnknownValue), err
		}

		vals["secret"] = val

		if err := tftypes.ValidateValue(objectType, vals); err != nil {
			return tftypes.NewValue(objectType, tftypes.UnknownValue), err
		}

		return tftypes.NewValue(objectType, vals), nil
	case attr.ValueStateNull:
		return tftypes.NewValue(objectType, nil), nil
	case attr.ValueStateUnknown:
		return tftypes.NewValue(objectType, tftypes.UnknownValue), nil
	default:
		panic(fmt.Sprintf("unhandled Object state in ToTerraformValue: %s", v.state))
	}
}

func (v SigningSecretValue) IsNull() bool {
	return v.state == attr.ValueStateNull
}

func (v SigningSecretValue) IsUnknown() bool {
	return v.state == attr.ValueStateUnknown
}

func (v SigningSecretValue) String() string {
	return "SigningSecretValue"
}

func (v SigningSecretValue) ToObjectValue(ctx context.Context) (basetypes.ObjectValue, diag.Diagnostics) {
	var diags diag.Diagnostics

	objVal, diags := types.ObjectValue(
		map[string]attr.Type{
			"algorithm": basetypes.StringType{},
			"secret":    basetypes.StringType{},
		},
		map[string]attr.Value{
			"algorithm": v.Algorithm,
			"secret":    v.Secret,
		})

	return objVal, diags
}

func (v SigningSecretValue) Equal(o attr.Value) bool {
	other, ok := o.(SigningSecretValue)

	if !ok {
		return false
	}

	if v.state != other.state {
		return false
	}

	if v.state != attr.ValueStateKnown {
		return true
	}

	if !v.Algorithm.Equal(other.Algorithm) {
		return false
	}

	if !v.Secret.Equal(other.Secret) {
		return false
	}

	return true
}

func (v SigningSecretValue) Type(ctx context.Context) attr.Type {
	return SigningSecretType{
		basetypes.ObjectType{
			AttrTypes: v.AttributeTypes(ctx),
		},
	}
}

func (v SigningSecretValue) AttributeTypes(ctx context.Context) map[string]attr.Type {
	return map[string]attr.Type{
		"algorithm": basetypes.StringType{},
		"secret":    basetypes.StringType{},
	}
}
