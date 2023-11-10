---
page_title: "cpln_gvc Resource - terraform-provider-cpln"
subcategory: "Global Virtual Cloud"
description: |-
---

# cpln_gvc (Resource)

Manages an org's [Global Virtual Cloud (GVC)](https://docs.controlplane.com/reference/gvc).

## Declaration

### Required

- **name** (String) Name of the GVC.
- **locations** (List of String) A list of [locations](https://docs.controlplane.com/reference/location#current) making up the Global Virtual Cloud.

### Optional

- **description** (String) Description of the GVC.
- **tags** (Map of String) Key-value map of resource tags.
- **domain** (String) Custom domain name used by associated workloads.
- **pull_secrets** (List of String) A list of [pull secret](https://docs.controlplane.com/reference/gvc#pull-secrets) names used to authenticate to any private image repository referenced by Workloads within the GVC.
- **env** (Array of Name-Value Pair) Key-value array of resource env variables.
- **load_balancer** (Block List, Max: 1) ([see below](#nestedblock--load_balancer))
- **lightstep_tracing** (Block List, Max: 1) ([see below](#nestedblock--lightstep_tracing)).
- **otel_tracing** (Block List, Max: 1) ([see below](#nestedblock--otel_tracing)).

~> **Note** Only one of the tracing blocks can be defined.

<a id="nestedblock--lightstep_tracing"></a>

### `lightstep_tracing`

Required:

- **sampling** (Int) Sampling percentage.
- **endpoint** (String) Tracing Endpoint Workload. Either the canonical endpoint or the internal endpoint.

Optional:

- **credentials** (String) Full link to referenced Opaque Secret.

<a id="nestedblock--otel_tracing"></a>

### `otel_tracing`

Required:

- **sampling** (Int) Sampling percentage.
- **endpoint** (String) Tracing Endpoint Workload. Either the canonical endpoint or the internal endpoint.

<a id="nestedblock--load_balancer"></a>

### `load_balancer`

Required:

- **dedicated** (Boolean) Creates a dedicated load balancer in each location and enables additional Domain features: custom ports, protocols and wildcard hostnames. Charges apply for each location.

## Outputs

The following attributes are exported:

- **self_link** (String) Full link to this resource. Can be referenced by other resources.

## Example Usage

```terraform

resource "cpln_secret" "docker" {
  name        = "docker-secret"
  description = "docker secret"

  tags = {
    terraform_generated = "true"
    acceptance_test     = "true"
    secret_type         = "docker"
  }

  docker = "{\"auths\":{\"your-registry-server\":{\"username\":\"your-name\",\"password\":\"your-pword\",\"email\":\"your-email\",\"auth\":\"<Secret>\"}}}"
}

resource "cpln_secret" "opaque" {

  name        = "opaque-random-tbd"
  description = "description opaque-random-tbd"

  tags = {
    terraform_generated = "true"
    acceptance_test     = "true"
    secret_type         = "opaque"
  }

  opaque {
    payload  = "opaque_secret_payload"
    encoding = "plain"
  }
}

resource "cpln_gvc" "example" {

  name        = "gvc-example"
  description = "Example GVC"

  # Example Locations: `aws-eu-central-1`, `aws-us-west-2`, `azure-east2`, `gcp-us-east1`
  locations = ["aws-eu-central-1", "aws-us-west-2"]

  # domain = "app.example.com"
  pull_secrets = [cpln_secret.docker.name]

  tags = {
    terraform_generated = "true"
    example             = "true"
  }

  env = {
    env_var_key          = "env_var_value"
    workload_can_inherit = "true"
  }

  lightstep_tracing {

    sampling = 50
    endpoint = "test.cpln.local:8080"

    // Opaque Secret Only
    credentials = cpln_secret.opaque.self_link
  }

  load_balancer {
    dedicated = true
  }

}

```

## Import Syntax

The `terraform import` command is used to bring existing infrastructure resources, created outside of Terraform, into the Terraform state file, enabling their management through Terraform going forward.

To update a statefile with an existing GVC resource, execute the following import command:

```terraform
terraform import cpln_gvc.RESOURCE_NAME GVC_NAME
```

-> 1. Substitute RESOURCE_NAME with the same string that is defined in the HCL file.<br/>2. Substitute GVC_NAME with the corresponding GVC defined in the resource.
