package workspaceprovider

import (
	"context"
	"fmt"

	"terraform-provider-fabric/internal/fabricapi"
	"terraform-provider-fabric/internal/fabricapi/fabricClientModels"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
)

var (
	_ datasource.DataSource                   = &WorkspaceDataSource{} // Ensure that WorkspaceDataSource implements the DataSource interface.
	_ datasource.DataSourceWithValidateConfig = &WorkspaceDataSource{} // Ensure that WorkspaceDataSource implements the DataSourceWithValidateConfig interface.
	_ datasource.DataSourceWithConfigure      = &WorkspaceDataSource{} // Ensure that WorkspaceDataSource implements the DataSourceWithConfigure interface.
)

// New is a function that creates a new instance of the Fabric DataSource.
func NewWorkspaceDataSource() datasource.DataSource {
	return &WorkspaceDataSource{}
}

// Struct that represents the Fabric DataSource.
type WorkspaceDataSource struct {
	client *fabricapi.FabricClient
}

// Metadata is a method that sets the metadata for the DataSource.
func (d *WorkspaceDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + itemName
}

// Schema is a method that sets the schema for the Source.
func (d *WorkspaceDataSource) Schema(_ context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{

		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Fabric " + itemName + " data source",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "ID (GUID) of the " + itemName,
				Optional:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "name of the " + itemName,
				Computed:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "description of the " + itemName,
				Computed:            true,
			},
			"fabric_capacity_id": schema.StringAttribute{
				MarkdownDescription: "ID (GUID) of the Fabric Capacity of the " + itemName + ".",
				Optional:            false,
				Required:            true,
			},
		},
	}
}

// ValidateConfig validates the configuration for the DataSource.
func (d *WorkspaceDataSource) ValidateConfig(ctx context.Context, req datasource.ValidateConfigRequest, resp *datasource.ValidateConfigResponse) {

}

// Configure is a method that configures the DataSource.
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

// Read is a method that reads the data from the Fabric service and returns the result.
func (d *WorkspaceDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config, state WorkspaceProviderModel
	var workspace *fabricClientModels.WorkspaceReadModel
	var err error

	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if !config.Id.IsNull() {

		workspace, err = fabricapi.GetItem[fabricClientModels.WorkspaceReadModel](config.Id.ValueString(), apiItemName, "", *d.client)
		if err != nil {
			resp.Diagnostics.AddError(fmt.Sprintf("Cannot retrieve "+itemName+" with Id %s", config.Id), err.Error())
			return
		}
	}

	if workspace == nil {
		resp.Diagnostics.AddError("No "+itemName+" found", "No "+itemName+" found with the specified id.")
		return
	}

	state = ConvertApiModelToTerraformModel(workspace)

	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
