package workspaceprovider

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type WorkspaceProviderModel struct {
	Id          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
}
