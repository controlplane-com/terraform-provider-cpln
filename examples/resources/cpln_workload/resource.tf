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

  type = "cron"
	  
  container {
    name  = "container-01"
    image = "gcr.io/knative-samples/helloworld-go"
    memory = "128Mi"
    cpu = "50m"

    command = "override-command"
    working_directory = "/usr"
  
    env = {
    env-name-01 = "env-value-01",
    env-name-02 = "env-value-02",
    }
  
    args = ["arg-01", "arg-02"]

    volume {
    uri  = "s3://bucket"
    path = "/testpath01"
    }

    volume {
    uri  = "azureblob://storageAccount/container"
    path = "/testpath02"
    }

    metrics {
    path = "/metrics"
    port = 8181
    }
  }
          
  options {
    capacity_ai = false
    spot = true
    timeout_seconds = 5
  
    autoscaling {
    target = 100
    max_scale = 1
    min_scale = 1
    max_concurrency = 0
    scale_to_zero_delay = 300
    }
  }
  
  firewall_spec {
    external {
    inbound_allow_cidr =  ["0.0.0.0/0"]
    // outbound_allow_cidr =  []
    outbound_allow_hostname =  ["*.controlplane.com", "*.cpln.io"]
    }
    internal { 
    # Allowed Types: "none", "same-gvc", "same-org", "workload-list"
    inbound_allow_type = "none"
    inbound_allow_workload = []
    }
  }

  job {
    schedule = "* * * * *"
    concurrency_policy = "Forbid"
    history_limit = 5
    restart_policy = "Never"
    active_deadline_seconds = 1200
  }
}