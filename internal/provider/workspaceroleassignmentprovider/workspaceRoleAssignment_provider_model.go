package workspaceRoleAssignmentProvider

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type Principal struct {
	Id   types.String `tfsdk:"id"`
	Type types.String `tfsdk:"type"`
}

type WorkspaceRoleAssignmentProviderModel struct {
	Id           types.String `tfsdk:"id"`
	Workspace_Id types.String `tfsdk:"workspace_id"`
	Principal    Principal    `tfsdk:"principal"`
	Role         types.String `tfsdk:"role"`
}
