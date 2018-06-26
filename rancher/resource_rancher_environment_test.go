package rancher

import (
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	rancherClient "github.com/rancher/go-rancher/v2"
)

func TestAccRancherEnvironment_basic(t *testing.T) {
	var environment rancherClient.Project

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRancherEnvironmentDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccRancherEnvironmentConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancherEnvironmentExists("rancher_environment.foo", &environment),
					resource.TestCheckResourceAttr("rancher_environment.foo", "name", "foo"),
					resource.TestCheckResourceAttr("rancher_environment.foo", "description", "Terraform acc test group"),
					resource.TestCheckResourceAttr("rancher_environment.foo", "orchestration", "cattle"),
				),
			},
			resource.TestStep{
				Config: testAccRancherEnvironmentUpdateConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancherEnvironmentExists("rancher_environment.foo", &environment),
					resource.TestCheckResourceAttr("rancher_environment.foo", "name", "foo2"),
					resource.TestCheckResourceAttr("rancher_environment.foo", "description", "Terraform acc test group - updated"),
					resource.TestCheckResourceAttr("rancher_environment.foo", "orchestration", "swarm"),
				),
			},
		},
	})
}

func TestAccRancherEnvironment_disappears(t *testing.T) {
	var environment rancherClient.Project

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRancherEnvironmentDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccRancherEnvironmentConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancherEnvironmentExists("rancher_environment.foo", &environment),
					testAccRancherEnvironmentDisappears(&environment),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccRancherEnvironment_members(t *testing.T) {
	var environment rancherClient.Project

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRancherEnvironmentDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccRancherEnvironmentMembersConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancherEnvironmentExists("rancher_environment.foo", &environment),
					resource.TestCheckResourceAttr("rancher_environment.foo", "name", "foo"),
					resource.TestCheckResourceAttr("rancher_environment.foo", "description", "Terraform acc test group"),
					resource.TestCheckResourceAttr("rancher_environment.foo", "orchestration", "cattle"),
					resource.TestCheckResourceAttr("rancher_environment.foo", "member.#", "2"),
				),
			},
			resource.TestStep{
				Config: testAccRancherEnvironmentMembersUpdateConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancherEnvironmentExists("rancher_environment.foo", &environment),
					resource.TestCheckResourceAttr("rancher_environment.foo", "name", "foo2"),
					resource.TestCheckResourceAttr("rancher_environment.foo", "description", "Terraform acc test group - updated"),
					resource.TestCheckResourceAttr("rancher_environment.foo", "orchestration", "cattle"),
					resource.TestCheckResourceAttr("rancher_environment.foo", "member.#", "1"),
				),
			},
		},
	})
}

func testAccRancherEnvironmentDisappears(env *rancherClient.Project) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client, err := testAccProvider.Meta().(*Config).GlobalClient()
		if err != nil {
			return err
		}
		if err := client.Project.Delete(env); err != nil {
			return fmt.Errorf("Error deleting Environment: %s", err)
		}
		stateConf := &resource.StateChangeConf{
			Pending:    []string{"active", "removed", "removing"},
			Target:     []string{"removed"},
			Refresh:    EnvironmentStateRefreshFunc(client, env.Id),
			Timeout:    10 * time.Minute,
			Delay:      1 * time.Second,
			MinTimeout: 3 * time.Second,
		}

		_, waitErr := stateConf.WaitForState()
		if waitErr != nil {
			return fmt.Errorf(
				"Error waiting for environment (%s) to be removed: %s", env.Id, waitErr)
		}
		return nil
	}
}

func testAccCheckRancherEnvironmentExists(n string, env *rancherClient.Project) resource.TestCheckFunc {
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

		foundEnv, err := client.Project.ById(rs.Primary.ID)
		if err != nil {
			return err
		}

		if foundEnv.Resource.Id != rs.Primary.ID {
			return fmt.Errorf("Environment not found")
		}

		*env = *foundEnv

		return nil
	}
}

func testAccCheckRancherEnvironmentDestroy(s *terraform.State) error {
	client, err := testAccProvider.Meta().(*Config).GlobalClient()
	if err != nil {
		return err
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "rancher_environment" {
			continue
		}
		env, err := client.Project.ById(rs.Primary.ID)

		if err == nil {
			if env != nil &&
				env.Resource.Id == rs.Primary.ID &&
				env.State != "removed" {
				return fmt.Errorf("Environment still exists")
			}
		}

		return nil
	}
	return nil
}

