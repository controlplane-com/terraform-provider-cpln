---
page_title: "cpln_identity Resource - terraform-provider-cpln"
subcategory: "Identity"
description: |-
---

# cpln_identity (Resource)

Manages a GVC's [Identities](https://docs.controlplane.com/reference/identity).

## Declaration

### Required

- **name** (String) Name of the Identity.
- **gvc** (String) Name of the GVC.

### Optional

- **description** (String) Description of the Identity.
- **tags** (Map of String) Key-value map of resource tags.
- **aws_access_policy** (Block List, Max: 1) ([see below](#nestedblock--aws_access_policy)).
- **azure_access_policy** (Block List, Max: 1) ([see below](#nestedblock--azure_access_policy)).
- **gcp_access_policy** (Block List, Max: 1) ([see below](#nestedblock--gcp_access_policy)).
- **ngs_access_policy** (Block List, Max: 1) ([see below](#nestedblock--ngs_access_policy)).

- **network_resource** (Block List) ([see below](#nestedblock--network_resource)).
- **native_network_resource** (Block List) ([see below](#nestedblock--native_network_resource)

<a id="nestedblock--aws_access_policy"></a>

### `aws_access_policy`

Required:

- **cloud_account_link** (String) Full link to referenced cloud account.

Optional:

~> **Note** AWS Identity can either contain an existing `role_name` or multiple `policy_refs`.

- **policy_refs** (List of String) List of policies.
- **role_name** (String) Role name.

<a id="nestedblock--azure_access_policy"></a>

### `azure_access_policy`

Required:

- **cloud_account_link** (String) Full link to referenced cloud account.

Optional:

- **role_assignment** (Block List) ([see below](#nestedblock--azure_access_policy--role_assignment)).

<a id="nestedblock--azure_access_policy--role_assignment"></a>

### `azure_access_policy.role_assignment`

Optional:

- **roles** (List of String) List of assigned roles.
- **scope** (String) Scope of roles.

<a id="nestedblock--gcp_access_policy"></a>

### `gcp_access_policy`

~> **Note** The GCP access policy can either contain an existing service_account or multiple bindings.

Required:

- **cloud_account_link** (String) Full link to referenced cloud account.

Optional:

- **scopes** (String) Comma delimited list of GCP scope URLs.
- **service_account** (String) Name of existing GCP service account.
- **binding** (Block List) ([see below](#nestedblock--gcp_access_policy--binding)).

<a id="nestedblock--gcp_access_policy--binding"></a>

### `gcp_access_policy.binding`

Optional:

- **resource** (String) Name of resource for binding.
- **roles** (List of String) List of allowed roles.

<a id="nestedblock--ngs_access_policy"></a>

### `ngs_access_policy`

Required:

- **cloud_account_link** (String) Full link to referenced cloud account.

Optional:

- **pub** (Block List, Max: 1) Pub Permission. ([see below](#nestedblock--ngs_access_policy--perm)).
- **sub** (Block List, Max: 1) Sub Permission. ([see below](#nestedblock--ngs_access_policy--perm)).
- **resp** (Block List, Max: 1) Reponses. ([see below](#nestedblock--ngs_access_policy--resp)).
- **subs** (Number) Max number of subscriptions per connection. Default: -1
- **data** (Number) Max number of bytes a connection can send. Default: -1
- **payload** (Number) Max message payload. Default: -1

<a id="nestedblock--ngs_access_policy--perm"></a>

### `ngs_access_policy.pub` / `ngs_access_policy.sub`

Optional:

- **allow** (List of String) List of allow subjects.
- **deny** (List of String) List of deny subjects.

<a id="nestedblock--ngs_access_policy--resp"></a>

### `ngs_access_policy.resp`

Optional:

- **max** (Number) Number of responses allowed on the replyTo subject, -1 means no limit. Default: -1
- **ttl** (String) Deadline to send replies on the replyTo subject [#ms(millis) | #s(econds) | m(inutes) | h(ours)]. -1 means no restriction.

<a id="nestedblock--network_resource"></a>

### `network_resource`

A network resource can be configured with:

- A fully qualified domain name (FQDN) and ports.
- An FQDN, resolver IP, and ports.
- IP's and ports.

Required:

- **name** (String) Name of the Network Resource.
- **agent_link** (String) Full link to referenced Agent.

Optional:

- **fqdn** (String) Fully qualified domain name.
- **resolver_ip** (String) Resolver IP.
- **ips** (Set of String) List of IP addresses.

- **ports** (Set of Number) Ports to expose.

<a id="nestedblock--native_network_resource"></a>

### `native_network_resource`

~> **Note** The configuration of a native network resource requires the assistance of Control Plane support.

Required:

- **name** (String) Name of the Native Network Resource.
- **fqdn** (String) Fully qualified domain name.
- **ports** (Set of Number) Ports to expose. At least one port is required.

Optional:

Exactly one of:

- **aws_private_link** (Block List, Max: 1) ([see below](#nestedblock--native_network_resource--aws_private_link))
- **gcp_service_connect** (Block List, Max: 1) ([see below](#nestedblock--native_network_resource--gcp_service_connect))

<a id="nestedblock--native_network_resource--aws_private_link"></a>

### `aws_private_link`

Required:

- **endpoint_service_name** (String) Endpoint service name.

<a id="nestedblock--native_network_resource--gcp_service_connect"></a>

### `gcp_service_connect`

Required:

- **target_service** (String) Target service name.

## Outputs

The following attributes are exported:

- **cpln_id** (String) ID, in GUID format, of the Identity.
- **self_link** (String) Full link to this resource. Can be referenced by other resources.

## Example Usage

```terraform
resource "cpln_gvc" "example" {

  name        = "gvc-example"
  description = "Example GVC"

  locations = ["aws-us-west-2"]

  tags = {
    terraform_generated = "true"
    example             = "true"
  }
}

resource "cpln_agent" "example" {

  name        = "agent-example"
  description = "Example Agent"
}

resource "cpln_cloud_account" "example_aws" {

  name        = "aws-example"
  description = "Example AWS Cloud Account"

  tags = {
    terraform_generated = "true"
    example             = "true"
  }

  aws {
    role_arn = "arn:aws:iam::1234:role/example_role"
  }
}

resource "cpln_cloud_account" "example_azure" {

  name        = "azure-example"
  description = "Example Azure Cloud Account"

  tags = {
    terraform_generated = "true"
    example             = "true"
  }

  azure {
    // Use the full link to the Azure secret
    secret_link = "/org/ORG_NAME/secret/AZURE_SECRET"
  }
}

resource "cpln_cloud_account" "example-gcp" {

  name        = "gcp-example"
  description = "Example GCP Cloud Account"

  tags = {
    terraform_generated = "true"
    example             = "true"
  }

  gcp {
    project_id = "cpln_gcp_project_1234"
  }
}

resource "cpln_cloud_account" "example-ngs" {

  name        = "ngs-example"
  description = "Example NGS Cloud Account"

  tags = {
    terraform_generated = "true"
    acceptance_test     = "true"
  }

  ngs {
    // Use the full link for now
    // secret_link = "//secret/tf_secret_nats_account"
    secret_link = "/org/ORG_NAME/secret/NATS_ACCOUNT_SECRET"
  }
}

resource "cpln_identity" "example" {

  gvc = cpln_gvc.example.name

  name        = "identity-example"
  description = "Example Identity"

  tags = {
    terraform_generated = "true"
    example             = "true"
  }

  # Network Resource with FQDN
  network_resource {
    name       = "test-network-resource-fqdn"
    agent_link = cpln_agent.example.self_link
    fqdn       = "domain.example.com"
    ports      = [1234, 5432]
  }

  # Network Resource with FQDN and Resolver IP
  network_resource {
    name        = "test-network-resource-fqdn-rip"
    agent_link  = cpln_agent.example.self_link
    fqdn        = "domain2.example.com"
    resolver_ip = "192.168.1.1"
    ports       = [12345, 54321]
  }

  # Network Resource with IP
  network_resource {
    name       = "test-network-resource-ip"
    agent_link = cpln_agent.example.self_link
    ips        = ["192.168.1.1", "192.168.1.250"]
    ports      = [3099, 7890]
  }

  # Native Network Resource with AWS Private Link
  native_network_resource {
    name  = "test-native-network-resource-aws"
    fqdn  = "aws.test.com"
    ports = [12345, 54321]
    aws_private_link {
      endpoint_service_name = "com.amazonaws.vpce.us-west-2.vpce-svc-01af6c4c9260ac550"
    }
  }
  # Native Network Resource with GCP Service Connect
  native_network_resource {
    name  = "test-native-network-resource-gcp"
    fqdn  = "gcp.test.com"
    ports = [12345, 54321]
    gcp_service_connect {
      target_service = "projects/example-project/regions/example-region/serviceAttachments/example-service-attachments"
    }
  }

  aws_access_policy {

    cloud_account_link = cpln_cloud_account.example_aws.self_link

    # The AWS access policy can either contain an existing role_name or multiple policy_refs

    // role_name = "rds-monitoring-role"

    policy_refs = ["aws::/job-function/SupportUser", "aws::AWSSupportAccess"]
  }

  azure_access_policy {

    cloud_account_link = cpln_cloud_account.example_azure.self_link

    role_assignment {
      scope = "/subscriptions/d0d1e522-0825-415a-8b07-f7759b5c8a7e/resourceGroups/CP-Test-Resource-Group"
      roles = ["AcrPull", "AcrPush"]
    }

    role_assignment {
      scope = "/subscriptions/d0d1e522-0825-415a-8b07-f7759b5c8a7e/resourceGroups/CP-Test-Resource-Group/providers/Microsoft.Storage/storageAccounts/cplntest"
      roles = ["Support Request Contributor"]
    }
  }

  gcp_access_policy {

    cloud_account_link = cpln_cloud_account.example_gcp.self_link
    scopes             = "https://www.googleapis.com/auth/cloud-platform"

    # The GCP access policy can either contain an existing service_account or multiple bindings

    // service_account = "cpln-tf@cpln-test.iam.gserviceaccount.com"

    binding {
      resource = "//cloudresourcemanager.googleapis.com/projects/cpln-test"
      roles    = ["roles/appengine.appViewer", "roles/actions.Viewer"]
    }

    binding {
      resource = "//iam.googleapis.com/projects/cpln-test/serviceAccounts/cpln-tf@cpln-test.iam.gserviceaccount.com"
      roles    = ["roles/editor", "roles/iam.serviceAccountUser"]
    }
  }

  ngs_access_policy {

    cloud_account_link = cpln_cloud_account.example_ngs.self_link

    pub {
      allow = ["pa1", "pa2"]
      deny  = ["pd1", "pd2"]
    }

    sub {
      allow = ["sa1", "sa2"]
      deny  = ["sd1", "sd2"]
    }

    resp {
      max = 1
      ttl = "5m"
    }

    subs    = 1
    data    = 2
    payload = 3
  }
}
```

## Import Syntax

The `terraform import` command is used to bring existing infrastructure resources, created outside of Terraform, into the Terraform state file, enabling their management through Terraform going forward.

To update a statefile with an existing identity resource, execute the following import command:

```terraform
terraform import cpln_identity.RESOURCE_NAME GVC_NAME:IDENTITY_NAME
```

-> 1. Substitute RESOURCE_NAME with the same string that is defined in the HCL file.<br/>2. Substitute GVC_NAME and IDENTITY_NAME with the corresponding GVC and identity name defined in the resource.
