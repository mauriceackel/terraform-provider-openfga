package store

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"

	"github.com/openfga/go-sdk/client"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &StoresDataSource{}
var _ datasource.DataSourceWithConfigure = &StoresDataSource{}

func NewStoresDataSource() datasource.DataSource {
	return &StoresDataSource{}
}

type StoresDataSource struct {
	client *StoreClient
}

type StoresDataSourceModel struct {
	Stores []StoreModel `tfsdk:"stores"`
}

func (d *StoresDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_stores"
}

func (d *StoresDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Provides the ability to list and retrieve details of existing OpenFGA stores.",

		Attributes: map[string]schema.Attribute{
			"stores": schema.ListNestedAttribute{
				MarkdownDescription: "List of existing stores.",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							MarkdownDescription: "The unique ID of the store.",
							Computed:            true,
						},
						"name": schema.StringAttribute{
							MarkdownDescription: "The name of the store.",
							Computed:            true,
						},
					},
				},
			},
		},
	}
}

func (d *StoresDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*client.OpenFgaClient)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *client.OpenFgaClient, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	d.client = NewStoreClient(client)
}

func (d *StoresDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state StoresDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	storeModels, err := d.client.ListStores(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read stores, got error: %s", err))
		return
	}

	state.Stores = *storeModels

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
