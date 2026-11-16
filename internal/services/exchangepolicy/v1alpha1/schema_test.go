package v1alpha1

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
)

func validateObject(v validator.Object, obj types.Object) []string {
	req := validator.ObjectRequest{
		Path:        path.Root("test"),
		ConfigValue: obj,
	}
	resp := &validator.ObjectResponse{}
	v.ValidateObject(context.Background(), req, resp)
	details := make([]string, 0, resp.Diagnostics.ErrorsCount())
	for _, d := range resp.Diagnostics.Errors() {
		details = append(details, d.Detail())
	}
	return details
}

func TestExactlyOneVariantValidator(t *testing.T) {
	attrTypes := map[string]attr.Type{
		"a": types.StringType,
		"b": types.StringType,
	}

	tests := []struct {
		name    string
		obj     types.Object
		wantErr string
	}{
		{
			name: "null object is allowed",
			obj:  types.ObjectNull(attrTypes),
		},
		{
			name: "unknown object is allowed",
			obj:  types.ObjectUnknown(attrTypes),
		},
		{
			name: "exactly one variant set",
			obj: types.ObjectValueMust(attrTypes, map[string]attr.Value{
				"a": types.StringValue("x"),
				"b": types.StringNull(),
			}),
		},
		{
			name: "no variant set",
			obj: types.ObjectValueMust(attrTypes, map[string]attr.Value{
				"a": types.StringNull(),
				"b": types.StringNull(),
			}),
			wantErr: "Exactly one thing variant must be set, but got 0.",
		},
		{
			name: "both variants set",
			obj: types.ObjectValueMust(attrTypes, map[string]attr.Value{
				"a": types.StringValue("x"),
				"b": types.StringValue("y"),
			}),
			wantErr: "Exactly one thing variant must be set, but got 2.",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := validateObject(exactlyOneVariantValidator{name: "thing"}, tt.obj)
			if tt.wantErr == "" {
				assert.Empty(t, got)
				return
			}
			assert.Equal(t, []string{tt.wantErr}, got)
		})
	}
}

func TestOAuthASValidator(t *testing.T) {
	oauthAs := func(grantType, issuerURL, tokenURL attr.Value) types.Object {
		return types.ObjectValueMust(oauthAsAttrTypes, map[string]attr.Value{
			"grant_type": grantType,
			"issuer_url": issuerURL,
			"token_url":  tokenURL,
			"audiences":  types.ListNull(types.StringType),
			"timeout":    types.Int64Null(),
		})
	}
	null := types.StringNull()
	grant := types.StringValue("client_credentials")
	url := types.StringValue("https://example.com")

	const missingURL = "At least one of issuer_url or token_url is required when oauth_as is set."

	tests := []struct {
		name    string
		obj     types.Object
		wantErr string
	}{
		{
			name: "null object is allowed",
			obj:  types.ObjectNull(oauthAsAttrTypes),
		},
		{
			name: "unknown object is allowed",
			obj:  types.ObjectUnknown(oauthAsAttrTypes),
		},
		{
			name: "issuer url only",
			obj:  oauthAs(grant, url, null),
		},
		{
			name: "token url only",
			obj:  oauthAs(grant, null, url),
		},
		{
			name: "both urls",
			obj:  oauthAs(grant, url, url),
		},
		{
			name: "unknown url counts as set",
			obj:  oauthAs(grant, types.StringUnknown(), null),
		},
		{
			name:    "neither url set",
			obj:     oauthAs(grant, null, null),
			wantErr: missingURL,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := validateObject(oauthAsValidator(), tt.obj)
			if tt.wantErr == "" {
				assert.Empty(t, got)
				return
			}
			assert.Equal(t, []string{tt.wantErr}, got)
		})
	}
}

func TestClientAssertionValidator(t *testing.T) {
	clientAssertion := func(issuerURL, audiences attr.Value) types.Object {
		return types.ObjectValueMust(clientAssertionAttrTypes, map[string]attr.Value{
			"client_id":  types.StringValue("client"),
			"issuer_url": issuerURL,
			"audiences":  audiences,
			"timeout":    types.Int64Null(),
		})
	}
	nullURL := types.StringNull()
	nullAud := types.ListNull(types.StringType)
	url := types.StringValue("https://example.com")
	aud := types.ListValueMust(types.StringType, []attr.Value{types.StringValue("https://example.com/token")})

	const missing = "At least one of issuer_url or audiences is required when client_assertion is set."

	tests := []struct {
		name    string
		obj     types.Object
		wantErr string
	}{
		{name: "null object is allowed", obj: types.ObjectNull(clientAssertionAttrTypes)},
		{name: "unknown object is allowed", obj: types.ObjectUnknown(clientAssertionAttrTypes)},
		{name: "issuer url only", obj: clientAssertion(url, nullAud)},
		{name: "audiences only", obj: clientAssertion(nullURL, aud)},
		{name: "issuer url and audiences non-null", obj: clientAssertion(url, aud)},
		{name: "unknown audiences count as non-null", obj: clientAssertion(nullURL, types.ListUnknown(types.StringType))},
		{name: "issuer url and audiences both null", obj: clientAssertion(nullURL, nullAud), wantErr: missing},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := validateObject(clientAssertionValidator(), tt.obj)
			if tt.wantErr == "" {
				assert.Empty(t, got)
				return
			}
			assert.Equal(t, []string{tt.wantErr}, got)
		})
	}
}
