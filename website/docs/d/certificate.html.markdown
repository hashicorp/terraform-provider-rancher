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

### Simple datasource declaration

```hcl
data "rancher_certificate" "foo" {
  name           = "foo"
  environment_id = "1a5"
}
```

### Let's encrypt with DNS challenge

This setup will ensure that the Load Balancer stack is not created before the Let's Encrypt's certificate is actually present in Rancher's certificates manager.

```hcl
locals {
  environment_id = "1a5"
}

resource "rancher_stack" "letsencrypt" {
  name            = "letsencrypt"
  environment_id  = "${local.environment_id}"
  catalog_id      = "community:letsencrypt:4"

  environment {
    CERT_NAME      = "letsencrypt"
    DOMAINS        = "foo.example.com"
    PROVIDER       = "Route53"
    AWS_ACCESS_KEY = "${var.aws_access_key}"
    AWS_SECRET_KEY = "${var.aws_secret_key}"
    ...
  }
}

data "rancher_certificate" "letsencrypt" {
  environment_id = "${local.environment_id}"
  name           = "${rancher_stack.letsencrypt.environment["CERT_NAME"]}"
}

resource "rancher_stack" "lb" {
  name           = "lb"
  environment_id = "${local.environment_id}"

  docker_compose = <<EOF
version: '2'
services:
  lb:
    image: rancher/lb-service-haproxy:v0.7.9
    ports:
    - 443:443/tcp
    labels:
      io.rancher.container.agent.role: environmentAdmin
      io.rancher.container.create_agent: 'true'
EOF

  rancher_compose = <<EOF
version: '2'
services:
  lb:
    scale: 1
    start_on_create: true
    lb_config:
      certs: []
      default_cert: ${data.rancher_certificate.letsencrypt.name}
      port_rules:
      - protocol: https
        service: mystack/myservice
        source_port: 443
        target_port: 80
    health_check:
      healthy_threshold: 2
      response_timeout: 2000
      port: 42
      unhealthy_threshold: 3
      interval: 2000
      strategy: recreate
EOF
}
```

### Let's encrypt with HTTP challenge

This setup will ensure that the HTTPS Load Balancer stack is not created before the Let's Encrypt's certificate is actually present in Rancher's certificates manager.

```hcl
locals {
  environment_id = "1a5"
}

resource "rancher_stack" "letsencrypt" {
  name            = "letsencrypt"
  environment_id  = "${local.environment_id}"
  catalog_id      = "community:letsencrypt:4"

  environment {
    CERT_NAME      = "letsencrypt"
    DOMAINS        = "foo.example.com"
    PROVIDER       = "HTTP"
    ...
  }
}

resource "rancher_stack" "lb-http" {
  name           = "lb-http"
  environment_id = "${local.environment_id}"

  docker_compose = <<EOF
version: '2'
services:
  lb:
    image: rancher/lb-service-haproxy:v0.7.9
    ports:
    - 80:80/tcp
    labels:
      io.rancher.container.agent.role: environmentAdmin
      io.rancher.container.create_agent: 'true'
EOF

  rancher_compose = <<EOF
version: '2'
services:
  lb:
    scale: 1
    start_on_create: true
    lb_config:
      certs: []
      - hostname: ''
        path: /.well-known/acme-challenge
        priority: 1
        protocol: http
        service: letsencrypt/letsencrypt
        source_port: 80
        target_port: 80
    health_check:
      healthy_threshold: 2
      response_timeout: 2000
      port: 42
      unhealthy_threshold: 3
      interval: 2000
      strategy: recreate
EOF
}

data "rancher_certificate" "letsencrypt" {
  environment_id = "${local.environment_id}"
  name           = "${rancher_stack.letsencrypt.environment["CERT_NAME"]}"
}

resource "rancher_stack" "lb-https" {
  name           = "lb-https"
  environment_id = "${local.environment_id}"

  docker_compose = <<EOF
version: '2'
services:
  lb:
    image: rancher/lb-service-haproxy:v0.7.9
    ports:
    - 443:443/tcp
    labels:
      io.rancher.container.agent.role: environmentAdmin
      io.rancher.container.create_agent: 'true'
EOF

  rancher_compose = <<EOF
version: '2'
services:
  lb:
    scale: 1
    start_on_create: true
    lb_config:
      certs: []
      default_cert: ${data.rancher_certificate.letsencrypt.name}
      port_rules:
      - protocol: https
        service: mystack/myservice
        source_port: 443
        target_port: 80
    health_check:
      healthy_threshold: 2
      response_timeout: 2000
      port: 42
      unhealthy_threshold: 3
      interval: 2000
      strategy: recreate
EOF
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
