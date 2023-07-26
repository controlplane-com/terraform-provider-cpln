---
page_title: "cpln_audit_context Resource - terraform-provider-cpln"
subcategory: "Audit Context"
description: |-
---

# cpln_audit_context (Resource)

Manages an org's [Audit Context](https://docs.controlplane.com/reference/auditctx).

~> **Note** Audit Contexts are immutable and can not be deleted. When destroying using Terraform, it will only be removed from the state. Existing audit contexts must be imported.

## Declaration

### Required

- **name** (String) Name of the Audit Context.

### Optional

- **description** (String) Description of the Audit Context.
- **tags** (Map of String) Key-value map of resource tags.

## Outputs

The following attributes are exported:

- **self_link** (String) Full link to this resource. Can be referenced by other resources.

## Example Usage

```terraform
resource "cpln_audit_context" "example" {
    name = "audit-context-example"
    description = "audit context description"

    tags = {
        terraform_generated = "true"
        example = "true"
    }
}
```

## Import Syntax

The `terraform import` command is used to bring existing infrastructure resources, created outside of Terraform, into the Terraform state file, enabling their management through Terraform going forward.

To update a statefile with an existing audit context resource, execute the following import command:

```terraform
terraform import cpln_audit_context.RESOURCE_NAME AUDIT_CONTEXT_NAME
```

-> 1. Substitute RESOURCE_NAME with the same string that is defined in the HCL file.<br/>2. Substitute AUDIT_CONTEXT_NAME with the corresponding audit context defined in the resource.
