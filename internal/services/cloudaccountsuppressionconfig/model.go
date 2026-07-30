package cloudaccountsuppressionconfig

import "github.com/hashicorp/terraform-plugin-framework/types"

type Model struct {
	ID             types.String `tfsdk:"id"`
	CloudAccountID types.String `tfsdk:"cloud_account_id"`
	Suppressed     types.Bool   `tfsdk:"suppressed"`
}
