## 1.5.1 (Unreleased)
## 1.5.0 (January 21, 2020)

IMPROVEMENTS:

* Use new Terraform plugin SDK ([#100](https://github.com/terraform-providers/terraform-provider-rancher/pull/100))
* Paginate list actions ([#102](https://github.com/terraform-providers/terraform-provider-rancher/pull/102))
* Fix acceptance tests ([#103](https://github.com/terraform-providers/terraform-provider-rancher/pull/103))

## 1.4.0 (July 26, 2019)

FEATURES:

* Add option to skip config validation ([#97](https://github.com/terraform-providers/terraform-provider-rancher/pull/97))

## 1.3.0 (June 04, 2019)

FEATURES:

* Prepare provider for Terraform v0.12 ([#96](https://github.com/terraform-providers/terraform-provider-rancher/pull/96))

IMPROVEMENTS:

* Upgrade to GO 1.11 ([#91](https://github.com/terraform-providers/terraform-provider-rancher/pull/91))
* Switch to Go Modules ([#94](https://github.com/terraform-providers/terraform-provider-rancher/pull/94))
* Use SVG badge in README ([#95](https://github.com/terraform-providers/terraform-provider-rancher/pull/95))

BUG FIXES:

* Recreate rancher_host resource when hostname changes ([#82](https://github.com/terraform-providers/terraform-provider-rancher/issues/82))
* Avoid crashing when registry ID cannot be found ([#92](https://github.com/terraform-providers/terraform-provider-rancher/pull/92))

## 1.2.1 (June 08, 2018)

BUG FIXES:

* resource/rancher_volume: changing environment_id or driver forces resource renewal ([#79](https://github.com/terraform-providers/terraform-provider-rancher/pull/79))

## 1.2.0 (December 19, 2017)

FEATURES:

* **New Resource:** `rancher_volume` ([#63](https://github.com/terraform-providers/terraform-provider-rancher/issues/63))

IMPROVEMENTS:

* doc/rancher_volume: add resource usage documentation ([#64](https://github.com/terraform-providers/terraform-provider-rancher/issues/64))
* resource/rancher_registry_credential: Deprecate email attribute ([#61](https://github.com/terraform-providers/terraform-provider-rancher/issues/61))

BUG FIXES:

* doc: add missing entries in doc sidebar ([#59](https://github.com/terraform-providers/terraform-provider-rancher/issues/59))

## 1.1.1 (November 10, 2017)

BUG FIXES:

* doc/rancher_certificate: add datasource usage documentation ([#58](https://github.com/terraform-providers/terraform-provider-rancher/issues/58))
* doc/rancher_registration_token: fix documentation ([#50](https://github.com/terraform-providers/terraform-provider-rancher/issues/50))
* datasource/rancher_certificate: increase NotFoundChecks ([#57](https://github.com/terraform-providers/terraform-provider-rancher/issues/57))
* resource/rancher_environment: add missing pending states ([#53](https://github.com/terraform-providers/terraform-provider-rancher/issues/53))

## 1.1.0 (October 11, 2017)

FEATURES:

* **New Data Source**: `rancher_certificate` ([#47](https://github.com/terraform-providers/terraform-provider-rancher/issues/47))
* **New Data Source**: `rancher_environment` ([#48](https://github.com/terraform-providers/terraform-provider-rancher/issues/48))

BUG FIXES:

* resource/rancher_host: remove rancher_host's importer ([#43](https://github.com/terraform-providers/terraform-provider-rancher/issues/43))
* resource/rancher_host: increase NotFoundChecks ([#46](https://github.com/terraform-providers/terraform-provider-rancher/issues/46))

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
