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

Optional:

- **cloud_account_link** (String) Full link to referenced cloud account. 
- **role_assignment** (Block List) ([see below](#nestedblock--azure_access_policy--role_assignment)).

<a id="nestedblock--azure_access_policy--role_assignment"></a>
 ### `azure_access_policy.role_assignment`

Optional:

- **roles** (List of String) List of assigned roles.
- **scope** (String) Scope of roles.



<a id="nestedblock--gcp_access_policy"></a>
 ### `gcp_access_policy`

~> **Note** The GCP access policy can either contain an existing service_account or multiple bindings.

- **cloud_account_link** (String) Full link to referenced Cloud Account. 
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

- **cloud_account_link** (String) Full link to referenced Cloud Account.

Optional:

- **pub** (Block List, Max: 1) ([see below](#nestedblock--ngs_access_policy--premissions_set)).
- **sub** (Block List, Max: 1) ([see below](#nestedblock--ngs_access_policy--premissions_set)).
- **resp** (Block List, Max: 1) ([see below](#nestedblock--ngs_access_policy--resp)).
- **subs** (Number) Max number of subscriptions per connection.
- **data** (Number) Max number of bytes a connection can send.
- **payload** (Number) Max message payload.

<a id="nestedblock--ngs_access_policy--premissions_set"></a>
### `ngs_access_policy.premissions_set`

Optional:

- **allow** (Set of String) // TODO: missing description.
- **deny** (Set of String) // TODO: missing description.

<a id="nestedblock--ngs_access_policy--resp"></a>
 ### `ngs_access_policy.resp`

Optional:

- **max** (Number) Number of responses allowed on the replyTo subject, -1 means no limit.
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

Required:
-- **name** (String) Name of the Native Network Resource.
-- **fqln** (String)  // TODO: Missing description
-- **ports** (Set of Number) Ports to expose. At least one is required.

Optional:

Exactly one of:
-- **aws_private_link** (Block List) ([see below](#nestedblock--native_network_resource--aws_private_link))
-- **gcp_service_connect**  (Block List) ([see below](#nestedblock--native_network_resource--gcp_service_connect))

<a id="nestedblock--native_network_resource--aws_private_link"></a>
### `aws_private_link`

Required:
-- **endpoint_service_name** (String)  // TODO: Missing description

<a id="nestedblock--native_network_resource--gcp_service_connect"></a>
### `gcp_service_connect`

Required:
-- **target_service** (String) // TODO: Missing description

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
    agent_link = cpln_agent.test_agent.self_link
    ips        = ["192.168.1.1", "192.168.1.250"]
    ports      = [3099, 7890]
  }

  # Native Network Resource with AWS Private Link
  native_network_resource {
			name = "test-native-network-resource"
			fqdn = "test.com"
			ports = [12345, 54321]

			aws_private_link {
				endpoint_service_name = "" // TODO: Add endpoint service name as an example
			}
		}

  # Native Network Resource with GCP Service Connect
  /*native_network_resource {
    name = "test-native-network-resource"
    fqdn = "test.com"
    ports = [12345, 54321]

    gcp_service_connect {
      target_service = "" // TODO: Add target service as an example
    }
  }*/

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
    scopes             = ["https://www.googleapis.com/auth/cloud-platform"]

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
    cloud_account_link = "/org/{org_name}/cloudaccount/{nats_account_secret}"

    pub {
      allow = ["allow_1", "allow_2"]
      deny = ["deny_1", "deny_2"]
    }

    sub {
      allow = ["allow_3", "allow_4"]
      deny = ["deny_3", "deny_4"]
    }

    resp {
      max = 10
      ttl = "500ms"
    }

    subs = -1
    data = -1
    payload = -1
  }
}
```