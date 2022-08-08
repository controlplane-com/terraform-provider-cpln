---
page_title: "cpln_org_Tracing Resource - terraform-provider-cpln"
subcategory: "Org Tracing"
description: |-
  
---

# cpln_org_tracing (Resource)

Manages an Org's tracing configuration.


## Declaration

### Required 

At least one of the following tracing blocks are required:
  
- **lightstep_tracing** (Block List, Max: 1) ([see below](#nestedblock--lightstep_tracing)).


<a id="nestedblock--lightstep_tracing"></a>
 ### `lightstep_tracing`


Required:

- **sampling** (Int) Sampling percentage.
- **endpoint** (String) Tracing Endpoint Workload. Either the canonical endpoint or the internale endpoint. 

Optional:

- **credentials** (String) Full link to referenced Opaque Secret. 

~> **Note** The workload that the endpoint is pointing to must have the tag `cpln/tracingDisabled` set to  `true`.


## Example Usage

### Lightstep

```terraform

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

	resource "cpln_org_tracing" "new" {

		lightstep_tracing {

			sampling = 50
			endpoint = "test.cpln.local:8080"

			// Opaque Secret Only
			credentials = cpln_secret.opaque.self_link
		}	
	}

```

