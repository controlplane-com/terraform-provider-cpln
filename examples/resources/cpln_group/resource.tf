resource "cpln_service_account" "example" {

  name        = "service-account-example"
  description = "Service Account to be used as a member of a group"

  tags = {
    terraform_generated = "true"
    example             = "true"
  }
}

resource "cpln_group" "example" {

  name        = "group-example"
  description = "group description ${var.random-name}"

  tags = {
    terraform_generated = "true"
    example             = "true"
  }

  user_id_email   = ["valid_user@example.com"]
  service_account = [cpln_service_account.example.name]

  member_query {

    fetch = "items"

    spec {
      match = "all"

      terms {
        op    = "="
        tag   = "firebase/sign_in_provider"
        value = "microsoft.com"
      }
    }
  }
}
