package rancher

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccRancherSettingDataSource_accessLog(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckRancherSettingDataSourceConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.rancher_setting.access_log", "value", "/dev/null"),
				),
			},
		},
	})
}

// Testing owner parameter
const testAccCheckRancherSettingDataSourceConfig = `
data "rancher_setting" "access_log" {
	name = "access.log"
}
`
