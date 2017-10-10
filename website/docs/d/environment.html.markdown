---
layout: "rancher"
page_title: "Rancher: rancher_environment"
sidebar_current: "docs-rancher-datasource-environment"
description: |-
  Get information on a Rancher environment.
---

# rancher\_environment

Use this data source to retrieve information about a Rancher environment.

## Example Usage

```hcl
data "rancher_environment" "foo" {
  name = "foo"
}
```

## Argument Reference

 * `name` - (Required) The setting name.

## Attributes Reference

* `id` - The ID of the resource.
* `description` - The environment description.
* `orchestration` - The environment orchestration engine.
* `project_template_id` - The environment project template ID.
* `member` - The environment members.
