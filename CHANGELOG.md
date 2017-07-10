## 0.1.2 (Unreleased)

BUG FIXES:

* resource/rancher_stack: wait when stack is in 'upgrading' state [GH-21]
* resource/rancher_stack: wait when stack is in 'registering' state [GH-22]

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
