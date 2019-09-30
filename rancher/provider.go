package rancher

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/url"
	"os"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

// CLIConfig used to store data from file.
type CLIConfig struct {
	AccessKey   string `json:"accessKey"`
	SecretKey   string `json:"secretKey"`
	URL         string `json:"url"`
	Environment string `json:"environment"`
	Path        string `json:"path,omitempty"`
}

// Provider returns a terraform.ResourceProvider.
func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"api_url": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("RANCHER_URL", ""),
				Description: descriptions["api_url"],
			},
			"access_key": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("RANCHER_ACCESS_KEY", ""),
				Description: descriptions["access_key"],
			},
			"secret_key": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("RANCHER_SECRET_KEY", ""),
				Description: descriptions["secret_key"],
			},
			"config": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("RANCHER_CLIENT_CONFIG", ""),
				Description: descriptions["config"],
			},
			"skip_config_validation": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: descriptions["skip_config_validation"],
			},
		},

		ResourcesMap: map[string]*schema.Resource{
			"rancher_certificate":         resourceRancherCertificate(),
			"rancher_environment":         resourceRancherEnvironment(),
			"rancher_host":                resourceRancherHost(),
			"rancher_registration_token":  resourceRancherRegistrationToken(),
			"rancher_registry":            resourceRancherRegistry(),
			"rancher_registry_credential": resourceRancherRegistryCredential(),
			"rancher_secret":              resourceRancherSecret(),
			"rancher_stack":               resourceRancherStack(),
			"rancher_volume":              resourceRancherVolume(),
		},

		DataSourcesMap: map[string]*schema.Resource{
			"rancher_certificate": dataSourceRancherCertificate(),
			"rancher_environment": dataSourceRancherEnvironment(),
			"rancher_setting":     dataSourceRancherSetting(),
		},

		ConfigureFunc: providerConfigure,
	}
}

var descriptions map[string]string

func init() {
	descriptions = map[string]string{
		"access_key": "API Key used to authenticate with the rancher server",

		"secret_key": "API secret used to authenticate with the rancher server",

		"api_url": "The URL to the rancher API, must include version uri (ie. v1 or v2-beta)",

		"config": "Path to the Rancher client cli.json config file",

		"skip_config_validation": "Skip the configuration parameters validation.",
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	apiURL := d.Get("api_url").(string)
	accessKey := d.Get("access_key").(string)
	secretKey := d.Get("secret_key").(string)

	if scv := d.Get("skip_config_validation").(bool); scv {
		config := &Config{
			APIURL:    apiURL,
			AccessKey: accessKey,
			SecretKey: secretKey,
		}
		return config, nil
	}

	if configFile := d.Get("config").(string); configFile != "" {
		config, err := loadConfig(configFile)
		if err != nil {
			return config, err
		}

		if apiURL == "" && config.URL != "" {
			u, err := url.Parse(config.URL)
			if err != nil {
				return config, err
			}
			apiURL = u.Scheme + "://" + u.Host
		}

		if accessKey == "" {
			accessKey = config.AccessKey
		}

		if secretKey == "" {
			secretKey = config.SecretKey
		}
	}

	if apiURL == "" {
		return &Config{}, fmt.Errorf("No api_url provided")
	}

	config := &Config{
		APIURL:    apiURL,
		AccessKey: accessKey,
		SecretKey: secretKey,
	}

	client, err := config.GlobalClient()
	if err != nil {
		return &Config{}, err
	}
	// Let Rancher Client normalizes the URL making it reliable as a base.
	config.APIURL = client.GetOpts().Url

	return config, nil
}

func loadConfig(path string) (CLIConfig, error) {
	config := CLIConfig{
		Path: path,
	}

	content, err := ioutil.ReadFile(path)
	if os.IsNotExist(err) {
		return config, nil
	} else if err != nil {
		return config, err
	}

	err = json.Unmarshal(content, &config)
	config.Path = path

	return config, err
}
