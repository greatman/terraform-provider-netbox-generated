package provider

import (
	"context"
	"fmt"
	"github.com/greatman/go-netbox/v4"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"terraform-provider-netbox-generated/internal/resource_site"

	"github.com/hashicorp/terraform-plugin-framework/resource"
)

var _ resource.Resource = (*siteResource)(nil)

func NewSiteResource() resource.Resource {
	return &siteResource{}
}

type siteResource struct {
	provider *netboxProvider
}

func (r *siteResource) Configure(ctx context.Context, request resource.ConfigureRequest, response *resource.ConfigureResponse) {
	if request.ProviderData == nil {
		return
	}

	provider, ok := request.ProviderData.(*netboxProvider)
	if !ok {
		response.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *client.NetBoxAPI, got: %T, Please report this issue to the provider developers.", request.ProviderData),
		)
		return
	}

	r.provider = provider
}

func (r *siteResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_site"
}

func (r *siteResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = resource_site.SiteResourceSchema(ctx)
}

func (r *siteResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data resource_site.SiteModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Create API call logic
	site := generateWritableSite(data, resp.Diagnostics)

	_, _, err := r.provider.client.DcimAPI.DcimSitesCreate(ctx).WritableSiteRequest(*site).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error while creating the site",
			err.Error())
	}
	//TODO: Some handling of the data back from Netbox

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *siteResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data resource_site.SiteModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Read API call logic

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *siteResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data resource_site.SiteModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	//resp.Diagnostics.AddWarning()

	if resp.Diagnostics.HasError() {
		return
	}

	// Update API call logic

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *siteResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data resource_site.SiteModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Delete API call logic
}

func generateWritableSite(data resource_site.SiteModel, diag diag.Diagnostics) *netbox.WritableSiteRequest {
	site := netbox.NewWritableSiteRequest(data.Name.ValueString(), data.Slug.ValueString())

	status, err := netbox.NewLocationStatusValueFromValue(data.Status.ValueString())

	if err != nil {
		diag.AddWarning("Invalid status value.", "This is invalid.")
	}
	site.Status = status

	if !data.Description.IsNull() {
		site.Description = data.Description.ValueStringPointer()
	}

	if !data.Facility.IsNull() {
		site.Facility = data.Facility.ValueStringPointer()
	}

	if !data.Longitude.IsNull() {
		site.Longitude = *netbox.NewNullableFloat64(data.Longitude.ValueFloat64Pointer())
	}

	if !data.Latitude.IsNull() {
		site.Latitude = *netbox.NewNullableFloat64(data.Latitude.ValueFloat64Pointer())
	}

	if !data.PhysicalAddress.IsNull() {
		site.PhysicalAddress = data.PhysicalAddress.ValueStringPointer()
	}

	if !data.ShippingAddress.IsNull() {
		site.ShippingAddress = data.ShippingAddress.ValueStringPointer()
	}

	/*if !data.Region.IsNull() {
		site.Region = data.RegionID.ValueInt64Pointer()
	}*/

	/*if !data.Group.IsNull() {

		site.Group = netbox.NewNullableBriefSiteGroupRequest(netbox.NewBriefSiteGroupRequest(data.Group.String())).Get()
	}*/
	/*if !data.Tenant.IsNull() {
		site.Tenant = data.TenantID.ValueInt64Pointer()
	}

	if !data.Tags.IsNull() {
		site.Tags = getNestedTagListFromResourceDataSetV6(api, data.Tags, diag)
	} else {
		site.Tags = []*models.NestedTag{}
	}

	if !data.TimeZone.IsNull() {
		site.TimeZone = data.Timezone.ValueStringPointer()
	}
	data.Asns.
	if !data.ASNIDs.IsNull() {
		site.Asns = toInt64ListV6(data.ASNIDs)
	} else {
		site.Asns = []int64{}
	}*/
	return site
}

func handle_api(data resource_site.SiteModel) {

}
