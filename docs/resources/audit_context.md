---
page_title: "cpln_audit_context Resource - terraform-provider-cpln"
subcategory: "Audit Context"
description: |-

---

# cpln_audit_context (Resource)

Manages an org's [Audit Context](https://docs.controlplane.com/reference/auditctx).

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