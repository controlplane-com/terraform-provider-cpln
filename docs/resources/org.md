---
page_title: "cpln_org Resource - terraform-provider-cpln"
subcategory: "Org"
description: |-
---

# cpln_org (Resource)

Manage an [organization](https://docs.controlplane.com/reference/org) (org).

~> **Note** The target org name for this resource is referenced from the `provider` block.

~> **Note** Since Terraform creates/updates resources in parallel, org resources (such as GVC, Workload, etc) must have the `depends_on` property when using Terraform to create an org. This allows the org to exist before other resources are created or updated. See example at the bottom for proper usage.

## Declaration

### Required

- **observability** (Block List, Max: 1) ([see below](#nestedblock--observability)).

### Optional

- **account_id** (String) The associated account ID that will be used when creating the org. Only used on org creation. The account ID can be obtained from the `Org Management & Billing` page.
- **org_invitees** (List of String) When an org is created, the list of email addresses which will receive an invitation to join the org and be assigned to the `superusers` group. The user account used when creating the org will be included in this list.
- **session_timeout_seconds** (Int) The idle time (in seconds) in which the console UI will automatically sign-out the user. Default: 900 (15 minutes)
- **auth_config** (Block List, Max: 1) ([see below](#nestedblock--auth_config)).
- **security** (Block List, Max: 1) ([see below](#nestedblock--security)).


~> **Note** To create an org, the provider **must** [authenticate](https://registry.terraform.io/providers/controlplane-com/cpln/latest/docs#authentication) with the `CLI` or `refresh_token` using a user account that has the `org_creator` role for the associated account.

~> **Note** To update the `auth_config` property, the provider **must** [authenticate](https://registry.terraform.io/providers/controlplane-com/cpln/latest/docs#authentication) with the `CLI` or `refresh_token` using a user account that was authenticated using a SAML provider.

~> **Note** Executing `terraform destroy` on a `cpln_org` resource will not delete the org. The properties of this resources will be restored to their default values.

<a id="nestedblock--observability"></a>

### `observability`

The retention period (in days) for logs, metrics, and traces.

Charges apply for storage beyond the 30 day default.

Optional:

- **logs_retention_days** (Int) Log retention days. Default: 30
- **metrics_retention_days** (Int) Metrics retention days. Default: 30
- **traces_retention_days** (Int) Traces retention days. Default: 30

~> **Note** The `observability` block is required, but the sub-properties are optional and will use the default value if not provided.

<a id="nestedblock--auth_config"></a>

### `auth_config`

Required:

- **domain_auto_members** (List of String) List of domains which will auto-provision users when authenticating using SAML.
- **saml_only** (Boolean) Enforce SAML only authentication.

<a id="nestedblock--security"></a>

### `security`

Optional:

- **threat_detection** (Block List, Max: 1) ([see below](#nestedblock--security--threat_detection))

<a id="nestedblock--security--threat_detection"></a>

### `security.threat_detection`

Optional:

- **enabled** (Boolean)
- **minimum_severity** (String)
- **syslog** (Block List, Max: 1) ([see below](#nestedblock--security--threat_detection--syslog))

<a id="nestedblock--security--threat_detection--syslog"></a>

### `security.threat_detection.syslog`

Required:

- **port** (Int)

Optional:

- **transport** (String) Default: `tcp`.
- **host** (String)

## Outputs

The following attributes are exported:

- **cpln_id** (String) ID, in GUID format, of the org.
- **name** (String) The name of the org.
- **self_link** (String) Full link to this resource. Can be referenced by other resources.
- **status** (List of Object) ([see below](#nestedblock--status)).

<a id="nestedblock--status"></a>

### `status`

Status of the org.

Read-Only:

- **account_link** (String) The link of the account the org belongs to.
- **active** (Boolean) Indicates whether the org is active or not.

## Example Usage

```terraform

resource "cpln_org" "example" {

    account_id = "a1b23456-7cd8-901e-fgh2-3i456j7k89lm"
    invitees   = ["example-1@mail.com", "example-2@mail.com"]

    description = "Example Org"

    tags = {
        terraform_generated = "true"
        example             = "true"
    }

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

    security {
        threat_detection {
            enabled          = true
            minimum_severity = "warning"
            syslog {
                transport = "tcp"
                host 	  = "example.com"
                port  	  = 8080
            }
        }
    }
}

resource "cpln_gvc" "example" {

  depends_on = [cpln_org.example]

  name        = "gvc-example"
  description = "Example GVC"

  locations = ["aws-eu-central-1", "aws-us-west-2"]

}

```

## Import Syntax

The `terraform import` command is used to bring existing infrastructure resources, created outside of Terraform, into the Terraform state file, enabling their management through Terraform going forward.

To update a statefile with an existing org resource, execute the following import command:

```terraform
terraform import cpln_org.RESOURCE_NAME ORG_NAME
```

-> 1. Substitute RESOURCE_NAME with the same string that is defined in the HCL file.<br/>2. Substitute ORG_NAME with the corresponding org name defined in the resource.
