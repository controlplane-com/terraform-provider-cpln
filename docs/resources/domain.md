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

The APEX domain is required to be added to one of the orgs. Any subdomain within that org will not need the TXT records. Any subdomain added to another org will require the TXT records be added.

## Declaration

### Required

- **name** (String) Domain name. (e.g., example.com / test.example.com). Control Plane will validate the existence of the domain with DNS. Create and Update will fail if the required DNS entries cannot be validated.

~> **Note** For a subdomain, include a `depends_on` property that points to the APEX domain declaration if the APEX was created in the same org.

- **spec** (Block List, Max: 1) ([see below](#nestedblock--spec))

~> **Note** If no spec properties are configured, an empty spec declaration (e.g., **spec { }**) is required to allow the default properties to exist in the state file.

### Optional

- **description** (String) Description for the domain name.
- **tags** (Map of String) Key-value map of resource tags.

<a id="nestedblock--spec"></a>

### `spec`

Required:

- **ports** (Block List) ([see below](#nestedblock--spec-ports))

~> **Note** If no ports are configured, an empty ports declaration (e.g., **ports { }**) is required to allow the default properties to exist in the state file.

Optional:

- **dns_mode** (String) In 'cname' dnsMode, Control Plane will configure workloads to accept traffic for the domain but will not manage DNS records for the domain. End users must configure CNAME records in their own DNS pointed to the canonical workload endpoint. Currently 'cname' dnsMode requires that a TLS server certificate be configured when subdomain based routing is used. In 'ns' dnsMode, Control Plane will manage the subdomains and create all necessary DNS records. End users configure NS records to forward DNS requests to the Control Plane managed DNS servers. Valid values: "cname", "ns". Default: "cname".
- **gvc_link** (String) This value is set to a target GVC (using a full link) for use by subdomain based routing. Each workload in the GVC will receive a subdomain in the form ${workload.name}.${domain.name}. No not include if path based routing is used.

<a id="nestedblock--spec-ports"></a>

### `spec.ports`

Required:

- **tls** (Block List, Max: 1) ([see below](#nestedblock--spec--ports--tls))

~> **Note** If no tls properties are configured, an empty tls declaration (e.g., **tls { }**) is required to allow for the default properties to exist in the state file.

Optional:

- **number** (Number)
- **protocol** (String)
- **cors** (Block List, Max: 1) ([see below](#nestedblock--spec--ports--cors))

<a id="nestedblock--spec--ports--cors"></a>

### `spec.ports.cors`

Optional:

- **allow_origins** (Block List) ([see below](#nestedblock--spec--ports--cors--allow_origins))
- **allow_methods** (List of Strings)
- **allow_headers** (List of Strings)
- **max_age** (String)
- **allow_credentials** (Boolean)

<a id="nestedblock--spec--ports--cors--allow_origins"></a>

### `spec.ports.cors.allow_origins`

Optional:

- **exact** (String)

<a id="nestedblock--spec--ports--tls"></a>

### `spec.ports.tls`

- **min_protocol_version** (String)
- **cipher_suites** (String)
- **client_certificate** (Block List, Max: 1) ([see below](#nestedblock--spec--ports--tls--certificate))
- **server_certificate** (Block List, Max: 1) ([see below](#nestedblock--spec--ports--tls--certificate))

<a id="nestedblock--spec--ports--tls--certificate"></a>

### `spec.ports.tls.certificate`

Optional:

- **secret_link** (String) Full link to a TLS secret.

## Outputs

The following attributes are exported:

- **self_link** (String) Full link to this resource. Can be referenced by other resources.

## Import Syntax

To update a statefile with an existing domain resource, execute the following import command:

```terraform
terraform import cpln_domain.RESOURCE_NAME DOMAIN_NAME
```

-> 1. Substitute RESOURCE_NAME with the same string that is defined in the HCL file.<br/>2. Substitute DOMAIN_NAME with the corresponding domain defined in the resource.

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
