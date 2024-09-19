package workspaceprovider

import (
	"context"
	"fmt"
	"terraform-provider-fabric/internal/fabricapi"
	"terraform-provider-fabric/internal/fabricapi/fabricClientModels"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var _ resource.Resource = &WorkspaceResource{}                // Ensure that WorkspaceResource implements the Resource interface.
var _ resource.ResourceWithImportState = &WorkspaceResource{} // Ensure that WorkspaceResource implements the ResourceWithImportState interface.
var _ resource.ResourceWithConfigure = &WorkspaceResource{}   // Ensure that WorkspaceResource implements the ResourceWithConfigure interface.

// NewWorkspaceResource is a function that creates a new instance of the WorkspaceResource.
func NewWorkspaceResource() resource.Resource {
	return &WorkspaceResource{}
}

// WorkspaceResource is a struct that represents the Power BI workspace resource.
type WorkspaceResource struct {
	client *fabricapi.FabricClient
}

// Configure configures the WorkspaceResource.
func (r *WorkspaceResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*fabricapi.FabricClient)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			"Expected *powerbiapi.Client, got: %T. Please report this issue to the provider developers.",
		)

		return
	}

	r.client = client
}

// Create creates a new Power BI workspace.
func (r *WorkspaceResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {

	var config WorkspaceProviderModel
	var state WorkspaceProviderModel
	var workspaceCreated *fabricClientModels.WorkspaceReadModel
	var err error

	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)

	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Creating workspace with name: %s", config.Name.ValueString()))

	var workspaceToCreate = fabricClientModels.WorkspaceCreateRequestModel{
		Description: config.Description.ValueString(),
		DisplayName: config.Name.ValueString(),
	}

	// workspaceCreated, err = r.client.CreateWorkspace(workspaceToCreate)

	workspaceCreated, err = fabricapi.CreateItem[fabricClientModels.WorkspaceCreateRequestModel, fabricClientModels.WorkspaceReadModel](workspaceToCreate, "workspaces", "", *r.client)

	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Cannot create workspace with name %s", config.Name.ValueString()), err.Error())
		return
	}
	tflog.Debug(ctx, "Workspace created successfully")

	tflog.Debug(ctx, "Populate the response with the workspace data")

	state.Id = types.StringValue(workspaceCreated.Id)
	state.Name = types.StringValue(workspaceCreated.DisplayName)
	state.Description = types.StringValue(workspaceCreated.Description)

	diags := resp.State.Set(ctx, &state)

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the Power BI workspace.
func (r *WorkspaceResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {

	var state WorkspaceProviderModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)

	var err error

	if resp.Diagnostics.HasError() {
		return
	}

	// tflog.Debug(ctx, fmt.Sprintf("Deleting workspace with name: %s", entitityToDelete.Name.ValueString()))

	err = fabricapi.DeleteItem(state.Id.ValueString(), "workspaces", "", *r.client)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting Workspace",
			"Could not delete, unexpected error: "+err.Error(),
		)
		return
	}

}

// ImportState implements resource.ResourceWithImportState.
func (r *WorkspaceResource) ImportState(_ context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	panic("unimplemented")
}

// Metadata sets the metadata for the WorkspaceResource.
func (r *WorkspaceResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_workspace"
}

// Read updates the state with the data from the Power BI service.
func (r *WorkspaceResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state WorkspaceProviderModel
	var workspace *fabricClientModels.WorkspaceReadModel
	var err error

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Reading workspace with name: %s", state.Name.ValueString()))

	workspace, err = fabricapi.GetItem[fabricClientModels.WorkspaceReadModel](state.Id.ValueString(), "workspaces", "", *r.client)

	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Cannot retrieve workspace with Id %s", state.Id.ValueString()), err.Error())
		return
	}

	state.Id = types.StringValue(workspace.Id)
	state.Name = types.StringValue(workspace.DisplayName)
	state.Description = types.StringValue(workspace.Description)

	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Schema sets the schema for the WorkspaceResource.
func (r *WorkspaceResource) Schema(_ context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{

		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Fabric workspace data source",

		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				MarkdownDescription: "name of the workspace. Max Size is 200 characters.",
				Optional:            false,
				Required:            true,
				Computed:            false,
			},
			"id": schema.StringAttribute{
				MarkdownDescription: "ID (GUID) of the workspace. This value is read-only and will be set by the fabric service itself.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "description of the workspace. Max size is 2000 characters",
				Optional:            true,
				Required:            false,
				Computed:            false,
			},
		},
	}
}

// Update updates the Power BI workspace.
func (r *WorkspaceResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {

	var plan WorkspaceProviderModel
	var state WorkspaceProviderModel
	var workspaceUpdated *fabricClientModels.WorkspaceReadModel
	var updateRequest fabricClientModels.WorkspaceUpdateRequestModel
	var err error

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	updateRequest = fabricClientModels.WorkspaceUpdateRequestModel{
		DisplayName: plan.Name.ValueString(),
		Description: plan.Description.ValueString(),
	}

	tflog.Debug(ctx, fmt.Sprintf("Updating workspace with name: %s", state.Name.ValueString()))

	err = fabricapi.UpdateItem[fabricClientModels.WorkspaceUpdateRequestModel](state.Id.ValueString(), "workspaces", updateRequest, "", *r.client)

	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Cannot update workspace with Id %s", state.Id.ValueString()), err.Error())
		return
	}

	tflog.Debug(ctx, "Workspace updated successfully")
	tflog.Debug(ctx, "Populate the response with the workspace data")
	tflog.Debug(ctx, fmt.Sprintf("Reading workspace with name: %s", plan.Name.ValueString()))

	workspaceUpdated, err = fabricapi.GetItem[fabricClientModels.WorkspaceReadModel](state.Id.ValueString(), "workspaces", "", *r.client)
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Cannot retrieve workspace with Id %s", state.Id.ValueString()), err.Error())
		return
	}

	state.Id = types.StringValue(workspaceUpdated.Id)
	state.Name = types.StringValue(workspaceUpdated.DisplayName)
	state.Description = types.StringValue(workspaceUpdated.Description)

	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
