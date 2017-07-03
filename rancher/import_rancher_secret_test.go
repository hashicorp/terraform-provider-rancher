package rancher

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccRancherSecret_importBasic(t *testing.T) {
	resourceName := "rancher_secret.foo"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRancherSecretDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccRancherSecretConfig,
			},

			resource.TestStep{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"secret_value"},
			},
		},
	})
}
