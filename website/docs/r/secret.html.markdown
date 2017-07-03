---
layout: "rancher"
page_title: "Rancher: rancher_secret"
sidebar_current: "docs-rancher-resource-secret"
description: |-
  Provides a Rancher Secret resource. This can be used to create secrets for rancher environments and retrieve their information.
---

# rancher\_secrets

Provides a Rancher Secret resource. This can be used to create secrets for rancher environments and retrieve their information.

## Example Usage

```hcl
# Create a new Rancher Secret
resource rancher_secret "foo" {
  name           = "foo"
  environment_id = "${rancher_environment.test.id}"
  value          = "my great password"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the secret.
* `description` - (Optional) A description of the secret.
* `environment_id` - (Required) The ID of the environment to create the secret for.
* `value` - (Required) The secret value.


## Import

Secrets can be imported using the Secret ID in the format
`<environment_id>/<secret_id>`

```
$ terraform import rancher_secret.mysec 1a5/1se10
```

If the credentials for the Rancher provider have access to the global API,
then `environment_id` can be omitted e.g.

```
$ terraform import rancher_secret.mysec 1se10
```
