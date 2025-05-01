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

~> **Note** Only one of `prefix` OR `regex` may be provided in a single resource.

- **domain_link** (String) The self link of the domain to add the route to.
- **domain_port** (int) The port the route corresponds to. Default: 443
- **prefix** (String) The path will match any unmatched path prefixes for the subdomain.
- **regex** (String) Used to match URI paths. Uses the google re2 regex syntax.
- **workload_link** (String) The link of the workload to map the prefix to.

### Optional

~> **Note** Only one of `host_prefix` OR `host_regex` may be provided in a single resource.

- **replace_prefix** (String) A path prefix can be configured to be replaced when forwarding the request to the Workload.
- **port** (Number) For the linked workload, the port to route traffic to.
- **host_prefix** (String) This option allows forwarding traffic for different host headers to different workloads. This will only be used when the target GVC has dedicated load balancing enabled and the Domain is configured for wildcard support. Please contact us on Slack or at support@controlplane.com for additional details.
- **host_regex** (String) A regex to match the host header. This will only be used when the target GVC has dedicated load balancing enabled and the Domain is configure for wildcard support. Contact your account manager for details.
- **headers** (Block List, Max: 1) ([see below](#nestedblock--headers))

<a id="nestedblock--headers"></a>

### `headers`

Modify the headers for all http requests for this route.

Optional:

- **request** (Block List, Max: 1) ([see below](#nestedblock--headers-request))

<a id="nestedblock--headers-request"></a>

### `headers.request`

Manipulates HTTP headers.

Optional:

- **set** (Map of String) Sets or overrides headers to all http requests for this route.

## Example Usage

### Prefix

#### With Host Prefix

```terraform
resource "cpln_domain" "apex" {
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

resource "cpln_domain" "subdomain" {

  depends_on = [cpln_domain.apex]

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

resource "cpln_domain_route" "first-route" {

  // The first route depends on the domain being created first
  depends_on  = [cpln_domain.subdomain]

  domain_link = cpln_domain.subdomain.self_link
  domain_port = 443

  prefix = "/example-1"
  replace_prefix = "/replace_example"
  host_prefix = "www.example.com"
  workload_link = "LINK_TO_WORKLOAD"
  port = 80

  headers {
    request {
      set = {
        Host = "example.com"
        "Content-Type" = "application/json"
      }
    }
  }
}

resource "cpln_domain_route" "second-route" {

  // The second route depends on the first route
  depends_on  = [cpln_domain_route.first-route]

  domain_link = cpln_domain.subdomain.self_link
  domain_port = 443

  prefix = "/example-2"
  replace_prefix = "/"
  host_prefix = "www.foo.com"
  workload_link = "LINK_TO_WORKLOAD"
  port = 80

  headers {
    request {
      set = {
        Host = "example.com"
        "Content-Type" = "application/json"
      }
    }
  }
}
```

#### With Host Regex

```terraform
resource "cpln_domain" "apex" {
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

resource "cpln_domain" "subdomain" {

  depends_on = [cpln_domain.apex]

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

resource "cpln_domain_route" "first-route" {

  // The first route depends on the domain being created first
  depends_on  = [cpln_domain.subdomain]

  domain_link = cpln_domain.subdomain.self_link
  domain_port = 443

  prefix = "/example-1"
  replace_prefix = "/replace_example"
  host_prefix = "www.example.com"
  workload_link = "LINK_TO_WORKLOAD"
  port = 80

  headers {
    request {
      set = {
        Host = "example.com"
        "Content-Type" = "application/json"
      }
    }
  }
}

resource "cpln_domain_route" "second-route" {

  // The second route depends on the first route
  depends_on  = [cpln_domain_route.first-route]

  domain_link = cpln_domain.subdomain.self_link
  domain_port = 443

  prefix = "/example-2"
  replace_prefix = "/"
  host_regex = "req"
  workload_link = "LINK_TO_WORKLOAD"
  port = 80

  headers {
    request {
      set = {
        Host = "example.com"
        "Content-Type" = "application/json"
      }
    }
  }
}
```

### Regex

```terraform
resource "cpln_domain" "apex" {
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

resource "cpln_domain" "subdomain" {

  depends_on = [cpln_domain.apex]

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

resource "cpln_domain_route" "new" {

    // The first route depends on the domain being created first
    depends_on  = [cpln_domain.subdomain]

    domain_link = cpln_domain.subdomain.self_link
    domain_port = 443

    regex          = "/user/.*/profile"
    replace_prefix = "/replace-example"
    host_prefix    = "www.example.com"
    workload_link  = "LINK_TO_WORKLOAD"
    port           = 80

    headers {
      request {
        set = {
          Host = "example.com"
          "Content-Type" = "application/json"
        }
      }
    }
}
```

## Import Syntax

The `terraform import` command is used to bring existing infrastructure resources, created outside of Terraform, into the Terraform state file, enabling their management through Terraform going forward.

To update a statefile with an existing domain route resource, execute the following import command:

```terraform
terraform import cpln_domain_route.RESOURCE_NAME DOMAIN_LINK:DOMAIN_PORT:[PREFIX|REGEX]
```

-> 1. Substitute RESOURCE_NAME with the same string that is defined in the HCL file.<br/>2. Substitute DOMAIN_LINK with the corresponding domain link defined in the resource.<br/>3. Substitute DOMAIN_PORT with the corresponding domain port defined in the resource.<br/>4. Substitute PREFIX with the corresponding prefix defined in the resource.
