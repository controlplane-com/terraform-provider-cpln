resource "cpln_gvc" "example" {

  name        = "gvc-example"
  description = "Example GVC"

  locations = ["aws-us-west-2"]

  tags = {
    terraform_generated = "true"
    example             = "true"
  }
}

resource "cpln_agent" "example" {

  name        = "agent-example"
  description = "Example Agent"
}

resource "cpln_cloud_account" "example_aws" {

  name        = "aws-example"
  description = "Example AWS Cloud Account"

  tags = {
    terraform_generated = "true"
    example             = "true"
  }

  aws {
    role_arn = "arn:aws:iam::1234:role/example_role"
  }
}

resource "cpln_cloud_account" "example_azure" {

  name        = "azure-example"
  description = "Example Azure Cloud Account"

  tags = {
    terraform_generated = "true"
    example             = "true"
  }

  azure {
    // Use the full link
    secret_link = "/org/ORG_NAME/secret/AZURE_SECRET"
  }
}

resource "cpln_cloud_account" "example-gcp" {

  name        = "gcp-example"
  description = "Example GCP Cloud Account"

  tags = {
    terraform_generated = "true"
    example             = "true"
  }

  gcp {
    project_id = "cpln_gcp_project_1234"
  }
}

resource "cpln_cloud_account" "test_ngs_cloud_account" {

  name        = "tf-ca-ngs-${var.random_name}"
  description = "cloud account description tf-ca-ngs"

  tags = {
    terraform_generated = "true"
    acceptance_test     = "true"
  }

  ngs {
    // Use the full link for now
    // secret_link = "//secret/tf_secret_nats_account"
    secret_link = "/org/${var.org_name}/secret/tf_secret_nats_account"
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

  # Network Resource with FQDN
  network_resource {
    name       = "test-network-resource-fqdn"
    agent_link = cpln_agent.example.self_link
    fqdn       = "domain.example.com"
    ports      = [1234, 5432]
  }

  # Network Resource with FQDN and Resolver IP
  network_resource {
    name        = "test-network-resource-fqdn-rip"
    agent_link  = cpln_agent.example.self_link
    fqdn        = "domain2.example.com"
    resolver_ip = "192.168.1.1"
    ports       = [12345, 54321]
  }

  # Network Resource with IP
  network_resource {
    name       = "test-network-resource-ip"
    agent_link = cpln_agent.test_agent.self_link
    ips        = ["192.168.1.1", "192.168.1.250"]
    ports      = [3099, 7890]
  }

  aws_identity {

    cloud_account_link = cpln_cloud_account.example_aws.self_link

    # AWS Identity can either contain an existing role_name or multiple policy_refs

    // role_name = "rds-monitoring-role"

    policy_refs = ["aws::/job-function/SupportUser", "aws::AWSSupportAccess"]
  }

  azure_identity {

    cloud_account_link = cpln_cloud_account.example_azure.self_link

    role_assignment {
      scope = "/subscriptions/d0d1e522-0825-415a-8b07-f7759b5c8a7e/resourceGroups/CP-Test-Resource-Group"
      roles = ["AcrPull", "AcrPush"]
    }

    role_assignment {
      scope = "/subscriptions/d0d1e522-0825-415a-8b07-f7759b5c8a7e/resourceGroups/CP-Test-Resource-Group/providers/Microsoft.Storage/storageAccounts/cplntest"
      roles = ["Support Request Contributor"]
    }
  }

  gcp_identity {

    cloud_account_link = cpln_cloud_account.example_gcp.self_link
    scopes             = ["https://www.googleapis.com/auth/cloud-platform"]

    # GCP Identity can either contain an existing service_account or multiple bindings

    // service_account = "cpln-tf@cpln-test.iam.gserviceaccount.com"

    binding {
      resource = "//cloudresourcemanager.googleapis.com/projects/cpln-test"
      roles    = ["roles/appengine.appViewer", "roles/actions.Viewer"]
    }

    binding {
      resource = "//iam.googleapis.com/projects/cpln-test/serviceAccounts/cpln-tf@cpln-test.iam.gserviceaccount.com"
      roles    = ["roles/editor", "roles/iam.serviceAccountUser"]
    }
  }

  ngs_access_policy {

    cloud_account_link = cpln_cloud_account.test_ngs_cloud_account.self_link

    pub {
      allow = ["pa1", "pa2"]
      deny  = ["pd1", "pd2"]
    }

    sub {
      allow = ["sa1", "sa2"]
      deny  = ["sd1", "sd2"]
    }

    resp {
      max = 1
      ttl = "5m"
    }

    subs    = 1
    data    = 2
    payload = 3
  }
}
