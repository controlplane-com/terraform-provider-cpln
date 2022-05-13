resource "cpln_agent" "example" {
  name        = "agent-example"
  description = "Example Agent"

  tags = {
    terraform_generated = "true"
    example             = "true"
  }
}
