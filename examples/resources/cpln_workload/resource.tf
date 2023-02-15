resource "cpln_gvc" "example" {
  name        = "gvc-example"
  description = "Example GVC"

  locations = ["aws-eu-central-1", "aws-us-west-2"]

  tags = {
    terraform_generated = "true"
    example             = "true"
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
}

resource "cpln_workload" "new" {

  gvc = cpln_gvc.example.name

  name        = "workload-example"
  description = "Example Workload"

  tags = {
    terraform_generated = "true"
    example             = "true"
  }

  identity_link = cpln_identity.example.self_link

  container {
    name        = "container-01"
    image       = "gcr.io/knative-samples/helloworld-go"
    port        = 8080
    memory      = "128Mi"
    cpu         = "50m"
    inherit_env = false

    env = {
      env-name-01 = "env-value-01",
      env-name-02 = "env-value-02",
    }

    args = ["arg-01", "arg-02"]

    readiness_probe {

      tcp_socket {
        port = 8181
      }

      period_seconds        = 11
      timeout_seconds       = 2
      failure_threshold     = 4
      success_threshold     = 2
      initial_delay_seconds = 1
    }

    liveness_probe {

      http_get {
        path   = "/path"
        port   = 8282
        scheme = "HTTPS"
        http_headers = {
          header-name-01 = "header-value-01"
          header-name-02 = "header-value-02"
        }
      }

      period_seconds        = 10
      timeout_seconds       = 3
      failure_threshold     = 5
      success_threshold     = 1
      initial_delay_seconds = 2
    }
  }

  options {
    capacity_ai     = false
    timeout_seconds = 30
    suspend         = false

    autoscaling {
      metric          = "concurrency"
      target          = 100
      max_scale       = 3
      min_scale       = 2
      max_concurrency = 500
    }
  }

  firewall_spec {
    external {
      inbound_allow_cidr      = ["0.0.0.0/0"]
      outbound_allow_cidr     = []
      outbound_allow_hostname = ["*.controlplane.com", "*.cpln.io"]
    }
    internal {
      # Allowed Types: "none", "same-gvc", "same-org", "workload-list"
      inbound_allow_type     = "none"
      inbound_allow_workload = []
    }
  }
}
