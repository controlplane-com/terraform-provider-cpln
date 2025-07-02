---
page_title: "cpln_cloud_account Resource - terraform-provider-cpln"
subcategory: "Cloud Account"
description: |-
---

# cpln_cloud_account (Resource)

[Cloud Accounts](https://docs.controlplane.com/reference/cloudaccount) enable Control Plane identities (workload identities) to define least-privilege access rules so that workloads can seamlessly consume resources from one or multiple clouds. The `cpln_cloud_account` resource defines a cloud account of either AWS, Azure, GCP or NGS.

The configuration at the associated cloud provider must exist before using Terraform to manage a `Cloud Account`.

Refer to the [Cloud Account Reference Page](https://docs.controlplane.com/reference/cloudaccount)
for additional details.

## Declaration

### Required

- **name** (String) Name of the Cloud Account.

### Optional

- **description** (String) Description of the Cloud Account.
- **tags** (Map of String) Key-value map of resource tags.

~> **Note** Only one of the cloud providers listed below can be included in a resource. Create resources for each additional cloud provider.

- **aws** (Block List, Max: 1) ([see below](#nestedblock--aws)).
- **azure** (Block List, Max: 1) ([see below](#nestedblock--azure)).
- **gcp** (Block List, Max: 1) ([see below](#nestedblock--gcp)).
- **ngs** (Block List, Max: 1) ([see below](#nestedblock--ngs)).

<a id="nestedblock--aws"></a>

### `aws`

Required:

- **role_arn** (String) Amazon Resource Name (ARN) Role.

<a id="nestedblock--azure"></a>

### `azure`

Required:

- **secret_link** (String) Full link to an Azure secret. (e.g., /org/ORG_NAME/secret/AZURE_SECRET).

<a id="nestedblock--gcp"></a>

### `gcp`

Required:

- **project_id** (String) GCP project ID. Obtained from the GCP cloud console.

<a id="nestedblock--ngs"></a>

### `ngs`

Required:

- **secret_link** (String) Full link to a NATS Account Secret secret. (e.g., /org/ORG_NAME/secret/NATS_ACCOUNT_SECRET).

## Outputs

The following attributes are exported:

- **cpln_id** (String) ID, in GUID format, of the Cloud Account.
- **self_link** (String) Full link to this resource. Can be referenced by other resources.
- **gcp_service_account_name** (String) GCP service account name used during the configuration of the cloud account at GCP.
- **gcp_roles** (List of String) GCP roles used during the configuration of the cloud account at GCP.

## Example Usage

```terraform
# AWS Cloud Account
resource "cpln_cloud_account" "aws" {
  name        = "cloud-account-aws"
  description = "AWS cloud account"

  tags = {
    terraform_generated = "true"
    example             = "true"
  }

  aws {
    role_arn = "arn:aws:iam::1234:role/test_role"
  }
}

# Azure Cloud Account
resource "cpln_cloud_account" "azure" {
  name        = "cloud-account-azure"
  description = "Azure cloud account"

  tags = {
    terraform_generated = "true"
    example             = "true"
  }

  azure {
    // Use full link
    secret_link = "/org/ORG_NAME/secret/AZURE_SECRET"
  }
}

# GCP Cloud Account
resource "cpln_cloud_account" "gcp" {
  name        = "cloud-account-gcp"
  description = "GCP cloud account"

  tags = {
    terraform_generated = "true"
    example             = "true"
  }

  gcp {
    project_id = "cpln_gcp_project_1234"
  }
}

output "cloud_account_gcp_service_account_name" {
  value = cpln_cloud_account.gcp.gcp_service_account_name
}

output "cloud_account_gcp_roles" {
  value = cpln_cloud_account.gcp.gcp_roles
}


# NGS Cloud Account
resource "cpln_cloud_account" "ngs" {
  name        = "cloud-account-ngs"
  description = "NGS cloud account "

  tags = {
    terraform_generated = "true"
    example             = "true"
  }

  ngs {
    // Use full link
    secret_link = "/org/ORG_NAME/secret/NATS_ACCOUNT_SECRET"
  }
}
```

## Import Syntax

The `terraform import` command is used to bring existing infrastructure resources, created outside of Terraform, into the Terraform state file, enabling their management through Terraform going forward.

To update a statefile with an existing cloud account resource, execute the following import command:

```terraform
terraform import cpln_cloud_account.RESOURCE_NAME CLOUD_ACCOUNT_NAME
```

-> 1. Substitute RESOURCE_NAME with the same string that is defined in the HCL file.<br/>2. Substitute CLOUD_ACCOUNT_NAME with the corresponding cloud account defined in the resource.
