---
page_title: "cpln_group Resource - terraform-provider-cpln"
subcategory: "Group"
description: |-
---

# cpln_group (Resource)

Manages an org's [Groups](https://docs.controlplane.com/reference/group).

Groups membership can contain [users](https://docs.controlplane.com/reference/user), [service accounts](https://docs.controlplane.com/reference/serviceaccount), and [custom queries](https://docs.controlplane.com/reference/group#query).

## Declaration

### Required

- **name** (String) Name of the Group.

### Optional

- **description** (String) Description of Group.
- **tags** (Map of String) Key-value map of resource tags.
- **service_accounts** (List of String) List of service accounts that exists within the configured org. Group membership will fail if the service account does not exits within the org.
- **user_ids_and_emails** (List of String) List of either the user ID or email address for a user that exists within the configured org. Group membership will fail if the user ID / email does not exist within the org.

- **member_query** (Block List, Max: 1) ([see below](#nestedblock--member_query)).
- **identity_matcher** (Block List, Max: 1) ([see below](#nestedblock--identity_matcher)).

<a id="nestedblock--member_query"></a>

### `member_query`

Optional:

- **fetch** (String) Type of fetch. Specify either: `links` or `items`. Default: `items`.
- **spec** (Block List, Max: 1) ([see below](#nestedblock--member_query--spec)).

<a id="nestedblock--member_query--spec"></a>

### `member_query.spec`

Optional:

- **match** (String) Type of match. Available values: `all`, `any`, `none`. Default: `all`.
- **terms** (Block List) ([see below](#nestedblock--member_query--spec--terms)).

<a id="nestedblock--member_query--spec--terms"></a>

### `member_query.spec.terms`

<!-- Terms can only contain one of the following attributes: `property`, `rel`, `tag`. -->

Terms can only contain one of the following attributes: `property`, `tag`.

Optional:

- **op** (String) Type of query operation. Available values: `=`, `>`, `>=`, `<`, `<=`, `!=`, `exists`, `!exists`. Default: `=`.

- **property** (String) Property to use for query evaluation.
<!-- - **rel** (String) Rel to use use for query evaluation. -->
- **tag** (String) Tag key to use for query evaluation.
- **value** (String) Testing value for query evaluation.

<a id="nestedblock--identity_matcher"></a>

### `identity_matcher`

Required:

- **expression** (String) Executes the expression against the users' claims to decide whether a user belongs to this group. This method is useful for managing the grouping of users logged in with SAML providers.

Optional:

- **language** (String) Language of the expression. Either `jmespath` or `javascript`. Default: `jmespath`.

## Outputs

The following attributes are exported:

- **origin** (String) Origin of the service account. Either `builtin` or `default`.
- **self_link** (String) Full link to this resource. Can be referenced by other resources.

## Example Usage

-> The `identity_matcher` expressions (evaluated when using SAML authentication) for the JMESPath and JavaScript example resources below are based on the following JSON that was provided to Control Plane from the domain IdP (identity provider).

**Provided JSON from SAML IdP**

```javascript
{
  "identities": {
    "saml.example.com": [
      "user@example.com"
    ],
    "email": [
      "user@example.com"
    ]
  },
  "sign_in_provider": "saml.example.com",
  "sign_in_attributes": {
    "orgPath": "/",
    "memberOf": "developers"
  }
```

**Example Terraform**

```terraform
resource "cpln_service_account" "example" {

  name        = "service-account-example"
  description = "Service Account to be used as a member of a group"

  tags = {
    terraform_generated = "true"
    example             = "true"
  }
}

resource "cpln_group" "example" {

  name        = "group-example"
  description = "group example"

  tags = {
    terraform_generated = "true"
    example             = "true"
  }

  user_ids_and_emails   = ["valid_user@example.com"]
  service_accounts      = [cpln_service_account.example.name]

  member_query {

    fetch = "items"

    spec {
      match = "all"

      terms {
        op    = "="
        tag   = "firebase/sign_in_provider"
        value = "microsoft.com"
      }
    }
  }
}

resource "cpln_group" "example_jmespath" {

  name        = "group-example"
  description = "group jmespath"

  tags = {
    terraform_generated = "true"
    example             = "true"
  }

  identity_matcher {
    expression = "contains(sign_in_attributes.memberOf, 'developers')"
    language = "jmespath"
  }
}

resource "cpln_group" "example_javascript" {

  name        = "group-example"
  description = "group javascript"

  tags = {
    terraform_generated = "true"
    example             = "true"
  }

  identity_matcher {
    expression = "if ($.sign_in_attributes) { $.sign_in_attributes.memberOf.includes('developers'); }"
    language = "javascript"
  }
}
```

## Import Syntax

The `terraform import` command is used to bring existing infrastructure resources, created outside of Terraform, into the Terraform state file, enabling their management through Terraform going forward.

To update a statefile with an existing group resource, execute the following import command:

```terraform
terraform import cpln_group.RESOURCE_NAME GROUP_NAME
```

-> 1. Substitute RESOURCE_NAME with the same string that is defined in the HCL file.<br/>2. Substitute GROUP_NAME with the corresponding group defined in the resource.
