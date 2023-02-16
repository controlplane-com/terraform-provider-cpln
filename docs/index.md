---
page_title: "Provider: cpln"
subcategory: ""
description: |-
  
---

# Control Plane Terraform Provider

The Control Plane Terraform Provider Plugin enables the scaffolding of any Control Plane object as code using [HCL](https://www.terraform.io/docs/language/syntax/configuration.html). It enables infrastructure as code with all the added benefit of the global virtual cloud (GVC). You can build your VPCs, subnets, databases, queues, caches, etc. and overlay them with a multi-cloud/multi-region universal compute workloads that span regions and clouds. Nearly everything you can do using the Control Plane CLI, UI or API is available using Terraform.

Each header below (i.e., `cpln_agent`) corresponds to a resource within the Control Plane Terraform provider.

## Authentication

Authenticate using one of the following methods:

1. CLI
  - [Install the CLI](https://docs.controlplane.com/reference/cli#installation) and execute the command `cpln login`. After a successful login, the Terraform provider will use the `default` profile to authenticate. To use a different profile, set the `profile` variable when initializing the provider or set the `CPLN_PROFILE` environment variable.

2. Token
  - The `token` variable can be set when initializing the provider or by setting the `CPLN_TOKEN` environment variable.
  - The value of `token` can be either:
      - The output of running the command `cpln profile token PROFILE_NAME`, or
      - In the case of a [Service Account](https://docs.controlplane.com/reference/serviceaccount), the value of one of it's [keys](https://docs.controlplane.com/reference/serviceaccount#keys)

~> **Note** To perform automated tasks using Terraform, the preferred method is to use a `Service Account` and one of it's `keys` as the `token` value. 


## Provider Declaration

### Required

- **org** (String) The Control Plane org that this provider will perform actions against. Can be specified with the `CPLN_ORG` environment variable.
### Optional

- **endpoint** (String) The Control Plane Data Service API endpoint. Default is: "https://api.cpln.io". Can be specified with the `CPLN_ENDPOINT` environment variable.
- **profile** (String) The user/service account profile that this provider will use to authenticate to the data service. Can be specified with the `CPLN_PROFILE` environment variable.
- **token** (String) The generated token that can be used to authenticate to the data service API. Can be specified with the `CPLN_TOKEN` environment variable.

~> **Note** If the `token` value is empty, the Control Plane CLI (cpln) must be installed and the command `cpln login` must be used to authenticate.

## Example Usage

```terraform
terraform {
  required_providers {
    cpln = {
      source = "controlplane-com/cpln"
      version = "1.1.0"
    }
  }
}

provider "cpln" {

  # Required
  # Can use CPLN_ORG Environment Variable
  org = var.org

  # Optional
  # Default Value: https://api.cpln.io
  # Can use CPLN_ENDPOINT Environment Variable
  endpoint = var.endpoint

  # Optional
  # Can use CPLN_PROFILE Environment Variable  
  profile = var.profile

  # Optional
  # Can use CPLN_TOKEN Environment Variable 
  token = var.token
}
```

