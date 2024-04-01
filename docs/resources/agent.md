---
page_title: "cpln_agent Resource - terraform-provider-cpln"
subcategory: "Agent"
description: |-
---

# cpln_agent (Resource)

Supports the creation of an [Agent](https://docs.controlplane.com/reference/agent). Multiple agents can be created for an org.

Agents allow secure communication between workloads running on the Control Plane platform and TCP endpoints inside private networks such as VPCs.

## Declaration

### Required

- **name** (String) Name of the Agent.

### Optional

- **description** (String) Description of the Agent.
- **tags** (Map of String) Key-value map of resource tags.

## Outputs

The following attributes are exported:

- **cpln_id** (String) The ID, in GUID format, of the Agent.
- **self_link** (String) Full link to this resource. Can be referenced by other resources.
- **user_data** (String, Sensitive) The JSON output needed when [creating an agent](https://docs.controlplane.com/guides/agent).

**Note:** The `user_data` output value is only generated when the resource is created. Because of its sensitive nature, the `user_data` value will not be displayed.

To use the `user_data` output:

1. After the initial apply, the `cpln_agent` output can either be directed to a file using the command `terraform output -json > ./cpln_agent.json`, or,
2. During the apply, used as a resource in a Terraform script to instantiate the agent at a cloud provider.

** Only the `user_data` value is required when configuring an agent, not the entire output. **

Refer to this [example](https://github.com/controlplane-com/examples/blob/main/terraform/poc/example-postgres/main.tf) in which
one of the steps creates an Agent at AWS using the `user_data` output.

## Example Usage

```terraform
resource "cpln_agent" "example" {

  name        = "agent-example"
  description = "Example Agent"

  tags = {
    terraform_generated = "true"
    example             = "true"
  }
}
```

## Import Syntax

The `terraform import` command is used to bring existing infrastructure resources, created outside of Terraform, into the Terraform state file, enabling their management through Terraform going forward.

To update a statefile with an existing agent resource, execute the following import command:

```terraform
terraform import cpln_agent.RESOURCE_NAME AGENT_NAME
```

-> 1. Substitute RESOURCE_NAME with the same string that is defined in the HCL file.<br/>2. Substitute AGENT_NAME with the corresponding agent defined in the resource.
