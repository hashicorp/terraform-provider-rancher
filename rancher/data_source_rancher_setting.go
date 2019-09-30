package rancher

import (
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceRancherSetting() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceRancherSettingRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"value": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceRancherSettingRead(d *schema.ResourceData, meta interface{}) error {
	name := d.Get("name").(string)
	log.Printf("[INFO] Refreshing Rancher Setting: %s", name)

	client, err := meta.(*Config).GlobalClient()
	if err != nil {
		return err
	}

	setting, err := client.Setting.ById(name)
	if err != nil {
		return err
	}

	d.SetId(name)
	d.Set("value", setting.Value)

	return nil
}
