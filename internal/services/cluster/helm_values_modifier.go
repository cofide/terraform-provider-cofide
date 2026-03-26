package cluster

import (
	"github.com/cofide/terraform-provider-cofide/internal/util"
	tftypes "github.com/hashicorp/terraform-plugin-framework/types"
	"google.golang.org/protobuf/types/known/structpb"
)

// helmValuesForState returns the tftypes.String value to store in state for
// extra_helm_values after a Read. If the existing state value is semantically
// equivalent to the API response, it is preserved unchanged so that the
// original user-provided format (YAML or JSON) is not replaced with JSON.
func helmValuesForState(apiValues *structpb.Struct, stateValue tftypes.String) (tftypes.String, error) {
	return util.HelmValuesForState(apiValues, stateValue)
}
