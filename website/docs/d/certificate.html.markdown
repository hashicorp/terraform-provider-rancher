---
layout: "rancher"
page_title: "Rancher: rancher_certificate"
sidebar_current: "docs-rancher-datasource-certificate"
description: |-
  Get information on a Rancher certificate.
---

# rancher\_certificate

Use this data source to retrieve information about a Rancher certificate.

## Example Usage

```hcl
data "rancher_certificate" "foo" {
  name = "foo"
  environment_id = "1a5"
}
```

## Argument Reference

 * `name` - (Required) The setting name.
 * `environment_id` - (Required) The ID of the environment.

## Attributes Reference

* `id` - The ID of the resource.
* `cn` - The certificate CN.
* `algorithm` - The certificate algorithm.
* `cert_fingerprint` - The certificate fingerprint.
* `expires_at` - The certificate expiration date.
* `issued_at` - The certificate creation date.
* `issuer` - The certificate issuer.
* `serial_number` - The certificate serial number.
* `subject_alternative_names` - The list of certificate Subject Alternative Names.
* `version` - The certificate version.
