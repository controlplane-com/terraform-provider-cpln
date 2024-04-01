---
page_title: "cpln_org_Tracing Resource - terraform-provider-cpln"
subcategory: "Org Tracing"
description: |-
---

# cpln_org_tracing (Resource)

Manages an org's tracing configuration.

## Declaration

### Required

Only one of the following tracing blocks can be defined:

- **lightstep_tracing** (Block List, Max: 1) ([see below](#nestedblock--lightstep_tracing)).
- **otel_tracing** (Block List, Max: 1) ([see below](#nestedblock--otel_tracing)).
- **controlplane_tracing** (Block List, Max: 1) ([see below](#nestedblock--controlplane_tracing)).

<a id="nestedblock--lightstep_tracing"></a>

### `lightstep_tracing`

Required:

- **sampling** (Int) Determines what percentage of requests should be traced.
- **endpoint** (String) Tracing Endpoint Workload. Either the canonical endpoint or the internal endpoint.

Optional:

- **credentials** (String) Full link to referenced Opaque Secret.
- **custom_tags** (Map of String) Key-value map of custom tags.

<a id="nestedblock--otel_tracing"></a>

### `otel_tracing`

Required:

- **sampling** (Int) Determines what percentage of requests should be traced.
- **endpoint** (String) Tracing Endpoint Workload. Either the canonical endpoint or the internal endpoint.

Optional:

- **custom_tags** (Map of String) Key-value map of custom tags.

<a id="nestedblock--controlplane_tracing"></a>

### `controlplane_tracing`

Required:

- **sampling** (Int) Determines what percentage of requests should be traced.

Optional:

- **custom_tags** (Map of String) Key-value map of custom tags.

## Outputs

The following attributes are exported:

- **cpln_id** (String) ID, in GUID format, of the org.
- **name** (String) The name of org.
- **description** (String) The description of org.
- **tags** (Map of String) Key-value map of the org's tags.

## Example Usage

### Lightstep

```terraform
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

resource "cpln_org_tracing" "new" {

  lightstep_tracing {

    sampling = 50
    endpoint = "test.cpln.local:8080"

    // Opaque Secret Only
    credentials = cpln_secret.opaque.self_link

    custom_tags = {
      key = "value"
    }
  }
}
```

### Otel

```terraform

resource "cpln_org_tracing" "new" {

  otel_tracing {

    sampling = 50
    endpoint = "test.cpln.local:8080"

    custom_tags = {
      key = "value"
    }
  }
}
```

### Control Plane

```terraform

resource "cpln_org_tracing" "new" {

  controlplane_tracing {

    sampling = 50

    custom_tags = {
      key = "value"
    }
  }
}
```

## Import Syntax

The `terraform import` command is used to bring existing infrastructure resources, created outside of Terraform, into the Terraform state file, enabling their management through Terraform going forward.

To update a statefile with an existing org tracing resource, execute the following import command:

```terraform
terraform import cpln_org_tracing.RESOURCE_NAME ORG_NAME
```

-> 1. Substitute RESOURCE_NAME with the same string that is defined in the HCL file.<br/>2. Substitute ORG_NAME with the target org.
