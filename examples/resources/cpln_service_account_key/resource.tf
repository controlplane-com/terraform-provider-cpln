resource "cpln_service_account" "example" {

  name        = "service-account-example"
  description = "Example Service Account"

  tags = {
    terraform_generated = "true"
    example             = "true"
  }
}

resource "cpln_service_account_key" "example" {

  service_account_name = cpln_service_account.example.name
  description          = "Service Account Key"
}


resource "cpln_service_account_key" "example_02" {

  // When adding another key, use `depends_on` to add the keys synchronously 
  depends_on = [cpln_service_account_key.example]

  service_account_name = cpln_service_account.example.name
  description          = "Service Account Key #2"
}

output "key_01" {
  value = cpln_service_account_key.example.key
}

output "key_02" {
  value = cpln_service_account_key.example_02.key
}
