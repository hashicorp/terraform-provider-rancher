package rancher

import (
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	rancherClient "github.com/rancher/go-rancher/v2"
)

var (
	defaultProjectTemplates = map[string]string{
		"mesos":      "",
		"kubernetes": "",
		"windows":    "",
		"swarm":      "",
		"cattle":     "",
	}

	// ErrNetworkPolicy network policy validation error
	ErrNetworkPolicy = errors.New("only one of `within` `between` `to|from` can be set")
)

func resourceRancherEnvironment() *schema.Resource {
	return &schema.Resource{
		Create: resourceRancherEnvironmentCreate,
		Read:   resourceRancherEnvironmentRead,
		Update: resourceRancherEnvironmentUpdate,
		Delete: resourceRancherEnvironmentDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"orchestration": &schema.Schema{
				Type:          schema.TypeString,
				Optional:      true,
				ValidateFunc:  validation.StringInSlice([]string{"cattle", "kubernetes", "mesos", "swarm", "windows"}, true),
				Computed:      true,
				ConflictsWith: []string{"project_template_id"},
				ForceNew:      true,
			},
			"project_template_id": &schema.Schema{
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"orchestration"},
				ForceNew:      true,
			},
			"description": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"default_policy": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice([]string{"allow", "deny"}, true),
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
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"external_id_type": {
							Type:     schema.TypeString,
							Required: true,
						},
						"external_id": {
							Type:     schema.TypeString,
							Required: true,
						},
						"role": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
		},
	}
}

func resourceRancherEnvironmentCreate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Creating Environment: %s", d.Id())
	populateProjectTemplateIDs(meta.(*Config))

	client, err := meta.(*Config).GlobalClient()
	if err != nil {
		return err
	}

	name := d.Get("name").(string)
	description := d.Get("description").(string)
	orchestration := d.Get("orchestration").(string)
	projectTemplateID := d.Get("project_template_id").(string)

	// validate policy arguments
	polices := d.Get("policy").(*schema.Set).List()
	if err := validateNetworkPolicies(polices); err != nil {
		return err
	}

	projectTemplateID, err = getProjectTemplateID(orchestration, projectTemplateID)
	if err != nil {
		return err
	}

	data := map[string]interface{}{
		"name":              &name,
		"description":       &description,
		"projectTemplateId": &projectTemplateID,
	}

	var newEnv rancherClient.Project
	if err := client.Create("project", data, &newEnv); err != nil {
		return err
	}

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"active", "activating", "registering", "removed", "removing"},
		Target:     []string{"active"},
		Refresh:    EnvironmentStateRefreshFunc(client, newEnv.Id),
		Timeout:    10 * time.Minute,
		Delay:      1 * time.Second,
		MinTimeout: 3 * time.Second,
	}
	_, waitErr := stateConf.WaitForState()
	if waitErr != nil {
		return fmt.Errorf(
			"Error waiting for environment (%s) to be created: %s", newEnv.Id, waitErr)
	}

	d.SetId(newEnv.Id)
	log.Printf("[INFO] Environment ID: %s", d.Id())

	if err := projectNetworkPolicyCreateOrUpdate(d, meta); err != nil {
		return err
	}

	if err := projectMembersCreateOrUpdate(d, meta); err != nil {
		return err
	}

	return resourceRancherEnvironmentRead(d, meta)
}

func resourceRancherEnvironmentRead(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Refreshing Environment: %s", d.Id())
	client, err := meta.(*Config).GlobalClient()
	if err != nil {
		return err
	}

	env, err := client.Project.ById(d.Id())
	if err != nil {
		return err
	}

	if env == nil {
		log.Printf("[INFO] Environment %s not found", d.Id())
		d.SetId("")
		return nil
	}

	if removed(env.State) {
		log.Printf("[INFO] Environment %s was removed on %v", d.Id(), env.Removed)
		d.SetId("")
		return nil
	}

	log.Printf("[INFO] Environment Name: %s", env.Name)

	d.Set("description", env.Description)
	d.Set("name", env.Name)
	d.Set("orchestration", getActiveOrchestration(env))
	d.Set("project_template_id", env.ProjectTemplateId)

	envClient, err := meta.(*Config).EnvironmentClient(d.Id())
	if err != nil {
		return err
	}

	network, err := envClient.Network.ById(env.DefaultNetworkId)
	if err != nil {
		return errors.New("Error creating environment, no default network found")
	}
	d.Set("default_policy", network.DefaultPolicyAction)

	normalizedNetworkPolicies := normalizeNetworkPolicies(network.Policy)
	d.Set("policy", normalizedNetworkPolicies)
	log.Printf("[LOG] policies found %+v", normalizedNetworkPolicies)

	members, _ := envClient.ProjectMember.List(NewListOpts())
	normalizedMembers := normalizeMembers(members.Data)
	if len(normalizedMembers) > 0 {
		d.Set("member", normalizedMembers)
	}

	return nil
}

