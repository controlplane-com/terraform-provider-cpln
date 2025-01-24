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

`1. CLI`
- [Install the CLI](https://docs.controlplane.com/reference/cli#installation) and execute the command `cpln login`. After a successful login, the Terraform provider will use the `default` profile to authenticate. To use a different profile, set the `profile` variable when initializing the provider or set the `CPLN_PROFILE` environment variable.

`2. Token`
- The `token` variable can be set when initializing the provider or by setting the `CPLN_TOKEN` environment variable.
- The value of `token` can be either:
  - The output of running the command `cpln profile token PROFILE_NAME`, or
  - In the case of a [Service Account](https://docs.controlplane.com/reference/serviceaccount), the value of one of it's [keys](https://docs.controlplane.com/reference/serviceaccount#keys)

`3. Refresh Token`
- The `refresh_token` variable is used when the provider is required to create an org or update the `auth_config` property using the `cpln_org` resource. The `refresh_token` variable can be set when initializing the provider or by setting the `CPLN_REFRESH_TOKEN` environment variable.
- When creating an org, the `refresh_token` **must** belong to a user that has the `org_creator` role for the associated account.
- When updating the org `auth_config` property, the `refresh_token` **must** belong to a user that was authenticated using SAML.
- The `refresh_token` can be obtained by following these steps:
  - Using the CLI, authenticate with a user account by executing `cpln login`.
  - Browser to the path `~/.config/cpln/profiles`. This path will contain JSON files corresponding to the name of the profile (i.e., `default.json`).
  - The contents of the JSON file will contain a key named `refreshToken`. Use the value of this key for the `refresh_token` variable.
  
~> **Note** To perform automated tasks using Terraform, the preferred method is to use a `Service Account` and one of it's `keys` as the `token` value.

## Provider Declaration

### Required

- **org** (String) The Control Plane org that this provider will perform actions against. Can be specified with the `CPLN_ORG` environment variable.

### Optional

- **endpoint** (String) The Control Plane Data Service API endpoint. Default is: `https://api.cpln.io`. Can be specified with the `CPLN_ENDPOINT` environment variable.
- **profile** (String) The user/service account profile that this provider will use to authenticate to the data service. Can be specified with the `CPLN_PROFILE` environment variable.
- **token** (String) A generated token that can be used to authenticate to the data service API. Can be specified with the `CPLN_TOKEN` environment variable.
- **refresh_token** (String) A generated token that can be used to authenticate to the data service API. Can be specified with the `CPLN_REFRESH_TOKEN` environment variable. Used when the provider is required to create an org or update the `auth_config` property. Refer to the section above on how to obtain the refresh token.

~> **Note** If the `token` or `refresh_token` value is empty, the Control Plane CLI (cpln) must be installed and the `cpln login` command must be used to authenticate.

## Example Usage

```terraform
terraform {
  required_providers {
    cpln = {
      source = "controlplane-com/cpln"
      version = "1.1.43"
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

  # Optional
  # Can use CPLN_REFRESH_TOKEN Environment Variable
  refresh_token = var.refresh_token
}
```
