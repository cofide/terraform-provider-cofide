package cluster

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type optionalComputedModifier struct{}

var _ planmodifier.String = optionalComputedModifier{}

func (m optionalComputedModifier) Description(_ context.Context) string {
	return "Handles optional+computed attributes: preserves state on update if config is removed, and marks as unknown on create if not configured."
}

func (m optionalComputedModifier) MarkdownDescription(ctx context.Context) string {
	return m.Description(ctx)
}

func (m optionalComputedModifier) PlanModifyString(ctx context.Context, req planmodifier.StringRequest, resp *planmodifier.StringResponse) {
	// If the user is not setting a value in config, we can't know the final
	// value until apply. It could be a new value from the API, or the
	// existing state value. Mark it as unknown.
	if req.ConfigValue.IsNull() {
		resp.PlanValue = types.StringUnknown()
	}
}
