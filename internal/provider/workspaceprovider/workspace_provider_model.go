package workspaceprovider

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Workspace is a struct that represents the workspace data model.
type WorkspaceProviderModel struct {
	Id          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	CapacityId  types.String `tfsdk:"capacity_id"`
}