func resourceRancherEnvironmentUpdate(d *schema.ResourceData, meta interface{}) error {
	populateProjectTemplateIDs(meta.(*Config))

	client, err := meta.(*Config).GlobalClient()
	if err != nil {
		return err
	}

	// validate policy arguments
	polices := d.Get("policy").(*schema.Set).List()
	if err := validateNetworkPolicies(polices); err != nil {
		return err
	}

	name := d.Get("name").(string)
	description := d.Get("description").(string)

	data := map[string]interface{}{
		"name":        &name,
		"description": &description,
	}

	var newEnv rancherClient.Project
	env, err := client.Project.ById(d.Id())
	if err != nil {
		return err
	}

	if err := client.Update("project", &env.Resource, data, &newEnv); err != nil {
		return err
	}

	if err := projectNetworkPolicyCreateOrUpdate(d, meta); err != nil {
		return err
	}

	if err := projectMembersCreateOrUpdate(d, meta); err != nil {
		return err
	}

	return resourceRancherEnvironmentRead(d, meta)
}

func resourceRancherEnvironmentDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Deleting Environment: %s", d.Id())
	id := d.Id()
	client, err := meta.(*Config).GlobalClient()
	if err != nil {
		return err
	}

	env, err := client.Project.ById(id)
	if err != nil {
		return err
	}

	if err := client.Project.Delete(env); err != nil {
		return fmt.Errorf("Error deleting Environment: %s", err)
	}

	log.Printf("[DEBUG] Waiting for environment (%s) to be removed", id)

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"active", "removed", "removing"},
		Target:     []string{"removed"},
		Refresh:    EnvironmentStateRefreshFunc(client, id),
		Timeout:    10 * time.Minute,
		Delay:      1 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, waitErr := stateConf.WaitForState()
	if waitErr != nil {
		return fmt.Errorf(
			"Error waiting for environment (%s) to be removed: %s", id, waitErr)
	}

	d.SetId("")
	return nil
}

func projectNetworkPolicyCreateOrUpdate(d *schema.ResourceData, meta interface{}) error {
	client, err := meta.(*Config).GlobalClient()
	if err != nil {
		return err
	}

	env, err := client.Project.ById(d.Id())
	if err != nil {
		return err
	}

	envClient, err := meta.(*Config).EnvironmentClient(d.Id())
	if err != nil {
		return err
	}

	network, err := envClient.Network.ById(env.DefaultNetworkId)
	if err != nil {
		return fmt.Errorf("Error failed retrive default network interface with id: %s", env.DefaultNetworkId)
	}

	desiredNetworkPolicy := d.Get("default_policy").(string)
	if desiredNetworkPolicy != "" && desiredNetworkPolicy != network.DefaultPolicyAction {
		_, err := envClient.Network.Update(network, map[string]interface{}{
			"defaultPolicyAction": desiredNetworkPolicy,
		})

		if err != nil {
			return fmt.Errorf("Error failed to set default network policy, network id %s", env.DefaultNetworkId)
		}
	}

	policies := d.Get("policy").(*schema.Set).List()
	_, err = envClient.Network.Update(network, map[string]interface{}{
		"policy": makeNetworkPolicies(policies),
	})

	if err != nil {
		return fmt.Errorf("Error failed to set network policies, network id %s", env.DefaultNetworkId)
	}

	return nil
}

