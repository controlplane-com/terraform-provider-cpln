resource "cpln_domain" "example" {

  name        = "app.example.com"
  description = "Custom domain that can be set on a GVC and used by associated workloads"

  tags = {
    terraform_generated = "true"
    example             = "true"
  }
  
  spec {
    dns_mode         = "string"
    gvc_link         = "string"
    accept_all_hosts = "true"

    ports {
      number   = 443
      protocol = "http"

      routes {
        prefix         = "string"
        replace_prefix = "string"
        workload_link  = "string"
        port           = 80
      }

      cors {
        allow_origins {
          exact = "string"
        }

        allow_methods     = ["allow_method_1", "allow_method_2", "allow_method_3"]
        allow_headers     = ["allow_header_1", "allow_header_2", "allow_header_3"]
        max_age           = "24h"
        allow_credentials = "true"
      }

      tls {
        min_protocol_version = "string"
        cipher_suites        = "string"

        client_certificate {
          secret_link = "string"
        }

        server_certificate {
          secret_link = "string"
        }
      }
    }

    status {
      end_points {
        url           = "string"
        workload_link = "string"
      }

      status  = "string"
      warning = "string"

      locations {
        name               = "string"
        certificate_status = "string"
      }

      fingerprint = "string"
    }
  }
}