package workspaceProvider

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type WorkspaceProviderModel struct {
	Id               types.String `tfsdk:"id"`
	Name             types.String `tfsdk:"name"`
	Description      types.String `tfsdk:"description"`
	FabricCapacityId types.String `tfsdk:"fabric_capacity_id"`
}
