package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
)

var testAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
	// newProvider is an example function that returns a provider.Provider
	"netbox": providerserver.NewProtocol6WithError(NewProvider()()),
}
