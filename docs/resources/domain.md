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

## Declaration

### Required

- **name** (String) Domain name. Must be a valid domain name with at least three segments (e.g., test.example.com). Control Plane will validate the existence of the domain with DNS. Create and Update will fail if the required DNS entries cannot be validated.

### Optional

- **description** (String) Description for the domain name.
- **tags** (Map of String) Key-value map of resource tags.
- **spec** (Block List, Max: 1) ([see below](#nestedblock--spec))
- **status** (Block List, Max: 1) ([see below](#nestedblock--status))

<a id="nestedblock--spec"></a>
### `spec`

Optional:

-- **dns_mode** (String) In 'cname' dnsMode, Control Plane will configure workloads to accept traffic for the domain but will not manage DNS records for the domain. End users configure CNAME records in their own DNS pointed to the canonical workload endpoint. Currently 'cname' dnsMode requires that a tls.serverCertificate is configured when subdomain based routing is used. In 'ns' dnsMode, Control Plane will manage the subdomains and create all necessary DNS records. End users configure an NS record to forward DNS requests to the Control Plane managed DNS servers.
-- **gvc_link** (String) One of gvcLink and routes may be provided. When gvcLink is configured each workload in the GVC will receive a subdomain in the form ${workload.name}.${domain.name}
-- **accept_all_hosts** (Boolean)
-- **ports** (Block List) ([see below](#nestedblock--spec-ports))

<a id="nestedblock--spec-ports"></a>
### `spec.ports`

Optional:

-- **number** (Number)
-- **protocol** (String)
-- **cors** (Block List, Max: 1) ([see below](#nestedblock--spec--ports--cors))
-- **tls** (Block List, Max: 1) ([see below](#nestedblock--spec--ports--tls))

<a id="nestedblock--spec--ports--cors"></a>
### `spec.ports.cors`

Optional:

-- **allow_origins** (Block List) ([see below](#nestedblock--spec--ports--cors--allow_origins))
-- **allow_methods** (List of Strings)
-- **allow_headers** (List of Strings)
-- **max_age** (String)
-- **allow_credentials** (Boolean)

<a id="nestedblock--spec--ports--cors--allow_origins"></a>
### `spec.ports.cors.allow_origins`

Optional:

-- **exact** (String)

<a id="nestedblock--spec--ports--tls"></a>
### `spec.ports.tls`

-- **min_protocol_version** (String)
-- **cipher_suites** (String)
-- **client_certificate** (Block List, Max: 1) ([see below](#nestedblock--spec--ports--tls--certificate))
-- **server_certificate** (Block List, Max: 1) ([see below](#nestedblock--spec--ports--tls--certificate))

<a id="nestedblock--spec--ports--tls--certificate"></a>
### `spec.ports.tls.certificate`

Optional:
-- **secret_link** (String)

<a id="nestedblock--status"></a>
### `status`

Optional

-- **endpoints** (Block List) ([see below](#nestedblock--status--endpoints))
-- **status** (String)
-- **warning** (String)
-- **locations** (Block List) ([see below](#nestedblock--status--locations))
-- **fingerprint** (String)

<a id="nestedblock--status--endpoints"></a>
### `status.endpoints`

Optional:

-- **url** (String)
-- **workload_link** (String)

<a id="nestedblock--status--locations"></a>
### `status.locations`

Optionals:

-- **name** (String)
-- **certificate_status** (String)

## Outputs

The following attributes are exported:

- **self_link** (String) Full link to this resource. Can be referenced by other resources. 

## Example Usage

```terraform
resource "cpln_domain" "example_ns_subdomain" {

  name        = "app.example.com"
  description = "Custom domain that can be set on a GVC and used by associated workloads"

  tags = {
    terraform_generated = "true"
    example             = "true"
  }

  spec {
    dns_mode         = "ns"
    gvc_link         = "/org/myorg/gvc/mygvc"
    accept_all_hosts = "true"

    ports {
      number   = 443
      protocol = "http"

      cors {
        allow_origins {
          exact = "example.com"
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

        client_certificate {}
      }
    }
  }
}
```