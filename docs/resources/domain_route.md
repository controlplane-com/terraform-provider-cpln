---
page_title: "cpln_service_account_key Resource - terraform-provider-cpln"
subcategory: "Domain Route"
description: |-

---

# cpln_domain_route (Resource)

Manages a domain's [Routes](https://docs.controlplane.com/reference/domain#path-based-routing).
Used in conjunction with a Domain.

## Declaration

### Required

- **domain_link** (String) The self link of the domain to add the route to.
- **prefix** (String) The path will match any unmatched path prefixes for the subdomain. Default: `/`
- **workload_link** (String) The link of the workload to map the prefix to.

### Optional

- **replace_prefix** (String) A path prefix can be configured to be replaced when forwarding the request to the Workload..
- **port** (Number) // TODO: Add description.

## Outputs

The following attributes are exported:

// TODO: Add outputs.

## Example Usage

```terraform
resource "cpln_domain" "example_cname_routes" {

  name        = "app.example.com"
  description = "Custom domain that can be set on a GVC and used by associated workloads"

  tags = {
    terraform_generated = "true"
    example             = "true"
  }

  spec {
    dns_mode         = "ns"
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

resource "cpln_domain_route" "example_route" {
    domain_link = cpln_domain.example_cname_routes.self_link

    prefix = "/example"
    replace_prefix = "/replace_example"
    workload_link = "self/link/to/workload"
    port = 80
}
```