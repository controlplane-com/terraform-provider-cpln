# Terraform Provider for Control Plane

## Overview

This is the official [Control Plane](https://controlplane.com) Terraform provider. Terraform usage and reference documentation is available on the [Terraform Registry](https://registry.terraform.io/providers/controlplane-com/cpln/latest/docs).

## Prerequisites:

1. [Control Plane CLI](https://docs.controlplane.com/reference/cli#command-line-interface)
2. [Terraform CLI](https://www.terraform.io/downloads.html)
3. [Go](https://golang.org/doc/install)

## Installation

### From Source

Install `terraform-provider` locally.

- From a path under '${GOPATH}', clone the provider code:

```
git clone https://github.com/controlplane-com/terraform-provider-cpln.git
```

- cd into the cloned direction and install using make (Default OS Architecture is linux_amd64):

```
cd terraform-provider
make install

For macOS (Apple Silicion):
make install OS_ARCH=darwin_arm64

For macOS:
make install OS_ARCH=darwin_amd64
```

The provider is installed under the `~/.terraform.d/plugins/` directory.

### Alternative Local Execution of the Provider

Considering the provider is installed at path below after running the `make install` command:

```
/Users/<username>/.terraform.d/plugins/controlplane.com/com/cpln/1.0.10/darwin_arm64
```

1. Add a `.terraform.rc` file on root with the contents below

   ```hcl
   provider_installation {
     filesystem_mirror {
       path    = "/Users/<username>/.terraform.d/plugins"
     }
     direct {
       exclude = ["controlplane.com/*/*"]
     }
   }
   ```

2. In your Terraform config (`main.tf`), set:

   ```hcl
   terraform {
     required_providers {
       cpln = {
         source  = "controlplane.com/com/cpln"
         version = "<version>"
       }
     }
   }
   ```

3. Run `terraform init`. You should see the warning below if it succeeded:

    ```
    Warning: Incomplete lock file information for providers

    Due to your customized provider installation methods, Terraform was forced to
    calculate lock file checksums locally for the following providers:
    - controlplane.com/com/cpln
    ```

### Debugging

When debugging, create the file at `.vscode/launch.json` by copying and editing it from the `launch.json.example` file.

## Examples

See the [HCL](https://www.terraform.io/docs/configuration/syntax.html) examples within the `/examples` directory.

## Quickstart

```bash
# 1. Log in via CLI
cpln login

# (For test environment)
cpln profile update default --login --endpoint https://api.test.cpln.io

# 2. Initialize Terraform
terraform init

# 3. Apply a sample configuration
terraform apply -var="org=<YOUR_ORG>"

# 4. When you’re done:
terraform destroy -var="org=<YOUR_ORG>"
```

For full examples, see the [examples/](./examples) directory.

## Testing

- **Unit Tests**:
  ```bash
  make test
  ```
- **Acceptance Tests**:
  ```bash
  make testacc
  ```

Make sure to configure a valid `org` in `internal/provider/config.go` before running acceptance tests.

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

## Generate Reference Documentation

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

macOS (Apple Silicon)

```
tar -cvzf terraform-provider-cpln_1.0.0_darwin_amd64.tgz terraform-provider-cpln_1.0.0_darwin_arm64
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
- git tag vX.X.X
- git push origin vX.X.X

## Notes

- Needed to add CGO_ENABLED=0 for linux build when running within an alpine image. See: https://github.com/Mastercard/terraform-provider-restapi/issues/65
- CLI Config File: https://www.terraform.io/docs/cli/config/config-file.html

## Development

When adding new resources or data sources:

1. Update the API client schema in `internal/provider/client/<item>.go`.
2. Define or adjust the Terraform schema in:
   - `internal/provider/resource_<item>.go` (resources)
   - `internal/provider/data_source_<item>.go` (data sources)
3. Update the CRUD/context functions under each resource.
4. Add or update tests in `*_test.go`.
5. Regenerate the reference docs via:
   ```bash
   cd scripts
   ./cpln_docs.sh
   # then copy terraform_reference.mdx into your docs site
   ```
6. Update resource / data source examples.

## CHANGELOG

All version history and release notes have been moved to [CHANGELOG.md](./CHANGELOG.md).
