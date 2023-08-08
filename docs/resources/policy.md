---
page_title: "cpln_policy Resource - terraform-provider-cpln"
subcategory: "Policy"
description: |-
---

# cpln_policy (Resource)

Manages an org's [Policy](https://docs.controlplane.com/reference/policy).

## Declaration

### Required

- **name** (String) Name of the Policy.

### Optional

- **description** (String) Description of the Policy.
- **tags** (Map of String) Key-value map of resource tags.
- **target_kind** (String) The kind of resource to target (e.g., gvc, serviceaccount, etc.).
- **gvc** (String) The GVC for 'identity', 'workload' and 'volumeset' target kinds only.

- **target** (String) Set this value of this attribute to `all` if this policy should target all objects of the given target_kind. Otherwise, do not include the attribute.
- **target_links** (List of String) List of the targets this policy will be applied to. Not used if `target` is set to `all`.
- **target_query** (Block List, Max: 1) ([see below](#nestedblock--target_query)).

- **binding** (Block Set, Max: 50) ([see below](#nestedblock--binding)).

<a id="nestedblock--binding"></a>

### `binding`

Optional:

- **permissions** (Set of String) List of permissions to allow.
- **principal_links** (Set of String) List of the principals this binding will be applied to. Principal links format: `group/GROUP_NAME`, `user/USER_EMAIL`, `gvc/GVC_NAME/identity/IDENTITY_NAME`, `serviceaccount/SERVICE_ACCOUNT_NAME`.

<a id="nestedblock--target_query"></a>

### `target_query`

Optional:

- **fetch** (String) Type of fetch. Either: `links` or `items`. Default: `items`.
- **spec** (Block List, Max: 1) ([see below](#nestedblock--target_query--spec)).

<a id="nestedblock--target_query--spec"></a>

### `target_query.spec`

Optional:

- **match** (String) Type of match. Available values: `all`, `any`, `none`. Default: `all`.
- **terms** (Block List) ([see below](#nestedblock--target_query--spec--terms)).

<a id="nestedblock--target_query--spec--terms"></a>

### `target_query.spec.terms`

Terms can only contain one of the following attributes: `property`, `tag`.

Optional:

- **op** (String) Type of query operation. Available values: `=`, `>`, `>=`, `<`, `<=`, `!=`, `exists`, `!exists`. Default: `=`.

- **property** (String) Property to use for query evaluation.
<!-- - **rel** (String) Rel to use use for query evaluation. -->
- **tag** (String) Tag key to use for query evaluation.
- **value** (String) Testing value for query evaluation.

## Outputs

The following attributes are exported:

- **cpln_id** (String) ID, in GUID format, of the Policy.
- **origin** (String) Origin of the Policy. Either `builtin` or `default`.
- **self_link** (String) Full link to this resource. Can be referenced by other resources.

## Example Usage

- Org Resources

```terraform
resource "cpln_service_account" "example" {

  name        = "service-account-example"
  description = "Example Service Account"

  tags = {
    terraform_generated = "true"
    example             = "true"
  }
}

resource "cpln_policy" "example" {

  name        = "policy-example"
  description = "Example Policy"

  tags = {
    terraform_generated = "true"
    example             = "true"
  }

  target_kind = "serviceaccount"

  # Policy can either target `all` or specific target links

  // target = "all"
  target_links = [cpln_service_account.example.name]

  target_query {

    spec {
      # match is either "all", "any", or "none"
      match = "all"

      terms {
        op    = "="
        tag   = "firebase/sign_in_provider"
        value = "microsoft.com"
      }
    }
  }

  binding {

    # Available permissions are based on the target kind
    permissions = ["manage", "edit"]

    # Principal links format: `group/GROUP_NAME`, `user/USER_EMAIL`, `gvc/GVC_NAME/identity/IDENTITY_NAME`, `serviceaccount/SERVICE_ACCOUNT_NAME`
    principal_links = ["user/email@example.com", "group/viewers"]
  }
}
```

- GVC Resources (i.e. identity, workload and volumeset)

```terraform
resource "cpln_gvc" "example" {

  name        = "gvc-example"
  description = "Example GVC"

  locations = ["aws-eu-central-1"]

  tags = {
    terraform_generated = "true"
  }
}

resource "cpln_identity" "example" {

  gvc = cpln_gvc.exmaple.name

  name        = "identity-example"
  description = "Example Identity"

  tags = {
    terraform_generated = "true"
  }
}

  resource "cpln_policy" "example" {

  name = "policy-example"
  description = "Example Policy for GVC resources"

  tags = {
    terraform_generated = "true"
  }

  target_kind = "identity"

  # gvc required for 'identity', 'workload' and 'volumeset' target kinds
  gvc = cpln_gvc.terraform_gvc.name

  target_links = [cpln_identity.example.name]


  binding {
    permissions = ["manage", "edit"]
    principal_links = ["user/support@controlplane.com", "group/viewers", "serviceaccount/service-account-${var.random-name}","gvc/${cpln_gvc.terraform_gvc.name}/identity/${cpln_identity.terraform_identity.name}"]
  }
}
```

## Import Syntax

The `terraform import` command is used to bring existing infrastructure resources, created outside of Terraform, into the Terraform state file, enabling their management through Terraform going forward.

To update a statefile with an existing policy resource, execute the following import command:

```terraform
terraform import cpln_policy.RESOURCE_NAME POLICY_NAME
```

-> 1. Substitute RESOURCE_NAME with the same string that is defined in the HCL file.<br/>2. Substitute POLICY_NAME with the corresponding policy defined in the resource.
