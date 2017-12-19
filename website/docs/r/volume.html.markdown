---
layout: "rancher"
page_title: "Rancher: rancher_volume"
sidebar_current: "docs-rancher-resource-volume"
description: |-
  Provides a Rancher Volume resource. This can be used to create volumes for rancher environments and retrieve their information.
---

# rancher\_volumes

Provides a Rancher Volume resource. This can be used to create volumes for rancher environments and retrieve their information.

## Example Usage

```hcl
# Create a new Rancher Volume
resource rancher_volume "foo" {
  name           = "foo"
  environment_id = "${rancher_environment.test.id}"
  driver         = "rancher-nfs"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the volume.
* `description` - (Optional) A description of the volume.
* `environment_id` - (Required) The ID of the environment to create the volume for.
* `driver` - (Required) The volume driver.


## Import

Volumes can be imported using the Volume ID in the format
`<environment_id>/<volume_id>`

```
$ terraform import rancher_volume.mysec 1a5/1v123456
```

If the credentials for the Rancher provider have access to the global API,
then `environment_id` can be omitted e.g.

```
$ terraform import rancher_volume.mysec 1se10
```
