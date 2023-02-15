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
- **domain** (String) Custom domain name used by associated workloads.
- **pull_secrets** (List of String) A list of [pull secret](https://docs.controlplane.com/reference/gvc#pull-secrets) names used to authenticate to any private image repository referenced by Workloads within the GVC.
- **tags** (Map of String) Key-value map of resource tags.
- **env** (Array of Name-Value Pair) Key-value array of resource env variables.
- **lightstep_tracing** (Block List, Max: 1) ([see below](#nestedblock--lightstep_tracing)).

<a id="nestedblock--lightstep_tracing"></a>

### `lightstep_tracing`

Required:

- **sampling** (Int) Sampling percentage.
- **endpoint** (String) Tracing Endpoint Workload. Either the canonical endpoint or the internal endpoint.

Optional:

- **credentials** (String) Full link to referenced Opaque Secret.

~> **Note** The workload that the endpoint is pointing to must have the tag `cpln/tracingDisabled` set to `true`.

## Outputs

The following attributes are exported:

- **self_link** (String) Full link to this resource. Can be referenced by other resources.

## Example Usage

```terraform
resource "cpln_secret" "docker" {
	name = "docker-secret"
	description = "docker secret"

	tags = {
		terraform_generated = "true"
		acceptance_test = "true"
		secret_type = "docker"
	}

	env = {
		env_var_key = "env_var_value"
		workload_can_inherit = "true"
	}

	docker = "{\"auths\":{\"your-registry-server\":{\"username\":\"your-name\",\"password\":\"your-pword\",\"email\":\"your-email\",\"auth\":\"<Secret>\"}}}"
}

resource "cpln_secret" "opaque" {

	name = "opaque-random-tbd"
	description = "description opaque-random-tbd"

	tags = {
		terraform_generated = "true"
		acceptance_test = "true"
		secret_type = "opaque"
	}

	opaque {
		payload = "opaque_secret_payload"
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

  lightstep_tracing {

		sampling = 50
		endpoint = "test.cpln.local:8080"

		// Opaque Secret Only
		credentials = cpln_secret.opaque.self_link
	}
}
```
