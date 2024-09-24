package roleAssignmentProvider

import (
	"context"
	"fmt"
	"terraform-provider-fabric/internal/fabricapi"
	"terraform-provider-fabric/internal/fabricapi/fabricClientModels"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var _ resource.Resource = &WorkspaceRoleAssignmentResource{}                // Ensure that WorkspaceResource implements the Resource interface.
var _ resource.ResourceWithImportState = &WorkspaceRoleAssignmentResource{} // Ensure that WorkspaceResource implements the ResourceWithImportState interface.
var _ resource.ResourceWithConfigure = &WorkspaceRoleAssignmentResource{}   // Ensure that WorkspaceResource implements the ResourceWithConfigure interface.

// New is a function that creates a new instance of the Resource.
func NewWorkspaceRoleAssignmentResource() resource.Resource {
	return &WorkspaceRoleAssignmentResource{}
}

// Struct that represents the Fabric workspace resource.
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
			"Expected *fabricapi.Client, got: %T. Please report this issue to the provider developers.",
		)

		return
	}

	r.client = client
}

// ImportState implements resource.ResourceWithImportState.
func (r *WorkspaceRoleAssignmentResource) ImportState(_ context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	panic("unimplemented")
}

// Metadata sets the metadata for the resource.
func (r *WorkspaceRoleAssignmentResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + itemName
}

// Schema sets the schema for the resource.
func (r *WorkspaceRoleAssignmentResource) Schema(_ context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Fabric " + itemName + " data source",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "ID (GUID) of the " + itemName + ". This value is read-only.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"workspace_id": schema.StringAttribute{
				MarkdownDescription: "ID (GUID) of the " + itemName + ". This is required.",
				Required:            true,
				Optional:            false,
				Computed:            false,
			},
			"principal": schema.SingleNestedAttribute{
				MarkdownDescription: "Details of the principal (user or group) associated with the " + itemName + ".",
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

// Read updates the state with the data from the Fabric service.
func (r *WorkspaceRoleAssignmentResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state, newState WorkspaceRoleAssignmentProviderModel
	var workspaceRoleAssignmentCreated *fabricClientModels.WorkspaceRoleAssignmentReadModel
	var err error

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Reading "+itemName+" with ID: %s", state.Id.ValueString()))

	workspaceRoleAssignmentCreated, err = fabricapi.GetItem[fabricClientModels.WorkspaceRoleAssignmentReadModel](state.Id.ValueString(), apiItemName, state.Workspace_Id.ValueString(), *r.client)

	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Cannot retrieve "+itemName+" with Id %s", state.Id.ValueString()), err.Error())
		return
	}

	// Mapping the Values from the REST API to the provider model
	newState = ConvertApiModelToTerraformModel(workspaceRoleAssignmentCreated, state.Workspace_Id.ValueString())

	diags := resp.State.Set(ctx, &newState)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Create creates a new Fabric item.
func (r *WorkspaceRoleAssignmentResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan, state WorkspaceRoleAssignmentProviderModel
	var workspaceRoleAssignmentCreated *fabricClientModels.WorkspaceRoleAssignmentReadModel
	var err error

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)

	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Creating "+itemName+" with id: %s", plan.Id.ValueString()))

	var workspaceReleAssignmentToCreate = ConvertTerraformModelToApiCreateModel(plan)

	workspaceRoleAssignmentCreated, err = fabricapi.CreateItem[fabricClientModels.WorkspaceRoleAssignmentCreateRequestModel, fabricClientModels.WorkspaceRoleAssignmentReadModel](workspaceReleAssignmentToCreate, "roleAssignments", plan.Workspace_Id.ValueString(), *r.client)
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Cannot create "+itemName+" with id %s", plan.Id.ValueString()), err.Error())
		return
	}

	state = ConvertApiModelToTerraformModel(workspaceRoleAssignmentCreated, plan.Workspace_Id.ValueString())
	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the Fabric item.
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

	updateRequest = ConvertTerraformModelToApiUpdateModel(plan)

	tflog.Debug(ctx, fmt.Sprintf("Updating "+itemName+" with name: %s", state.Id.ValueString()))

	err = fabricapi.UpdateItem[fabricClientModels.WorkspaceRoleAssignmentUpdateRequestModel](state.Id.ValueString(), apiItemName, updateRequest, plan.Workspace_Id.ValueString(), *r.client)

	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Cannot update "+itemName+" with Id %s", state.Id.ValueString()), err.Error())
		return
	}

	workspaceRoleAssignmentUpdated, err = fabricapi.GetItem[fabricClientModels.WorkspaceRoleAssignmentReadModel](state.Id.ValueString(), apiItemName, state.Workspace_Id.ValueString(), *r.client)
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Cannot retrieve "+itemName+" with Id %s", state.Id.ValueString()), err.Error())
		return
	}

	state = ConvertApiModelToTerraformModel(workspaceRoleAssignmentUpdated, plan.Workspace_Id.ValueString())

	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the Fabric item.
func (r *WorkspaceRoleAssignmentResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state WorkspaceRoleAssignmentProviderModel
	var err error

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Deleting "+itemName+" with id: %s", state.Id.ValueString()))

	err = fabricapi.DeleteItem(state.Id.ValueString(), apiItemName, state.Workspace_Id.ValueString(), *r.client)
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Cannot delete "+itemName+" with Id %s", state.Id.ValueString()), err.Error())
		return
	}

	tflog.Debug(ctx, ""+itemName+" deleted successfully")
}
