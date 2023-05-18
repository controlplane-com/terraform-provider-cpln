---
page_title: "cpln_domain_route Resource - terraform-provider-cpln"
subcategory: "Domain"
description: |-

---

# cpln_domain_route (Resource)

Manages a domain's [Routes](https://docs.controlplane.com/reference/domain#path-based-routing).

Used in conjunction with a Domain.

~> **Note** It is mandatory to use the `depends_on` clause on each `cpln_domain_route`. Each route needs to depend on its upstream route (e.g., the second route depends on the first route, etc.). The first route needs to depend on the domain it is linked to.

## Declaration

### Required

- **domain_link** (String) The self link of the domain to add the route to.
- **domain_port** (int) The port the route corresponds to. Default: 443
- **prefix** (String) The path will match any unmatched path prefixes for the subdomain. 
- **workload_link** (String) The link of the workload to map the prefix to.

### Optional

- **replace_prefix** (String) A path prefix can be configured to be replaced when forwarding the request to the Workload.
- **port** (Number) For the linked workload, the port to route traffic to.

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

resource "cpln_domain" "example_cname_routes" {

  depends_on = [cpln_domain.domain_apex]

  name        = "app.example.com"
  description = "Custom domain that can be set on a GVC and used by associated workloads"

  tags = {
    terraform_generated = "true"
    example             = "true"
  }

  spec {
    dns_mode = "ns"

    ports {
      number   = 443
      protocol = "http2"

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
      }
    }
  }
}

resource "cpln_domain_route" "example_route" {

    // The first route depends on the domain being created first
    depends_on  = [cpln_domain.example_cname_routes]

    domain_link = cpln_domain.example_cname_routes.self_link
    domain_port = 443

    prefix = "/example"
    replace_prefix = "/replace_example"
    workload_link = "LINK_TO_WORKLOAD"
    port = 80
}

resource "cpln_domain_route" "example_second_route" {

    // The second route depends on the first route
    depends_on  = [cpln_domain_route.example_route]

    domain_link = cpln_domain.example_cname_routes.self_link
    domain_port = 443

    prefix = "/example_second_route"
    replace_prefix = "/"
    workload_link = "LINK_TO_WORKLOAD"
    port = 80
}
```