func TestAccRancherEnvironmentDefaultPolicy(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRancherEnvironmentDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccRancherEnvironmentDefaultPolicyConfig,
				Check: func(s *terraform.State) error {
					rs, ok := s.RootModule().Resources["rancher_environment.foo"]

					if !ok {
						return fmt.Errorf("Environment not found")
					}

					if rs.Primary.ID == "" {
						return fmt.Errorf("No App Name is set")
					}

					client, err := testAccProvider.Meta().(*Config).GlobalClient()
					if err != nil {
						return err
					}

					env, err := client.Project.ById(rs.Primary.ID)
					if err != nil {
						return err
					}

					envClient, err := testAccProvider.Meta().(*Config).EnvironmentClient(rs.Primary.ID)
					if err != nil {
						return fmt.Errorf("Failed to create project scoped client")
					}

					network, err := envClient.Network.ById(env.DefaultNetworkId)
					if err != nil {
						return fmt.Errorf("Error failed retrive default network interface with id: %s", env.DefaultNetworkId)
					}

					if want := rs.Primary.Attributes["default_policy"]; want != network.DefaultPolicyAction {
						return fmt.Errorf("Mismatch network policy want: %s received: %s", want, network.DefaultPolicyAction)
					}

					return nil
				},
			},
		},
	},
	)
}

func TestEnviromentPolicyRules(t *testing.T) {
	testCases := []struct {
		desc     string
		in       []interface{}
		expected error
	}{
		{
			desc: "to-from",
			in: []interface{}{
				map[string]interface{}{
					"from":  "a",
					"to":    "b",
					"ports": []string{"1000"},
				},
			},
			expected: nil,
		},
		{
			desc: "within",
			in: []interface{}{
				map[string]interface{}{
					"action": "allow",
					"within": "foo",
				},
			},
			expected: nil,
		},
		{
			desc: "between",
			in: []interface{}{
				map[string]interface{}{
					"between": "stack",
					"action":  "allow",
				},
			},
			expected: nil,
		},
		{
			desc: "incompatible to/from within",
			in: []interface{}{
				map[string]interface{}{
					"from":   "a",
					"to":     "b",
					"within": "service",
				},
			},
			expected: ErrNetworkPolicy,
		},
		{
			desc: "incompatible to/from between",
			in: []interface{}{
				map[string]interface{}{
					"from":    "a",
					"to":      "b",
					"between": "service",
				},
			},
			expected: ErrNetworkPolicy,
		},
		{
			desc: "incompatible to/from between",
			in: []interface{}{
				map[string]interface{}{
					"from":   "a",
					"to":     "b",
					"within": "service",
				},
			},
			expected: ErrNetworkPolicy,
		},
		{
			desc: "missing to",
			in: []interface{}{
				map[string]interface{}{
					"from": "a",
				},
			},
			expected: ErrNetworkPolicy,
		},
		{
			desc: "missing from",
			in: []interface{}{
				map[string]interface{}{
					"to": "a",
				},
			},
			expected: ErrNetworkPolicy,
		},
	}

	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			if recv := validateNetworkPolicies(tC.in); recv != tC.expected {
				t.Errorf("rule validation failed, expected %+v, received %+v", tC.expected, recv)
			}
		})
	}
}

const testAccRancherEnvironmentConfig = `
resource "rancher_environment" "foo" {
	name = "foo"
	description = "Terraform acc test group"
	orchestration = "cattle"
}
`

const testAccRancherEnvironmentUpdateConfig = `
resource "rancher_environment" "foo" {
	name = "foo2"
	description = "Terraform acc test group - updated"
	orchestration = "swarm"
}
`
const testAccRancherEnvironmentMembersConfig = `
resource "rancher_environment" "foo" {
	name = "foo"
	description = "Terraform acc test group"
	orchestration = "cattle"

	member {
		external_id = "1234"
		external_id_type = "github_user"
		role = "owner"
	}

	member {
		external_id = "8765"
		external_id_type = "github_team"
		role = "member"
	}
}
`

const testAccRancherEnvironmentMembersUpdateConfig = `
resource "rancher_environment" "foo" {
	name = "foo2"
	description = "Terraform acc test group - updated"
	orchestration = "cattle"

	member {
		external_id = "1235"
		external_id_type = "github_user"
		role = "owner"
	}
}
`

const testAccRancherEnvironmentDefaultPolicyConfig = `
resource "rancher_environment" "foo" {
	name = "foo"
	description = "Terraform acc test group"
	orchestration = "cattle"
	default_policy = "deny"
}
`

const testAccRancherInvalidEnvironmentConfig = `
resource "rancher_environment_invalid_config" "bar" {
	name = "bar"
	description = "Terraform acc test group - failure"
	orchestration = "cattle"
	project_template_id = "1pt1"
}
`
