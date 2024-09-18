package workspaceprovider

import (
	"context"
	"fmt"

	"terraform-provider-fabric/internal/fabricapi"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource                   = &WorkspaceDataSource{} // Ensure that WorkspaceDataSource implements the DataSource interface.
	_ datasource.DataSourceWithValidateConfig = &WorkspaceDataSource{} // Ensure that WorkspaceDataSource implements the DataSourceWithValidateConfig interface.
	_ datasource.DataSourceWithConfigure      = &WorkspaceDataSource{} // Ensure that WorkspaceDataSource implements the DataSourceWithConfigure interface.
)

// NewWorkspaceDataSource is a function that creates a new instance of the WorkspaceDataSource.
func NewWorkspaceDataSource() datasource.DataSource {
	return &WorkspaceDataSource{}
}

// WorkspaceDataSource is a struct that represents the Power BI workspace data source.
type WorkspaceDataSource struct {
	client *fabricapi.FabricClient
}

// Metadata is a method that sets the metadata for the WorkspaceDataSource.
func (d *WorkspaceDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_workspace"
}

// Schema is a method that sets the schema for the WorkspaceDataSource.
func (d *WorkspaceDataSource) Schema(_ context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{

		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Fabric workspace data source",

		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				MarkdownDescription: "The name of the workspace",
				Optional:            true,
				Computed:            true,
			},
			"id": schema.StringAttribute{
				MarkdownDescription: "The id of the workspace",
				Optional:            true,
				Computed:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "the description of the workspace.",
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
func (d *WorkspaceDataSource) ValidateConfig(ctx context.Context, req datasource.ValidateConfigRequest, resp *datasource.ValidateConfigResponse) {

	var data WorkspaceProviderModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	bothNull := data.Name.IsNull() && data.Id.IsNull()
	bothSet := !data.Name.IsNull() && !data.Id.IsNull()

	if bothNull || bothSet {
		resp.Diagnostics.AddAttributeError(
			path.Root("name"),
			"Invalid attribute configuration",
			"one of 'name' or 'id' must be set",
		)
	}
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
func (d *WorkspaceDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
func (d *WorkspaceDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data WorkspaceProviderModel
	var workspace *fabricapi.WorkspaceReadModel
	var err error

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	if !data.Id.IsNull() {
		workspace, err = d.client.GetWorkspace(data.Id.ValueString())
		if err != nil {
			resp.Diagnostics.AddError(fmt.Sprintf("Cannot retrieve workspace with Id %s", data.Id), err.Error())
			return
		}
	}

	if !data.Name.IsNull() {
		workspaces, err := d.client.GetWorkspaces(fmt.Sprintf("name eq '%s'", data.Name), 0, 0)
		if err != nil {
			resp.Diagnostics.AddError(fmt.Sprintf("Cannot retrieve workspace with name %s", data.Name), err.Error())
			return
		}

		if len(workspaces.Value) == 0 {
			resp.Diagnostics.AddError(fmt.Sprintf("Cannot retrieve workspace with name %s", data.Name), "No groups found")
			return
		}

		if len(workspaces.Value) > 1 {
			resp.Diagnostics.AddError(fmt.Sprintf("Cannot retrieve workspace with name %s", data.Name), "Multiple groups found")
			return
		}

		workspace = &workspaces.Value[0]
	}

	if workspace == nil {
		resp.Diagnostics.AddError("No workspace found", "No workspace found with the specified name or id.")
		return
	}

	// Mapping the Values from the REST API to the provider model
	data.Id = types.StringValue(workspace.Id)
	data.Name = types.StringValue(workspace.DisplayName)
	data.Description = types.StringValue(workspace.Description)

	diags := resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}
