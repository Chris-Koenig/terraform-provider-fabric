package workspaceprovider

import (
	"context"
	"fmt"
	"terraform-provider-fabric/internal/fabricapi"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
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
	var workspaceCreated *fabricapi.WorkspaceReadModel
	var err error

	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)

	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Creating workspace with name: %s", config.Name.ValueString()))

	workspaceCreated, err = r.client.CreateWorkspace(config.Name.ValueString())

	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Cannot create workspace with name %s", config.Name.ValueString()), err.Error())
		return
	}
	tflog.Debug(ctx, "Workspace created successfully")

	tflog.Debug(ctx, "Populate the response with the workspace data")
	state.Id = types.StringValue(workspaceCreated.Id)
	state.Name = types.StringValue(workspaceCreated.DisplayName)
	state.CapacityId = types.StringValue("")
	state.Description = types.StringValue("")

	//state.IsOnDedicatedCapacity = types.BoolValue(workspaceCreated.IsOnDedicatedCapacity)

	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the Power BI workspace.
func (r *WorkspaceResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// var state models.Workspace
	// var err error

	// resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	// if resp.Diagnostics.HasError() {
	// 	return
	// }

	// tflog.Debug(ctx, fmt.Sprintf("Deleting workspace with name: %s", state.Name.ValueString()))
	// err = r.client.DeleteGroup(state.Id.ValueString())
	// if err != nil {
	// 	resp.Diagnostics.AddError(fmt.Sprintf("Cannot delete workspace with Id %s", state.Id.ValueString()), err.Error())
	// 	return
	// }

	// tflog.Debug(ctx, "Workspace deleted successfully")
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
	var workspace *fabricapi.WorkspaceReadModel
	var err error

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Reading workspace with name: %s", state.Name.ValueString()))
	workspace, err = r.client.GetWorkspace(state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Cannot retrieve workspace with Id %s", state.Id.ValueString()), err.Error())
		return
	}

	state.Id = types.StringValue(workspace.Id)
	state.Name = types.StringValue(workspace.DisplayName)

	//state.IsOnDedicatedCapacity = types.BoolValue(workspace.IsOnDedicatedCapacity)

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
				MarkdownDescription: "The name of the workspace",
				Optional:            false,
				Required:            true,
				Computed:            false,
			},
			"id": schema.StringAttribute{
				MarkdownDescription: "The id of the workspace",
				Computed:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "the description of the workspace.",
				Optional:            true,
				Required:            false,
				Computed:            false,
			},
			"capacity_id": schema.StringAttribute{
				MarkdownDescription: "the id of the capacity to which the workspace belongs.",
				Optional:            true,
				Required:            false,
				Computed:            false,
			},
		},
	}
}

// Update updates the Power BI workspace.
func (r *WorkspaceResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {

	// var plan models.Workspace
	// var state models.Workspace
	// var workspace *pbiModels.Group
	// var updateRequest *pbiModels.UpdateGroupRequest
	// var err error

	// resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	// resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	// if resp.Diagnostics.HasError() {
	// 	return
	// }

	// updateRequest = &pbiModels.UpdateGroupRequest{
	// 	Name: plan.Name.ValueString(),
	// }

	// tflog.Debug(ctx, fmt.Sprintf("Updating workspace with name: %s", state.Name.ValueString()))
	// err = r.client.UpdateGroup(state.Id.ValueString(), updateRequest)
	// if err != nil {
	// 	resp.Diagnostics.AddError(fmt.Sprintf("Cannot update workspace with Id %s", state.Id.ValueString()), err.Error())
	// 	return
	// }

	// tflog.Debug(ctx, "Workspace updated successfully")

	// tflog.Debug(ctx, "Populate the response with the workspace data")
	// tflog.Debug(ctx, fmt.Sprintf("Reading workspace with name: %s", plan.Name.ValueString()))
	// workspace, err = r.client.GetGroup(state.Id.ValueString())
	// if err != nil {
	// 	resp.Diagnostics.AddError(fmt.Sprintf("Cannot retrieve workspace with Id %s", state.Id.ValueString()), err.Error())
	// 	return
	// }

	// state.Id = types.StringValue(workspace.Id)
	// state.Name = types.StringValue(workspace.Name)
	// state.IsReadOnly = types.BoolValue(workspace.IsReadOnly)
	// state.IsOnDedicatedCapacity = types.BoolValue(workspace.IsOnDedicatedCapacity)

	// diags := resp.State.Set(ctx, &state)
	// resp.Diagnostics.Append(diags...)
	// if resp.Diagnostics.HasError() {
	// 	return
	// }
}
