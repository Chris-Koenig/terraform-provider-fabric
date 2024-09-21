package workspaceroleassignmentprovider

import (
	"context"
	"fmt"

	"terraform-provider-fabric/internal/fabricapi"
	"terraform-provider-fabric/internal/fabricapi/fabricClientModels"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource                   = &WorkspaceRoleAssignmentDataSource{} // Ensure that WorkspaceDataSource implements the DataSource interface.
	_ datasource.DataSourceWithValidateConfig = &WorkspaceRoleAssignmentDataSource{} // Ensure that WorkspaceDataSource implements the DataSourceWithValidateConfig interface.
	_ datasource.DataSourceWithConfigure      = &WorkspaceRoleAssignmentDataSource{} // Ensure that WorkspaceDataSource implements the DataSourceWithConfigure interface.
)

// New is a function that creates a new instance of the Fabric DataSource.
func NewWorkspaceRoleAssignmentDataSource() datasource.DataSource {
	return &WorkspaceRoleAssignmentDataSource{}
}

// Struct that represents the Fabric DataSource.
type WorkspaceRoleAssignmentDataSource struct {
	client *fabricapi.FabricClient
}

// Metadata is a method that sets the metadata for the DataSource.
func (d *WorkspaceRoleAssignmentDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + itemName
}

// Schema is a method that sets the schema for the DataSource.
func (d *WorkspaceRoleAssignmentDataSource) Schema(_ context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{

		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Fabric " + itemName + " data source",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "ID (GUID) of the " + itemName,
				Computed:            false,
				Required:            true,
			},
			"workspace_id": schema.StringAttribute{
				MarkdownDescription: "Workspace ID (GUID) of the workspace.",
				Computed:            false,
				Required:            true,
			},
			"principal": schema.ObjectAttribute{

				Optional: true,
				AttributeTypes: map[string]attr.Type{
					"id":   types.StringType,
					"type": types.StringType,
				},
				Description: "The principal assigned to the role.",
			},
			"role": schema.StringAttribute{
				MarkdownDescription: "Role assigned to the principal ('Member', 'Admin', 'Contributor', 'Viewer').",
				Computed:            true,
			},
		},
	}
}

// ValidateConfig validates the configuration for the DataSource.
func (d *WorkspaceRoleAssignmentDataSource) ValidateConfig(ctx context.Context, req datasource.ValidateConfigRequest, resp *datasource.ValidateConfigResponse) {

}

// Configure is a method that configures the DataSource.
func (d *WorkspaceRoleAssignmentDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*fabricapi.FabricClient)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			"Expected *fabricapi.FabricClient, got: %T. Please report this issue to the provider developers.",
		)

		return
	}

	d.client = client
}

// Read is a method that reads the data from the Fabric service and returns the result.
func (d *WorkspaceRoleAssignmentDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data, state WorkspaceRoleAssignmentProviderModel
	var roleAssignment *fabricClientModels.WorkspaceRoleAssignmentReadModel
	var err error

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if !data.Id.IsNull() {

		roleAssignment, err = fabricapi.GetItem[fabricClientModels.WorkspaceRoleAssignmentReadModel](data.Id.ValueString(), apiItemName, data.Workspace_Id.ValueString(), *d.client)

		if err != nil {
			resp.Diagnostics.AddError(fmt.Sprintf("Cannot retrieve "+itemName+" with Id %s", data.Id), err.Error())
			return
		}
	}

	if roleAssignment == nil {
		resp.Diagnostics.AddError("No "+itemName+" found", "No "+itemName+" found with the specified id.")
		return
	}

	// Mapping the Values from the REST API to the provider model
	state = ConvertApiModelToTerraformModel(roleAssignment, data.Workspace_Id.ValueString())

	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}
