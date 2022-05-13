resource "cpln_service_account" "example" {

  name        = "service-account-example"
  description = "Example Service Account"

  tags = {
    terraform_generated = "true"
    example             = "true"
  }
}

resource "cpln_policy" "example" {

  name        = "policy-example"
  description = "Example Policy"

  tags = {
    terraform_generated = "true"
    example             = "true"
  }

  target_kind = "serviceaccount"

  # Policy can eith target `all` or specific target links

  // target = "all"
  target_links = [cpln_service_account.example.name]

  target_query {

    spec {
      # match is either "all", "any", or "none"
      match = "all"

      terms {
        op    = "="
        tag   = "firebase/sign_in_provider"
        value = "microsoft.com"
      }
    }
  }

  binding {

    # Available permissions are based on the target kind
    permissions = ["manage", "edit"]

    # Principal links format: `group/GROUP_NAME`, `user/USER_EMAIL`, `gvc/GVC_NAME/identity/IDENTITY_NAME`, `serviceaccount/SERVICE_ACCOUNT_NAME`
    principal_links = ["user/email@example.com", "group/viewers"]
  }

}
