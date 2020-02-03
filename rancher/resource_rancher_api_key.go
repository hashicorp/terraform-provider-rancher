package rancher

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	rancherClient "github.com/rancher/go-rancher/v2"
)

func resourceRancherApiKey() *schema.Resource {
	return &schema.Resource{
		Create: resourceRancherApiKeyCreate,
		Read:   resourceRancherApiKeyRead,
		Delete: resourceRancherApiKeyDelete,
		Update: resourceRancherApiKeyUpdate,
		Importer: &schema.ResourceImporter{
			State: resourceRancherApiKeyImport,
		},

		Schema: map[string]*schema.Schema{
			"id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: false,
			},
			"description": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: false,
			},
			"reference_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"kind": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"public_value": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"secret_value": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"uuid": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceRancherApiKeyCreate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Creating ApiKey: %s", d.Id())
	client, err := meta.(*Config).ApiKeyClient(d.Get("reference_id").(string))
	if err != nil {
		return err
	}

	name := d.Get("name").(string)
	description := d.Get("description").(string)
	reference_id := d.Get("reference_id").(string)

	data := rancherClient.ApiKey{
		Name:        name,
		Description: description,
		AccountId:   reference_id,
	}

	// var newApiKey rancherClient.ApiKey
	newApiKey, err := client.ApiKey.Create(&data)
	if err != nil {
		return err
	}

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"active", "activating", "removed", "removing"},
		Target:     []string{"active"},
		Refresh:    ApiKeyStateRefreshFunc(client, newApiKey.Id),
		Timeout:    10 * time.Minute,
		Delay:      1 * time.Second,
		MinTimeout: 3 * time.Second,
	}
	_, waitErr := stateConf.WaitForState()
	if waitErr != nil {
		return fmt.Errorf(
			"Error waiting for API Key (%s) to be created: %s", newApiKey.Id, waitErr)
	}

	d.SetId(newApiKey.Id)
	log.Printf("[INFO] ApiKey ID: %s", d.Id())
	d.Set("name", newApiKey.Name)
	d.Set("description", newApiKey.Description)
	d.Set("reference_id", newApiKey.AccountId)
	d.Set("public_value", newApiKey.PublicValue)
	d.Set("secret_value", newApiKey.SecretValue)
	d.Set("kind", newApiKey.Kind)
	log.Printf("[INFO] ApiKey PublicValue: %s", newApiKey.PublicValue)
	log.Printf("[INFO] ApiKey SecretValue: %s", newApiKey.SecretValue)

	return nil
}

func resourceRancherApiKeyRead(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Refreshing ApiKey: %s", d.Id())
	log.Printf("[INFO] ApiKey reference_id: %s", d.Get("reference_id"))
	client, err := meta.(*Config).ApiKeyClient(d.Get("reference_id").(string))
	if err != nil {
		return err
	}

	apiKey, err := client.ApiKey.ById(d.Id())
	if err != nil {
		return err
	}

	if apiKey == nil {
		log.Printf("[INFO] ApiKey %s not found", d.Id())
		d.SetId("")
		return nil
	}

	if removed(apiKey.State) {
		log.Printf("[INFO] API Key %s was removed on %v", d.Id(), apiKey.Removed)
		d.SetId("")
		return nil
	}

	log.Printf("[INFO] ApiKey Name: %s", apiKey.Name)

	d.Set("name", apiKey.Name)
	d.Set("description", apiKey.Description)
	d.Set("reference_id", apiKey.AccountId)
	d.Set("public_value", apiKey.PublicValue)
	if apiKey.SecretValue != "" && d.Get("secret_value") != nil {
		d.Set("secret_value", apiKey.SecretValue)
	}
	d.Set("kind", apiKey.Kind)

	return nil
}

func resourceRancherApiKeyDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Deleting ApiKey: %s", d.Id())
	id := d.Id()
	client, err := meta.(*Config).ApiKeyClient(d.Get("reference_id").(string))
	if err != nil {
		return err
	}

	apiKey, err := client.ApiKey.ById(id)
	if err != nil {
		return err
	}

	// Step 1: Deactivate
	if _, e := client.ApiKey.ActionDeactivate(apiKey); e != nil {
		return fmt.Errorf("Error deactivating ApiKey: %s", err)
	}

	log.Printf("[DEBUG] Waiting for API Key (%s) to be deactivated", id)

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"active", "inactive", "deactivating"},
		Target:     []string{"inactive"},
		Refresh:    ApiKeyStateRefreshFunc(client, id),
		Timeout:    10 * time.Minute,
		Delay:      1 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, waitErr := stateConf.WaitForState()
	if waitErr != nil {
		return fmt.Errorf(
			"Error waiting for API Key (%s) to be deactivated: %s", id, waitErr)
	}

	// Update resource to reflect its state
	apiKey, err = client.ApiKey.ById(id)
	if err != nil {
		return fmt.Errorf("Failed to refresh state of deactivated API Key (%s): %s", id, err)
	}

	// Step 2: Remove
	if _, err := client.ApiKey.ActionRemove(apiKey); err != nil {
		return fmt.Errorf("Error removing ApiKey: %s", err)
	}

	log.Printf("[DEBUG] Waiting for API Key (%s) to be removed", id)

	stateConf = &resource.StateChangeConf{
		Pending:    []string{"inactive", "removed", "removing"},
		Target:     []string{"removed"},
		Refresh:    ApiKeyStateRefreshFunc(client, id),
		Timeout:    10 * time.Minute,
		Delay:      1 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, waitErr = stateConf.WaitForState()
	if waitErr != nil {
		return fmt.Errorf(
			"Error waiting for API Key (%s) to be removed: %s", id, waitErr)
	}

	d.SetId("")
	return nil
}

func resourceRancherApiKeyUpdate(d *schema.ResourceData, meta interface{}) error {
	//if d.HasChange("host_labels") {
	//newCommand := addHostLabels(
	//d.Get("command").(string),
	//d.Get("host_labels").(map[string]interface{}))
	//d.Set("command", newCommand)
	//}
	return resourceRancherApiKeyRead(d, meta)
}

func resourceRancherApiKeyImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	envID, resourceID := splitID(d.Id())
	d.SetId(resourceID)
	if envID != "" {
		d.Set("reference_id", envID)
	} else {
		client, err := meta.(*Config).GlobalClient()
		if err != nil {
			return []*schema.ResourceData{}, err
		}
		token, err := client.ApiKey.ById(d.Id())
		if err != nil {
			return []*schema.ResourceData{}, err
		}
		d.Set("reference_id", token.AccountId)
	}
	return []*schema.ResourceData{d}, nil
}

// ApiKeyStateRefreshFunc returns a resource.StateRefreshFunc that is used to watch
// a Rancher ApiKey.
func ApiKeyStateRefreshFunc(client *rancherClient.RancherClient, apiKeyID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		apiKey, err := client.ApiKey.ById(apiKeyID)

		if err != nil {
			return nil, "", err
		}

		return apiKey, apiKey.State, nil
	}
}
