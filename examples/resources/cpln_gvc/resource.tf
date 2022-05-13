resource "cpln_gvc" "example" {

  name        = "gvc-example"
  description = "Example GVC"

  # Example Locations: `aws-eu-central-1`, `aws-us-west-2`, `azure-east2`, `gcp-us-east1`
  locations = ["aws-eu-central-1", "aws-us-west-2"]

  # domain = "app.example.com"

  tags = {
    terraform_generated = "true"
    example             = "true"
  }
}
