package workspaceProvider

import (
	"terraform-provider-fabric/internal/fabricapi/fabricClientModels"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

var itemName = "workspace"
var apiItemName = "workspaces"

func ConvertApiModelToTerraformModel(apiModel *fabricClientModels.WorkspaceReadModel) WorkspaceProviderModel {
	var terraformModel WorkspaceProviderModel

	terraformModel.Id = types.StringValue(apiModel.Id)
	terraformModel.Name = types.StringValue(apiModel.DisplayName)

	if apiModel.Description == "" {
		terraformModel.Description = types.StringNull()
	} else {
		terraformModel.Description = types.StringValue(apiModel.Description)
	}

	if apiModel.FabricCapacityId == "" {

		terraformModel.FabricCapacityId = types.StringNull()
	} else {
		terraformModel.FabricCapacityId = types.StringValue(apiModel.FabricCapacityId)
	}
	return terraformModel
}

func ConvertTerraformModelToApiCreateModel(terraformModel WorkspaceProviderModel) fabricClientModels.WorkspaceCreateRequestModel {
	var apiCreateModel fabricClientModels.WorkspaceCreateRequestModel
	apiCreateModel.DisplayName = terraformModel.Name.ValueString()
	apiCreateModel.Description = terraformModel.Description.ValueString()
	return apiCreateModel
}

func ConvertTerraformModelToApiUpdateModel(terraformModel WorkspaceProviderModel) fabricClientModels.WorkspaceUpdateRequestModel {
	var apiCreateModel fabricClientModels.WorkspaceUpdateRequestModel
	apiCreateModel.DisplayName = terraformModel.Name.ValueString()
	apiCreateModel.Description = terraformModel.Description.ValueString()
	return apiCreateModel
}
