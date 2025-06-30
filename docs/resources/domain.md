---
page_title: "cpln_domain Resource - terraform-provider-cpln"
subcategory: "Domain"
description: |-
---

# cpln_domain (Resource)

Manages an org's custom [Domain](https://docs.controlplane.com/reference/domain).

The required DNS entries must exist before using Terraform to manage a `Domain`.

Refer to the [Configure a Domain](https://docs.controlplane.com/guides/configure-domain#dns-entries)
page for additional details.

During the creation of a domain, Control Plane will verify that the DNS entries exists. If they do
not exist, the Terraform script will fail.

The APEX domain is required to be added to one of the orgs. Any subdomain within that org will not need the `TXT` validation record. Any subdomain added to another org will require the `TXT` validation record be added.

## Declaration

### Required

- **name** (String) Domain name. (e.g., `example.com` / `test.example.com`). Control Plane will validate the existence of the domain with DNS. Create and Update will fail if the required DNS entries cannot be validated.

~> **Note** For a subdomain, include a `depends_on` property that points to the APEX domain declaration if the APEX was created in the same org.

- **spec** (Block List, Max: 1) ([see below](#nestedblock--spec))

~> **Note** If no spec properties are configured, an empty spec declaration (e.g., **spec { }**) is required to allow the default properties to exist in the state file.

### Optional

- **description** (String) Description of the domain name.
- **tags** (Map of String) Key-value map of resource tags.

<a id="nestedblock--spec"></a>

### `spec`

Required:

- **ports** (Block List) ([see below](#nestedblock--spec-ports))

~> **Note** If no ports are configured, an empty ports declaration (e.g., **ports { }**) is required to allow the default properties to exist in the state file.

Optional:

- **dns_mode** (String) In `cname` dnsMode, Control Plane will configure workloads to accept traffic for the domain but will not manage DNS records for the domain. End users must configure CNAME records in their own DNS pointed to the canonical workload endpoint. Currently `cname` dnsMode requires that a TLS server certificate be configured when subdomain based routing is used. In `ns` dnsMode, Control Plane will manage the subdomains and create all necessary DNS records. End users configure NS records to forward DNS requests to the Control Plane managed DNS servers. Valid values: `cname`, `ns`. Default: `cname`.
- **gvc_link** (String) This value is set to a target GVC (using a full link) for use by subdomain based routing. Each workload in the GVC will receive a subdomain in the form ${workload.name}.${domain.name}. **Do not include if path based routing is used.**
- **accept_all_hosts** (Boolean) Allows domain to accept wildcards. The associated GVC must have dedicated load balancing enabled.

<a id="nestedblock--spec-ports"></a>

### `spec.ports`

Required:

- **tls** (Block List, Max: 1) ([see below](#nestedblock--spec--ports--tls))

~> **Note** If no tls properties are configured, an empty tls declaration (e.g., **tls { }**) is required to allow for the default properties to exist in the state file.

Optional:

- **number** (Number) Port to expose externally. Values: `80`, `443`. Default: `443`.
- **protocol** (String) Allowed protocol. Valid values: `http`, `http2`, `tcp`. Default: `http2`.
- **cors** (Block List, Max: 1) ([see below](#nestedblock--spec--ports--cors))

<a id="nestedblock--spec--ports--cors"></a>

### `spec.ports.cors`

Optional:

- **allow_origins** (Block List) Determines which origins are allowed to access a particular resource on a server from a web browser. ([see below](#nestedblock--spec--ports--cors--allow_origins))
- **allow_methods** (List of Strings) Specifies the HTTP methods (such as `GET`, `POST`, `PUT`, `DELETE`, etc.) that are allowed for a cross-origin request to a specific resource.
- **allow_headers** (List of Strings) Specifies the custom HTTP headers that are allowed in a cross-origin request to a specific resource.
- **expose_headers** (List of Strings) The HTTP headers that a server allows to be exposed to the client in response to a cross-origin request. These headers provide additional information about the server's capabilities or requirements, aiding in proper handling of the request by the client's browser or application.
- **max_age** (String) Maximum amount of time that a preflight request result can be cached by the client browser. Input is expected as a duration string (i.e, 24h, 20m, etc.).
- **allow_credentials** (Boolean) Determines whether the client-side code (typically running in a web browser) is allowed to include credentials (such as cookies, HTTP authentication, or client-side SSL certificates) in cross-origin requests.

<a id="nestedblock--spec--ports--cors--allow_origins"></a>

### `spec.ports.cors.allow_origins`

Optional:

- **exact** (String) Value of allowed origin.
- **regex** (String)

<a id="nestedblock--spec--ports--tls"></a>

### `spec.ports.tls`

- **min_protocol_version** (String) Minimum TLS version to accept. Minimum is `1.0`. Default: `1.2`.
- **cipher_suites** (List of Strings) Allowed cipher suites. Refer to the [Domain Reference](https://docs.controlplane.com/reference/domain#cipher-suites) for details.
- **client_certificate** (Block List, Max: 1) The certificate authority PEM, stored as a TLS Secret, used to verify the authority of the client certificate. The only verification performed checks that the CN of the PEM matches the Domain (i.e., CN=*.DOMAIN). ([see below](#nestedblock--spec--ports--tls--certificate))
- **server_certificate** (Block List, Max: 1) Custom Server Certificate. ([see below](#nestedblock--spec--ports--tls--certificate))

~> **Note** If a custom server certificate is configured on a domain, it is the responsibility of the user to ensure that the certificate is valid and not expired.

<a id="nestedblock--spec--ports--tls--certificate"></a>

### `spec.ports.tls.certificate`

Optional:

- **secret_link** (String) Full link to a TLS secret.

## Outputs

The following attributes are exported:

- **cpln_id** (String) The ID, in GUID format, of the Domain.
- **self_link** (String) Full link to this resource. Can be referenced by other resources.
- **status** (Block List, Max: 1) ([see below](#nestedblock--status))

<a id="nestedblock--status"></a>

### `status`

Status of the domain.

Read-Only:

- **status** (String) Status of Domain. Possible values: `initializing`, `ready`, `pendingDnsConfig`, `pendingCertificate`, `usedByGvc`.
- **warning** (String) Warning message.
- **fingerprint** (String)
- **endpoints** (Block List) ([see below](#nestedblock--status--endpoints))
- **locations** (Block List) ([see below](#nestedblock--status--locations))
- **dns_config** (Block List) ([see below](#nestedblock--status--dns_config))

<a id="nestedblock--status--endpoints"></a>

### `status.endpoints`

List of configured domain endpoints.

- **url** (String) URL of endpoint.
- **workload_link** (String) Full link to associated workload.

<a id="nestedblock--status--locations"></a>

### `status.locations`

List of locations where domain is deployed.

- **name** (String) Location name.
- **certificate_status** (String) Status of certificate. Valud values: `initializing`, `ready`, `pendingDnsConfig`, `pendingCertificate`, `ignored`.

<a id="nestedblock--status--dns_config"></a>

### `status.dns_config`

List of required DNS record entries.

- **type** (String) The DNS record type specifies the type of data the DNS record contains. Valid values: `CNAME`, `NS`, `TXT`.
- **ttl** (Number) Time to live (TTL) is a value that signifies how long (in seconds) a DNS record should be cached by a resolver or a browser before a new request should be sent to refresh the data. Lower TTL values mean records are updated more frequently, which is beneficial for dynamic DNS configurations or during DNS migrations. Higher TTL values reduce the load on DNS servers and improve the speed of name resolution for end users by relying on cached data.
- **host** (String) The host in DNS terminology refers to the domain or subdomain that the DNS record is associated with. It's essentially the name that is being queried or managed. For example, in a DNS record for `www.example.com`, `www` is a host in the domain `example.com`.
- **value** (String) The value of a DNS record contains the data the record is meant to convey, based on the type of the record.

## Example Usage

```terraform
resource "cpln_domain" "domain_apex" {
		name        = "example.com"
		description = "APEX domain example"

		tags = {
		  terraform_generated = "true"
		}

		spec {
			ports {
				tls { }
			 }
		}
}

resource "cpln_domain" "example_ns_subdomain" {

  depends_on  = [cpln_domain.domain_apex]

  name        = "app.example.com"
  description = "Custom domain that can be set on a GVC and used by associated workloads"

  tags = {
    terraform_generated = "true"
    example             = "true"
  }

  spec {
    dns_mode         = "ns"
    gvc_link         = "/org/myorg/gvc/mygvc"

    ports {
      number   = 443
      protocol = "http2"

      cors {

        allow_origins {
          exact = "example.com"
        }

         allow_origins {
          exact = "*"
        }

        allow_methods     = ["allow_method_1", "allow_method_2", "allow_method_3"]
        allow_headers     = ["allow_header_1", "allow_header_2", "allow_header_3"]
        max_age           = "24h"
        allow_credentials = "true"
      }

      tls {
        min_protocol_version = "TLSV1_2"
        cipher_suites = [
          "ECDHE-ECDSA-AES256-GCM-SHA384",
          "ECDHE-ECDSA-CHACHA20-POLY1305",
          "ECDHE-ECDSA-AES128-GCM-SHA256",
          "ECDHE-RSA-AES256-GCM-SHA384",
          "ECDHE-RSA-CHACHA20-POLY1305",
          "ECDHE-RSA-AES128-GCM-SHA256",
          "AES256-GCM-SHA384",
          "AES128-GCM-SHA256",
        ]

        server_certificate {
          secret_link = "LINK_TO_TLS_CERTIFICATE"
        }
      }
    }
  }
}
```

## Import Syntax

The `terraform import` command is used to bring existing infrastructure resources, created outside of Terraform, into the Terraform state file, enabling their management through Terraform going forward.

To update a statefile with an existing domain resource, execute the following import command:

```terraform
terraform import cpln_domain.RESOURCE_NAME DOMAIN_NAME
```

-> 1. Substitute RESOURCE_NAME with the same string that is defined in the HCL file.<br/>2. Substitute DOMAIN_NAME with the corresponding domain defined in the resource.
