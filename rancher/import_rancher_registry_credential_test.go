package rancher

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccRancherRegistryCredential_importBasic(t *testing.T) {
	resourceName := "rancher_registry_credential.foo"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRancherRegistryCredentialDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRancherRegistryCredentialConfig,
			},

			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"secret_value"},
			},
		},
	})
}
