package roleAssignmentProvider

import (
	"terraform-provider-fabric/internal/fabricapi/fabricClientModels"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

var itemName = "role_assignment"
var apiItemName = "roleassignments"

func ConvertApiModelToTerraformModel(apiModel *fabricClientModels.RoleAssignmentReadModel, targetWorkspaceId string) RoleAssignmentProviderModel {
	var terraformModel RoleAssignmentProviderModel

	terraformModel.Id = types.StringValue(apiModel.Id)
	terraformModel.Workspace_Id = types.StringValue(targetWorkspaceId)
	terraformModel.Role = types.StringValue(apiModel.Role)

	// The principal must be instantiate
	terraformModel.Principal = Principal{
		Id:   types.StringValue(apiModel.Principal.Id),
		Type: types.StringValue(apiModel.Principal.Type),
	}

	return terraformModel
}

func ConvertTerraformModelToApiCreateModel(terraformModel RoleAssignmentProviderModel) fabricClientModels.RoleAssignmentCreateRequestModel {
	return fabricClientModels.RoleAssignmentCreateRequestModel{
		Role: terraformModel.Role.ValueString(),
		Principal: fabricClientModels.Principal{
			Id:   terraformModel.Principal.Id.ValueString(),
			Type: terraformModel.Principal.Type.ValueString(),
		},
	}
}

func ConvertTerraformModelToApiUpdateModel(terraformModel RoleAssignmentProviderModel) fabricClientModels.RoleAssignmentUpdateRequestModel {
	return fabricClientModels.RoleAssignmentUpdateRequestModel{
		Role: terraformModel.Role.ValueString(),
		Principal: fabricClientModels.Principal{
			Id:   terraformModel.Principal.Id.ValueString(),
			Type: terraformModel.Principal.Type.ValueString(),
		},
	}
}
