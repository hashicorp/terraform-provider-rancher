package rancher

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	rancher "github.com/rancher/go-rancher/v2"
)

func dataSourceRancherEnvironment() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceRancherEnvironmentRead,

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"orchestration": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"project_template_id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"default_policy": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"policy": &schema.Schema{
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"action": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringInSlice([]string{"allow", "deny"}, true),
						},
						"within": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringInSlice([]string{"stack", "service", "linked"}, true),
						},
						"between": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"from": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"to": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"ports": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"member": &schema.Schema{
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"external_id_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"external_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"role": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceRancherEnvironmentRead(d *schema.ResourceData, meta interface{}) error {
	name := d.Get("name").(string)
	log.Printf("[INFO] Refreshing Rancher Environment: %s", name)

	client, err := meta.(*Config).GlobalClient()
	if err != nil {
		return err
	}

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"active", "removed", "removing", "not found"},
		Target:     []string{"active"},
		Refresh:    findEnv(client, name),
		Timeout:    10 * time.Minute,
		Delay:      1 * time.Second,
		MinTimeout: 3 * time.Second,
	}
	env, waitErr := stateConf.WaitForState()
	if waitErr != nil {
		return fmt.Errorf(
			"Error waiting for environment (%s) to be found: %s", name, waitErr)
	}

	environment := env.(rancher.Project)
	d.SetId(environment.Id)

	d.Set("description", environment.Description)
	d.Set("name", environment.Name)

	// Computed values
	d.Set("orchestration", getActiveOrchestration(&environment))
	d.Set("project_template_id", environment.ProjectTemplateId)

	envClient, err := meta.(*Config).EnvironmentClient(d.Id())
	if err != nil {
		return err
	}

	network, err := envClient.Network.ById(environment.DefaultNetworkId)
	if err != nil {
		return errors.New("Error creating environment, no default network found")
	}
	d.Set("default_policy", network.DefaultPolicyAction)

	normalizedNetworkPolicies := normalizeNetworkPolicies(network.Policy)
	if len(normalizedNetworkPolicies) > 0 {
		d.Set("policy", normalizedNetworkPolicies)
	}

	members, _ := envClient.ProjectMember.List(NewListOpts())
	normalizedMembers := normalizeMembers(members.Data)
	if len(normalizedMembers) > 0 {
		d.Set("member", normalizedMembers)
	}

	return nil
}

func findEnv(client *rancher.RancherClient, envname string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		envs, err := client.Project.List(NewListOpts())
		if err != nil {
			return nil, "", err
		}

		for _, env := range envs.Data {
			if env.Name == envname {
				return env, env.State, nil
			}
		}

		return nil, "not found", nil
	}
}
