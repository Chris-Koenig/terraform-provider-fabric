package workspaceRoleAssignmentProvider

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type Principal struct {
	Id   string `tfsdk:"id"`
	Type string `tfsdk:"type"`
}

type WorkspaceRoleAssignmentProviderModel struct {
	Id           types.String `tfsdk:"id"`
	Workspace_Id types.String `tfsdk:"workspace_id"`
	Principal    Principal    `tfsdk:"principal"`
	Role         string       `tfsdk:"role"`
}
