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
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var _ resource.Resource = &WorkspaceResource{}                // Ensure that WorkspaceResource implements the Resource interface.
var _ resource.ResourceWithImportState = &WorkspaceResource{} // Ensure that WorkspaceResource implements the ResourceWithImportState interface.
var _ resource.ResourceWithConfigure = &WorkspaceResource{}   // Ensure that WorkspaceResource implements the ResourceWithConfigure interface.

// New is a function that creates a new instance of the Resource.
func NewWorkspaceResource() resource.Resource {
	return &WorkspaceResource{}
}

// Struct that represents the Fabric workspace resource.
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
			"Expected *fabricapi.Client, got: %T. Please report this issue to the provider developers.",
		)

		return
	}

	r.client = client
}

// ImportState implements resource.ResourceWithImportState.
func (r *WorkspaceResource) ImportState(_ context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	panic("unimplemented")
}

// Metadata sets the metadata for the resource.
func (r *WorkspaceResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + itemName
}

// Schema sets the schema for the Resource.
func (r *WorkspaceResource) Schema(_ context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {

	resp.Schema = schema.Schema{

		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Fabric " + itemName + " data source",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "ID (GUID) of the " + itemName + ". This value is read-only and will be set by the fabric service itself.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "name of the " + itemName + ". Max Size is 200 characters.",
				Optional:            false,
				Required:            true,
				Computed:            false,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "description of the " + itemName + ". Max size is 2000 characters",
				Optional:            false,
				Required:            true,
			},
			"fabric_capacity_id": schema.StringAttribute{
				MarkdownDescription: "ID (GUID) of the Fabric Capacity of the " + itemName + ".",
				Optional:            false,
				Required:            true,
			},
		},
	}
}

// Read updates the state with the data from the Fabric service.
func (r *WorkspaceResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var stateCurrent, stateNew WorkspaceProviderModel
	var workspace *fabricClientModels.WorkspaceReadModel
	var err error

	resp.Diagnostics.Append(req.State.Get(ctx, &stateCurrent)...)

	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Reading "+itemName+" with name: %s", stateCurrent.Name.ValueString()))

	workspace, err = fabricapi.GetItem[fabricClientModels.WorkspaceReadModel](stateCurrent.Id.ValueString(), apiItemName, "", *r.client)

	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Cannot retrieve "+itemName+" with Id %s", stateCurrent.Id.ValueString()), err.Error())
		return
	}

	// convert the api data model in the terraform data model
	stateNew = ConvertApiModelToTerraformModel(workspace)

	diags := resp.State.Set(ctx, &stateNew)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Create creates a new Fabric item.
func (r *WorkspaceResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan, state WorkspaceProviderModel
	var workspaceCreated *fabricClientModels.WorkspaceReadModel
	var err error

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)

	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Creating "+itemName+" with id: %s", plan.Id.ValueString()))

	var workspaceToCreate = ConvertTerraformModelToApiCreateModel(plan)

	workspaceCreated, err = fabricapi.CreateWorkspace(workspaceToCreate, plan.FabricCapacityId.ValueStringPointer(), *r.client)

	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Cannot create "+itemName+" with name %s", plan.Name.ValueString()), err.Error())
		return
	}

	// convert the api data model in the terraform data model
	state = ConvertApiModelToTerraformModel(workspaceCreated)

	diags := resp.State.Set(ctx, &state)

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the Fabric item.
func (r *WorkspaceResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {

	var plan, state WorkspaceProviderModel

	var workspaceUpdated *fabricClientModels.WorkspaceReadModel
	var updateRequest fabricClientModels.WorkspaceUpdateRequestModel
	var err error

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	updateRequest = ConvertTerraformModelToApiUpdateModel(plan)

	tflog.Debug(ctx, fmt.Sprintf("Updating "+itemName+" with name: %s", state.Name.ValueString()))

	err = fabricapi.UpdateItem[fabricClientModels.WorkspaceUpdateRequestModel](state.Id.ValueString(), apiItemName, updateRequest, "", *r.client)

	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Cannot update "+itemName+" with Id %s", state.Id.ValueString()), err.Error())
		return
	}

	workspaceUpdated, err = fabricapi.GetItem[fabricClientModels.WorkspaceReadModel](state.Id.ValueString(), "workspaces", "", *r.client)
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Cannot retrieve "+itemName+" with Id %s", state.Id.ValueString()), err.Error())
		return
	}

	// convert the api data model in the terraform data model
	state = ConvertApiModelToTerraformModel(workspaceUpdated)

	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the Fabric item.
func (r *WorkspaceResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state WorkspaceProviderModel
	var err error

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Deleting "+itemName+" with name: %s", state.Name.ValueString()))

	err = fabricapi.DeleteItem(state.Id.ValueString(), apiItemName, "", *r.client)
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Cannot delete "+itemName+" with Id %s", state.Id.ValueString()), err.Error())
		return
	}

	tflog.Debug(ctx, ""+itemName+" deleted successfully")
}
