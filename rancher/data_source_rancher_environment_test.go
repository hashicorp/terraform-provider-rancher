package rancher

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccRancherEnvironmentDataSource_foo(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckRancherEnvironmentDataSourceConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.rancher_environment.default", "orchestration", "cattle"),
				),
			},
		},
	})
}

// Testing owner parameter
const testAccCheckRancherEnvironmentDataSourceConfig = `
data "rancher_environment" "default" {
	name = "Default"
}
`
