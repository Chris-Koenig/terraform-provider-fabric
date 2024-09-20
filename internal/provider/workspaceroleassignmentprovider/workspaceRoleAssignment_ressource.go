package workspaceRoleAssignmentProvider

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

var _ resource.Resource = &WorkspaceRoleAssignmentResource{}                // Ensure that WorkspaceResource implements the Resource interface.
var _ resource.ResourceWithImportState = &WorkspaceRoleAssignmentResource{} // Ensure that WorkspaceResource implements the ResourceWithImportState interface.
var _ resource.ResourceWithConfigure = &WorkspaceRoleAssignmentResource{}   // Ensure that WorkspaceResource implements the ResourceWithConfigure interface.

// NewWorkspaceResource is a function that creates a new instance of the WorkspaceResource.
func NewWorkspaceRoleAssignmentResource() resource.Resource {
	return &WorkspaceRoleAssignmentResource{}
}

// WorkspaceResource is a struct that represents the Power BI workspace resource.
type WorkspaceRoleAssignmentResource struct {
	client *fabricapi.FabricClient
}

// Configure configures the WorkspaceResource.
func (r *WorkspaceRoleAssignmentResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

// ImportState implements resource.ResourceWithImportState.
func (r *WorkspaceRoleAssignmentResource) ImportState(_ context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	panic("unimplemented")
}

// Metadata sets the metadata for the WorkspaceResource.
func (r *WorkspaceRoleAssignmentResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_workspaceroleassignment"
}

// Schema sets the schema for the WorkspaceResource.
func (r *WorkspaceRoleAssignmentResource) Schema(_ context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Schema for Workspace Role Assignment",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "ID (GUID) of the workspace role assignment. This value is read-only.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"workspace_id": schema.StringAttribute{
				MarkdownDescription: "ID (GUID) of the workspace. This is required.",
				Required:            true,
				Optional:            false,
				Computed:            false,
			},
			"principal": schema.SingleNestedAttribute{
				MarkdownDescription: "Details of the principal (user or group) associated with the workspace role assignment.",
				Required:            true,
				Attributes: map[string]schema.Attribute{
					"id": schema.StringAttribute{
						MarkdownDescription: "ID of the principal. This is required.",
						Required:            true,
					},
					"type": schema.StringAttribute{
						MarkdownDescription: "Type of the principal (e.g., 'User', 'Group'). This is required.",
						Required:            true,
					},
				},
			},
			"role": schema.StringAttribute{
				MarkdownDescription: "Role assigned to the principal ('Member', 'Admin', 'Contributor', 'Viewer'). This is required.",
				Required:            true,
			},
		},
	}
}

// Read updates the state with the data from the Power BI service.
func (r *WorkspaceRoleAssignmentResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state WorkspaceRoleAssignmentProviderModel
	var workspaceRoleAssignmentCreated *fabricClientModels.WorkspaceRoleAssignmentReadModel
	var err error

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Reading workspace with ID: %s", state.Id.ValueString()))

	workspaceRoleAssignmentCreated, err = fabricapi.GetItem[fabricClientModels.WorkspaceRoleAssignmentReadModel](state.Id.ValueString(), "roleassignments", state.Workspace_Id.ValueString(), *r.client)

	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Cannot retrieve workspace with Id %s", state.Id.ValueString()), err.Error())
		return
	}

	state.Id = types.StringValue(workspaceRoleAssignmentCreated.Id)
	state.Role = workspaceRoleAssignmentCreated.Role
	state.Principal.Id = workspaceRoleAssignmentCreated.Principal.Id
	state.Principal.Type = workspaceRoleAssignmentCreated.Principal.Type

	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Create creates a new Power BI workspace.
func (r *WorkspaceRoleAssignmentResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {

	var config WorkspaceRoleAssignmentProviderModel
	var state WorkspaceRoleAssignmentProviderModel
	var workspaceRoleAssignmentCreated *fabricClientModels.WorkspaceRoleAssignmentReadModel
	var err error

	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)

	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Creating roleassignments with id: %s", config.Id.ValueString()))

	var workspaceReleAssignmentToCreate = fabricClientModels.WorkspaceRoleAssignmentCreateRequestModel{
		Role:      config.Role,
		Principal: fabricClientModels.Principal(config.Principal),
	}

	workspaceRoleAssignmentCreated, err = fabricapi.CreateItem[fabricClientModels.WorkspaceRoleAssignmentCreateRequestModel, fabricClientModels.WorkspaceRoleAssignmentReadModel](workspaceReleAssignmentToCreate, "roleAssignments", config.Workspace_Id.ValueString(), *r.client)

	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Cannot create roleassignments with id %s", config.Id.ValueString()), err.Error())
		return
	}
	tflog.Debug(ctx, "Workspace created successfully")

	tflog.Debug(ctx, "Populate the response with the workspace data")

	state.Id = types.StringValue(workspaceRoleAssignmentCreated.Id)
	state.Role = workspaceRoleAssignmentCreated.Role
	state.Principal.Id = workspaceRoleAssignmentCreated.Principal.Id
	state.Principal.Type = workspaceRoleAssignmentCreated.Principal.Type
	state.Workspace_Id = config.Workspace_Id

	diags := resp.State.Set(ctx, &state)

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the Power BI workspace.
func (r *WorkspaceRoleAssignmentResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {

	var plan, state WorkspaceRoleAssignmentProviderModel
	var workspaceRoleAssignmentUpdated *fabricClientModels.WorkspaceRoleAssignmentReadModel
	var updateRequest fabricClientModels.WorkspaceRoleAssignmentUpdateRequestModel
	var err error

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	updateRequest = fabricClientModels.WorkspaceRoleAssignmentUpdateRequestModel{
		Principal: fabricClientModels.Principal(plan.Principal),
		Role:      plan.Role,
	}

	tflog.Debug(ctx, fmt.Sprintf("Updating roleAssignments with name: %s", state.Id.ValueString()))

	err = fabricapi.UpdateItem[fabricClientModels.WorkspaceRoleAssignmentUpdateRequestModel](state.Id.ValueString(), "roleAssignments", updateRequest, plan.Workspace_Id.ValueString(), *r.client)

	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Cannot update workspace with Id %s", state.Id.ValueString()), err.Error())
		return
	}

	tflog.Debug(ctx, "Workspace updated successfully")
	tflog.Debug(ctx, "Populate the response with the workspace data")
	tflog.Debug(ctx, fmt.Sprintf("Reading workspace with name: %s", plan.Id.ValueString()))

	workspaceRoleAssignmentUpdated, err = fabricapi.GetItem[fabricClientModels.WorkspaceRoleAssignmentReadModel](state.Id.ValueString(), "roleAssignments", state.Workspace_Id.ValueString(), *r.client)
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Cannot retrieve workspace with Id %s", state.Id.ValueString()), err.Error())
		return
	}

	state.Id = types.StringValue(workspaceRoleAssignmentUpdated.Id)
	state.Role = workspaceRoleAssignmentUpdated.Role
	state.Principal.Id = workspaceRoleAssignmentUpdated.Principal.Id
	state.Principal.Type = workspaceRoleAssignmentUpdated.Principal.Type
	state.Workspace_Id = plan.Workspace_Id

	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the Power BI workspace.
func (r *WorkspaceRoleAssignmentResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state WorkspaceRoleAssignmentProviderModel
	var err error

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Deleting roleassignments with id: %s", state.Id.ValueString()))

	err = fabricapi.DeleteItem(state.Id.ValueString(), "roleassignments", state.Workspace_Id.ValueString(), *r.client)
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Cannot delete roleassignments with Id %s", state.Id.ValueString()), err.Error())
		return
	}

	tflog.Debug(ctx, "roleAssignments deleted successfully")
}
