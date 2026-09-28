package v1alpha1

import (
	"testing"

	trustzonepb "github.com/cofide/cofide-api-sdk/gen/go/proto/trust_zone/v1alpha1"
	tftypes "github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"
)

func TestNewUpdateRequest(t *testing.T) {
	t.Parallel()

	state := TrustZoneModel{
		ID:                    tftypes.StringValue("tz-123"),
		Name:                  tftypes.StringValue("old-name"),
		TrustDomain:           tftypes.StringValue("example.cofide.dev"),
		OrgID:                 tftypes.StringValue("org-123"),
		IsManagementZone:      tftypes.BoolValue(false),
		BundleEndpointURL:     tftypes.StringValue("https://bundle.example.com"),
		BundleEndpointProfile: tftypes.StringValue("BUNDLE_ENDPOINT_PROFILE_HTTPS_WEB"),
		JWTIssuer:             tftypes.StringValue("https://issuer.example.com"),
	}

	tests := []struct {
		name    string
		state   func(TrustZoneModel) TrustZoneModel
		want    func(*trustzonepb.TrustZone)
		wantErr bool
	}{
		{
			name: "carries computed fields over from state",
		},
		{
			name: "omits empty computed fields",
			state: func(s TrustZoneModel) TrustZoneModel {
				s.BundleEndpointURL = tftypes.StringValue("")
				s.BundleEndpointProfile = tftypes.StringNull()
				s.JWTIssuer = tftypes.StringUnknown()
				return s
			},
			want: func(tz *trustzonepb.TrustZone) {
				tz.BundleEndpointUrl = nil
				tz.BundleEndpointProfile = nil
				tz.JwtIssuer = nil
			},
		},
		{
			name: "rejects unknown bundle endpoint profile",
			state: func(s TrustZoneModel) TrustZoneModel {
				s.BundleEndpointProfile = tftypes.StringValue("BOGUS")
				return s
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			s := state
			if tt.state != nil {
				s = tt.state(s)
			}
			plan := s
			plan.Name = tftypes.StringValue("new-name")

			got, err := newUpdateRequest(plan, s)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)

			want := &trustzonepb.TrustZone{
				Id:                    proto.String("tz-123"),
				Name:                  "new-name",
				TrustDomain:           "example.cofide.dev",
				OrgId:                 proto.String("org-123"),
				BundleEndpointUrl:     proto.String("https://bundle.example.com"),
				BundleEndpointProfile: trustzonepb.BundleEndpointProfile_BUNDLE_ENDPOINT_PROFILE_HTTPS_WEB.Enum(),
				JwtIssuer:             proto.String("https://issuer.example.com"),
			}
			if tt.want != nil {
				tt.want(want)
			}
			assert.True(t, proto.Equal(want, got), "want %v, got %v", want, got)
		})
	}
}
