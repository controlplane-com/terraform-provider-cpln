module terraform-provider-cpln

go 1.15

require (
	github.com/go-test/deep v1.0.7
	github.com/hashicorp/terraform-plugin-docs v0.4.0
	github.com/hashicorp/terraform-plugin-sdk/v2 v2.7.1
)

replace github.com/hashicorp/terraform-plugin-sdk/v2 => github.com/controlplane-com/terraform-plugin-sdk/v2 v2.7.2
