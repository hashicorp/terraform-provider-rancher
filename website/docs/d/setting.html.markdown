---
layout: "rancher"
page_title: "Rancher: rancher_setting"
sidebar_current: "docs-rancher-datasource-setting"
description: |-
  Get information on a Rancher setting.
---

# rancher\_setting

Use this data source to retrieve information about a Rancher setting.

## Example Usage

```
data "rancher_setting" "cattle.cattle.version" {
  name = "cattle.cattle.version"
}
```

## Argument Reference

 * `name` - (Required) The setting name.

## Attributes Reference

 * `value` - the settting's value.
