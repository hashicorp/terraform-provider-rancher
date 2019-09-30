package rancher

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccRancherCertificateDataSource_foo(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckRancherCertificateDataSourceConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.rancher_certificate.foo", "cn", "foo"),
				),
			},
		},
	})
}

// Testing owner parameter
const testAccCheckRancherCertificateDataSourceConfig = `
data "rancher_certificate" "foo" {
	name = "foo"
	environment_id = "1a5"
}
`
