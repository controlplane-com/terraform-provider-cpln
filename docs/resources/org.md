---
page_title: "cpln_org Resource - terraform-provider-cpln"
subcategory: "Org"
description: |-
---

# cpln_org (Resource)

Manage an [organization](https://docs.controlplane.com/reference/org).

## Declaration

### Optional

- **account_id** (String) Only effective on creation, the account id that will be used to create the organization.
- **org_invitees** (List of String) Only effective on creation, the list of emails that will receive an invitation to the organization as superusers.
- **session_timeout_seconds** (Int) This timeout setting (in seconds) specifies when the console UI will automatically sign out. Default: 900. (15 minutes)
- **auth_config** (Block List, Max: 1) ([see below](#nestedblock--auth_config)).
- **observability** (Block List, Max: 1) ([see below](#nestedblock--observability)).

~> **Note** In order to create an organization, you will need a login token. A recommended practice would be to not set the `token` argument, have the [CLI](https://docs.controlplane.com/reference/cli) installed and have a profile logged into your account. This way, the provider will use the `CPLN_PROFILE` environment variable set by the [CLI](https://docs.controlplane.com/reference/cli). You can read more about managing [CLI](https://docs.controlplane.com/reference/cli) profiles [here](https://docs.controlplane.com/guides/manage-profile#prerequisites).

<a id="nestedblock--auth_config"></a>

### `auth_config`

Required:

- **domain_auto_members** (List of String) // TODO: Add description
- **saml_only** (String) // TODO: Add description

<a id="nestedblock--observability"></a>

### `observability`

The retention period for logs, metrics and traces defaults to 30 days and can be adjusted for each independently.

Charges apply for storage beyond the 30 day default.

Required:

- **logs_retention_days** (Int) // TODO: Add description
- **metrics_retention_days** (Int) // TODO: Add description
- **traces_retention_days** (Int) // TODO: Add description

## Outputs

The following attributes are exported:

- **cpln_id** (String) ID, in GUID format, of the organization.
- **self_link** (String) Full link to this resource. Can be referenced by other resources.
- **status** (List of Object) ([see below](#nestedblock--status)).

<a id="nestedblock--status"></a>

### `status`

Status of the organization.

Read-Only:

- **account_link** (String) The link of the account the organization belongs to.
- **active** (Boolean) Indicates whether the org is active or not.

## Example Usage

```terraform

resource "cpln_org" "example" {
    name       = "new-org"
    account_id = "a1b23456-7cd8-901e-fgh2-3i456j7k89lm"
    invitees   = ["example-1@mail.com", "example-2@mail.com"]

    session_timeout_seconds = 1200

    auth_config {
        domain_auto_members = ["example.com"]
        saml_only           = false
    }

    observability {
        logs_retention_days    = 30
        metrics_retention_days = 40
        traces_retention_days  = 50
    }
}
```

## Import Syntax

The `terraform import` command is used to bring existing infrastructure resources, created outside of Terraform, into the Terraform state file, enabling their management through Terraform going forward.

To update a statefile with an existing Org resource, execute the following import command:

```terraform
terraform import cpln_org.RESOURCE_NAME ORG_NAME
```

-> 1. Substitute RESOURCE_NAME with the same string that is defined in the HCL file.<br/>2. Substitute ORG_NAME with the corresponding Org defined in the resource.
