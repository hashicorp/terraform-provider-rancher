## 1.0.1 (Unreleased)

BUG FIXES:

* resource/rancher_host: remove rancher_host's importer [GH-43]

## 1.0.0 (September 20, 2017)

IMPROVEMENTS:

* resource/rancher_stack: improve error messages ([#40](https://github.com/terraform-providers/terraform-provider-rancher/issues/40))

BUG FIXES:

* resource/rancher_host: wait when host is in 'activating' state ([#31](https://github.com/terraform-providers/terraform-provider-rancher/issues/31))
* resource/rancher_environment: fix member creation and orchestration update
  ([#37](https://github.com/terraform-providers/terraform-provider-rancher/issues/37))
* resource/rancher_stack: get system scope from API ([#39](https://github.com/terraform-providers/terraform-provider-rancher/issues/39))
* resource/rancher_stack: use all valid suffixes to retrieve catalog template
  files ([#41](https://github.com/terraform-providers/terraform-provider-rancher/issues/41))

## 0.2.0 (August 24, 2017)

IMPROVEMENTS:

* resource/rancher_registration_token: add 'agent_ip' argument ([#23](https://github.com/terraform-providers/terraform-provider-rancher/issues/23))

BUG FIXES:

* resource/rancher_environment: fix setting membership on creation ([#29](https://github.com/terraform-providers/terraform-provider-rancher/issues/29))
* resource/rancher_stack: wait when stack is in 'upgrading' state ([#21](https://github.com/terraform-providers/terraform-provider-rancher/issues/21))
* resource/rancher_stack: wait when stack is in 'registering' state ([#22](https://github.com/terraform-providers/terraform-provider-rancher/issues/22))

## 0.1.1 (July 04, 2017)

FEATURES:

* **New Data Source**: `rancher_setting` ([#13](https://github.com/terraform-providers/terraform-provider-rancher/issues/13))
* **New Resource:** `rancher_secret` ([#11](https://github.com/terraform-providers/terraform-provider-rancher/issues/11))

BUG FIXES:

* tests: Add orchestration parameter to fix acceptance tests ([#10](https://github.com/terraform-providers/terraform-provider-rancher/issues/10))
* resource/rancher_certificate: fix doc ([#12](https://github.com/terraform-providers/terraform-provider-rancher/issues/12))
* resource/rancher_environment: fix members casting when creating ([#18](https://github.com/terraform-providers/terraform-provider-rancher/issues/18))
* resource/rancher_host: wait for host to be created ([#17](https://github.com/terraform-providers/terraform-provider-rancher/issues/17))
* resource/rancher_host: deactivate host before deleting it ([#20](https://github.com/terraform-providers/terraform-provider-rancher/issues/20))
* fix missing attributes in doc ([#16](https://github.com/terraform-providers/terraform-provider-rancher/issues/16))

## 0.1.0 (June 21, 2017)

IMPROVEMENTS:

* Move to Rancher V2 API [[#13908](https://github.com/terraform-providers/terraform-provider-rancher/issues/13908)](https://github.com/hashicorp/terraform/pull/13908)
