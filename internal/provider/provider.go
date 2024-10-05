package provider

import (
	"context"
	"fmt"
	"github.com/greatman/go-netbox/v4"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"net/http"
	"os"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

var _ provider.Provider = (*netboxProvider)(nil)

func NewProvider() func() provider.Provider {
	return func() provider.Provider {
		return &netboxProvider{}
	}
}

type netboxProvider struct {
	client *netbox.APIClient
}

type netboxProviderModel struct {
	ApiToken                    types.String `tfsdk:"api_token"`
	ServerUrl                   types.String `tfsdk:"server_url"`
	SkipVersionCheck            types.Bool   `tfsdk:"skip_version_check"`
	AllowInsecureHttps          types.Bool   `tfsdk:"allow_insecure_https"`
	Headers                     types.Map    `tfsdk:"headers"`
	StripTrailingSlashesFromUrl types.Bool   `tfsdk:"strip_trailing_slashes_from_url"`
	RequestTimeout              types.Int32  `tfsdk:"request_timeout"`
}

func (p *netboxProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"api_token": schema.StringAttribute{
				MarkdownDescription: "Netbox API authentication token. Can be set via the `NETBOX_API_TOKEN` environment variable.",
				Optional:            true,
			},
			"server_url": schema.StringAttribute{
				MarkdownDescription: "Location of Netbox server including scheme (http or https) and optional port. Can be set via the `NETBOX_SERVER_URL` environment variable.",
				Optional:            true,
			},
			"skip_version_check": schema.BoolAttribute{
				MarkdownDescription: "If true, do not try to determine the running Netbox version at provider startup. Disables warnings about possibly unsupported Netbox version. Also useful for local testing on terraform plans. Can be set via the `NETBOX_SKIP_VERSION_CHECK` environment variable. Defaults to `false`.",
				Optional:            true,
			},
			"allow_insecure_https": schema.BoolAttribute{
				Optional:            true,
				MarkdownDescription: "Flag to set whether to allow https with invalid certificates. Can be set via the `NETBOX_ALLOW_INSECURE_HTTPS` environment variable. Defaults to `false`.",
			},
			"headers": schema.MapAttribute{
				Optional:            true,
				MarkdownDescription: "Set these header on all requests to Netbox. Can be set via the `NETBOX_HEADERS` environment variable.",
				ElementType:         types.StringType,
			},
			"strip_trailing_slashes_from_url": schema.BoolAttribute{
				Optional:            true,
				MarkdownDescription: "If true, strip trailing slashes from the `server_url` parameter and print a warning when doing so. Note that using trailing slashes in the `server_url` parameter will usually lead to errors. Can be set via the `NETBOX_STRIP_TRAILING_SLASHES_FROM_URL` environment variable. Defaults to `true`.",
			},
			"request_timeout": schema.Int32Attribute{
				Optional:            true,
				MarkdownDescription: "Netbox API HTTP request timeout in seconds. Can be set via the `NETBOX_REQUEST_TIMEOUT` environment variable.",
			},
		},
	}
}

const authHeaderName = "Authorization"
const authHeaderFormat = "Token %v"
const languageHeaderName = "Accept-Language"
const languageHeaderValue = "en-US"

func (p *netboxProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var data netboxProviderModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	apiToken := os.Getenv("NETBOX_API_TOKEN")
	serverUrl := os.Getenv("NETBOX_SERVER_URL")
	requestTimeoutString := os.Getenv("NETBOX_REQUEST_TIMEOUT")
	var requestTimeout int32
	if requestTimeoutString == "" {
		requestTimeout = 10
	} else {
		requestTimeoutConversion, err := strconv.ParseInt(requestTimeoutString, 10, 32)

		if err != nil {
			resp.Diagnostics.AddError("Invalid Request Timeout environment variable.", "TODO")
		}
		requestTimeout = int32(requestTimeoutConversion)
	}

	if !data.ApiToken.IsNull() {
		apiToken = data.ApiToken.ValueString()
	}

	if !data.ServerUrl.IsNull() {
		serverUrl = data.ServerUrl.ValueString()
	}

	if !data.RequestTimeout.IsNull() {
		requestTimeout = data.RequestTimeout.ValueInt32()
	}

	if apiToken == "" {
		resp.Diagnostics.AddError(
			"Missing API Token Configuration",
			"TODO DETAIL")
	}

	if serverUrl == "" {
		resp.Diagnostics.AddError(
			"Missing server URL configuration.",
			"TODO details")
	}
	httpClient := &http.Client{
		Timeout: time.Second * time.Duration(requestTimeout),
	}
	conf := netbox.NewConfiguration()
	conf.Servers[0].URL = serverUrl
	conf.AddDefaultHeader(
		authHeaderName,
		fmt.Sprintf(authHeaderFormat, apiToken),
	)

	conf.AddDefaultHeader(
		languageHeaderName,
		languageHeaderValue,
	)
	conf.HTTPClient = httpClient
	netboxClient := netbox.NewAPIClient(conf)

	res, _, err := netboxClient.StatusAPI.StatusRetrieve(ctx).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error while checking Netbox version...",
			err.Error(),
		)
	}
	resp.ResourceData = netboxClient
	supportedVersions := []string{"4.1.3"}
	if !slices.Contains(supportedVersions, res["netbox-version"].(string)) {
		resp.Diagnostics.AddWarning("Possibly unsupported Netbox version", fmt.Sprintf("Your Netbox version is v%v. The provider was successfully tested against the following versions:\n\n  %v\n\nUnexpected errors may occur.", res["netbox-version"], strings.Join(supportedVersions, ", ")))
	}
	p.client = netboxClient
	resp.ResourceData = p
	resp.DataSourceData = p
}

func (p *netboxProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "netbox"
}

func (p *netboxProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{}
}

func (p *netboxProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewSiteResource,
	}
}
