provider "cpln" {

  # Required
  # Can use CPLN_ORG Environment Variable
  org = var.org

  # Optional
  # Default Value: https://api.cpln.io
  # Can use CPLN_ENDPOINT Environment Variable
  endpoint = var.endpoint

  # Optional
  # Can use CPLN_PROFILE Environment Variable  
  profile = var.profile

  # Optional
  # Can use CPLN_TOKEN Environment Variable 
  token = var.token
}
