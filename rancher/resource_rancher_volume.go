package rancher

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	rancher "github.com/rancher/go-rancher/v2"
)

func resourceRancherVolume() *schema.Resource {
	return &schema.Resource{
		Create: resourceRancherVolumeCreate,
		Read:   resourceRancherVolumeRead,
		Update: resourceRancherVolumeUpdate,
		Delete: resourceRancherVolumeDelete,
		Importer: &schema.ResourceImporter{
			State: resourceRancherVolumeImport,
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
				ForceNew: true,
			},
			"driver": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceRancherVolumeCreate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO][rancher] Creating Volume: %s", d.Id())
	client, err := meta.(*Config).EnvironmentClient(d.Get("environment_id").(string))
	if err != nil {
		return err
	}

	name := d.Get("name").(string)
	description := d.Get("description").(string)
	driver := d.Get("driver").(string)

	volume := rancher.Volume{
		Name:        name,
		Description: description,
		Driver:      driver,
	}
	newVolume, err := client.Volume.Create(&volume)
	if err != nil {
		return err
	}

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"active", "inactive", "removed", "removing"},
		Target:     []string{"active", "inactive"},
		Refresh:    VolumeStateRefreshFunc(client, newVolume.Id),
		Timeout:    10 * time.Minute,
		Delay:      1 * time.Second,
		MinTimeout: 3 * time.Second,
	}
	_, waitErr := stateConf.WaitForState()
	if waitErr != nil {
		return fmt.Errorf(
			"Error waiting for volume (%s) to be created: %s", newVolume.Id, waitErr)
	}

	d.SetId(newVolume.Id)
	log.Printf("[INFO] Volume ID: %s", d.Id())

	return resourceRancherVolumeUpdate(d, meta)
}

func resourceRancherVolumeRead(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Refreshing Volume: %s", d.Id())
	client, err := meta.(*Config).EnvironmentClient(d.Get("environment_id").(string))
	if err != nil {
		return err
	}

	volume, err := client.Volume.ById(d.Id())
	if err != nil {
		return err
	}

	log.Printf("[INFO] Volume Name: %s", volume.Name)

	d.Set("description", volume.Description)
	d.Set("name", volume.Name)
	d.Set("driver", volume.Driver)

	return nil
}

func resourceRancherVolumeUpdate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Updating Volume: %s", d.Id())
	client, err := meta.(*Config).EnvironmentClient(d.Get("environment_id").(string))
	if err != nil {
		return err
	}

	volume, err := client.Volume.ById(d.Id())
	if err != nil {
		return err
	}

	name := d.Get("name").(string)
	description := d.Get("description").(string)
	driver := d.Get("driver").(string)

	data := map[string]interface{}{
		"name":        &name,
		"description": &description,
		"driver":      &driver,
	}

	var newVolume rancher.Volume
	if err := client.Update("volume", &volume.Resource, data, &newVolume); err != nil {
		return err
	}

	return resourceRancherVolumeRead(d, meta)
}

func resourceRancherVolumeDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Deleting Volume: %s", d.Id())
	id := d.Id()
	client, err := meta.(*Config).EnvironmentClient(d.Get("environment_id").(string))
	if err != nil {
		return err
	}

	volume, err := client.Volume.ById(id)
	if err != nil {
		return err
	}

	if err := client.Volume.Delete(volume); err != nil {
		return fmt.Errorf("Error deleting Volume: %s", err)
	}

	log.Printf("[DEBUG] Waiting for volume (%s) to be removed", id)

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"active", "removed", "removing"},
		Target:     []string{"removed"},
		Refresh:    VolumeStateRefreshFunc(client, id),
		Timeout:    10 * time.Minute,
		Delay:      1 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, waitErr := stateConf.WaitForState()
	if waitErr != nil {
		return fmt.Errorf(
			"Error waiting for volume (%s) to be removed: %s", id, waitErr)
	}

	d.SetId("")
	return nil
}

func resourceRancherVolumeImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	envID, resourceID := splitID(d.Id())
	d.SetId(resourceID)
	if envID != "" {
		d.Set("environment_id", envID)
	} else {
		client, err := meta.(*Config).GlobalClient()
		if err != nil {
			return []*schema.ResourceData{}, err
		}
		stack, err := client.Stack.ById(d.Id())
		if err != nil {
			return []*schema.ResourceData{}, err
		}
		d.Set("environment_id", stack.AccountId)
	}
	return []*schema.ResourceData{d}, nil
}

// VolumeStateRefreshFunc returns a resource.StateRefreshFunc that is used to watch
// a Rancher Volume.
func VolumeStateRefreshFunc(client *rancher.RancherClient, volumeID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		cert, err := client.Volume.ById(volumeID)

		if err != nil {
			return nil, "", err
		}

		return cert, cert.State, nil
	}
}
