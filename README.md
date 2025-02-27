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

Add a `.terraform.rc` file on root with the contents below

```
provider_installation {
  filesystem_mirror {
    path    = "/Users/<username>/.terraform.d/plugins"
  }
  direct {
    exclude = ["terraform.local/*/*"]
  }
}
```

Then change the `required_provider` source to `controlplane.com/com/cpln` within `main.tf`

Init Terraform project with the local provider now by running `terraform init`

You should see the warning below if it succeeded:

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

`cpln profile update default --login --endpoint https://api.test.cpln.io`

Note: In provider section, the endpoint for test is `https://api.test.cpln.io`.

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

## Notes for Developing New Features

- Update item schema under internal/provider/client/<item>.go (Needs to follow data-service api structure)
- Update item's Terraform resource schema under "internal/provider/resource\_<item>.go"
- Update ResourceCreate Context
- Update ResourceRead Context
- Update ResourceUpdate Context
- Update ResourceDelete Context (If needed)
- Update/Add Terraform resource test
- Update dataSource of the item
- Update resource example
- Update resource documentation

### Notes

- Flatten methods transform api object to terraform resource

## Version Info

- v1.0.2 - Publish To Terraform Registry.
- v1.0.3 - Update docs.
- v1.0.4 - Add standard workload type.
- v1.0.5 - Add Org External Logging.
- v1.0.7 - Add Org/Gvc Tracing (lightstep).
- v1.0.8 - Add Gvc Data Source.
- v1.0.9 - Only remove certain `cpln/*` server generated tags. Increase max containers.
- v1.0.10 - Add Location / Locations Data Source.
- v1.0.11 - Fix issue with secrets having json types. Remove built-in server generated secret tags.
- v1.0.12 - Update GitHub Action.
- v1.0.121 - HotFix for removal of workload option spot property.
- v1.0.122 - HotFix for new and missing workload properties.
- v1.0.123 - Updates for -refresh-only flag.
- v1.0.13 - Add workload lifecycle hooks (post start / pre stor). Add GVC Environment Variables. Add workload suspend.
- v1.1.0 - Update to Go 1.1.8 and Terraform SDK 2.25.0. Add workload lifecycle hooks (post start / pre stor). Add GVC Environment Variables. Add workload suspend.
- v1.1.1 - Add identity manager to group.
- v1.1.2 - Added `NATS Account` secret. Added NGS cloud account. Added NGS cloud access policy. Sync'ed GVC schema (env will be returned). Added GCP service account name and roles output to cloud account. Added elastic logging to org external logging. Added audit context resource. Added native network resources.
- v1.1.3 - Fixed issue with identity and workload policies. 'gvc' property now required for those policy kinds.
- v1.1.4 - Added domain and domain route.
- v1.1.5 - Updated Terraform SDK to v2.26.1. Added CRON job workload type. Add workload rollout and security options. Add disabled scaling strategy. Add GVC load balancer. Add workload support dynamic tags.
- v1.1.6 - Added volume sets. Added GPU to workload. Added external firewall outbound allow ports. Add domain host prefix.
- v1.1.7 - Update all resources to allow import. Updated docs with import details and syntax.
- v1.1.8 - Add to secret the output `dictionary_as_env`. Updated import docs. Initial logic for the deprecation of the `port` container attribute. Updated SDK to v2.27.0.
- v1.1.9 - Fix issue with volume set status.locations. Fix import domain route.
- v1.1.10 - Add to org logging to multiple external provider. Update volumeset performance classes. Add cloud account data source that has aws identifiers as output. Fix issue with workload import and legacy port.
- v1.1.11 - Add Otel tracing to org and gvc.
- v1.1.12 - Fix bug when passing null container command args.
- v1.1.13 - Update dependencies. Add to a volume set policy that a gvc is required.
- v1.1.14 - Fix bug with secrets when updating. Added `trusted_proxies` to GVC load balancer.
- v1.1.15 - Fix issue with tag values that were stored as number types.
- v1.1.16 - Add org creation. Add org properties.
- v1.1.17 - Add schedule to volume set.  Add status to domain output. Add Control Plane tracing and custom tags to tracing. Add external ID to ECR secret. Add generic elastic to org logging. Add geo properties to locations. Add missing properties to workload.
- v1.1.18 - Add min CPU and Memory to container. 
- v1.1.19 - Add external ID to AWS secret.
- v1.1.20 - Update GVC data-source and resource docs. Add descriptions in the schema for all data-sources and resources. Add image and images data-sources.
- v1.1.21 - Fix workload validation. Add default to autoscaling metric property. Add stateful workload example in docs. Update dependencies and Go to v1.21
- v1.1.22 - Add formatted link for secret and volume set. Update workload autoscaling to be optional.
- v1.1.23 - Update docs to indicate the locations are optional for a GVC. Update image data source to return the latest image if there is not tag. pdate images data source to accept a query. Add storage class suffix to volume set.
- v1.1.24 - Update images data source to fetch all images.
- v1.1.25 - Handle case when no images are found. Add cloudwatch, fluentd, and stack driver to org logging.
- v1.1.26 - Add mk8s resource. Add location resource. Add syslog to org logging.
- v1.1.27 - Add additional descriptions. Add sysbox, hetznerLabels, awsTags to mk8s. Add regex to domain-route resource. Add threat-detection to org.
- v1.1.28 - Fix mk8s add-on update flow.
- v1.1.29 - Fix mk8s addons boolean value issue
- v1.1.30 - Fix mk8s addons being removed on update
- v1.1.31
    - Add a custom location resource.
    - Fix typo in org doc.
    - Add floating ip selector to mk8s resource.
    - Add external fields to cloud watching logging in org logging resource.
    - Update policy resource to allow self-link references.
    - Add headers to domain route resource.
    - Force a tag indicating resource was created by terraform.
    - Set networking as required in mk8s resource.
- v1.1.32
    - Update docs to include import syntax.
    - Ignore server-side tags for location resource.
    - Enable term rel in query.
    - Fix bug in tracing.
    - Fix plan not empty issue in cloud account.
    - Fix tag format.
    - Add geo locations to security options in workload.
    - Add load balancer to workload resource.
- v1.1.33
    - Add import support to mk8s resource.
- v1.1.34
    - Make workload options optional.
- v1.1.35
    - Fix empty options bug in workload resource.
- v1.1.36
    - Add secret data source.
    - Add autoscaling by memory.
- v1.1.37
    - Add cpln_ipset resource.
    - Add redirect to GVC load balancer.
    - Add lambdalabs to mk8s.
    - Add linode to mk8s.
    - Add oblivus to mk8s.
    - Add paperspace to mk8s.
    - Add triton to mk8s.
    - Add digital ocean to mk8s.
- v1.1.38
    - Fix empty object issue in mk8s.
- v1.1.39
    - Add deploy role chain to AWS mk8s provider.
- v1.1.40
    - Add extraNodePolicies to aws provider.
    - Add multi to workload autoscaling.
    - Allow more empty objects and fix bugs.
    - Add validators to volume set resource.
- v1.1.41
    - Update Mk8s namespaces example.
    - Add load balancer to mk8s triton.
- v1.1.42
    - Fix issue with policies that require a GVC reference.
- v1.1.43
    - Add privateNetworkIds, metadata and tags to Triton load balancer manual attribute
    - Update dependencies
- v1.1.44
    - Fix secret resource update.
- v1.1.45
    - Replace _sentinel with placeholder_attribute to fix an issue in pulumi-go.
