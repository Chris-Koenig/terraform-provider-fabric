package workspaceRoleAssignmentProvider

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

// NewWorkspaceDataSource is a function that creates a new instance of the WorkspaceDataSource.
func NewWorkspaceRoleAssignmentDataSource() datasource.DataSource {
	return &WorkspaceRoleAssignmentDataSource{}
}

// WorkspaceDataSource is a struct that represents the Power BI workspace data source.
type WorkspaceRoleAssignmentDataSource struct {
	client *fabricapi.FabricClient
}

// Metadata is a method that sets the metadata for the WorkspaceDataSource.
func (d *WorkspaceRoleAssignmentDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_workspaceroleassignment"
}

// Schema is a method that sets the schema for the WorkspaceDataSource.
func (d *WorkspaceRoleAssignmentDataSource) Schema(_ context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{

		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Fabric workspace data source",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "ID (GUID) of the workspace role assignment. This value is read-only.",
				Computed:            false,
				Required:            true,
			},
			"workspace_id": schema.StringAttribute{
				MarkdownDescription: "ID (GUID) of the workspace. This is required.",
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
				MarkdownDescription: "Role assigned to the principal ('Member', 'Admin', 'Contributor', 'Viewer'). This is required.",
				Computed:            true,
			},
		},
	}
}

// ValidateConfig validates the configuration for the WorkspaceDataSource.
// It checks if either the 'name' or 'id' attribute is set, and adds an error to the response diagnostics if both are empty or both are non-empty.
// If there are any errors in the response diagnostics, the function returns without further processing.
// Parameters:
//   - ctx: The context.Context object for the request.
//   - req: The ValidateConfigRequest object containing the configuration to validate.
//   - resp: The ValidateConfigResponse object to store the validation results.
//
// Returns: None.
func (d *WorkspaceRoleAssignmentDataSource) ValidateConfig(ctx context.Context, req datasource.ValidateConfigRequest, resp *datasource.ValidateConfigResponse) {

}

// Configure is a method that configures the WorkspaceDataSource.
// It creates a new instance of the powerbiapi.Client with the specified base URL.
// If an error occurs while creating the client, an error is added to the response diagnostics.
// Parameters:
//   - ctx: The context.Context object for the request.
//   - req: The ConfigureRequest object containing the configuration to configure.
//   - resp: The ConfigureResponse object to store the configuration results.
//
// Returns: None.
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

// Read is a method that reads the data from the Power BI service and returns the result.
// It takes a ReadRequest and a ReadResponse as input and output parameters, respectively.
// It reads the data from the Power BI service and returns the result in the response.
// Parameters:
//   - ctx: The context.Context object for the request.
//   - req: The ReadRequest object containing the configuration to read.
//   - resp: The ReadResponse object to store the read results.
//
// Returns: None.
func (d *WorkspaceRoleAssignmentDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data, state WorkspaceRoleAssignmentProviderModel
	var roleAssignment *fabricClientModels.WorkspaceRoleAssignmentReadModel
	var err error

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if !data.Id.IsNull() {

		roleAssignment, err = fabricapi.GetItem[fabricClientModels.WorkspaceRoleAssignmentReadModel](data.Id.ValueString(), "roleassignments", data.Workspace_Id.ValueString(), *d.client)

		if err != nil {
			resp.Diagnostics.AddError(fmt.Sprintf("Cannot retrieve workspace with Id %s", data.Id), err.Error())
			return
		}
	}

	if roleAssignment == nil {
		resp.Diagnostics.AddError("No RoleAssignment found", "No workspace found with the specified name or id.")
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
