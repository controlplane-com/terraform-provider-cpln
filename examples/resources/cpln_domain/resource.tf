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

resource "cpln_domain_route" "first_route" {
  depends_on = [cpln_domain.example_cname_routes]
  domain_name = "app.example.com"
  domain_port = 443

  prefix         = "/app"
  replace_prefix = "/replaceApp"
  workload_link  = "/org/myorg/gvc/mygvc/workload_two"
  port           = 80
}
