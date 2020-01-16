package rancher

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	rancherClient "github.com/rancher/go-rancher/v2"
)

func TestAccRancherSecret_basic(t *testing.T) {
	var secret rancherClient.Secret

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRancherSecretDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRancherSecretConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancherSecretExists("rancher_secret.foo", &secret),
					resource.TestCheckResourceAttr("rancher_secret.foo", "name", "foo"),
					resource.TestCheckResourceAttr("rancher_secret.foo", "description", "Terraform acc test group"),
				),
			},
		},
	})
}

func testAccCheckRancherSecretExists(n string, env *rancherClient.Secret) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No App Name is set")
		}

		client, err := testAccProvider.Meta().(*Config).GlobalClient()
		if err != nil {
			return err
		}

		foundEnv, err := client.Secret.ById(rs.Primary.ID)
		if err != nil {
			return err
		}

		if foundEnv.Resource.Id != rs.Primary.ID {
			return fmt.Errorf("Secret not found")
		}

		*env = *foundEnv

		return nil
	}
}

func testAccCheckRancherSecretDestroy(s *terraform.State) error {
	client, err := testAccProvider.Meta().(*Config).GlobalClient()
	if err != nil {
		return err
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "rancher_secret" {
			continue
		}
		env, err := client.Secret.ById(rs.Primary.ID)

		if err == nil {
			if env != nil &&
				env.Resource.Id == rs.Primary.ID &&
				env.State != "removed" {
				return fmt.Errorf("Secret still exists")
			}
		}

		return nil
	}
	return nil
}

const testAccRancherSecretConfig = `
resource "rancher_secret" "foo" {
	name = "foo"
	description = "Terraform acc test group"
	environment_id = "1a5"
	value = "mypasswd"
}
`
