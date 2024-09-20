// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"terraform-provider-fabric/internal/fabricapi"
	workspaceProvider "terraform-provider-fabric/internal/provider/workspaceProvider"
	"terraform-provider-fabric/internal/provider/workspaceRoleAssignmentProvider"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure FabricProvider satisfies various provider interfaces.
var _ provider.Provider = &FabricProvider{}

// FabricProvider defines the provider implementation.
type FabricProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// FabricProviderModel describes the provider data model.
type FabricProviderModel struct {
	Client_id       types.String `tfsdk:"client_id"`
	Client_secret   types.String `tfsdk:"client_secret"`
	Tenant_id       types.String `tfsdk:"tenant_id"`
	Subscription_id types.String `tfsdk:"subscription_id"`
}

func (p *FabricProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "fabric"
	resp.Version = p.version
}

// Schema defines the provider-level schema for configuration data.
func (p *FabricProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{

			"client_id": schema.StringAttribute{
				MarkdownDescription: "ID (GUID) of the service principal.",
				Description:         "ID (GUID) of the service principal.",
				Optional:            false,
				Required:            true,
			},
			"client_secret": schema.StringAttribute{
				MarkdownDescription: "Secret of the service principal.",
				Description:         "Secret of the service principal.",
				Optional:            false,
				Required:            true,
			},
			"tenant_id": schema.StringAttribute{
				MarkdownDescription: "Your Entra ID Tenant ID (GUID)",
				Description:         "Your Entra ID Tenant ID (GUID)",
				Optional:            false,
				Required:            true,
			},
			"subscription_id": schema.StringAttribute{
				MarkdownDescription: "Your Subscription ID (GUID)",
				Description:         "Your Subscription ID (GUID)",
				Optional:            false,
				Required:            true,
			},
		},
	}
}

func (p *FabricProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var data FabricProviderModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	client, err := fabricapi.NewFabricClient()
	if err != nil {
		resp.Diagnostics.AddError("failed to create client", err.Error())
		return
	}

	resp.DataSourceData = client
	resp.ResourceData = client
}

func (p *FabricProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		workspaceProvider.NewWorkspaceResource,
		workspaceRoleAssignmentProvider.NewWorkspaceRoleAssignmentResource,
	}
}

func (p *FabricProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		workspaceProvider.NewWorkspaceDataSource,
		workspaceRoleAssignmentProvider.NewWorkspaceRoleAssignmentDataSource,
	}
}

func NewFabricProvider(version string) func() provider.Provider {
	return func() provider.Provider {
		return &FabricProvider{
			version: version,
		}
	}
}
