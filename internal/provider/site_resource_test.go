package provider

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestSiteResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testSiteResourceConfig("test"),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("netbox_site.test", tfjsonpath.New("name"), knownvalue.StringExact("test")),
				},
			},
		},
	})
}

func testSiteResourceConfig(attribute string) string {
	return fmt.Sprintf(`
resource "netbox_site" "test" {
	name = "%[1]s"
	slug = "testslug"
}
`, attribute)
}
