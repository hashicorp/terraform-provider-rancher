package rancher

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccRancherApiKey_importBasic(t *testing.T) {
	resourceName := "rancher_registration_token.foo"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRancherApiKeyDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccRancherApiKeyConfig,
			},

			resource.TestStep{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
