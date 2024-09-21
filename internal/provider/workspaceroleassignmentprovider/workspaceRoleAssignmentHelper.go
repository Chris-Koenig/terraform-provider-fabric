package workspaceroleassignmentprovider

import (
	"terraform-provider-fabric/internal/fabricapi/fabricClientModels"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

var itemName = "workspaceroleassignment"
var apiItemName = "roleassignments"

func ConvertApiModelToTerraformModel(apiModel *fabricClientModels.WorkspaceRoleAssignmentReadModel, targetWorkspaceId string) WorkspaceRoleAssignmentProviderModel {
	var terraformModel WorkspaceRoleAssignmentProviderModel

	terraformModel.Id = types.StringValue(apiModel.Id)
	terraformModel.Workspace_Id = types.StringValue(targetWorkspaceId)
	terraformModel.Role = types.StringValue(apiModel.Role)

	// The principal must be instantiate
	terraformModel.Principal = &Principal{
		Id:   types.StringValue(apiModel.Principal.Id),
		Type: types.StringValue(apiModel.Principal.Type),
	}

	return terraformModel
}

func ConvertTerraformModelToApiCreateModel(terraformModel WorkspaceRoleAssignmentProviderModel) fabricClientModels.WorkspaceRoleAssignmentCreateRequestModel {
	return fabricClientModels.WorkspaceRoleAssignmentCreateRequestModel{
		Role: terraformModel.Role.ValueString(),
		Principal: fabricClientModels.Principal{
			Id:   terraformModel.Principal.Id.ValueString(),
			Type: terraformModel.Principal.Type.ValueString(),
		},
	}
}

func ConvertTerraformModelToApiUpdateModel(terraformModel WorkspaceRoleAssignmentProviderModel) fabricClientModels.WorkspaceRoleAssignmentUpdateRequestModel {
	return fabricClientModels.WorkspaceRoleAssignmentUpdateRequestModel{
		Role: terraformModel.Role.ValueString(),
		Principal: fabricClientModels.Principal{
			Id:   terraformModel.Principal.Id.ValueString(),
			Type: terraformModel.Principal.Type.ValueString(),
		},
	}
}
