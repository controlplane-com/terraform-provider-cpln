## 1.2.15

- Fix org resource SpecReplace removing logging and tracing configured by separate resources, add mutex synchronization to prevent race conditions between org, org_logging, and org_tracing resources.

## 1.2.14

- Improve error message on retried API requests.
- Add byok addon to mk8s resource.
- Add gcp provider to mk8s resource.
- Add `capacity_ai_update_minutes` to workload resource options.
- Add `drop_metrics` to workload resource container.
- Add workload data source.
- Add `custom_encryption` to volume set resource.
- Add missing status and origin attributes to resources and data sources.
- Add `none` load balancer option to mk8s resource triton provider.
- Fix trust policy statement format in identity resource.
- Allow autoscaling target to be defined alongside keda in workload options.

## 1.2.13

- Fix order of values in command array in workload resource.
- Deprecate 'default' endpoint_naming_format in gvc resource; add 'legacy' as replacement.
- Fix order of values in deploy_role_chain in mk8s resource.
- Fix order of values in container ports in workload resource.

## 1.2.12

- Fix noisy diffs for `dictionary_as_envs` across runs in secret resource.

## 1.2.11

- Fix serivce account key preservation.

## 1.2.10

- Add acceptAllSubdomains, certChallengeType and workloadLink to domain resource.
- Add authentication_ref to workload KEDA.
- Add secrets to GVC KEDA.
- Switch between List & Set usage when necessary.

## 1.2.9

- Add KEDA to workload and GVC resources.

## 1.2.8

- Use types.List for args instead of types.Set.
- Fix Addon config with false values.

## 1.2.7

- Unbreak breaking change done to scopes in identity resource.

## 1.2.6

- Add logging to MK8s Triton Load Balancer.
- Add registry_mirror to MK8s addons.
- Replace slice fields with types.List to handle unknown-value decoding.

## 1.2.5

- Inject a tag that represents the current provider version.
- Prevent GVC replacement unless cpln/terraformReplace tag is set to true.
- Show a warning when updating domain spec.
- Add a donwtime warning on workload rename.
- Fix endpoint_naming_format requires replace after migration to `1.2.x`.

## 1.2.4
- Move the replace plan modifier from replace_prefix to regex attribute in cpln_domain_route.
- Fix domain routes removal upon update.

## 1.2.3
- Force user_data in cpln_agent to persist it value set at creation.
- Fix nil dereference issue with empty/null environment variables in cpln_workload.

## 1.2.2
- Fix an env flatten bug that caused a panic on nil value.

## 1.2.1
- Add `cpln_mk8s_kubeconfig` resource back.

## 1.2.0
- Migrate to the Terraform Plugin Framework.
- Add docker, kubelet, kernel and events to the logs addon in `cpln_mk8s` resource.
- Add azure provider to `cpln_mk8s` resource.
- Allow for giving a full target link in `cpln_policy` resource.
- Allow for using capacity AI in cron workloads.
- Attempt to delete resource again after receiving 409.
- Do not retry on 429 if the message contains "Quota".
- Use endpoint from profile if user didn't specify an endpoint.
- Fix some bugs.

## 1.1.48
- Fix org import with empty `saml only`.
- Fix empty IP Set issue in GVC load balancer.

## 1.1.47
- Notify users to use the custom location resource for BYOK locations.
- Allow resources to have 2 character names.
- Add Custom GPU to workload resource.
- Add inbound & outbound blocked cidr to workload firewall.
- Add Request Retry Policy to workload resource.
- Add Default Alert Emails to Org resource.
- Add ipset to load balancer in GVC resource.
- Specify that an IP Set resource can target a GVC link.
- Add status401 to load balancer in GVC resource.
- Add endpoint naming format to GVC resource.
- Add host regex to domain resource.
- Update Go dependencies.

## 1.1.46
- Add MK8s Kubeconfig resource.
- Add extras to workload resource.

## 1.1.45
- Replace `_sentinel` with `placeholder_attribute` to fix an issue in pulumi-go.

## 1.1.44
- Fix secret resource update.

## 1.1.43
- Add `privateNetworkIds`, metadata and tags to Triton load balancer manual attribute.
- Update dependencies.

## 1.1.42
- Fix issue with policies that require a GVC reference.

## 1.1.41
- Update Mk8s namespaces example.
- Add load balancer to mk8s triton.

## 1.1.40
- Add extraNodePolicies to aws provider.
- Add multi to workload autoscaling.
- Allow more empty objects and fix bugs.
- Add validators to volume set resource.

## 1.1.39
- Add deploy role chain to AWS mk8s provider.

## 1.1.38
- Fix empty object issue in mk8s.

## 1.1.37
- Add cpln_ipset resource.
- Add redirect to GVC load balancer.
- Add lambdalabs to mk8s.
- Add linode to mk8s.
- Add oblivus to mk8s.
- Add paperspace to mk8s.
- Add triton to mk8s.
- Add digital ocean to mk8s.

## 1.1.36
- Add secret data source.
- Add autoscaling by memory.

## 1.1.35
- Fix empty options bug in workload resource.

## 1.1.34
- Make workload options optional.

## 1.1.33
- Add import support to mk8s resource.

