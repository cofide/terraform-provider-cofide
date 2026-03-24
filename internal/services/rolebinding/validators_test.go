package rolebinding

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// roleBindingTFType is the tftypes representation of the role binding schema.
var roleBindingTFType = tftypes.Object{
	AttributeTypes: map[string]tftypes.Type{
		"id":      tftypes.String,
		"role_id": tftypes.String,
		"user": tftypes.Object{
			AttributeTypes: map[string]tftypes.Type{
				"subject": tftypes.String,
			},
		},
		"group": tftypes.Object{
			AttributeTypes: map[string]tftypes.Type{
				"claim_value": tftypes.String,
			},
		},
		"resource": tftypes.Object{
			AttributeTypes: map[string]tftypes.Type{
				"type": tftypes.String,
				"id":   tftypes.String,
			},
		},
	},
}

var (
	userTFType  = roleBindingTFType.AttributeTypes["user"]
	groupTFType = roleBindingTFType.AttributeTypes["group"]

	nullUser  = tftypes.NewValue(userTFType, nil)
	nullGroup = tftypes.NewValue(groupTFType, nil)
)

func knownUserVal(subject string) tftypes.Value {
	return tftypes.NewValue(userTFType, map[string]tftypes.Value{
		"subject": tftypes.NewValue(tftypes.String, subject),
	})
}

func unknownUserVal() tftypes.Value {
	return tftypes.NewValue(userTFType, map[string]tftypes.Value{
		"subject": tftypes.NewValue(tftypes.String, tftypes.UnknownValue),
	})
}

func knownGroupVal(claimValue string) tftypes.Value {
	return tftypes.NewValue(groupTFType, map[string]tftypes.Value{
		"claim_value": tftypes.NewValue(tftypes.String, claimValue),
	})
}

func unknownGroupVal() tftypes.Value {
	return tftypes.NewValue(groupTFType, map[string]tftypes.Value{
		"claim_value": tftypes.NewValue(tftypes.String, tftypes.UnknownValue),
	})
}

func makeValidateRequest(t *testing.T, user, group tftypes.Value) resource.ValidateConfigRequest {
	t.Helper()
	raw := tftypes.NewValue(roleBindingTFType, map[string]tftypes.Value{
		"id":      tftypes.NewValue(tftypes.String, tftypes.UnknownValue),
		"role_id": tftypes.NewValue(tftypes.String, "admin"),
		"user":    user,
		"group":   group,
		"resource": tftypes.NewValue(roleBindingTFType.AttributeTypes["resource"], map[string]tftypes.Value{
			"type": tftypes.NewValue(tftypes.String, "Organization"),
			"id":   tftypes.NewValue(tftypes.String, "org-id"),
		}),
	})
	return resource.ValidateConfigRequest{
		Config: tfsdk.Config{
			Schema: resourceSchema(),
			Raw:    raw,
		},
	}
}

func TestOneOfValidator(t *testing.T) {
	tests := []struct {
		name           string
		user           tftypes.Value
		group          tftypes.Value
		wantErrCount   int
		wantErrSummary string
	}{
		{
			name:         "user set with known subject passes validation",
			user:         knownUserVal("user-subject-1"),
			group:        nullGroup,
			wantErrCount: 0,
		},
		{
			name:         "group set with known claim value passes validation",
			user:         nullUser,
			group:        knownGroupVal("group-claim-1"),
			wantErrCount: 0,
		},
		{
			name:           "both user and group set returns conflict errors",
			user:           knownUserVal("user-subject-1"),
			group:          knownGroupVal("group-claim-1"),
			wantErrCount:   2,
			wantErrSummary: "Conflicting Attributes",
		},
		{
			name:           "neither user nor group set returns missing errors",
			user:           nullUser,
			group:          nullGroup,
			wantErrCount:   2,
			wantErrSummary: "Missing Required Attribute",
		},
		{
			// Regression: validator must defer when subject is unknown, which occurs
			// during the for_each pre-expansion validation pass where each.value is
			// not yet resolved to a concrete value.
			name:         "user set with unknown subject defers validation",
			user:         unknownUserVal(),
			group:        nullGroup,
			wantErrCount: 0,
		},
		{
			// Regression: same deferral required when group claim_value is unknown.
			name:         "group set with unknown claim_value defers validation",
			user:         nullUser,
			group:        unknownGroupVal(),
			wantErrCount: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := &oneOfValidator{}
			req := makeValidateRequest(t, tt.user, tt.group)
			resp := &resource.ValidateConfigResponse{}

			v.ValidateResource(context.Background(), req, resp)

			require.Len(t, resp.Diagnostics, tt.wantErrCount)
			if tt.wantErrSummary != "" {
				for _, d := range resp.Diagnostics {
					assert.Equal(t, tt.wantErrSummary, d.Summary())
				}
			}
		})
	}
}
