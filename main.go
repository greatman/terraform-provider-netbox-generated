package main

import (
	"context"
	"log"

	"terraform-provider-netbox-generated/internal/provider"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
)

//go:generate tfplugingen-openapi generate --config api/provider_config.yml --output api/provider_code_spec.json api/openapi.yaml
//go:generate tfplugingen-framework generate resources --input api/provider_code_spec.json --output internal
func main() {
	opts := providerserver.ServeOpts{
		Address: "hashicorp.com/greatman/netbox",
	}

	err := providerserver.Serve(context.Background(), provider.NewProvider(), opts)
	if err != nil {
		log.Fatal(err.Error())
	}
}
