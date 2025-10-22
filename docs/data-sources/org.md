---
page_title: "cpln_org Data Source - terraform-provider-cpln"
subcategory: "Org"
description: |-
  
---
# cpln_org (Data Source)

Use this data source to access details about the current [org](https://docs.controlplane.com/reference/org) targeted by the provider configuration.

## Outputs

The following attributes are exported:

- **id** (String) The unique identifier for this org.
- **cpln_id** (String) The ID, in GUID format, of the org.
- **name** (String) The name of org.
- **description** (String) Description of the org.
- **tags** (Map of String) Key-value map of resource tags.
- **self_link** (String) Full link to this resource. Can be referenced by other resources.
- **observability** (Block List, Max: 1) ([see below](#nestedblock--observability)).
- **account_id** (String) The associated account ID that was used when creating the org.
- **invitees** (Set of String) Email addresses that received invitations to join the org and were assigned to the `superusers` group.
- **session_timeout_seconds** (Int) The idle time (in seconds) after which the console UI signs out the user. Default: `900`.
- **auth_config** (Block List, Max: 1) ([see below](#nestedblock--auth_config)).
- **security** (Block List, Max: 1) ([see below](#nestedblock--security)).
- **status** (List of Object) ([see below](#nestedblock--status)).

<a id="nestedblock--observability"></a>

### `observability`

The retention period (in days) for logs, metrics, and traces. Charges apply for storage beyond the 30 day default.

Read-Only:

- **logs_retention_days** (Number) Log retention days. Default: `30`.
- **metrics_retention_days** (Number) Metrics retention days. Default: `30`.
- **traces_retention_days** (Number) Traces retention days. Default: `30`.
- **default_alert_emails** (Set of String) These emails are configured as alert recipients in Grafana when the `grafana-default-email` contact delivery type is `Email`.

<a id="nestedblock--auth_config"></a>

### `auth_config`

Configuration settings related to authentication within the org.

Read-Only:

- **domain_auto_members** (Set of String) List of domains that auto-provision users when authenticating using SAML.
- **saml_only** (Boolean) Enforces SAML-only authentication.

<a id="nestedblock--security"></a>

### `security`

Read-Only:

- **threat_detection** (Block List, Max: 1) ([see below](#nestedblock--security--threat_detection)).

<a id="nestedblock--security--threat_detection"></a>

### `security.threat_detection`

Read-Only:

- **enabled** (Boolean) Indicates whether threat detection information is forwarded.
- **minimum_severity** (String) Any threats with this severity and more severe are sent. Others are ignored. Valid values: `warning`, `error`, or `critical`.
- **syslog** (Block List, Max: 1) ([see below](#nestedblock--security--threat_detection--syslog)).

<a id="nestedblock--security--threat_detection--syslog"></a>

### `security.threat_detection.syslog`

Read-Only:

- **port** (Number) The port to send syslog messages to.
- **transport** (String) The transport-layer protocol used for syslog messages. If `tcp` is chosen, messages are sent with TLS. Default: `tcp`.
- **host** (String) The hostname to send syslog messages to.

<a id="nestedblock--status"></a>

### `status`

Status of the org.

Read-Only:

- **account_link** (String) The link of the account the org belongs to.
- **active** (Boolean) Indicates whether the org is active or not.
- **endpoint_prefix** (String)

## Example Usage

```terraform
data "cpln_org" "org" {}

output "org_summary" {
  value = {
    name                = data.cpln_org.org.name
    cpln_id             = data.cpln_org.org.cpln_id
    account_id          = data.cpln_org.org.account_id
    session_timeout_sec = data.cpln_org.org.session_timeout_seconds
    observability       = data.cpln_org.org.observability
  }
}
```
