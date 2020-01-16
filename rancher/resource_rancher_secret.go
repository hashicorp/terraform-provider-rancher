package rancher

import (
	"encoding/base64"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	rancher "github.com/rancher/go-rancher/v2"
)

func resourceRancherSecret() *schema.Resource {
	return &schema.Resource{
		Create: resourceRancherSecretCreate,
		Read:   resourceRancherSecretRead,
		Update: resourceRancherSecretUpdate,
		Delete: resourceRancherSecretDelete,
		Importer: &schema.ResourceImporter{
			State: resourceRancherSecretImport,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"environment_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"value": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceRancherSecretCreate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO][rancher] Creating Secret: %s", d.Id())
	client, err := meta.(*Config).EnvironmentClient(d.Get("environment_id").(string))
	if err != nil {
		return err
	}

	name := d.Get("name").(string)
	description := d.Get("description").(string)
	value := d.Get("value").(string)

	secret := rancher.Secret{
		Name:        name,
		Description: description,
		Value:       base64.StdEncoding.EncodeToString([]byte(value)),
	}
	newSecret, err := client.Secret.Create(&secret)
	if err != nil {
		return err
	}

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"active", "removed", "removing"},
		Target:     []string{"active"},
		Refresh:    SecretStateRefreshFunc(client, newSecret.Id),
		Timeout:    10 * time.Minute,
		Delay:      1 * time.Second,
		MinTimeout: 3 * time.Second,
	}
	_, waitErr := stateConf.WaitForState()
	if waitErr != nil {
		return fmt.Errorf(
			"Error waiting for secret (%s) to be created: %s", newSecret.Id, waitErr)
	}

	d.SetId(newSecret.Id)
	log.Printf("[INFO] Secret ID: %s", d.Id())

	return resourceRancherSecretUpdate(d, meta)
}

func resourceRancherSecretRead(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Refreshing Secret: %s", d.Id())
	client, err := meta.(*Config).EnvironmentClient(d.Get("environment_id").(string))
	if err != nil {
		return err
	}

	secret, err := client.Secret.ById(d.Id())
	if err != nil {
		return err
	}

	if secret == nil {
		return fmt.Errorf("Failed to find secret %s", d.Id())
	}

	log.Printf("[INFO] Secret Name: %s", secret.Name)

	d.Set("description", secret.Description)
	d.Set("name", secret.Name)

	return nil
}

func resourceRancherSecretUpdate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Updating Secret: %s", d.Id())
	client, err := meta.(*Config).EnvironmentClient(d.Get("environment_id").(string))
	if err != nil {
		return err
	}

	secret, err := client.Secret.ById(d.Id())
	if err != nil {
		return err
	}

	name := d.Get("name").(string)
	description := d.Get("description").(string)
	value := d.Get("value").(string)

	data := map[string]interface{}{
		"name":        &name,
		"description": &description,
		"value":       base64.StdEncoding.EncodeToString([]byte(value)),
	}

	var newSecret rancher.Secret
	if err := client.Update("secret", &secret.Resource, data, &newSecret); err != nil {
		return err
	}

	return resourceRancherSecretRead(d, meta)
}

func resourceRancherSecretDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Deleting Secret: %s", d.Id())
	id := d.Id()
	client, err := meta.(*Config).EnvironmentClient(d.Get("environment_id").(string))
	if err != nil {
		return err
	}

	secret, err := client.Secret.ById(id)
	if err != nil {
		return err
	}

	if err := client.Secret.Delete(secret); err != nil {
		return fmt.Errorf("Error deleting Secret: %s", err)
	}

	log.Printf("[DEBUG] Waiting for secret (%s) to be removed", id)

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"active", "removed", "removing"},
		Target:     []string{"removed"},
		Refresh:    SecretStateRefreshFunc(client, id),
		Timeout:    10 * time.Minute,
		Delay:      1 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, waitErr := stateConf.WaitForState()
	if waitErr != nil {
		return fmt.Errorf(
			"Error waiting for secret (%s) to be removed: %s", id, waitErr)
	}

	d.SetId("")
	return nil
}

func resourceRancherSecretImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	envID, resourceID := splitID(d.Id())
	d.SetId(resourceID)
	if envID != "" {
		d.Set("environment_id", envID)
	} else {
		client, err := meta.(*Config).GlobalClient()
		if err != nil {
			return []*schema.ResourceData{}, err
		}
		sec, err := client.Secret.ById(d.Id())
		if err != nil {
			return []*schema.ResourceData{}, err
		}
		d.Set("environment_id", sec.AccountId)
	}
	return []*schema.ResourceData{d}, nil
}

// SecretStateRefreshFunc returns a resource.StateRefreshFunc that is used to watch
// a Rancher Secret.
func SecretStateRefreshFunc(client *rancher.RancherClient, secretID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		cert, err := client.Secret.ById(secretID)

		if err != nil {
			return nil, "", err
		}

		return cert, cert.State, nil
	}
}