## 1.1.32
- Update docs to include import syntax.
- Ignore server-side tags for location resource.
- Enable term rel in query.
- Fix bug in tracing.
- Fix plan not empty issue in cloud account.
- Fix tag format.
- Add geo locations to security options in workload.
- Add load balancer to workload resource.

## 1.1.31
- Add a custom location resource.
- Fix typo in org doc.
- Add floating ip selector to mk8s resource.
- Add external fields to cloud watching logging in org logging resource.
- Update policy resource to allow self-link references.
- Add headers to domain route resource.
- Force a tag indicating resource was created by terraform.
- Set networking as required in mk8s resource.

## 1.1.30
- Fix mk8s addons being removed on update.

## 1.1.29
- Fix mk8s addons boolean value issue.

## 1.1.28
- Fix mk8s add-on update flow.

## 1.1.27
- Add additional descriptions.
- Add sysbox, hetznerLabels, awsTags to mk8s.
- Add regex to domain-route resource.
- Add threat-detection to org.

## 1.1.26
- Add mk8s resource.
- Add location resource.
- Add syslog to org logging.

## 1.1.25
- Handle case when no images are found.
- Add cloudwatch, fluentd, and stack driver to org logging.

## 1.1.24
- Update images data source to fetch all images.

## 1.1.23
- Update docs to indicate the locations are optional for a GVC.
- Update image data source to return the latest image if there is not tag.
- Update images data source to accept a query.
- Add storage class suffix to volume set.

## 1.1.22
- Add formatted link for secret and volume set.
- Update workload autoscaling to be optional.

## 1.1.21
- Fix workload validation.
- Add default to autoscaling metric property.
- Add stateful workload example in docs.
- Update dependencies and Go to v1.21.

## 1.1.20
- Update GVC data-source and resource docs.
- Add descriptions in the schema for all data-sources and resources.
- Add image and images data-sources.

## 1.1.19
- Add external ID to AWS secret.

## 1.1.18
- Add min CPU and Memory to container.

## 1.1.17
- Add schedule to volume set.
- Add status to domain output.
- Add Control Plane tracing and custom tags to tracing.
- Add external ID to ECR secret.
- Add generic elastic to org logging.
- Add geo properties to locations.
- Add missing properties to workload.

## 1.1.16
- Add org creation.
- Add org properties.

## 1.1.15
- Fix issue with tag values that were stored as number types.

## 1.1.14
- Fix bug with secrets when updating.
- Add `trusted_proxies` to GVC load balancer.

## 1.1.13
- Update dependencies.
- Add to a volume set policy that a gvc is required.

## 1.1.12
- Fix bug when passing null container command args.

## 1.1.11
- Add Otel tracing to org and gvc.

## 1.1.10
- Add to org logging to multiple external provider.
- Update volumeset performance classes.
- Add cloud account data source that has aws identifiers as output.
- Fix issue with workload import and legacy port.

## 1.1.9
- Fix issue with volume set status.locations.
- Fix import domain route.

## 1.1.8
- Add to secret the output `dictionary_as_env`.
- Updated import docs.
- Initial logic for the deprecation of the `port` container attribute.
- Updated SDK to v2.27.0.

## 1.1.7
- Update all resources to allow import.
- Updated docs with import details and syntax.

## 1.1.6
- Added volume sets.
- Added GPU to workload.
- Added external firewall outbound allow ports.
- Add domain host prefix.

## 1.1.5
- Updated Terraform SDK to v2.26.1.
- Added CRON job workload type.
- Add workload rollout and security options.
- Add disabled scaling strategy.
- Add GVC load balancer.
- Add workload support dynamic tags.

## 1.1.4
- Added domain and domain route.

## 1.1.3
- Fixed issue with identity and workload policies.
- 'gvc' property now required for those policy kinds.

## 1.1.2
- Added NATS Account secret.
- Added NGS cloud account.
- Added NGS cloud access policy.
- Sync'ed GVC schema (env will be returned).
- Added GCP service account name and roles output to cloud account.
- Added elastic logging to org external logging.
- Added audit context resource.
- Added native network resources.

## 1.1.1
- Add identity manager to group.

## 1.1.0
- Update to Go 1.1.8 and Terraform SDK 2.25.0.
- Add workload lifecycle hooks (post start / pre stor).
- Add GVC Environment Variables.
- Add workload suspend.

## 1.0.13
- Add workload lifecycle hooks (post start / pre stor).
- Add GVC Environment Variables.
- Add workload suspend.

## 1.0.123
- Updates for -refresh-only flag.

## 1.0.122
- HotFix for new and missing workload properties.

## 1.0.121
- HotFix for removal of workload option spot property.

## 1.0.12
- Update GitHub Action.

## 1.0.11
- Fix issue with secrets having json types.
- Remove built-in server generated secret tags.

## 1.0.10
- Add Location / Locations Data Source.

## 1.0.9
- Only remove certain `cpln/*` server generated tags.
- Increase max containers.

## 1.0.8
- Add Gvc Data Source.

## 1.0.7
- Add Org/Gvc Tracing (lightstep).

## 1.0.5
- Add Org External Logging.

## 1.0.4
- Add standard workload type.

## 1.0.3
- Update docs.

## 1.0.2
- Publish To Terraform Registry.