func projectMembersCreateOrUpdate(d *schema.ResourceData, meta interface{}) error {
	client, err := meta.(*Config).GlobalClient()
	if err != nil {
		return err
	}

	env, err := client.Project.ById(d.Id())
	if err != nil {
		return err
	}

	// Create or Update members
	envClient, err := meta.(*Config).EnvironmentClient(d.Id())
	if err != nil {
		return err
	}
	members := d.Get("member").(*schema.Set).List()
	log.Printf("[INFO] Create or Update Project Members: %v", makeProjectMembers(members))
	if members != nil && len(members) > 0 {
		_, err = envClient.Project.ActionSetmembers(env, &rancherClient.SetProjectMembersInput{
			Members: makeProjectMembers(members),
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func getProjectTemplateID(orchestration, templateID string) (string, error) {
	id := templateID
	if templateID == "" && orchestration == "" {
		return "", fmt.Errorf("Need either 'orchestration' or 'project_template_id'")
	}

	if templateID == "" && orchestration != "" {
		ok := false
		id, ok = defaultProjectTemplates[orchestration]
		if !ok {
			return "", fmt.Errorf("Invalid orchestration: %s", orchestration)
		}
	}

	return id, nil
}

func normalizeMembers(in []rancherClient.ProjectMember) (out []interface{}) {
	for _, m := range in {
		mm := map[string]string{
			"external_id_type": m.ExternalIdType,
			"external_id":      m.ExternalId,
			"role":             m.Role,
		}
		out = append(out, mm)
	}
	return
}

func normalizeNetworkPolicies(in []rancherClient.NetworkPolicyRule) (out []interface{}) {
	for _, policy := range in {
		n := make(map[string]interface{})
		n["action"] = policy.Action

		if policy.Within != "" {
			n["within"] = policy.Within
		}

		if policy.Between.GroupBy != "" {
			n["between"] = policy.Between.GroupBy
		}

		if policy.From.Selector != "" {
			n["from"] = policy.From.Selector
		}

		if policy.To.Selector != "" {
			n["to"] = policy.To.Selector
		}

		if len(policy.Ports) > 0 {
			n["ports"] = strings.Join(policy.Ports, ",")
		}

		out = append(out, n)
	}

	return
}

func makeNetworkPolicies(in []interface{}) (out []map[string]interface{}) {
	for _, n := range in {
		nMap := n.(map[string]interface{})

		nn := map[string]interface{}{}
		within := nMap["within"].(string)
		between := nMap["between"].(string)
		to := nMap["to"].(string)
		from := nMap["from"].(string)

		nn["action"] = nMap["action"].(string)

		if within != "" {
			nn["within"] = within
		}

		if between != "" {
			nn["between"] = map[string]string{
				"groupBy": between,
			}
		}

		if from != "" {
			nn["from"] = map[string]string{
				"selector": from,
			}
		}

		if to != "" {
			nn["to"] = map[string]string{
				"selector": to,
			}
		}

		ports := nMap["ports"].(string)
		if ports != "" {
			nn["ports"] = strings.Split(ports, ",")
		}

		out = append(out, nn)
	}
	return
}

func validateNetworkPolicies(policies []interface{}) error {
	for _, n := range policies {
		nMap := n.(map[string]interface{})

		within, _ := nMap["within"].(string)
		between, _ := nMap["between"].(string)
		to, _ := nMap["to"].(string)
		from, _ := nMap["from"].(string)

		if within != "" && between != "" {
			return ErrNetworkPolicy
		}

		if (within != "" || between != "") && (from != "" || to != "") {
			return ErrNetworkPolicy
		}

		if (from != "" && to == "") || (from == "" && to != "") {
			return ErrNetworkPolicy
		}
	}

	return nil
}

func makeProjectMembers(in []interface{}) (out []rancherClient.ProjectMember) {
	for _, m := range in {
		mMap := m.(map[string]interface{})
		mm := rancherClient.ProjectMember{
			ExternalIdType: mMap["external_id_type"].(string),
			ExternalId:     mMap["external_id"].(string),
			Role:           mMap["role"].(string),
		}
		out = append(out, mm)
	}
	return
}

// EnvironmentStateRefreshFunc returns a resource.StateRefreshFunc that is used to watch
// a Rancher Environment.
func EnvironmentStateRefreshFunc(client *rancherClient.RancherClient, environmentID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		env, err := client.Project.ById(environmentID)

		if err != nil {
			return nil, "", err
		}

		// Wait until network is available
		if env != nil && env.DefaultNetworkId == "" {
			return nil, "", nil
		}

		// Env not returned, or State not set...
		if env == nil || env.State == "" {
			// This makes it so user level API keys can be used instead of just admin
			env = &rancherClient.Project{
				State: "removed",
			}
		}

		return env, env.State, nil
	}
}
