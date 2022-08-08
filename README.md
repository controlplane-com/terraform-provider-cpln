# Terraform Provider for Control Plane


## Prerequisites:

1. [Control Plane CLI](https://docs.controlplane.com/reference/cli#command-line-interface) 
2. [Terraform CLI](https://www.terraform.io/downloads.html) 
3. [Go](https://golang.org/doc/install) 


## Installation

Install `terraform-provider` locally.
- From a path under '${GOPATH}', clone the provider code:
```
git clone https://github.com/controlplane-com/terraform-provider-cpln.git
```
- cd into the cloned direction and install using make (Default OS Architecture is linux_amd64):
```
cd terraform-provider
make install

For macOS:
make install OS_ARCH=darwin_amd64
 ```
The provider is installed under the `~/.terraform.d/plugins/` directory.

## Examples

See the [HCL](https://www.terraform.io/docs/configuration/syntax.html) examples within the `/examples` directory.

## Example Usage

```
$ cd examples

Edit main.tf file.
$ vim main.tf 

Update the 'org' variable within the 'Provider' configuration with a valid organization you are authorized to modify.

Login To Control Plane via the CLI
$ cpln login

```
Login to test environment:

cpln profile update default --login --endpoint https://api.test.cpln.io
```

Initialize Terraform
$ terraform init

Create the infrastructure declared in the HCL (enter a valid organization name)
$ terraform apply -var 'org=ENTER-VALID-ORG'

Remove infrastructure
$ terraform destroy -var 'org=ENTER-VALID-ORG'
```

## Testing

Edit the `provider_test.go` file. Update the provider config on line 53 with a valid organization you are authorized to modify.

### [Unit Tests](https://www.terraform.io/docs/extend/testing/unit-testing.html)

```
$ make test
```

### [Acceptance Tests](https://www.terraform.io/docs/extend/testing/acceptance-tests/index.html)

```
$ make testacc
```


## Control Plane CLI Helper Notes:

1. Creating a new service account
```
    a) Create service account:
        'cpln serviceaccount create --name <service_account_name> --org <organization_name>'

    b) Edit group:
        'cpln group edit superusers --org <organization_name>'

        Using the editor, add the service account to the memberLinks element using the format:
        '/org/<organziation_name>/serviceaccount/<service_account_name>'

    c) Add key to service account:
        'cpln serviceaccount add-key <service_account_name> --org <organization_name> --description <key_description>'

    d) Create a profile with the token output from step c:
        'cpln profile create <profile_name> --token <token>'
```

2. To obtain a valid access token via the Control Plane CLI:
```
    a) Using the default profile:
        'cpln profile token'
        OR
        'cpln profile token default'

    b) Using a created profile:
        'cpln profile token <profile_name>'
```

3. Misc WSL Commands
```
    a) Allow WSL network through Windows Firewall:
        New-NetFirewallRule -DisplayName "WSL" -Direction Inbound  -InterfaceAlias "vEthernet (WSL)"  -Action Allow
    
```


## Terraform Documentation Helper Links

1. [Provider Documentation](https://www.terraform.io/docs/registry/providers/docs.html)
2. [Document Preview](https://registry.terraform.io/tools/doc-preview)


## Generate Reference Doc

1. Run the script `cpln_docs.sh` within the /scripts directory.
2. Copy the `terraform_reference.mdx` file to the directory `/pages/terraform` in the documentation project.


## Compress Commands

1. Run `make release`
2. Plugins will be in /bin directory
3. Use commands below to compress (update version)

macOS
```
tar -cvzf terraform-provider-cpln_1.0.0_darwin_amd64.tgz terraform-provider-cpln_1.0.0_darwin_amd64
```

Linux
```
tar -cvzf terraform-provider-cpln_1.0.0_linux_amd64.tgz terraform-provider-cpln_1.0.0_linux_amd64
```

Windows
```
tar -cvzf terraform-provider-cpln_1.0.0_windows_amd64.zip terraform-provider-cpln_1.0.0_windows_amd64.exe
```
## Terraform Registry Publishing

- https://www.terraform.io/registry/providers/publishing

## Notes

- Needed to add CGO_ENABLED=0 for linux build when running within an alpine image. See: https://github.com/Mastercard/terraform-provider-restapi/issues/65

- CLI Config File: https://www.terraform.io/docs/cli/config/config-file.html

## Version Info

- v1.0.2 - Publish To Terraform Registry.
- v1.0.3 - Update docs.
- v1.0.4 - Add standard workload type.
- v1.0.5 - Add Org External Logging.
- v1.0.6 - Add Org/Gvc Tracing (lightstep).
