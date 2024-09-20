package workspaceprovider

import (
	"terraform-provider-fabric/internal/fabricapi/fabricClientModels"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

func ConvertApiModelToTerraformModel(apiModel *fabricClientModels.WorkspaceReadModel) WorkspaceProviderModel {
	var terraformModel WorkspaceProviderModel

	terraformModel.Id = types.StringValue(apiModel.Id)
	terraformModel.Name = types.StringValue(apiModel.DisplayName)

	if apiModel.Description == "" {
		terraformModel.Description = types.StringValue("") //types.StringNull()
	} else {
		terraformModel.Description = types.StringValue(apiModel.Description)
	}

	return terraformModel
}

func ConvertTerraformModelToApiCreateModel(terraformModel WorkspaceProviderModel) fabricClientModels.WorkspaceCreateRequestModel {
	var apiCreateModel fabricClientModels.WorkspaceCreateRequestModel

	apiCreateModel.DisplayName = terraformModel.Name.ValueString()

	if terraformModel.Description == types.StringValue("") {
		apiCreateModel.Description = "" //types.StringNull() // Set to null
	} else if terraformModel.Description == types.StringNull() {
		apiCreateModel.Description = ""
	} else {
		apiCreateModel.Description = terraformModel.Description.ValueString()
	}

	return apiCreateModel
}

func ConvertTerraformModelToApiUpdateModel(terraformModel WorkspaceProviderModel) fabricClientModels.WorkspaceUpdateRequestModel {
	var apiCreateModel fabricClientModels.WorkspaceUpdateRequestModel

	apiCreateModel.DisplayName = terraformModel.Name.ValueString()

	if terraformModel.Description == types.StringValue("") {
		apiCreateModel.Description = "" //types.StringNull() // Set to null
	} else {
		apiCreateModel.Description = terraformModel.Description.ValueString()
	}

	return apiCreateModel
}
