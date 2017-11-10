package rancher

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	rancher "github.com/rancher/go-rancher/v2"
)

func dataSourceRancherCertificate() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceRancherCertificateRead,

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"environment_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"description": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"cn": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"algorithm": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"cert_fingerprint": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"expires_at": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"issued_at": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"issuer": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"key_size": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"serial_number": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"subject_alternative_names": &schema.Schema{
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"version": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceRancherCertificateRead(d *schema.ResourceData, meta interface{}) error {
	name := d.Get("name").(string)
	log.Printf("[INFO] Refreshing Rancher Certificate: %s", name)

	environmentId := d.Get("environment_id").(string)
	client, err := meta.(*Config).EnvironmentClient(environmentId)
	if err != nil {
		return err
	}

	stateConf := &resource.StateChangeConf{
		Pending:        []string{"active", "removed", "removing", "not found"},
		Target:         []string{"active"},
		Refresh:        findCert(client, name),
		Timeout:        10 * time.Minute,
		Delay:          1 * time.Second,
		MinTimeout:     3 * time.Second,
		NotFoundChecks: 50,
	}
	cert, waitErr := stateConf.WaitForState()
	if waitErr != nil {
		return fmt.Errorf(
			"Error waiting for certificate (%s) to be found: %s", name, waitErr)
	}

	certificate := cert.(rancher.Certificate)
	d.SetId(certificate.Id)

	d.Set("description", certificate.Description)
	d.Set("name", certificate.Name)

	// Computed values
	d.Set("cn", certificate.CN)
	d.Set("algorithm", certificate.Algorithm)
	d.Set("cert_fingerprint", certificate.CertFingerprint)
	d.Set("expires_at", certificate.ExpiresAt)
	d.Set("issued_at", certificate.IssuedAt)
	d.Set("issuer", certificate.Issuer)
	d.Set("serial_number", certificate.SerialNumber)
	d.Set("subject_alternative_names", certificate.SubjectAlternativeNames)
	d.Set("version", certificate.Version)

	return nil
}

func findCert(client *rancher.RancherClient, certname string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		certs, err := client.Certificate.List(NewListOpts())
		if err != nil {
			return nil, "", err
		}

		for _, cert := range certs.Data {
			if cert.Name == certname {
				return cert, cert.State, nil
			}
		}

		return nil, "not found", nil
	}
}
