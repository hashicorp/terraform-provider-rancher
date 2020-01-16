package rancher

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	rancherClient "github.com/rancher/go-rancher/v2"
)

func TestAccRancherVolume_basic(t *testing.T) {
	var volume rancherClient.Volume

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRancherVolumeDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRancherVolumeConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancherVolumeExists("rancher_volume.foo", &volume),
					resource.TestCheckResourceAttr("rancher_volume.foo", "name", "foo"),
					resource.TestCheckResourceAttr("rancher_volume.foo", "description", "Terraform acc test group"),
				),
			},
		},
	})
}

func testAccCheckRancherVolumeExists(n string, env *rancherClient.Volume) resource.TestCheckFunc {
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

		foundEnv, err := client.Volume.ById(rs.Primary.ID)
		if err != nil {
			return err
		}

		if foundEnv.Resource.Id != rs.Primary.ID {
			return fmt.Errorf("Volume not found")
		}

		*env = *foundEnv

		return nil
	}
}

func testAccCheckRancherVolumeDestroy(s *terraform.State) error {
	client, err := testAccProvider.Meta().(*Config).GlobalClient()
	if err != nil {
		return err
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "rancher_volume" {
			continue
		}
		env, err := client.Volume.ById(rs.Primary.ID)

		if err == nil {
			if env != nil &&
				env.Resource.Id == rs.Primary.ID &&
				env.State != "removed" {
				return fmt.Errorf("Volume still exists")
			}
		}

		return nil
	}
	return nil
}

const testAccRancherVolumeConfig = `
resource "rancher_environment" "foo" {
	name = "foo"
	orchestration = "cattle"
}

resource "rancher_volume" "foo" {
	name = "foo"
	description = "Terraform acc test group"
	environment_id = "${rancher_environment.foo.id}"
	driver = "local"
}
`
