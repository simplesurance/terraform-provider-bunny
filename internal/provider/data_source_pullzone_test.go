package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourcePullZone(t *testing.T) {
	resourceName := randResourceName()
	resource.Test(t, resource.TestCase{
		Providers: testProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourcePullZone(resourceName),
				Check: resource.TestCheckResourceAttr(
					"data.bunny_pullzone.dpz1",
					"name",
					resourceName,
				),
			},
		},
	})
}

func testAccDataSourcePullZone(rName string) string {
	return fmt.Sprintf(`
resource "bunny_pullzone" "pz1" {
  name       = %q
  origin_url = "https://terraform.io"
}

data "bunny_pullzone" "dpz1" {
  pull_zone_id = bunny_pullzone.pz1.id
}`, rName)
}
