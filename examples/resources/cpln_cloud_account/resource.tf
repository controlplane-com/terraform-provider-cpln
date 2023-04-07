# AWS Cloud Account 
resource "cpln_cloud_account" "aws" {

  name        = "cloud-account-aws"
  description = "AWS cloud account"

  tags = {
    terraform_generated = "true"
    example             = "true"
  }

  aws {
    role_arn = "arn:aws:iam::1234:role/test_role"
  }
}

# Azure Cloud Account 
resource "cpln_cloud_account" "azure" {

  name        = "cloud-account-azure"
  description = "Azure cloud account "

  tags = {
    terraform_generated = "true"
    example             = "true"
  }

  azure {
    // Use full link
    secret_link = "/org/ORG_NAME/secret/AZURE_SECRET"
  }
}

# GCP Cloud Account 
resource "cpln_cloud_account" "gcp" {

  name        = "cloud-account-gcp"
  description = "GCP cloud account"

  tags = {
    terraform_generated = "true"
    example             = "true"
  }

  gcp {
    project_id = "cpln_gcp_project_1234"
  }
}

# NGS Cloud Account 
resource "cpln_cloud_account" "ngs" {

  name        = "cloud-account-ngs"
  description = "NGS cloud account "

  tags = {
    terraform_generated = "true"
    example             = "true"
  }

  ngs {
    // Use full link
    secret_link = "/org/ORG_NAME/secret/NATS_ACCOUNT_SECRET"
  }
}