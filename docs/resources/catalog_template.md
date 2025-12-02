---
page_title: "cpln_catalog_template Resource - terraform-provider-cpln"
subcategory: "Catalog"
description: |-
---

# cpln_catalog_template (Resource)

Manages a Control Plane [Catalog Template](https://github.com/controlplane-com/templates). This resource allows you to install, update, and uninstall applications from the Control Plane marketplace catalog.

You can browse available catalog templates in the [Control Plane templates repository](https://github.com/controlplane-com/templates). These are the same templates available in the Control Plane console's catalog section. Once installed, catalog template releases will appear in the Releases page in the UI.

~> **Important** The token or service account used to manage catalog templates must have the `reveal` permission for secrets. This is required to read the release state from helm release secrets. For more information about secret permissions, see the [Reveal Permission documentation](https://docs.controlplane.com/reference/secret#reveal-permission).

## Declaration

### Required

- **name** (String) The release name for this catalog template release.
- **template** (String) The name of the catalog template to deploy (e.g., 'postgres', 'redis', 'nginx').
- **version** (String) The version of the catalog template to deploy.
- **values** (String) The values file content (YAML format) for customizing the template release.

### Optional

- **gvc** (String) The GVC where the template will be deployed. Leave empty if the template creates its own GVC (check template's createsGvc field).

~> **Note** The `name`, `template`, and `gvc` fields require resource replacement if changed.

## Outputs

The following attributes are exported:

- **id** (String) The unique identifier for this catalog template release (same as name).
- **resources** (List of Object) List of Control Plane resources created by this release. Each resource contains:
  - **kind** (String) The kind of resource (e.g., 'workload', 'secret', 'gvc').
  - **name** (String) The name of the resource.
  - **link** (String) The full Control Plane link to the resource.

## Example Usage

### Template with GVC

```terraform
resource "cpln_gvc" "catalog_test_gvc" {
  name = "my-gvc"
}

resource "cpln_catalog_template" "redis" {
  name     = "my-redis"
  template = "redis"
  version  = "3.0.1"
  gvc      = cpln_gvc.catalog_test_gvc.name

  values = <<-EOT
redis:
  image: redis/redis-stack:7.4.0-v3
  resources:
    cpu: 200m
    memory: 256Mi
    minCpu: 80m
    minMemory: 128Mi
  replicas: 3
  timeoutSeconds: 15
  auth:
    fromSecret:
      enabled: false
      name: example-redis-auth-password
      passwordKey: password
    password:
      enabled: false
      value: your-password
  serverCommand: redis-stack-server
  publicAccess:
    enabled: false
    address: redis-test.example-cpln.com
  firewall:
    internal_inboundAllowType: "same-gvc"
    external_inboundAllowCIDR: 0.0.0.0/0
    external_outboundAllowCIDR: "0.0.0.0/0"
  env: []
  dataDir: /data
  persistence:
    enabled: false
    volumes:
      data:
        initialCapacity: 10
        performanceClass: general-purpose-ssd
        fileSystemType: ext4
        snapshots:
          retentionDuration: 7d
          schedule: 0 0 * * *
        autoscaling:
          maxCapacity: 100
          minFreePercentage: 20
          scalingFactor: 1.2

sentinel:
  image: redis/redis-stack:7.4.0-v3
  resources:
    cpu: 200m
    memory: 256Mi
    minCpu: 80m
    minMemory: 128Mi
  replicas: 3
  timeoutSeconds: 10
  quorumAutoCalculation: true
  quorumOverride: null
  auth:
    fromSecret:
      enabled: false
      name: example-redis-auth-password
      passwordKey: password
    password:
      enabled: false
      value: your-password
  publicAccess:
    enabled: false
    address: redis-sentinel-test.example-cpln.com
  firewall:
    internal_inboundAllowType: "same-gvc"
    external_inboundAllowCIDR: 0.0.0.0/0
    external_outboundAllowCIDR: "0.0.0.0/0"
  env: []
  persistence:
    enabled: false
    volumes:
      data:
        initialCapacity: 10
        performanceClass: general-purpose-ssd
        fileSystemType: ext4
        snapshots:
          retentionDuration: 7d
          schedule: 0 0 * * *
        autoscaling:
          maxCapacity: 50
          minFreePercentage: 20
          scalingFactor: 1.2
EOT
}
```

### Template Without GVC (Creates Its Own)

```terraform
resource "cpln_catalog_template" "cockroach" {
  name     = "my-cockroach"
  template = "cockroach"
  version  = "1.0.0"

  values = <<-EOT
gvc:
  name: my-gvc
  locations:
    - name: aws-eu-central-1
      replicas: 3

resources:
  cpu: 2000m
  memory: 4096Mi

database:
  name: mydb
  user: myuser

cockroach_defaults:
  workload_name: cockroach
  sql_port: 26257
  http_port: 8080

internal_access:
  type: same-gvc # options: same-gvc, same-org, workload-list
  workloads:  # Note: can only be used if type is same-gvc or workload-list
    #- //gvc/GVC_NAME/workload/WORKLOAD_NAME
    #- //gvc/GVC_NAME/workload/WORKLOAD_NAME
EOT
}
```

### Accessing Created Resources

```terraform
resource "cpln_gvc" "catalog_test_gvc" {
  name = "my-gvc"
}

resource "cpln_catalog_template" "postgres" {
  name     = "my-postgres"
  template = "postgres"
  version  = "2.0.0"
  gvc      = cpln_gvc.catalog_test_gvc.name

  values = <<-EOT
resources:
  cpu: 500m
  memory: 1024Mi

config:
  username: username
  password: password
  database: test
EOT
}

# Output the resources created by the catalog template
output "postgres_resources" {
  value = cpln_catalog_template.postgres.resources
}
```

## Import Syntax

The `terraform import` command is used to bring existing infrastructure resources, created outside of Terraform, into the Terraform state file, enabling their management through Terraform going forward.

To update a statefile with an existing catalog template resource, execute the following import command:

```terraform
terraform import cpln_catalog_template.RESOURCE_NAME RELEASE_NAME
```

-> 1. Substitute RESOURCE_NAME with the same string that is defined in the HCL file.<br/>2. Substitute RELEASE_NAME with the name of the installed catalog template release.
