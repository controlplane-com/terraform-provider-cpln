---
page_title: "cpln_mk8s Resource - terraform-provider-cpln"
subcategory: "Mk8s"
description: |-
---

# cpln_mk8s (Resource)

Manages an org's [Managed K8s](https://docs.controlplane.com/mk8s/overview).

## Declaration

### Required

- **name** (String) Name of the Mk8s.
- **version** (String)

~> **Note** Only one of the providers listed below can be included in a resource.

- **generic_provider** (Block List, Max: 1) ([see below](#nestedblock--generic_provider))
- **hetzner_provider** (Block List, Max: 1) ([see below](#nestedblock--hetzner_provider))
- **aws_provider** (Block List, Max: 1) ([see below](#nestedblock--aws_provider))
- **linode_provider** (Block List, Max: 1) ([see below](#nestedblock--linode_provider))
- **oblivus_provider** (Block List, Max: 1) ([see below](#nestedblock--oblivus_provider))
- **lambdalabs_provider** (Block List, Max: 1) ([see below](#nestedblock--lambdalabs_provider))
- **paperspace_provider** (Block List, Max: 1) ([see below](#nestedblock--paperspace_provider))
- **ephemeral_provider** (Block List, Max: 1) ([see below](#nestedblock--ephemeral_provider))
- **triton_provider** (Block List, Max: 1) ([see below](#nestedblock--triton_provider))
- **azure_provider** (Block List, Max: 1) ([see below](#nestedblock--azure_provider))
- **digital_ocean_provider** (Block List, Max: 1) ([see below](#nestedblock--digital_ocean_provider))
- **gcp_provider** (Block List, Max: 1) ([see below](#nestedblock--gcp_provider))

### Optional

- **description** (String) Description of the Mk8s.
- **tags** (Map of String) Key-value map of resource tags.
- **firewall** (Block List, Max: 1) ([see below](#nestedblock--firewall))
- **add_ons** (Block List, Max: 1) ([see below](#nestedblock--add_ons))

<a id="nestedblock--generic_provider"></a>

### `generic_provider`

Required:

- **location** (String) Control Plane location that will host the K8s components. Prefer one that is closest to where the nodes are running.
- **networking** (Block List, Max: 1) ([see below](#nestedblock--generic_provider--networking))

Optional:

- **node_pool** (Block List) ([see below](#nestedblock--generic_provider--node_pool))

<a id="nestedblock--generic_provider--networking"></a>

### `generic_provider.networking`

Networking declaration is required even if networking is not utilized. Example usage: `networking {}`.

Optional:

- **service_network** (String) The CIDR of the service network.
- **pod_network** (String) The CIDR of the pod network.
- **dns_forwarder** (String) DNS forwarder used by the cluster. Can be a space-delimited list of dns servers. Default is /etc/resolv.conf when not specified.

<a id="nestedblock--generic_provider--node_pool"></a>

### `generic_provider.node_pool`

List of node pools.

Required:

- **name** (String)

Optional:

- **labels** (Map of String) Labels to attach to nodes of a node pool.
- **taint** (Block List) ([see below](#nestedblock--generic_provider--node_pool--taint))

<a id="nestedblock--generic_provider--node_pool--taint"></a>

### `generic_provider.node_pool.taint`

Taint for the nodes of a pool.

Optional:

- **key** (String)
- **value** (String)
- **effect** (String)

<a id="nestedblock--hetzner_provider"></a>

### `hetzner_provider`

Required:

- **region** (String) Hetzner region to deploy nodes to.
- **token_secret_link** (String) Link to a secret holding Hetzner access key.
- **network_id** (String) ID of the Hetzner network to deploy nodes to.
- **networking** (Block List, Max: 1) ([see below](#nestedblock--generic_provider--networking))

Optional:

- **hetzner_labels** (Map of String) Extra labels to attach to servers.
- **pre_install_script** (String) Optional shell script that will be run before K8s is installed.
- **firewall_id** (String) Optional firewall rule to attach to all nodes.
- **node_pool** (Block List) ([see below](#nestedblock--hetzner_provider--node_pool))
- **dedicated_server_node_pool** (Block List) ([see below](#nestedblock--hetzner_provider--dedicated_server_node_pool))
- **image** (String) Default image for all nodes.
- **ssh_key** (String) SSH key name for accessing deployed nodes.
- **autoscaler** (Block List, Max: 1) ([see below](#nestedblock--autoscaler))
- **floating_ip_selector** (Map of String) If supplied, nodes will get assigned a random floating ip matching the selector.

<a id="nestedblock--hetzner_provider--node_pool"></a>

### `hetzner_provider.node_pool`

List of node pools.

Required:

- **name** (String)
- **server_type** (String)

Optional:

- **labels** (Map of String) Labels to attach to nodes of a node pool.
- **taint** (Block List) ([see below](#nestedblock--generic_provider--node_pool--taint))
- **override_image** (String)
- **min_size** (Number)
- **max_size** (Number)

<a id="nestedblock--hetzner_provider--dedicated_server_node_pool"></a>

### `hetzner_provider.dedicated_server_node_pool`

Node pool that can configure dedicated Hetzner servers.

Required:

- **name** (String)

Optional:

- **labels** (Map of String) Labels to attach to nodes of a node pool.
- **taint** (Block List) ([see below](#nestedblock--generic_provider--node_pool--taint))

<a id="nestedblock--aws_provider"></a>

### `aws_provider`

Required:

- **region** (String) Region where the cluster nodes will live.
- **skip_create_roles** (Boolean) If true, Control Plane will not create any roles.
- **image** (Block List, Max: 1) ([see below](#nestedblock--aws_provider--ami))
- **deploy_role_arn** (String) Control Plane will set up the cluster by assuming this role.
- **vpc_id** (String) The vpc where nodes will be deployed. Supports SSM.
- **networking** (Block List, Max: 1) ([see below](#nestedblock--generic_provider--networking))

Optional:

- **deploy_role_chain** (Block List) ([see below](#nestedblock--aws_provider--deploy_role_chain))
- **aws_tags** (Map of String) Extra tags to attach to all created objects.
- **pre_install_script** (String) Optional shell script that will be run before K8s is installed. Supports SSM.
- **key_pair** (String) Name of keyPair. Supports SSM
- **disk_encryption_key_arn** (String) KMS key used to encrypt volumes. Supports SSM.
- **security_group_ids** (List of String) Security groups to deploy nodes to. Security groups control if the cluster is multi-zone or single-zon.
- **extra_node_policies** (List of String)
- **node_pool** (Block List) ([see below](#nestedblock--aws_provider--node_pool))
- **autoscaler** (Block List, Max: 1) ([see below](#nestedblock--autoscaler))

<a id="nestedblock--aws_provider--node_pool"></a>

### `aws_provider.node_pool`

List of node pools.

Required:

- **name** (String)
- **instance_types** (List of String)
- **override_image** (Block List, Max: 1) ([see below](#nestedblock--aws_provider--ami))
- **subnet_ids** (List of String)

Optional:

- **labels** (Map of String) Labels to attach to nodes of a node pool.
- **taint** (Block List) ([see below](#nestedblock--generic_provider--node_pool--taint))
- **boot_disk_size** (Number) Size in GB.
- **min_size** (Number)
- **max_size** (Number)
- **on_demand_base_capacity** (Number)
- **on_demand_percentage_above_base_capacity** (Number)
- **spot_allocation_strategy** (String)
- **extra_security_group_ids** (List of String)

<a id="nestedblock--aws_provider--ami"></a>

### `ami`

Default image for all nodes.

Required:

~> **Note** Only one of the following listed below can be included.

- **recommended** (String)
- **exact** (String) Support SSM.

<a id="nestedblock--aws_provider--deploy_role_chain"></a>

### `aws_provider.deploy_role_chain`

Required:

- **role_arn** (String)

Optional:

- **external_id** (String)
- **session_name_prefix** (String) Control Plane will append random.

<a id="nestedblock--linode_provider"></a>

### `linode_provider`

Required:

- **region** (String) Region where the cluster nodes will live.
- **token_secret_link** (String) Link to a secret holding Linode access key.
- **image** (String) Default image for all nodes.
- **vpc_id** (String) The vpc where nodes will be deployed. Supports SSM.

Optional:

- **firewall_id** (String) Optional firewall rule to attach to all nodes.
- **authorized_users** (List of String)
- **authorized_keys** (List of String)
- **pre_install_script** (String) Optional shell script that will be run before K8s is installed. Supports SSM.
- **networking** (Block List, Max: 1) ([see below](#nestedblock--generic_provider--networking))
- **node_pool** (Block List) ([see below](#nestedblock--linode_provider--node_pool))
- **autoscaler** (Block List, Max: 1) ([see below](#nestedblock--autoscaler))

<a id="nestedblock--linode_provider--node_pool"></a>

### `linode_provider.node_pool`

List of node pools.

Required:

- **name** (String)
- **subnet_id** (String)
- **server_type** (String)

Optional:

- **labels** (Map of String) Labels to attach to nodes of a node pool.
- **taint** (Block List) ([see below](#nestedblock--generic_provider--node_pool--taint))
- **override_image** (String)
- **min_size** (Number)
- **max_size** (Number)

<a id="nestedblock--oblivus_provider"></a>

### `oblivus_provider`

Required:

- **datacenter** (String)
- **token_secret_link** (String) Link to a secret holding Oblivus access key.

Optional:

- **node_pool** (Block List) ([see below](#nestedblock--oblivus_provider--node_pool))
- **ssh_keys** (List of String)
- **unmanaged_node_pool** (Block List) ([see below](#nestedblock--oblivus_provider--unmanaged_node_pool))
- **autoscaler** (Block List, Max: 1) ([see below](#nestedblock--autoscaler))
- **pre_install_script** (String) Optional shell script that will be run before K8s is installed. Supports SSM.

<a id="nestedblock--oblivus_provider--node_pool"></a>

### `oblivus_provider.node_pool`

List of node pools.

Required:

- **name** (String)

Optional:

- **labels** (Map of String) Labels to attach to nodes of a node pool.
- **taint** (Block List) ([see below](#nestedblock--generic_provider--node_pool--taint))
- **min_size** (Number)
- **max_size** (Number)
- **flavor** (String)

<a id="nestedblock--oblivus_provider--unmanaged_node_pool"></a>

### `oblivus_provider.unmanaged_node_pool`

Required:

- **name** (String)

Optional:

- **labels** (Map of String) Labels to attach to nodes of a node pool.
- **taint** (Block List) ([see below](#nestedblock--generic_provider--node_pool--taint))

<a id="nestedblock--lambdalabs_provider"></a>

### `lambdalabs_provider`

Required:

- **region** (String) Region where the cluster nodes will live.
- **token_secret_link** (String) Link to a secret holding Lambdalabs access key.
- **ssh_key** (String) SSH key name for accessing deployed nodes.

Optional:

- **node_pool** (Block List) ([see below](#nestedblock--lambdalabs_provider--node_pool))
- **unmanaged_node_pool** (Block List) ([see below](#nestedblock--lambdalabs_provider--unmanaged_node_pool))
- **autoscaler** (Block List, Max: 1) ([see below](#nestedblock--autoscaler))
- **file_systems** (List of String)
- **pre_install_script** (String) Optional shell script that will be run before K8s is installed. Supports SSM.

<a id="nestedblock--lambdalabs_provider--node_pool"></a>

### `lambdalabs_provider.node_pool`

List of node pools.

Required:

- **name** (String)
- **instance_type** (String)

Optional:

- **labels** (Map of String) Labels to attach to nodes of a node pool.
- **taint** (Block List) ([see below](#nestedblock--generic_provider--node_pool--taint))
- **min_size** (Number)
- **max_size** (Number)

<a id="nestedblock--lambdalabs_provider--unmanaged_node_pool"></a>

### `lambdalabs_provider.unmanaged_node_pool`

Required:

- **name** (String)

Optional:

- **labels** (Map of String) Labels to attach to nodes of a node pool.
- **taint** (Block List) ([see below](#nestedblock--generic_provider--node_pool--taint))

<a id="nestedblock--paperspace_provider"></a>

### `paperspace_provider`

Required:

- **region** (String) Region where the cluster nodes will live.
- **token_secret_link** (String) Link to a secret holding Paperspace access key.
- **network_id** (String)

Optional:

- **shared_drives** (List of String)
- **node_pool** (Block List) ([see below](#nestedblock--paperspace_provider--node_pool))
- **autoscaler** (Block List, Max: 1) ([see below](#nestedblock--autoscaler))
- **unmanaged_node_pool** (Block List) ([see below](#nestedblock--paperspace_provider--unmanaged_node_pool))
- **pre_install_script** (String) Optional shell script that will be run before K8s is installed. Supports SSM.
- **user_ids** (List of String)

<a id="nestedblock--paperspace_provider--node_pool"></a>

### `paperspace_provider.node_pool`

List of node pools.

Required:

- **name** (String)
- **public_ip_type** (String)
- **machine_type** (String)

Optional:

- **labels** (Map of String) Labels to attach to nodes of a node pool.
- **taint** (Block List) ([see below](#nestedblock--generic_provider--node_pool--taint))
- **min_size** (Number)
- **max_size** (Number)
- **boot_disk_size** (Number)

<a id="nestedblock--paperspace_provider--unmanaged_node_pool"></a>

### `paperspace_provider.unmanaged_node_pool`

Required:

- **name** (String)

Optional:

- **labels** (Map of String) Labels to attach to nodes of a node pool.
- **taint** (Block List) ([see below](#nestedblock--generic_provider--node_pool--taint))

<a id="nestedblock--ephemeral_provider"></a>

### `ephemeral_provider`

Required:

- **location** (String) Control Plane location that will host the K8s components. Prefer one that is closest to where the nodes are running.

Optional:

- **node_pool** (Block List) ([see below](#nestedblock--ephemeral_provider--node_pool))

<a id="nestedblock--ephemeral_provider--node_pool"></a>

### `ephemeral_provider.node_pool`

List of node pools.

Required:

- **name** (String)
- **count** (Int) Number of nodes to deploy.
- **arch** (String) CPU architecture of the nodes.
- **flavor** (String) Linux distro to use for ephemeral nodes.
- **cpu** (String) Allocated CPU.
- **memory** (String) Allocated memory.

Optional:

- **labels** (Map of String) Labels to attach to nodes of a node pool.
- **taint** (Block List) ([see below](#nestedblock--generic_provider--node_pool--taint))

<a id="nestedblock--triton_provider"></a>

### `triton_provider`

Required:

- **connection** (Block List, Max: 1) ([see below](#nestedblock--triton_provider--connection))
- **networking** (Block List, Max: 1) ([see below](#nestedblock--generic_provider--networking))
- **location** (String) Control Plane location that will host the K8s components. Prefer one that is closest to the Triton datacenter.
- **load_balancer** (Block List, Max: 1) ([see below](#nestedblock--triton_provider--load_balancer))
- **private_network_id** (String) ID of the private Fabric/Network.
- **image_id** (String) Default image for all nodes.

Optional:

- **pre_install_script** (String) Optional shell script that will be run before K8s is installed. Supports SSM.
- **firewall_enabled** (Boolean) Enable firewall for the instances deployed.
- **node_pool** (Block List) ([see below](#nestedblock--triton_provider--node_pool))
- **ssh_keys** (List of String) Extra SSH keys to provision for user root.
- **autoscaler** (Block List, Max: 1) ([see below](#nestedblock--autoscaler))

<a id="nestedblock--triton_provider--connection"></a>

### `triton_provider.connection`

Required:

- **url** (String)
- **account** (String)
- **user** (String)
- **private_key_secret_link** (String) Link to a SSH or opaque secret.

Optional:

- **user** (String)

<a id="nestedblock--triton_provider--load_balancer"></a>

### `triton_provider.load_balancer`

Required:

~> **Note** Only one of the attributes listed below can be included in load balancer.

- **manual** (Block List, Max: 1) ([see below](#nestedblock--triton_provider--load_balancer--manual))
- **none** (Block List, Max: 1) Just an empty list. E.g. `none {}`.
- **gateway** (Block List, Max: 1) Just an empty list. E.g. `gateway {}`.

<a id="nestedblock--triton_provider--load_balancer--manual"></a>

### `triton_provider.load_balancer.manual`

Required:

- **package_id** (String)
- **image_id** (String)
- **public_network_id** (String) If set, machine will also get a public IP.
- **private_network_ids** (List of String) More private networks to join.
- **metadata** (Map of String) Extra tags to attach to instances from a node pool.
- **tags** (Map of String) Extra tags to attach to instances from a node pool.
- **logging** (Block List, Max: 1) ([see below](#nestedblock--triton_provider--load_balancer--manual--logging))
- **count** (Number)
- **cns_internal_domain** (String)
- **cns_public_domain** (String)

<a id="nestedblock--triton_provider--load_balancer--manual--logging"></a>

### `triton_provider.load_balancer.manual.logging`

Optional:

- **node_port** (Number)
- **external_syslog** (String)

<a id="nestedblock--triton_provider--node_pool"></a>

### `triton_provider.node_pool`

List of node pools.

Required:

- **name** (String)
- **package_id** (String)

Optional:

- **labels** (Map of String) Labels to attach to nodes of a node pool.
- **taint** (Block List) ([see below](#nestedblock--generic_provider--node_pool--taint))
- **override_image_id** (String)
- **public_network_id** (String) If set, machine will also get a public IP.
- **private_network_ids** (List of String) More private networks to join.
- **triton_tags** (Map of String) Extra tags to attach to instances from a node pool.
- **min_size** (Number)
- **max_size** (Number)

<a id="nestedblock--azure_provider"></a>

### `azure_provider`

Required:

- **location** (String) Region where the cluster nodes will live.
- **subscription_id** (String)
- **sdk_secret_link** (String)
- **resource_group** (String)
- **ssh_keys** (List of String) SSH keys to install for "azureuser" linux user.
- **network_id** (String) The vpc where nodes will be deployed.

Optional:

- **networking** (Block List, Max: 1) ([see below](#nestedblock--generic_provider--networking))
- **pre_install_script** (String) Optional shell script that will be run before K8s is installed. Supports SSM.
- **image** (Block List, Max: 1) ([see below](#nestedblock--azure_provider--image))
- **tags** (Map of String) Extra tags to attach to all created objects.
- **node_pool** (Block List) ([see below](#nestedblock--azure_provider--node_pool))
- **autoscaler** (Block List, Max: 1) ([see below](#nestedblock--autoscaler))

<a id="nestedblock--azure_provider--image"></a>

### `azure_provider.image`

Default image for all nodes.

Required:

~> **Note** Only one of the following listed below can be included.

- **recommended** (String)
- **reference** (Block List, Max: 1) ([see below](#nestedblock--azure_provider--image--reference))

<a id="nestedblock--azure_provider--image--reference"></a>

### `azure_provider.image.reference`

Required:

- **publisher** (String)
- **offer** (String)
- **sku** (String)
- **version** (String)

<a id="nestedblock--azure_provider--node_pool"></a>

### `azure_provider.node_pool`

List of node pools.

Required:

- **name** (String)
- **size** (String)
- **subnet_id** (String)
- **zones** (List of Number)
- **boot_disk_size** (Number)

Optional:

- **labels** (Map of String) Labels to attach to nodes of a node pool.
- **taint** (Block List) ([see below](#nestedblock--generic_provider--node_pool--taint))
- **override_image** (Block List, Max: 1) ([see below](#nestedblock--azure_provider--image))
- **min_size** (Number) Default: 0
- **max_size** (Number) Default: 0

<a id="nestedblock--gcp_provider"></a>

### `gcp_provider`

Required:

- **project_id** (String) GCP project ID that hosts the cluster infrastructure.
- **region** (String) Region where the cluster nodes will live.
- **network** (String) VPC network used by the cluster.
- **sa_key_link** (String) Link to a secret containing the service account JSON key.
- **networking** (Block List, Max: 1) ([see below](#nestedblock--generic_provider--networking))

Optional:

- **labels** (Map of String) Extra labels to attach to all created objects. Maximum: `10`.
- **tags** (List of String)
- **metadata** (Map of String)
- **pre_install_script** (String) Optional shell script that runs before K8s is installed.
- **image** (Block List, Max: 1) ([see below](#nestedblock--gcp_provider--image))
- **node_pool** (Block List) ([see below](#nestedblock--gcp_provider--node_pool))
- **autoscaler** (Block List, Max: 1) ([see below](#nestedblock--autoscaler))

<a id="nestedblock--gcp_provider--image"></a>

### `gcp_provider.image`

Default image for all nodes.

Optional:

- **recommended** (String) Recommended image alias. Valid values: `ubuntu/jammy-22.04`, `ubuntu/noble-24.04`, `debian/bookworm-12`, `debian/trixie-13`, `google/cos-stable`.
- **family** (Object) ([see below](#nestedblock--gcp_provider--image--family))
- **exact** (String)

<a id="nestedblock--gcp_provider--image--family"></a>

### `gcp_provider.image.family`

Required:

- **project** (String)
- **family** (String)

<a id="nestedblock--gcp_provider--node_pool"></a>

### `gcp_provider.node_pool`

List of node pools.

Required:

- **name** (String)
- **machine_type** (String) GCE machine type for nodes in this pool.
- **zone** (String) Zone where the pool nodes run.
- **boot_disk_size** (Number) Size in GB. Minimum: `20`.
- **subnet** (String) Subnet within the selected network.

Optional:

- **labels** (Map of String) Labels to attach to nodes of a node pool.
- **taint** (Block List) ([see below](#nestedblock--generic_provider--node_pool--taint))
- **assign_public_ip** (Boolean)
- **override_image** (Block List, Max: 1) ([see below](#nestedblock--gcp_provider--image))
- **min_size** (Number)
- **max_size** (Number)
- **preemptible** (Boolean)
- **local_persistent_disks** (Number)

<a id="nestedblock--digital_ocean_provider"></a>

### `digital_ocean_provider`

Required:

- **region** (String) Region to deploy nodes to.
- **networking** (Block List, Max: 1) ([see below](#nestedblock--generic_provider--networking))
- **token_secret_link** (String) Link to a secret holding personal access token.
- **vpc_id** (String) ID of the Hetzner network to deploy nodes to.
- **image** (String) Default image for all nodes.
- **ssh_keys** (List of String) SSH key name for accessing deployed nodes.

Optional:

- **digital_ocean_tags** (List of String) Extra tags to attach to droplets.
- **pre_install_script** (String) Optional shell script that will be run before K8s is installed. Supports SSM.
- **node_pool** (Block List) ([see below](#nestedblock--digital_ocean_provider--node_pool))
- **extra_ssh_keys** (List of String) Extra SSH keys to provision for user root that are not registered in the DigitalOcean.
- **autoscaler** (Block List, Max: 1) ([see below](#nestedblock--autoscaler))
- **reserved_ips** (List of String) Optional set of IPs to assign as extra IPs for nodes of the cluster.

<a id="nestedblock--digital_ocean_provider--node_pool"></a>

### `digital_ocean_provider.node_pool`

List of node pools.

Required:

- **name** (String)
- **droplet_size** (String)

Optional:

- **labels** (Map of String) Labels to attach to nodes of a node pool.
- **taint** (Block List) ([see below](#nestedblock--generic_provider--node_pool--taint))
- **override_image** (String)
- **min_size** (Number)
- **max_size** (Number)

<a id="nestedblock--autoscaler"></a>

### `autoscaler`

Optional:

- **expander** (List of String)
- **unneeded_time** (String)
- **unready_time** (String)
- **utilization_threshold** (Float64)

<a id="nestedblock--firewall"></a>

### `firewall`

Allow-list.

Required:

- **source_cidr** (String)

Optional:

- **description** (String)

<a id="nestedblock--add_ons"></a>

### `add_ons`

Optional:

- **dashboard** (Boolean)
- **azure_workload_identity** (Block List, Max: 1) ([see below](#nestedblock--add_ons--azure_workload_identity))
- **aws_workload_identity** (Boolean)
- **local_path_storage** (Boolean)
- **metrics** (Block List, Max: 1) ([see below](#nestedblock--add_ons--metrics))
- **logs** (Block List, Max: 1) ([see below](#nestedblock--add_ons--logs))
- **registry_mirror** (Block List, Max: 1) ([see below](#nestedblock--add_ons--registry_mirror))
- **nvidia** (Block List, Max: 1) ([see below](#nestedblock--add_ons--nvidia))
- **aws_efs** (Block List, Max: 1) ([see below](#nestedblock--add_ons--aws--efs))
- **aws_ecr** (Block List, Max: 1) ([see below](#nestedblock--add_ons--aws--ecr))
- **aws_elb** (Block List, Max: 1) ([see below](#nestedblock--add_ons--aws--elb))
- **azure_acr** (Block List, Max: 1) ([see below](#nestedblock--add_ons--azure_acr))
- **byok** (Object) ([see below](#nestedblock--add_ons--byok))
- **sysbox** (Boolean)

<a id="nestedblock--add_ons--azure_workload_identity"></a>

### `add_ons.azure_workload_identity`

Required:

- **tenant_id** (String) Tenant ID to use for workload identity.

<a id="nestedblock--add_ons--metrics"></a>

### `add_ons.metrics`

Optional:

- **kube_state** (Boolean) Enable kube-state metrics.
- **core_dns** (Boolean) Enable scraping of core-dns service.
- **kubelet** (Boolean) Enable scraping kubelet stats.
- **api_server** (Boolean) Enable scraping apiserver stats.
- **node_exporter** (Boolean) Enable collecting node-level stats (disk, network, filesystem, etc).
- **cadvisor** (Boolean) Enable CNI-level container stats.
- **scrape_annotated** (Block List, Max: 1) ([see below](#nestedblock--add_ons--metrics--scrape_annotated))

<a id="nestedblock--add_ons--metrics--scrape_annotated"></a>

### `add_ons.metrics.scrape_annotated`

Scrape pods annotated with prometheus.io/scrape=true.

Optional:

- **interval_seconds** (Number)
- **include_namespaces** (String)
- **exclude_namespaces** (String)
- **retain_labels** (String)

<a id="nestedblock--add_ons--logs"></a>

### `add_ons.logs`

Optional:

- **audit_enabled** (Boolean) Collect K8s audit log as log events.
- **include_namespaces** (String)
- **exclude_namespaces** (String)
- **docker** (Boolean)
- **kubelet** (Boolean)
- **kernel** (Boolean)
- **events** (Boolean)

<a id="nestedblock--add_ons--registry_mirror"></a>

### `add_ons.registry_mirror`

Optional:

- **mirror** (Block List) ([see below](#nestedblock--add_ons--registry_mirror--mirror))

<a id="nestedblock--add_ons--registry_mirror--mirror"></a>

### `add_ons.registry_mirror.mirror`

Required:

- **registry** (String)

Optional:

- **mirrors** (List of String)

<a id="nestedblock--add_ons--nvidia"></a>

### `add_ons.nvidia`

Required:

- **taint_gpu_nodes** (Boolean)

<a id="nestedblock--add_ons--aws--efs"></a>

### `add_ons.aws_efs`

Required:

- **role_arn** (String) Use this role for EFS interaction.

<a id="nestedblock--add_ons--aws--ecr"></a>

### `add_ons.aws_ecr`

Required:

- **role_arn** (String) Role to use when authorizing ECR pulls. Optional on AWS, in which case it will use the instance role to pull.

<a id="nestedblock--add_ons--aws--elb"></a>

### `add_ons.aws_elb`

Required:

- **role_arn** (String) Role to use when authorizing calls to EC2 ELB. Optional on AWS, when not provided it will create the recommended role.

<a id="nestedblock--add_ons--azure_acr"></a>

### `add_ons.azure_acr`

Required:

- **client_id** (String)

<a id="nestedblock--add_ons--byok"></a>

### `add_ons.byok`

Bring-your-own Kubernetes (BYOK) add-on settings.

Required:

- **location** (String) The full link of an existing BYOK location.

Optional:

- **ignore_updates** (Boolean) Disable Control Plane managed upgrades for BYOK components.
- **config** (Object) ([see below](#nestedblock--add_ons--byok--config))

<a id="nestedblock--add_ons--byok--config"></a>

Fine-grained configuration for the BYOK workloads.

### `add_ons.byok.config`

Optional:

- **actuator** (Object) ([see below](#nestedblock--add_ons--byok--config--actuator))
- **middlebox** (Object) ([see below](#nestedblock--add_ons--byok--config--middlebox))
- **common** (Object) ([see below](#nestedblock--add_ons--byok--config--common))
- **longhorn** (Object) ([see below](#nestedblock--add_ons--byok--config--longhorn))
- **ingress** (Object) ([see below](#nestedblock--add_ons--byok--config--ingress))
- **istio** (Object) ([see below](#nestedblock--add_ons--byok--config--istio))
- **log_splitter** (Object) ([see below](#nestedblock--add_ons--byok--config--log_splitter))
- **monitoring** (Object) ([see below](#nestedblock--add_ons--byok--config--monitoring))
- **redis** (Object) ([see below](#nestedblock--add_ons--byok--config--redis))
- **redis_ha** (Object) ([see below](#nestedblock--add_ons--byok--config--redis_ha))
- **redis_sentinel** (Object) ([see below](#nestedblock--add_ons--byok--config--redis_sentinel))
- **tempo_agent** (Object) ([see below](#nestedblock--add_ons--byok--config--tempo_agent))
- **internal_dns** (Object) ([see below](#nestedblock--add_ons--byok--config--internal_dns))

<a id="nestedblock--add_ons--byok--config--actuator"></a>

### `add_ons.byok.config.actuator`

Resource tuning for the actuator component.

Optional:

- **min_cpu** (String) Minimum CPU request applied to actuator pods (for example, "100m").
- **max_cpu** (String) CPU limit applied to actuator pods.
- **min_memory** (String) Minimum memory request applied to actuator pods (for example, "128Mi").
- **max_memory** (String) Memory limit applied to actuator pods.
- **log_level** (String) Log level override for actuator containers. Valid values are: `trace`, `info`, `error`.
- **env** (Map of String) Additional environment variables injected into actuator pods.

<a id="nestedblock--add_ons--byok--config--middlebox"></a>

### `add_ons.byok.config.middlebox`

Configuration for the optional middlebox traffic shaper.

Optional:

- **enabled** (Boolean) Whether to deploy the middlebox component.
- **bandwidth_alert_mbps** (Number) Alert threshold, in Mbps, for middlebox bandwidth usage.

<a id="nestedblock--add_ons--byok--config--common"></a>

### `add_ons.byok.config.common`

Shared rollout settings for BYOK workloads.

Optional:

- **deployment_replicas** (Number) Replica count shared by BYOK control plane deployments.
- **pdb** (Object) ([see below](#nestedblock--add_ons--byok--config--common--pdb))

<a id="nestedblock--add_ons--byok--config--common--pdb"></a>

### `add_ons.byok.config.common.pdb`

Pod disruption budget limits for BYOK workloads.

Optional:

- **max_unavailable** (Number) Maximum number of pods that can be unavailable during disruptions.

<a id="nestedblock--add_ons--byok--config--longhorn"></a>

### `add_ons.byok.config.longhorn`

Longhorn persistent volume settings.

Optional:

- **replicas** (Number) Replica factor for Longhorn volumes. Minimum: `1`.

<a id="nestedblock--add_ons--byok--config--ingress"></a>

### `add_ons.byok.config.ingress`

Ingress controller resource configuration.

Optional:

- **cpu** (String) CPU request/limit string applied to ingress pods.
- **memory** (String) Memory request/limit string applied to ingress pods.
- **target_percent** (Number) Target usage percentage that triggers ingress autoscaling.

<a id="nestedblock--add_ons--byok--config--istio"></a>

### `add_ons.byok.config.istio`

Istio service mesh configuration.

Optional:

- **istiod** (Object) ([see below](#nestedblock--add_ons--byok--config--istio--istiod))
- **ingress_gateway** (Object) ([see below](#nestedblock--add_ons--byok--config--istio--ingress_gateway))
- **sidecar** (Object) ([see below](#nestedblock--add_ons--byok--config--istio--sidecar))

<a id="nestedblock--add_ons--byok--config--istio--istiod"></a>

### `add_ons.byok.config.istio.istiod`

Control plane deployment settings for istiod.

Optional:

- **replicas** (Number) Number of istiod replicas.
- **min_cpu** (String) CPU request applied to istiod pods.
- **max_cpu** (String) CPU limit applied to istiod pods.
- **min_memory** (String) Memory request applied to istiod pods.
- **max_memory** (String) Memory limit applied to istiod pods.
- **pdb** (Number) Pod disruption budget `max_unavailable` for istiod.

<a id="nestedblock--add_ons--byok--config--istio--ingress_gateway"></a>

### `add_ons.byok.config.istio.ingress_gateway`

Istio ingress gateway deployment settings.

Optional:

- **replicas** (Number) Number of ingress gateway replicas.
- **max_cpu** (String) CPU limit applied to ingress gateway pods.
- **max_memory** (String) Memory limit applied to ingress gateway pods.

<a id="nestedblock--add_ons--byok--config--istio--sidecar"></a>

### `add_ons.byok.config.istio.sidecar`

Default resource requests for Istio sidecar injection.

Optional:

- **min_cpu** (String) CPU request applied to injected sidecars.
- **min_memory** (String) Memory request applied to injected sidecars.

<a id="nestedblock--add_ons--byok--config--log_splitter"></a>

### `add_ons.byok.config.log_splitter`

Log splitter deployment configuration.

Optional:

- **min_cpu** (String) CPU request applied to log splitter pods.
- **max_cpu** (String) CPU limit applied to log splitter pods.
- **min_memory** (String) Memory request applied to log splitter pods.
- **max_memory** (String) Memory limit applied to log splitter pods.
- **mem_buffer_size** (String) In-memory buffer size consumed by each log splitter pod.
- **per_pod_rate** (Number) Per-pod log processing rate limit.

<a id="nestedblock--add_ons--byok--config--monitoring"></a>

### `add_ons.byok.config.monitoring`

Monitoring stack configuration.

Optional:

- **min_memory** (String) Minimum memory request for monitoring components.
- **max_memory** (String) Maximum memory limit for monitoring components.
- **kube_state_metrics** (Object) ([see below](#nestedblock--add_ons--byok--config--monitoring--kube_state_metrics))
- **prometheus** (Object) ([see below](#nestedblock--add_ons--byok--config--monitoring--prometheus))

<a id="nestedblock--add_ons--byok--config--monitoring--kube_state_metrics"></a>

### `add_ons.byok.config.monitoring.kube_state_metrics`

Kube-state-metrics resource overrides.

Optional:

- **min_memory** (String) Memory request applied to kube-state-metrics pods.

<a id="nestedblock--add_ons--byok--config--monitoring--prometheus"></a>

### `add_ons.byok.config.monitoring.prometheus`

Prometheus deployment configuration.

Optional:

- **main** (Object) ([see below](#nestedblock--add_ons--byok--config--monitoring--prometheus--main))

<a id="nestedblock--add_ons--byok--config--monitoring--prometheus--main"></a>

### `add_ons.byok.config.monitoring.prometheus.main`

Primary Prometheus instance settings.

Optional:

- **storage** (String) Persistent volume size for Prometheus (for example, "50Gi").

<a id="nestedblock--add_ons--byok--config--redis"></a>

### `add_ons.byok.config.redis`

Redis cache configuration.

Optional:

- **min_cpu** (String) CPU request applied to the Redis pods.
- **max_cpu** (String) CPU limit applied to the Redis pods.
- **min_memory** (String) Memory request applied to the Redis pods.
- **max_memory** (String) Memory limit applied to the Redis pods.
- **storage** (String) Persistent storage size allocated to the Redis pods (for example, "8Gi").

<a id="nestedblock--add_ons--byok--config--redis_ha"></a>

### `add_ons.byok.config.redis_ha`

High-availability Redis configuration.

Optional:

- **min_cpu** (String) CPU request applied to the Redis pods.
- **max_cpu** (String) CPU limit applied to the Redis pods.
- **min_memory** (String) Memory request applied to the Redis pods.
- **max_memory** (String) Memory limit applied to the Redis pods.
- **storage** (Number) Persistent storage size allocated to the Redis pods, in GiB.

<a id="nestedblock--add_ons--byok--config--redis_sentinel"></a>

### `add_ons.byok.config.redis_sentinel`

Redis Sentinel configuration.

Optional:

- **min_cpu** (String) CPU request applied to the Redis pods.
- **max_cpu** (String) CPU limit applied to the Redis pods.
- **min_memory** (String) Memory request applied to the Redis pods.
- **max_memory** (String) Memory limit applied to the Redis pods.
- **storage** (Number) Persistent storage size allocated to the Redis pods, in GiB.

<a id="nestedblock--add_ons--byok--config--tempo_agent"></a>

### `add_ons.byok.config.tempo_agent`

Tempo agent resource configuration.

Optional:

- **min_cpu** (String) CPU request applied to tempo agent pods.
- **min_memory** (String) Memory request applied to tempo agent pods.

<a id="nestedblock--add_ons--byok--config--internal_dns"></a>

### `add_ons.byok.config.internal_dns`

Internal DNS deployment settings.

Optional:

- **min_cpu** (String) CPU request applied to internal DNS pods.
- **max_cpu** (String) CPU limit applied to internal DNS pods.
- **min_memory** (String) Memory request applied to internal DNS pods.
- **max_memory** (String) Memory limit applied to internal DNS pods.

## Outputs

The following attributes are exported:

- **cpln_id** (String) The ID, in GUID format, of the Mk8s.
- **alias** (String) The alias name of the Mk8s.
- **self_link** (String) Full link to this resource. Can be referenced by other resources.
- **status** (Block List, Max: 1) ([see below](#nestedblock--status)).

<a id="nestedblock--status"></a>

### `status`

Status of the mk8s.

Read-Only:

- **oidc_provider_url** (String)
- **server_url** (String)
- **home_location** (String)
- **add_ons** (Block List, Max: 1) ([see below](#nestedblock--status--add_ons))

<a id="nestedblock--status--add_ons"></a>

### `status.add_ons`

Read-Only:

- **dashboard** (Block List, Max: 1) ([see below](#nestedblock--status--add_ons--dashobard))
- **headlamp** (Block List, Max: 1) ([see below](#nestedblock--status--add_ons--dashobard))
- **aws_workload_identity** (Block List, Max: 1) ([see below](#nestedblock--status--add_ons--aws_workload_identity))
- **metrics** (Block List, Max: 1) ([see below](#nestedblock--status--add_ons--metrics))
- **logs** (Block List, Max: 1) ([see below](#nestedblock--status--add_ons--logs))
- **aws_ecr** (Block List, Max: 1) ([see below](#nestedblock--status--add_ons--aws))
- **aws_efs** (Block List, Max: 1) ([see below](#nestedblock--status--add_ons--aws))
- **aws_elb** (Block List, Max: 1) ([see below](#nestedblock--status--add_ons--aws))

<a id="nestedblock--status--add_ons--dashobard"></a>

### `status.add_ons.dashboard`

Read-Only:

- **url** (String) Access to dashboard.

<a id="nestedblock--status--add_ons--aws_workload_identity"></a>

### `status.add_ons.aws_workload_identity`

Read-Only:

- **oidc_provider_config** (Block List, Max: 1) ([see below](#nestedblock--status--add_ons--aws_workload_identity--oidc_provider_config))
- **trust_policy** (String)

<a id="nestedblock--status--add_ons--aws_workload_identity--oidc_provider_config"></a>

### `status.add_ons.aws_workload_identity.oidc_provider_config`

Read-Only:

- **provider_url** (String)
- **audience** (String)

<a id="nestedblock--status--add_ons--metrics"></a>

### `status.add_ons.metrics`

Read-Only:

- **prometheus_endpoint** (String)
- **remote_write_config** (String)

<a id="nestedblock--status--add_ons--logs"></a>

### `status.add_ons.logs`

Read-Only:

- **loki_address** (String) Loki endpoint to query logs from.

<a id="nestedblock--status--add_ons--aws"></a>

### `status.add_ons.aws`

Read-Only:

- **trust_policy** (String)

## Example Usage - Generic Provider

```terraform
resource "cpln_mk8s" "generic" {

    name        = "demo-mk8s-generic-provider"
    description = "demo-mk8s-generic-provider"

    tags = {
        terraform_generated = "true"
        acceptance_test     = "true"
    }

    version = "1.28.4"

    firewall {
        source_cidr = "192.168.1.255"
        description = "hello world"
    }
    
    generic_provider {
        location = "aws-eu-central-1"
        
        networking {
            service_network = "10.43.0.0/16"
            pod_network 	= "10.42.0.0/16"
        }
        
        node_pool {
            name = "my-node-pool"

            labels = {
                hello = "world"
            }

            taint {
                key    = "hello"
                value  = "world"
                effect = "NoSchedule"
            }
        }
    }

    add_ons {
        dashboard = true

        azure_workload_identity {
            tenant_id = "7f43458a-a34e-4bfa-9e56-e2289e49c4ec"
        }

        aws_workload_identity = true
        local_path_storage    = true

        metrics {
            kube_state    = true
            core_dns      = true
            kubelet       = true
            api_server    = true
            node_exporter = true
            cadvisor      = true

            scrape_annotated {
                interval_seconds   = 30
                include_namespaces = "^elastic"
                exclude_namespaces  = "^elastic"
                retain_labels      = "^\\w+$"
            }
        }

        logs {
            audit_enabled      = true
            include_namespaces = "^elastic"
            exclude_namespaces  = "^elastic"
        }

        registry_mirror {
            mirror {
                registry = "registry.mycompany.com"
                mirrors  = ["https://mirror1.mycompany.com"]
            }

            mirror {
                registry = "docker.io"
                mirrors  = ["https://us-mirror.gcr.io"]
            }
        }

        nvidia {
            taint_gpu_nodes = true
        }

        azure_acr {
            client_id = "4e25b134-160b-4a9d-b392-13b381ced5ef"
        }

        byok = {
            ignore_updates = false
            location       = "/org/terraform-test-org/location/test-byok"

            config = {
                actuator = {
                    min_cpu    = "50m"
                    max_cpu    = "8001m"
                    min_memory = "200Mi"
                    max_memory = "8000Mi"
                    log_level  = "info"
                    env = {
                        CACHE_PERIOD_DATA_SERVICE = "600"
                        LABEL_NODES               = "false"
                    }
                }

                middlebox = {
                    enabled              = false
                    bandwidth_alert_mbps = 650
                }

                common = {
                    deployment_replicas = 1

                    pdb = {
                        max_unavailable = 1
                    }
                }

                longhorn = {
                    replicas = 2
                }

                ingress = {
                    cpu            = "50m"
                    memory         = "200Mi"
                    target_percent = 6000
                }

                istio = {
                    istiod = {
                        replicas   = 2
                        min_cpu    = "50m"
                        max_cpu    = "8001m"
                        min_memory = "100Mi"
                        max_memory = "8000Mi"
                        pdb        = 0
                    }

                    ingress_gateway = {
                        replicas   = 2
                        max_cpu    = "1"
                        max_memory = "1Gi"
                    }

                    sidecar = {
                        min_cpu    = "0m"
                        min_memory = "200Mi"
                    }
                }

                log_splitter = {
                    min_cpu         = "1m"
                    max_cpu         = "1000m"
                    min_memory      = "10Mi"
                    max_memory      = "2000Mi"
                    mem_buffer_size = "128M"
                    per_pod_rate    = 10000
                }

                monitoring = {
                    min_memory = "100Mi"
                    max_memory = "20Gi"

                    kube_state_metrics = {
                        min_memory = "25Mi"
                    }

                    prometheus = {
                        main = {
                            storage = "10Gi"
                        }
                    }
                }

                redis = {
                    min_cpu    = "10m"
                    max_cpu    = "2001m"
                    min_memory = "100Mi"
                    max_memory = "1000Mi"
                    storage    = "8Gi"
                }

                redis_ha = {
                    min_cpu    = "50m"
                    max_cpu    = "2001m"
                    min_memory = "100Mi"
                    max_memory = "1000Mi"
                    storage    = 0
                }

                redis_sentinel = {
                    min_cpu    = "10m"
                    max_cpu    = "500m"
                    min_memory = "10Mi"
                    max_memory = "400Mi"
                    storage    = 0
                }

                tempo_agent = {
                    min_cpu    = "0m"
                    min_memory = "10Mi"
                }

                internal_dns = {
                    min_cpu    = "0m"
                    max_cpu    = "500m"
                    min_memory = "10Mi"
                    max_memory = "400Mi"
                }
            }
        }

        sysbox = true
    }
}
```

## Example Usage - Hetzner Provider

```terraform
resource "cpln_mk8s" "hetzner" {
    
    name        = "demo-mk8s-hetzner-provider"
    description = "demo-mk8s-hetzner-provider"

    tags = {
        terraform_generated = "true"
        acceptance_test     = "true"
    }

    version = "1.28.4"

    firewall {
        source_cidr = "192.168.1.255"
        description = "hello world"
    }
    
    hetzner_provider {
        
        region = "fsn1"

        hetzner_labels = {
            hello = "world"
        }

        networking {
            service_network = "10.43.0.0/16"
            pod_network 	= "10.42.0.0/16"
        }

        pre_install_script = "#! echo hello world"
        token_secret_link  = "/org/terraform-test-org/secret/hetzner"
        network_id 		   = "2808575"

        node_pool {
            name = "my-hetzner-node-pool"

            labels = {
                hello = "world"
            }

            taint {
                key    = "hello"
                value  = "world"
                effect = "NoSchedule"
            }

            server_type    = "cpx11"
            override_image = "debian-11"
            min_size 	   = 0
            max_size 	   = 0
        }

        dedicated_server_node_pool {
            name = "my-node-pool"

            labels = {
                hello = "world"
            }

            taint {
                key    = "hello"
                value  = "world"
                effect = "NoSchedule"
            }
        }

        image 	= "centos-7"
        ssh_key = "10925607"

        autoscaler {
            expander 	  		  = ["most-pods"]
            unneeded_time         = "10m"
            unready_time  		  = "20m"
            utilization_threshold = 0.7
        }

        floating_ip_selector = {
            floating_ip_1 = "123.45.67.89"
            floating_ip_2 = "98.76.54.32"
        }
    }

    add_ons {
        dashboard = true

        azure_workload_identity {
            tenant_id = "7f43458a-a34e-4bfa-9e56-e2289e49c4ec"
        }

        aws_workload_identity = true
        local_path_storage    = true

        metrics {
            kube_state    = true
            core_dns      = true
            kubelet       = true
            api_server    = true
            node_exporter = true
            cadvisor      = true

            scrape_annotated {
                interval_seconds   = 30
                include_namespaces = "^elastic"
                exclude_namespaces  = "^elastic"
                retain_labels      = "^\\w+$"
            }
        }

        logs {
            audit_enabled      = true
            include_namespaces = "^elastic"
            exclude_namespaces  = "^elastic"
        }

        registry_mirror {
            mirror {
                registry = "registry.mycompany.com"
                mirrors  = ["https://mirror1.mycompany.com"]
            }

            mirror {
                registry = "docker.io"
                mirrors  = ["https://us-mirror.gcr.io"]
            }
        }

        nvidia {
            taint_gpu_nodes = true
        }

        azure_acr {
            client_id = "4e25b134-160b-4a9d-b392-13b381ced5ef"
        }

        byok = {
            ignore_updates = false
            location       = "/org/terraform-test-org/location/test-byok"

            config = {
                actuator = {
                    min_cpu    = "50m"
                    max_cpu    = "8001m"
                    min_memory = "200Mi"
                    max_memory = "8000Mi"
                    log_level  = "info"
                    env = {
                        CACHE_PERIOD_DATA_SERVICE = "600"
                        LABEL_NODES               = "false"
                    }
                }

                middlebox = {
                    enabled              = false
                    bandwidth_alert_mbps = 650
                }

                common = {
                    deployment_replicas = 1

                    pdb = {
                        max_unavailable = 1
                    }
                }

                longhorn = {
                    replicas = 2
                }

                ingress = {
                    cpu            = "50m"
                    memory         = "200Mi"
                    target_percent = 6000
                }

                istio = {
                    istiod = {
                        replicas   = 2
                        min_cpu    = "50m"
                        max_cpu    = "8001m"
                        min_memory = "100Mi"
                        max_memory = "8000Mi"
                        pdb        = 0
                    }

                    ingress_gateway = {
                        replicas   = 2
                        max_cpu    = "1"
                        max_memory = "1Gi"
                    }

                    sidecar = {
                        min_cpu    = "0m"
                        min_memory = "200Mi"
                    }
                }

                log_splitter = {
                    min_cpu         = "1m"
                    max_cpu         = "1000m"
                    min_memory      = "10Mi"
                    max_memory      = "2000Mi"
                    mem_buffer_size = "128M"
                    per_pod_rate    = 10000
                }

                monitoring = {
                    min_memory = "100Mi"
                    max_memory = "20Gi"

                    kube_state_metrics = {
                        min_memory = "25Mi"
                    }

                    prometheus = {
                        main = {
                            storage = "10Gi"
                        }
                    }
                }

                redis = {
                    min_cpu    = "10m"
                    max_cpu    = "2001m"
                    min_memory = "100Mi"
                    max_memory = "1000Mi"
                    storage    = "8Gi"
                }

                redis_ha = {
                    min_cpu    = "50m"
                    max_cpu    = "2001m"
                    min_memory = "100Mi"
                    max_memory = "1000Mi"
                    storage    = 0
                }

                redis_sentinel = {
                    min_cpu    = "10m"
                    max_cpu    = "500m"
                    min_memory = "10Mi"
                    max_memory = "400Mi"
                    storage    = 0
                }

                tempo_agent = {
                    min_cpu    = "0m"
                    min_memory = "10Mi"
                }

                internal_dns = {
                    min_cpu    = "0m"
                    max_cpu    = "500m"
                    min_memory = "10Mi"
                    max_memory = "400Mi"
                }
            }
        }

        sysbox = true
    }
}
```

## Example Usage - AWS Provider

```terraform
resource "cpln_mk8s" "aws" {

    name        = "demo-mk8s-aws-provider"
    description = "demo-mk8s-aws-provider"

    tags = {
        terraform_generated = "true"
        acceptance_test     = "true"
    }

    version = "1.28.4"

    firewall {
        source_cidr = "192.168.1.255"
        description = "hello world"
    }

    aws_provider {

        region = "eu-central-1"

        aws_tags = {
            hello = "world"
        }

        skip_create_roles = false

        networking {}

        pre_install_script = "#! echo hello world"

        image {
            recommended = "amazon/al2023"
        }

        deploy_role_arn         = "arn:aws:iam::12345678901:role/cpln"
        vpc_id                  = "vpc-03105bd4dc058d3a8"
        key_pair                = "cem_uzak"
        disk_encryption_key_arn = "arn:aws:kms:eu-central-1:12345678901:key/0a1bcd23-4567-8901-e2fg-3h4i5jk678lm"

        security_group_ids = ["sg-031480aa7a1e6e38b"]
        extra_node_policies = ["arn:aws:iam::aws:policy/IAMFullAccess"]

        deploy_role_chain {
            role_arn            = "arn:aws:iam::483676437512:role/mk8s-chain-1"
            external_id         = "chain-1"
            session_name_prefix = "foo-"
        }

        node_pool {
            name = "my-aws-node-pool"

            labels = {
                hello = "world"
            }

            taint {
                key    = "hello"
                value  = "world"
                effect = "NoSchedule"
            }

            instance_types = ["t4g.nano"]

            override_image {
                exact = "ami-123"
            }

            boot_disk_size                           = 20
            min_size                                 = 0
            max_size                                 = 0
            on_demand_base_capacity                  = 0
            on_demand_percentage_above_base_capacity = 0
            spot_allocation_strategy                 = "lowest-price"

            subnet_ids               = ["subnet-0e564a042e2a45009"]
            extra_security_group_ids = ["sg-031480aa7a1e6e38b"]
        }

        autoscaler {
            expander 	  		  = ["most-pods"]
            unneeded_time         = "10m"
            unready_time  		  = "20m"
            utilization_threshold = 0.7
        }
    }

    add_ons {
        dashboard = true

        azure_workload_identity {
            tenant_id = "7f43458a-a34e-4bfa-9e56-e2289e49c4ec"
        }

        aws_workload_identity = true
        local_path_storage    = true

        metrics {
            kube_state    = true
            core_dns      = true
            kubelet       = true
            api_server    = true
            node_exporter = true
            cadvisor      = true

            scrape_annotated {
                interval_seconds   = 30
                include_namespaces = "^elastic"
                exclude_namespaces  = "^elastic"
                retain_labels      = "^\\w+$"
            }
        }

        logs {
            audit_enabled      = true
            include_namespaces = "^elastic"
            exclude_namespaces  = "^elastic"
        }

        registry_mirror {
            mirror {
                registry = "registry.mycompany.com"
                mirrors  = ["https://mirror1.mycompany.com"]
            }

            mirror {
                registry = "docker.io"
                mirrors  = ["https://us-mirror.gcr.io"]
            }
        }

        nvidia {
            taint_gpu_nodes = true
        }

        aws_efs {
            role_arn = "arn:aws:iam::123456789012:role/aws-efs-role"
        }

        aws_ecr {
            role_arn = "arn:aws:iam::123456789012:role/aws-ecr-role"
        }

        aws_elb {
            role_arn = "arn:aws:iam::123456789012:role/aws-elb-role"
        }

        azure_acr {
            client_id = "4e25b134-160b-4a9d-b392-13b381ced5ef"
        }

        byok = {
            ignore_updates = false
            location       = "/org/terraform-test-org/location/test-byok"

            config = {
                actuator = {
                    min_cpu    = "50m"
                    max_cpu    = "8001m"
                    min_memory = "200Mi"
                    max_memory = "8000Mi"
                    log_level  = "info"
                    env = {
                        CACHE_PERIOD_DATA_SERVICE = "600"
                        LABEL_NODES               = "false"
                    }
                }

                middlebox = {
                    enabled              = false
                    bandwidth_alert_mbps = 650
                }

                common = {
                    deployment_replicas = 1

                    pdb = {
                        max_unavailable = 1
                    }
                }

                longhorn = {
                    replicas = 2
                }

                ingress = {
                    cpu            = "50m"
                    memory         = "200Mi"
                    target_percent = 6000
                }

                istio = {
                    istiod = {
                        replicas   = 2
                        min_cpu    = "50m"
                        max_cpu    = "8001m"
                        min_memory = "100Mi"
                        max_memory = "8000Mi"
                        pdb        = 0
                    }

                    ingress_gateway = {
                        replicas   = 2
                        max_cpu    = "1"
                        max_memory = "1Gi"
                    }

                    sidecar = {
                        min_cpu    = "0m"
                        min_memory = "200Mi"
                    }
                }

                log_splitter = {
                    min_cpu         = "1m"
                    max_cpu         = "1000m"
                    min_memory      = "10Mi"
                    max_memory      = "2000Mi"
                    mem_buffer_size = "128M"
                    per_pod_rate    = 10000
                }

                monitoring = {
                    min_memory = "100Mi"
                    max_memory = "20Gi"

                    kube_state_metrics = {
                        min_memory = "25Mi"
                    }

                    prometheus = {
                        main = {
                            storage = "10Gi"
                        }
                    }
                }

                redis = {
                    min_cpu    = "10m"
                    max_cpu    = "2001m"
                    min_memory = "100Mi"
                    max_memory = "1000Mi"
                    storage    = "8Gi"
                }

                redis_ha = {
                    min_cpu    = "50m"
                    max_cpu    = "2001m"
                    min_memory = "100Mi"
                    max_memory = "1000Mi"
                    storage    = 0
                }

                redis_sentinel = {
                    min_cpu    = "10m"
                    max_cpu    = "500m"
                    min_memory = "10Mi"
                    max_memory = "400Mi"
                    storage    = 0
                }

                tempo_agent = {
                    min_cpu    = "0m"
                    min_memory = "10Mi"
                }

                internal_dns = {
                    min_cpu    = "0m"
                    max_cpu    = "500m"
                    min_memory = "10Mi"
                    max_memory = "400Mi"
                }
            }
        }

        sysbox = true
    }
}
```

## Example Usage - Linode Provider

```terraform
resource "cpln_mk8s" "linode" {

    name        = "demo-mk8s-linode-provider"
    description = "demo-mk8s-linode-provider"

    tags = {
        terraform_generated = "true"
        acceptance_test     = "true"
    }

    version = "1.28.4"

    firewall {
        source_cidr = "192.168.1.255"
        description = "hello world"
    }

    linode_provider {
        region             = "fr-par"
        token_secret_link  = "/org/terraform-test-org/secret/linode"
        image              = "linode/ubuntu24.04"
        vpc_id             = "93666"
        firewall_id        = "168425"
        pre_install_script = "#! echo hello world"

        authorized_users = ["cpln"]

        node_pool {
            name = "my-linode-node-pool"

            labels = {
                hello = "world"
            }

            taint {
                key    = "hello"
                value  = "world"
                effect = "NoSchedule"
            }

            server_type    = "g6-nanode-1"
            override_image = "linode/debian11"
            subnet_id      = "90623"
            min_size 	   = 0
            max_size 	   = 0
        }

        networking {
            service_network = "10.43.0.0/16"
            pod_network 	= "10.42.0.0/16"
        }

        autoscaler {
            expander 	  		  = ["most-pods"]
            unneeded_time         = "10m"
            unready_time  		  = "20m"
            utilization_threshold = 0.7
        }
    }

    add_ons {
        dashboard = true

        azure_workload_identity {
            tenant_id = "7f43458a-a34e-4bfa-9e56-e2289e49c4ec"
        }

        aws_workload_identity = true
        local_path_storage    = true

        metrics {
            kube_state    = true
            core_dns      = true
            kubelet       = true
            api_server    = true
            node_exporter = true
            cadvisor      = true

            scrape_annotated {
                interval_seconds   = 30
                include_namespaces = "^elastic"
                exclude_namespaces = "^elastic"
                retain_labels      = "^\\w+$"
            }
        }

        logs {
            audit_enabled      = true
            include_namespaces = "^elastic"
            exclude_namespaces = "^elastic"
        }

        registry_mirror {
            mirror {
                registry = "registry.mycompany.com"
                mirrors  = ["https://mirror1.mycompany.com"]
            }

            mirror {
                registry = "docker.io"
                mirrors  = ["https://us-mirror.gcr.io"]
            }
        }

        nvidia {
            taint_gpu_nodes = true
        }

        azure_acr {
            client_id = "4e25b134-160b-4a9d-b392-13b381ced5ef"
        }

        byok = {
            ignore_updates = false
            location       = "/org/terraform-test-org/location/test-byok"

            config = {
                actuator = {
                    min_cpu    = "50m"
                    max_cpu    = "8001m"
                    min_memory = "200Mi"
                    max_memory = "8000Mi"
                    log_level  = "info"
                    env = {
                        CACHE_PERIOD_DATA_SERVICE = "600"
                        LABEL_NODES               = "false"
                    }
                }

                middlebox = {
                    enabled              = false
                    bandwidth_alert_mbps = 650
                }

                common = {
                    deployment_replicas = 1

                    pdb = {
                        max_unavailable = 1
                    }
                }

                longhorn = {
                    replicas = 2
                }

                ingress = {
                    cpu            = "50m"
                    memory         = "200Mi"
                    target_percent = 6000
                }

                istio = {
                    istiod = {
                        replicas   = 2
                        min_cpu    = "50m"
                        max_cpu    = "8001m"
                        min_memory = "100Mi"
                        max_memory = "8000Mi"
                        pdb        = 0
                    }

                    ingress_gateway = {
                        replicas   = 2
                        max_cpu    = "1"
                        max_memory = "1Gi"
                    }

                    sidecar = {
                        min_cpu    = "0m"
                        min_memory = "200Mi"
                    }
                }

                log_splitter = {
                    min_cpu         = "1m"
                    max_cpu         = "1000m"
                    min_memory      = "10Mi"
                    max_memory      = "2000Mi"
                    mem_buffer_size = "128M"
                    per_pod_rate    = 10000
                }

                monitoring = {
                    min_memory = "100Mi"
                    max_memory = "20Gi"

                    kube_state_metrics = {
                        min_memory = "25Mi"
                    }

                    prometheus = {
                        main = {
                            storage = "10Gi"
                        }
                    }
                }

                redis = {
                    min_cpu    = "10m"
                    max_cpu    = "2001m"
                    min_memory = "100Mi"
                    max_memory = "1000Mi"
                    storage    = "8Gi"
                }

                redis_ha = {
                    min_cpu    = "50m"
                    max_cpu    = "2001m"
                    min_memory = "100Mi"
                    max_memory = "1000Mi"
                    storage    = 0
                }

                redis_sentinel = {
                    min_cpu    = "10m"
                    max_cpu    = "500m"
                    min_memory = "10Mi"
                    max_memory = "400Mi"
                    storage    = 0
                }

                tempo_agent = {
                    min_cpu    = "0m"
                    min_memory = "10Mi"
                }

                internal_dns = {
                    min_cpu    = "0m"
                    max_cpu    = "500m"
                    min_memory = "10Mi"
                    max_memory = "400Mi"
                }
            }
        }

        sysbox = true
    }
}
```

## Example Usage - Oblivus Provider

```terraform
resource "cpln_mk8s" "oblivus" {
    
    name        = "demo-mk8s-oblivus-provider"
    description = "demo-mk8s-oblivus-provider"

    tags = {
        terraform_generated = "true"
        acceptance_test     = "true"
    }

    version = "1.28.4"

    firewall {
        source_cidr = "192.168.1.255"
        description = "hello world"
    }

    oblivus_provider {
        datacenter         = "OSL1"
        token_secret_link  = "/org/terraform-test-org/secret/oblivus"
        pre_install_script = "#! echo hello world"

        node_pool {
            name     = "my-oblivus-node-pool"
            min_size = 0
            max_size = 0
            flavor   = "INTEL_XEON_V3_x4"

            labels = {
                hello = "world"
            }

            taint {
                key    = "hello"
                value  = "world"
                effect = "NoSchedule"
            }
        }

        unmanaged_node_pool {
            name = "my-node-pool"

            labels = {
                hello = "world"
            }

            taint {
                key    = "hello"
                value  = "world"
                effect = "NoSchedule"
            }
        }

        autoscaler {
            expander 	  		  = ["most-pods"]
            unneeded_time         = "10m"
            unready_time  		  = "20m"
            utilization_threshold = 0.7
        }
    }

    add_ons {
        dashboard = true

        azure_workload_identity {
            tenant_id = "7f43458a-a34e-4bfa-9e56-e2289e49c4ec"
        }

        aws_workload_identity = true
        local_path_storage    = true

        metrics {
            kube_state    = true
            core_dns      = true
            kubelet       = true
            api_server    = true
            node_exporter = true
            cadvisor      = true

            scrape_annotated {
                interval_seconds   = 30
                include_namespaces = "^elastic"
                exclude_namespaces  = "^elastic"
                retain_labels      = "^\\w+$"
            }
        }

        logs {
            audit_enabled      = true
            include_namespaces = "^elastic"
            exclude_namespaces  = "^elastic"
        }

        registry_mirror {
            mirror {
                registry = "registry.mycompany.com"
                mirrors  = ["https://mirror1.mycompany.com"]
            }

            mirror {
                registry = "docker.io"
                mirrors  = ["https://us-mirror.gcr.io"]
            }
        }

        nvidia {
            taint_gpu_nodes = true
        }

        azure_acr {
            client_id = "4e25b134-160b-4a9d-b392-13b381ced5ef"
        }

        byok = {
            ignore_updates = false
            location       = "/org/terraform-test-org/location/test-byok"

            config = {
                actuator = {
                    min_cpu    = "50m"
                    max_cpu    = "8001m"
                    min_memory = "200Mi"
                    max_memory = "8000Mi"
                    log_level  = "info"
                    env = {
                        CACHE_PERIOD_DATA_SERVICE = "600"
                        LABEL_NODES               = "false"
                    }
                }

                middlebox = {
                    enabled              = false
                    bandwidth_alert_mbps = 650
                }

                common = {
                    deployment_replicas = 1

                    pdb = {
                        max_unavailable = 1
                    }
                }

                longhorn = {
                    replicas = 2
                }

                ingress = {
                    cpu            = "50m"
                    memory         = "200Mi"
                    target_percent = 6000
                }

                istio = {
                    istiod = {
                        replicas   = 2
                        min_cpu    = "50m"
                        max_cpu    = "8001m"
                        min_memory = "100Mi"
                        max_memory = "8000Mi"
                        pdb        = 0
                    }

                    ingress_gateway = {
                        replicas   = 2
                        max_cpu    = "1"
                        max_memory = "1Gi"
                    }

                    sidecar = {
                        min_cpu    = "0m"
                        min_memory = "200Mi"
                    }
                }

                log_splitter = {
                    min_cpu         = "1m"
                    max_cpu         = "1000m"
                    min_memory      = "10Mi"
                    max_memory      = "2000Mi"
                    mem_buffer_size = "128M"
                    per_pod_rate    = 10000
                }

                monitoring = {
                    min_memory = "100Mi"
                    max_memory = "20Gi"

                    kube_state_metrics = {
                        min_memory = "25Mi"
                    }

                    prometheus = {
                        main = {
                            storage = "10Gi"
                        }
                    }
                }

                redis = {
                    min_cpu    = "10m"
                    max_cpu    = "2001m"
                    min_memory = "100Mi"
                    max_memory = "1000Mi"
                    storage    = "8Gi"
                }

                redis_ha = {
                    min_cpu    = "50m"
                    max_cpu    = "2001m"
                    min_memory = "100Mi"
                    max_memory = "1000Mi"
                    storage    = 0
                }

                redis_sentinel = {
                    min_cpu    = "10m"
                    max_cpu    = "500m"
                    min_memory = "10Mi"
                    max_memory = "400Mi"
                    storage    = 0
                }

                tempo_agent = {
                    min_cpu    = "0m"
                    min_memory = "10Mi"
                }

                internal_dns = {
                    min_cpu    = "0m"
                    max_cpu    = "500m"
                    min_memory = "10Mi"
                    max_memory = "400Mi"
                }
            }
        }

        sysbox = true
    }
}
```

## Example Usage - Lambdalabs Provider

```terraform
resource "cpln_mk8s" "lambdalabs" {
    name        = "demo-mk8s-lambdalabs-provider"
    description = "demo-mk8s-lambdalabs-provider"

    tags = {
        terraform_generated = "true"
        acceptance_test     = "true"
    }

    version = "1.28.4"

    firewall {
        source_cidr = "192.168.1.255"
        description = "hello world"
    }

    lambdalabs_provider {
        region             = "europe-central-1"
        token_secret_link  = "/org/ORG_NAME/secret/SECRET_NAME"
        ssh_key            = "some-key"
        pre_install_script = "#! echo hello world"
        
        node_pool {
            name = "my-node-pool"

            labels = {
                hello = "world"
            }

            taint {
                key    = "hello"
                value  = "world"
                effect = "NoSchedule"
            }

            instance_type = "cpu_4x_general"
        }

        unmanaged_node_pool {
            name = "my-unmanaged-node-pool"

            labels = {
                hello = "world"
            }

            taint {
                key    = "hello"
                value  = "world"
                effect = "NoSchedule"
            }
        }

        autoscaler {
            expander 	  		  = ["most-pods"]
            unneeded_time         = "10m"
            unready_time  		  = "20m"
            utilization_threshold = 0.7
        }
    }

    add_ons {
        dashboard = true

        azure_workload_identity {
            tenant_id = "7f43458a-a34e-4bfa-9e56-e2289e49c4ec"
        }

        aws_workload_identity = true
        local_path_storage    = true

        metrics {
            kube_state    = true
            core_dns      = true
            kubelet       = true
            api_server    = true
            node_exporter = true
            cadvisor      = true

            scrape_annotated {
                interval_seconds   = 30
                include_namespaces = "^elastic"
                exclude_namespaces  = "^elastic"
                retain_labels      = "^\\w+$"
            }
        }

        logs {
            audit_enabled      = true
            include_namespaces = "^elastic"
            exclude_namespaces  = "^elastic"
        }

        registry_mirror {
            mirror {
                registry = "registry.mycompany.com"
                mirrors  = ["https://mirror1.mycompany.com"]
            }

            mirror {
                registry = "docker.io"
                mirrors  = ["https://us-mirror.gcr.io"]
            }
        }

        nvidia {
            taint_gpu_nodes = true
        }

        azure_acr {
            client_id = "4e25b134-160b-4a9d-b392-13b381ced5ef"
        }

        byok = {
            ignore_updates = false
            location       = "/org/terraform-test-org/location/test-byok"

            config = {
                actuator = {
                    min_cpu    = "50m"
                    max_cpu    = "8001m"
                    min_memory = "200Mi"
                    max_memory = "8000Mi"
                    log_level  = "info"
                    env = {
                        CACHE_PERIOD_DATA_SERVICE = "600"
                        LABEL_NODES               = "false"
                    }
                }

                middlebox = {
                    enabled              = false
                    bandwidth_alert_mbps = 650
                }

                common = {
                    deployment_replicas = 1

                    pdb = {
                        max_unavailable = 1
                    }
                }

                longhorn = {
                    replicas = 2
                }

                ingress = {
                    cpu            = "50m"
                    memory         = "200Mi"
                    target_percent = 6000
                }

                istio = {
                    istiod = {
                        replicas   = 2
                        min_cpu    = "50m"
                        max_cpu    = "8001m"
                        min_memory = "100Mi"
                        max_memory = "8000Mi"
                        pdb        = 0
                    }

                    ingress_gateway = {
                        replicas   = 2
                        max_cpu    = "1"
                        max_memory = "1Gi"
                    }

                    sidecar = {
                        min_cpu    = "0m"
                        min_memory = "200Mi"
                    }
                }

                log_splitter = {
                    min_cpu         = "1m"
                    max_cpu         = "1000m"
                    min_memory      = "10Mi"
                    max_memory      = "2000Mi"
                    mem_buffer_size = "128M"
                    per_pod_rate    = 10000
                }

                monitoring = {
                    min_memory = "100Mi"
                    max_memory = "20Gi"

                    kube_state_metrics = {
                        min_memory = "25Mi"
                    }

                    prometheus = {
                        main = {
                            storage = "10Gi"
                        }
                    }
                }

                redis = {
                    min_cpu    = "10m"
                    max_cpu    = "2001m"
                    min_memory = "100Mi"
                    max_memory = "1000Mi"
                    storage    = "8Gi"
                }

                redis_ha = {
                    min_cpu    = "50m"
                    max_cpu    = "2001m"
                    min_memory = "100Mi"
                    max_memory = "1000Mi"
                    storage    = 0
                }

                redis_sentinel = {
                    min_cpu    = "10m"
                    max_cpu    = "500m"
                    min_memory = "10Mi"
                    max_memory = "400Mi"
                    storage    = 0
                }

                tempo_agent = {
                    min_cpu    = "0m"
                    min_memory = "10Mi"
                }

                internal_dns = {
                    min_cpu    = "0m"
                    max_cpu    = "500m"
                    min_memory = "10Mi"
                    max_memory = "400Mi"
                }
            }
        }
    }
}
```


## Example Usage - Paperspace Provider

```terraform
resource "cpln_mk8s" "paperspace" {
		
    name        = "demo-mk8s-paperspace-provider"
    description = "demo-mk8s-paperspace-provider"

    tags = {
        terraform_generated = "true"
        acceptance_test     = "true"
    }

    version = "1.28.4"

    firewall {
        source_cidr = "192.168.1.255"
        description = "hello world"
    }

    paperspace_provider {
        region             = "CA1"
        token_secret_link  = "/org/terraform-test-org/secret/paperspace"
        pre_install_script = "#! echo hello world"
        shared_drives      = ["california"]
        network_id         = "nla0jotp"
        
        node_pool {
            name           = "my-paperspace-node-pool"
            min_size       = 0
            max_size       = 0
            public_ip_type = "dynamic"
            boot_disk_size = 50
            machine_type   = "GPU+"

            labels = {
                hello = "world"
            }

            taint {
                key    = "hello"
                value  = "world"
                effect = "NoSchedule"
            }
        }

        unmanaged_node_pool {
            name = "my-node-pool"

            labels = {
                hello = "world"
            }

            taint {
                key    = "hello"
                value  = "world"
                effect = "NoSchedule"
            }
        }

        autoscaler {
            expander 	  		  = ["most-pods"]
            unneeded_time         = "10m"
            unready_time  		  = "20m"
            utilization_threshold = 0.7
        }
    }

    add_ons {
        dashboard = true

        azure_workload_identity {
            tenant_id = "7f43458a-a34e-4bfa-9e56-e2289e49c4ec"
        }

        aws_workload_identity = true
        local_path_storage    = true

        metrics {
            kube_state    = true
            core_dns      = true
            kubelet       = true
            api_server    = true
            node_exporter = true
            cadvisor      = true

            scrape_annotated {
                interval_seconds   = 30
                include_namespaces = "^elastic"
                exclude_namespaces = "^elastic"
                retain_labels      = "^\\w+$"
            }
        }

        logs {
            audit_enabled      = true
            include_namespaces = "^elastic"
            exclude_namespaces = "^elastic"
        }

        registry_mirror {
            mirror {
                registry = "registry.mycompany.com"
                mirrors  = ["https://mirror1.mycompany.com"]
            }

            mirror {
                registry = "docker.io"
                mirrors  = ["https://us-mirror.gcr.io"]
            }
        }

        nvidia {
            taint_gpu_nodes = true
        }

        azure_acr {
            client_id = "4e25b134-160b-4a9d-b392-13b381ced5ef"
        }

        byok = {
            ignore_updates = false
            location       = "/org/terraform-test-org/location/test-byok"

            config = {
                actuator = {
                    min_cpu    = "50m"
                    max_cpu    = "8001m"
                    min_memory = "200Mi"
                    max_memory = "8000Mi"
                    log_level  = "info"
                    env = {
                        CACHE_PERIOD_DATA_SERVICE = "600"
                        LABEL_NODES               = "false"
                    }
                }

                middlebox = {
                    enabled              = false
                    bandwidth_alert_mbps = 650
                }

                common = {
                    deployment_replicas = 1

                    pdb = {
                        max_unavailable = 1
                    }
                }

                longhorn = {
                    replicas = 2
                }

                ingress = {
                    cpu            = "50m"
                    memory         = "200Mi"
                    target_percent = 6000
                }

                istio = {
                    istiod = {
                        replicas   = 2
                        min_cpu    = "50m"
                        max_cpu    = "8001m"
                        min_memory = "100Mi"
                        max_memory = "8000Mi"
                        pdb        = 0
                    }

                    ingress_gateway = {
                        replicas   = 2
                        max_cpu    = "1"
                        max_memory = "1Gi"
                    }

                    sidecar = {
                        min_cpu    = "0m"
                        min_memory = "200Mi"
                    }
                }

                log_splitter = {
                    min_cpu         = "1m"
                    max_cpu         = "1000m"
                    min_memory      = "10Mi"
                    max_memory      = "2000Mi"
                    mem_buffer_size = "128M"
                    per_pod_rate    = 10000
                }

                monitoring = {
                    min_memory = "100Mi"
                    max_memory = "20Gi"

                    kube_state_metrics = {
                        min_memory = "25Mi"
                    }

                    prometheus = {
                        main = {
                            storage = "10Gi"
                        }
                    }
                }

                redis = {
                    min_cpu    = "10m"
                    max_cpu    = "2001m"
                    min_memory = "100Mi"
                    max_memory = "1000Mi"
                    storage    = "8Gi"
                }

                redis_ha = {
                    min_cpu    = "50m"
                    max_cpu    = "2001m"
                    min_memory = "100Mi"
                    max_memory = "1000Mi"
                    storage    = 0
                }

                redis_sentinel = {
                    min_cpu    = "10m"
                    max_cpu    = "500m"
                    min_memory = "10Mi"
                    max_memory = "400Mi"
                    storage    = 0
                }

                tempo_agent = {
                    min_cpu    = "0m"
                    min_memory = "10Mi"
                }

                internal_dns = {
                    min_cpu    = "0m"
                    max_cpu    = "500m"
                    min_memory = "10Mi"
                    max_memory = "400Mi"
                }
            }
        }

        sysbox = true
    }
}
```

## Example Usage - Ephemeral Provider

```terraform
resource "cpln_mk8s" "ephemeral" {

    name        = "demo-mk8s-ephemeral-provider"
    description = "demo-mk8s-ephemeral-provider"

    tags = {
        terraform_generated = "true"
        acceptance_test     = "true"
    }

    version = "1.28.4"

    firewall {
        source_cidr = "192.168.1.255"
        description = "hello world"
    }
    
    ephemeral_provider {
        location = "aws-eu-central-1"
        
        node_pool {
            name = "my-node-pool"

            labels = {
                hello = "world"
            }

            taint {
                key    = "hello"
                value  = "world"
                effect = "NoSchedule"
            }

            count  = 1
            arch   = "arm64"
            flavor = "debian"
            cpu    = "50m"
            memory = "128Mi"
        }
    }

    add_ons {
        dashboard = true

        azure_workload_identity {
            tenant_id = "7f43458a-a34e-4bfa-9e56-e2289e49c4ec"
        }

        aws_workload_identity = true
        local_path_storage    = true

        metrics {
            kube_state    = true
            core_dns      = true
            kubelet       = true
            api_server    = true
            node_exporter = true
            cadvisor      = true

            scrape_annotated {
                interval_seconds   = 30
                include_namespaces = "^elastic"
                exclude_namespaces  = "^elastic"
                retain_labels      = "^\\w+$"
            }
        }

        logs {
            audit_enabled      = true
            include_namespaces = "^elastic"
            exclude_namespaces  = "^elastic"
        }

        registry_mirror {
            mirror {
                registry = "registry.mycompany.com"
                mirrors  = ["https://mirror1.mycompany.com"]
            }

            mirror {
                registry = "docker.io"
                mirrors  = ["https://us-mirror.gcr.io"]
            }
        }

        nvidia {
            taint_gpu_nodes = true
        }

        azure_acr {
            client_id = "4e25b134-160b-4a9d-b392-13b381ced5ef"
        }

        byok = {
            ignore_updates = false
            location       = "/org/terraform-test-org/location/test-byok"

            config = {
                actuator = {
                    min_cpu    = "50m"
                    max_cpu    = "8001m"
                    min_memory = "200Mi"
                    max_memory = "8000Mi"
                    log_level  = "info"
                    env = {
                        CACHE_PERIOD_DATA_SERVICE = "600"
                        LABEL_NODES               = "false"
                    }
                }

                middlebox = {
                    enabled              = false
                    bandwidth_alert_mbps = 650
                }

                common = {
                    deployment_replicas = 1

                    pdb = {
                        max_unavailable = 1
                    }
                }

                longhorn = {
                    replicas = 2
                }

                ingress = {
                    cpu            = "50m"
                    memory         = "200Mi"
                    target_percent = 6000
                }

                istio = {
                    istiod = {
                        replicas   = 2
                        min_cpu    = "50m"
                        max_cpu    = "8001m"
                        min_memory = "100Mi"
                        max_memory = "8000Mi"
                        pdb        = 0
                    }

                    ingress_gateway = {
                        replicas   = 2
                        max_cpu    = "1"
                        max_memory = "1Gi"
                    }

                    sidecar = {
                        min_cpu    = "0m"
                        min_memory = "200Mi"
                    }
                }

                log_splitter = {
                    min_cpu         = "1m"
                    max_cpu         = "1000m"
                    min_memory      = "10Mi"
                    max_memory      = "2000Mi"
                    mem_buffer_size = "128M"
                    per_pod_rate    = 10000
                }

                monitoring = {
                    min_memory = "100Mi"
                    max_memory = "20Gi"

                    kube_state_metrics = {
                        min_memory = "25Mi"
                    }

                    prometheus = {
                        main = {
                            storage = "10Gi"
                        }
                    }
                }

                redis = {
                    min_cpu    = "10m"
                    max_cpu    = "2001m"
                    min_memory = "100Mi"
                    max_memory = "1000Mi"
                    storage    = "8Gi"
                }

                redis_ha = {
                    min_cpu    = "50m"
                    max_cpu    = "2001m"
                    min_memory = "100Mi"
                    max_memory = "1000Mi"
                    storage    = 0
                }

                redis_sentinel = {
                    min_cpu    = "10m"
                    max_cpu    = "500m"
                    min_memory = "10Mi"
                    max_memory = "400Mi"
                    storage    = 0
                }

                tempo_agent = {
                    min_cpu    = "0m"
                    min_memory = "10Mi"
                }

                internal_dns = {
                    min_cpu    = "0m"
                    max_cpu    = "500m"
                    min_memory = "10Mi"
                    max_memory = "400Mi"
                }
            }
        }
    }
}
```

## Example Usage - Triton Provider

```terraform
resource "cpln_mk8s" "triton" {

    name        = "demo-mk8s-triton-provider"
    description = "demo-mk8s-triton-provider"

    tags = {
        terraform_generated = "true"
        acceptance_test     = "true"
    }

    version = "1.28.4"
	
    firewall {
        source_cidr = "192.168.1.255"
        description = "hello world"
    }

    triton_provider {
        pre_install_script = "#! echo hello world"
        location           = "aws-eu-central-1"
        private_network_id = "6704dae9-00f4-48b5-8bbf-1be538f20587"
        firewall_enabled   = false
        image_id           = "6b98a11c-53a4-4a62-99e7-cf3dcf150ab2"
        
        networking {}

        connection {
            url                     = "https://us-central-1.api.mnx.io"
            account                 = "eric_controlplane.com"
            user                    = "julian_controlplane.com"
            private_key_secret_link = "/org/terraform-test-org/secret/triton"
        }
      
        load_balancer {
            none {}
        }

        node_pool {
            name                = "my-triton-node-pool"
            package_id          = "da311341-b42b-45a8-9386-78ede625d0a4"
            override_image_id   = "e2f3f2aa-a833-49e0-94af-7a7e092cdd9e"
            public_network_id   = "5ff1fe03-075b-4e4c-b85b-73de0c452f77"
            min_size            = 0
            max_size            = 0

            private_network_ids = ["6704dae9-00f4-48b5-8bbf-1be538f20587"]

            labels = {
                hello = "world"
            }

            taint {
                key    = "hello"
                value  = "world"
                effect = "NoSchedule"
            }
            
            triton_tags = {
                drink = "water"
            }
        }
        
        autoscaler {
            expander 	  		  = ["most-pods"]
            unneeded_time         = "10m"
            unready_time  		  = "20m"
            utilization_threshold = 0.7
        }
    }

    add_ons {
        dashboard = true

        azure_workload_identity {
            tenant_id = "7f43458a-a34e-4bfa-9e56-e2289e49c4ec"
        }

        aws_workload_identity = true
        local_path_storage    = true

        metrics {
            kube_state    = true
            core_dns      = true
            kubelet       = true
            api_server    = true
            node_exporter = true
            cadvisor      = true

            scrape_annotated {
                interval_seconds   = 30
                include_namespaces = "^\\d+$"
                exclude_namespaces = "^elastic"
                retain_labels      = "^elastic"
            }
        }

        logs {
            audit_enabled      = true
            include_namespaces = "^elastic"
            exclude_namespaces = "^elastic"
        }

        registry_mirror {
            mirror {
                registry = "registry.mycompany.com"
                mirrors  = ["https://mirror1.mycompany.com"]
            }

            mirror {
                registry = "docker.io"
                mirrors  = ["https://us-mirror.gcr.io"]
            }
        }

        nvidia {
            taint_gpu_nodes = true
        }

        azure_acr {
            client_id = "4e25b134-160b-4a9d-b392-13b381ced5ef"
        }

        byok = {
            ignore_updates = false
            location       = "/org/terraform-test-org/location/test-byok"

            config = {
                actuator = {
                    min_cpu    = "50m"
                    max_cpu    = "8001m"
                    min_memory = "200Mi"
                    max_memory = "8000Mi"
                    log_level  = "info"
                    env = {
                        CACHE_PERIOD_DATA_SERVICE = "600"
                        LABEL_NODES               = "false"
                    }
                }

                middlebox = {
                    enabled              = false
                    bandwidth_alert_mbps = 650
                }

                common = {
                    deployment_replicas = 1

                    pdb = {
                        max_unavailable = 1
                    }
                }

                longhorn = {
                    replicas = 2
                }

                ingress = {
                    cpu            = "50m"
                    memory         = "200Mi"
                    target_percent = 6000
                }

                istio = {
                    istiod = {
                        replicas   = 2
                        min_cpu    = "50m"
                        max_cpu    = "8001m"
                        min_memory = "100Mi"
                        max_memory = "8000Mi"
                        pdb        = 0
                    }

                    ingress_gateway = {
                        replicas   = 2
                        max_cpu    = "1"
                        max_memory = "1Gi"
                    }

                    sidecar = {
                        min_cpu    = "0m"
                        min_memory = "200Mi"
                    }
                }

                log_splitter = {
                    min_cpu         = "1m"
                    max_cpu         = "1000m"
                    min_memory      = "10Mi"
                    max_memory      = "2000Mi"
                    mem_buffer_size = "128M"
                    per_pod_rate    = 10000
                }

                monitoring = {
                    min_memory = "100Mi"
                    max_memory = "20Gi"

                    kube_state_metrics = {
                        min_memory = "25Mi"
                    }

                    prometheus = {
                        main = {
                            storage = "10Gi"
                        }
                    }
                }

                redis = {
                    min_cpu    = "10m"
                    max_cpu    = "2001m"
                    min_memory = "100Mi"
                    max_memory = "1000Mi"
                    storage    = "8Gi"
                }

                redis_ha = {
                    min_cpu    = "50m"
                    max_cpu    = "2001m"
                    min_memory = "100Mi"
                    max_memory = "1000Mi"
                    storage    = 0
                }

                redis_sentinel = {
                    min_cpu    = "10m"
                    max_cpu    = "500m"
                    min_memory = "10Mi"
                    max_memory = "400Mi"
                    storage    = 0
                }

                tempo_agent = {
                    min_cpu    = "0m"
                    min_memory = "10Mi"
                }

                internal_dns = {
                    min_cpu    = "0m"
                    max_cpu    = "500m"
                    min_memory = "10Mi"
                    max_memory = "400Mi"
                }
            }
        }

        sysbox = true
    }
}
```

## Example Usage - Triton Provider - Load Balancer Gateway

```terraform
resource "cpln_mk8s" "triton" {

    name        = "demo-mk8s-triton-provider"
    description = "demo-mk8s-triton-provider"

    tags = {
        terraform_generated = "true"
        acceptance_test     = "true"
    }

    version = "1.28.4"
	
    firewall {
        source_cidr = "192.168.1.255"
        description = "hello world"
    }

    triton_provider {
        pre_install_script = "#! echo hello world"
        location           = "aws-eu-central-1"
        private_network_id = "6704dae9-00f4-48b5-8bbf-1be538f20587"
        firewall_enabled   = false
        image_id           = "6b98a11c-53a4-4a62-99e7-cf3dcf150ab2"
        
        networking {}

        connection {
            url                     = "https://us-central-1.api.mnx.io"
            account                 = "eric_controlplane.com"
            user                    = "julian_controlplane.com"
            private_key_secret_link = "/org/terraform-test-org/secret/triton"
        }
      
        load_balancer {
            gateway {}
        }

        node_pool {
            name                = "my-triton-node-pool"
            package_id          = "da311341-b42b-45a8-9386-78ede625d0a4"
            override_image_id   = "e2f3f2aa-a833-49e0-94af-7a7e092cdd9e"
            public_network_id   = "5ff1fe03-075b-4e4c-b85b-73de0c452f77"
            min_size            = 0
            max_size            = 0

            private_network_ids = ["6704dae9-00f4-48b5-8bbf-1be538f20587"]

            labels = {
                hello = "world"
            }

            taint {
                key    = "hello"
                value  = "world"
                effect = "NoSchedule"
            }
            
            triton_tags = {
                drink = "water"
            }
        }
        
        autoscaler {
            expander 	  		  = ["most-pods"]
            unneeded_time         = "10m"
            unready_time  		  = "20m"
            utilization_threshold = 0.7
        }
    }

    add_ons {
        dashboard = true

        azure_workload_identity {
            tenant_id = "7f43458a-a34e-4bfa-9e56-e2289e49c4ec"
        }

        aws_workload_identity = true
        local_path_storage    = true

        metrics {
            kube_state    = true
            core_dns      = true
            kubelet       = true
            api_server    = true
            node_exporter = true
            cadvisor      = true

            scrape_annotated {
                interval_seconds   = 30
                include_namespaces = "^elastic"
                exclude_namespaces = "^elastic"
                retain_labels      = "^\\w+$"
            }
        }

        logs {
            audit_enabled      = true
            include_namespaces = "^elastic"
            exclude_namespaces = "^elastic"
        }

        registry_mirror {
            mirror {
                registry = "registry.mycompany.com"
                mirrors  = ["https://mirror1.mycompany.com"]
            }

            mirror {
                registry = "docker.io"
                mirrors  = ["https://us-mirror.gcr.io"]
            }
        }

        nvidia {
            taint_gpu_nodes = true
        }

        azure_acr {
            client_id = "4e25b134-160b-4a9d-b392-13b381ced5ef"
        }

        byok = {
            ignore_updates = false
            location       = "/org/terraform-test-org/location/test-byok"

            config = {
                actuator = {
                    min_cpu    = "50m"
                    max_cpu    = "8001m"
                    min_memory = "200Mi"
                    max_memory = "8000Mi"
                    log_level  = "info"
                    env = {
                        CACHE_PERIOD_DATA_SERVICE = "600"
                        LABEL_NODES               = "false"
                    }
                }

                middlebox = {
                    enabled              = false
                    bandwidth_alert_mbps = 650
                }

                common = {
                    deployment_replicas = 1

                    pdb = {
                        max_unavailable = 1
                    }
                }

                longhorn = {
                    replicas = 2
                }

                ingress = {
                    cpu            = "50m"
                    memory         = "200Mi"
                    target_percent = 6000
                }

                istio = {
                    istiod = {
                        replicas   = 2
                        min_cpu    = "50m"
                        max_cpu    = "8001m"
                        min_memory = "100Mi"
                        max_memory = "8000Mi"
                        pdb        = 0
                    }

                    ingress_gateway = {
                        replicas   = 2
                        max_cpu    = "1"
                        max_memory = "1Gi"
                    }

                    sidecar = {
                        min_cpu    = "0m"
                        min_memory = "200Mi"
                    }
                }

                log_splitter = {
                    min_cpu         = "1m"
                    max_cpu         = "1000m"
                    min_memory      = "10Mi"
                    max_memory      = "2000Mi"
                    mem_buffer_size = "128M"
                    per_pod_rate    = 10000
                }

                monitoring = {
                    min_memory = "100Mi"
                    max_memory = "20Gi"

                    kube_state_metrics = {
                        min_memory = "25Mi"
                    }

                    prometheus = {
                        main = {
                            storage = "10Gi"
                        }
                    }
                }

                redis = {
                    min_cpu    = "10m"
                    max_cpu    = "2001m"
                    min_memory = "100Mi"
                    max_memory = "1000Mi"
                    storage    = "8Gi"
                }

                redis_ha = {
                    min_cpu    = "50m"
                    max_cpu    = "2001m"
                    min_memory = "100Mi"
                    max_memory = "1000Mi"
                    storage    = 0
                }

                redis_sentinel = {
                    min_cpu    = "10m"
                    max_cpu    = "500m"
                    min_memory = "10Mi"
                    max_memory = "400Mi"
                    storage    = 0
                }

                tempo_agent = {
                    min_cpu    = "0m"
                    min_memory = "10Mi"
                }

                internal_dns = {
                    min_cpu    = "0m"
                    max_cpu    = "500m"
                    min_memory = "10Mi"
                    max_memory = "400Mi"
                }
            }
        }

        sysbox = true
    }
}
```

## Example Usage - Triton Provider - Load Balancer Manual

```terraform
resource "cpln_mk8s" "triton" {

    name        = "demo-mk8s-triton-provider"
    description = "demo-mk8s-triton-provider"

    tags = {
        terraform_generated = "true"
        acceptance_test     = "true"
    }

    version = "1.28.4"
	
    firewall {
        source_cidr = "192.168.1.255"
        description = "hello world"
    }

    triton_provider {
        pre_install_script = "#! echo hello world"
        location           = "aws-eu-central-1"
        private_network_id = "6704dae9-00f4-48b5-8bbf-1be538f20587"
        firewall_enabled   = false
        image_id           = "6b98a11c-53a4-4a62-99e7-cf3dcf150ab2"
        
        networking {}

        connection {
            url                     = "https://us-central-1.api.mnx.io"
            account                 = "eric_controlplane.com"
            user                    = "julian_controlplane.com"
            private_key_secret_link = "/org/terraform-test-org/secret/triton"
        }
      
        load_balancer {
            manual {
                package_id          = "df26ba1d-1261-6fc1-b35c-f1b390bc06ff"
                image_id            = "8605a524-0655-43b9-adf1-7d572fe797eb"
                public_network_id   = "5ff1fe03-075b-4e4c-b85b-73de0c452f77"
                private_network_ids = ["6704dae9-00f4-48b5-8bbf-1be538f20587"]
                count               = 1
                cns_internal_domain = "example.com"
                cns_public_domain   = "example.com"

                metadata = {
                    key1 = "value1"
                    key2 = "value2"
                }

                tags = {
                    tag1 = "value1"
                    tag2 = "value2"
                }

                logging {
                  node_port       = 32000
                  external_syslog = "syslog.example.com:514"
                }
            }
        }

        node_pool {
            name                = "my-triton-node-pool"
            package_id          = "da311341-b42b-45a8-9386-78ede625d0a4"
            override_image_id   = "e2f3f2aa-a833-49e0-94af-7a7e092cdd9e"
            public_network_id   = "5ff1fe03-075b-4e4c-b85b-73de0c452f77"
            min_size            = 0
            max_size            = 0

            private_network_ids = ["6704dae9-00f4-48b5-8bbf-1be538f20587"]

            labels = {
                hello = "world"
            }

            taint {
                key    = "hello"
                value  = "world"
                effect = "NoSchedule"
            }
            
            triton_tags = {
                drink = "water"
            }
        }
        
        autoscaler {
            expander 	  		  = ["most-pods"]
            unneeded_time         = "10m"
            unready_time  		  = "20m"
            utilization_threshold = 0.7
        }
    }

    add_ons {
        dashboard = true

        azure_workload_identity {
            tenant_id = "7f43458a-a34e-4bfa-9e56-e2289e49c4ec"
        }

        aws_workload_identity = true
        local_path_storage    = true

        metrics {
            kube_state    = true
            core_dns      = true
            kubelet       = true
            api_server    = true
            node_exporter = true
            cadvisor      = true

            scrape_annotated {
                interval_seconds   = 30
                include_namespaces = "^elastic"
                exclude_namespaces = "^elastic"
                retain_labels      = "^\\w+$"
            }
        }

        logs {
            audit_enabled      = true
            include_namespaces = "^elastic"
            exclude_namespaces = "^elastic"
        }

        registry_mirror {
            mirror {
                registry = "registry.mycompany.com"
                mirrors  = ["https://mirror1.mycompany.com"]
            }

            mirror {
                registry = "docker.io"
                mirrors  = ["https://us-mirror.gcr.io"]
            }
        }

        nvidia {
            taint_gpu_nodes = true
        }

        azure_acr {
            client_id = "4e25b134-160b-4a9d-b392-13b381ced5ef"
        }

        byok = {
            ignore_updates = false
            location       = "/org/terraform-test-org/location/test-byok"

            config = {
                actuator = {
                    min_cpu    = "50m"
                    max_cpu    = "8001m"
                    min_memory = "200Mi"
                    max_memory = "8000Mi"
                    log_level  = "info"
                    env = {
                        CACHE_PERIOD_DATA_SERVICE = "600"
                        LABEL_NODES               = "false"
                    }
                }

                middlebox = {
                    enabled              = false
                    bandwidth_alert_mbps = 650
                }

                common = {
                    deployment_replicas = 1

                    pdb = {
                        max_unavailable = 1
                    }
                }

                longhorn = {
                    replicas = 2
                }

                ingress = {
                    cpu            = "50m"
                    memory         = "200Mi"
                    target_percent = 6000
                }

                istio = {
                    istiod = {
                        replicas   = 2
                        min_cpu    = "50m"
                        max_cpu    = "8001m"
                        min_memory = "100Mi"
                        max_memory = "8000Mi"
                        pdb        = 0
                    }

                    ingress_gateway = {
                        replicas   = 2
                        max_cpu    = "1"
                        max_memory = "1Gi"
                    }

                    sidecar = {
                        min_cpu    = "0m"
                        min_memory = "200Mi"
                    }
                }

                log_splitter = {
                    min_cpu         = "1m"
                    max_cpu         = "1000m"
                    min_memory      = "10Mi"
                    max_memory      = "2000Mi"
                    mem_buffer_size = "128M"
                    per_pod_rate    = 10000
                }

                monitoring = {
                    min_memory = "100Mi"
                    max_memory = "20Gi"

                    kube_state_metrics = {
                        min_memory = "25Mi"
                    }

                    prometheus = {
                        main = {
                            storage = "10Gi"
                        }
                    }
                }

                redis = {
                    min_cpu    = "10m"
                    max_cpu    = "2001m"
                    min_memory = "100Mi"
                    max_memory = "1000Mi"
                    storage    = "8Gi"
                }

                redis_ha = {
                    min_cpu    = "50m"
                    max_cpu    = "2001m"
                    min_memory = "100Mi"
                    max_memory = "1000Mi"
                    storage    = 0
                }

                redis_sentinel = {
                    min_cpu    = "10m"
                    max_cpu    = "500m"
                    min_memory = "10Mi"
                    max_memory = "400Mi"
                    storage    = 0
                }

                tempo_agent = {
                    min_cpu    = "0m"
                    min_memory = "10Mi"
                }

                internal_dns = {
                    min_cpu    = "0m"
                    max_cpu    = "500m"
                    min_memory = "10Mi"
                    max_memory = "400Mi"
                }
            }
        }

        sysbox = true
    }
}
```

## Example Usage - Azure Provider

```terraform
resource "cpln_mk8s" "azure" {

    name        = "demo-mk8s-azure-provider"
    description = "demo-mk8s-azure-provider"

    tags = {
        terraform_generated = "true"
        acceptance_test     = "true"
    }

    version = "1.32.1"

    firewall {
        source_cidr = "192.168.1.255"
        description = "hello world"
    }

    azure_provider {
        location        = "eastus"
        subscription_id = "12345678-1234-1234-1234-123456789abc"
        sdk_secret_link = "/org/my-org/secret/azure"
        resource_group  = "my-resource-group"
        ssh_keys        = ["ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAI... mk8s-key"]
        network_id      = "/subscriptions/12345678-1234-1234-1234-123456789abc/resourceGroups/my-resource-group/providers/Microsoft.Network/virtualNetworks/mk8s-vnet"

        tags = {
            hello = "world"
        }

        pre_install_script = "#! echo hello world"

        networking {}

        image {
            recommended = "ubuntu/jammy-22.04"
        }

        node_pool {
            name      = "my-azure-node-pool"
            size      = "Standard_B2s"
            subnet_id = "/subscriptions/12345678-1234-1234-1234-123456789abc/resourceGroups/my-resource-group/providers/Microsoft.Network/virtualNetworks/mk8s-vnet/subnets/default"

            zones          = [1]
            boot_disk_size = 30
            min_size       = 0
            max_size       = 0

            labels = {
                hello = "world"
            }

            taint {
                key    = "hello"
                value  = "world"
                effect = "NoSchedule"
            }

            override_image {
                reference {
                    publisher = "Canonical"
                    offer     = "0001-com-ubuntu-server-jammy"
                    sku       = "22_04-lts"
                    version   = "latest"
                }
            }
        }

        autoscaler {
            expander              = ["most-pods"]
            unneeded_time         = "10m"
            unready_time          = "20m"
            utilization_threshold = 0.7
        }
    }

    add_ons {
        dashboard = true

        azure_workload_identity {
            tenant_id = "7f43458a-a34e-4bfa-9e56-e2289e49c4ec"
        }

        aws_workload_identity = true
        local_path_storage    = true

        metrics {
            kube_state    = true
            core_dns      = true
            kubelet       = true
            api_server    = true
            node_exporter = true
            cadvisor      = true

            scrape_annotated {
                interval_seconds   = 30
                include_namespaces = "^elastic"
                exclude_namespaces = "^elastic"
                retain_labels      = "^\\w+$"
            }
        }

        logs {
            audit_enabled      = true
            include_namespaces = "^elastic"
            exclude_namespaces = "^elastic"
        }

        registry_mirror {
            mirror {
                registry = "registry.mycompany.com"
                mirrors  = ["https://mirror1.mycompany.com"]
            }

            mirror {
                registry = "docker.io"
                mirrors  = ["https://us-mirror.gcr.io"]
            }
        }

        nvidia {
            taint_gpu_nodes = true
        }

        azure_acr {
            client_id = "4e25b134-160b-4a9d-b392-13b381ced5ef"
        }

        sysbox = true
    }
}
```

## Example Usage - Digital Ocean Provider

```terraform
resource "cpln_mk8s" "digital-ocean-provider" {

    name        = "demo-mk8s-digital-ocean-provider"
    description = "demo-mk8s-digital-ocean-provider"

    tags = {
        terraform_generated = "true"
        acceptance_test     = "true"
    }

    version = "1.34.2"

    firewall {
        source_cidr = "0.0.0.0/0"
        description = "Default allow-all rule"
    }

    digital_ocean_provider {
        region            = "ams3"
        token_secret_link = "/org/my-org/secret/digitalocean"
        vpc_id            = "12345678-1234-1234-1234-123456789abc"
        image             = "almalinux-8-x64"
        ssh_keys          = ["12345678"]

        digital_ocean_tags = ["mk8s-test", "terraform"]

        pre_install_script = "#! echo hello world"

        networking {
            service_network = "10.43.0.0/16"
            pod_network     = "10.42.0.0/16"
        }

        node_pool {
            name           = "my-do-node-pool"
            droplet_size   = "s-1vcpu-1gb-intel"
            override_image = "ubuntu-22-04-x64"
            min_size       = 1
            max_size       = 1

            labels = {
                hello = "world"
            }

            taint {
                key    = "hello"
                value  = "world"
                effect = "NoSchedule"
            }
        }

        autoscaler {
            expander              = ["most-pods"]
            unneeded_time         = "10m"
            unready_time          = "20m"
            utilization_threshold = 0.7
        }
    }

    add_ons {
        dashboard = true

        aws_workload_identity = true
        local_path_storage    = true

        metrics {
            kube_state    = true
            core_dns      = true
            kubelet       = true
            api_server    = true
            node_exporter = true
            cadvisor      = true

            scrape_annotated {
                interval_seconds   = 30
                include_namespaces = "^elastic"
                exclude_namespaces = "^elastic"
                retain_labels      = "^\\w+$"
            }
        }

        logs {
            audit_enabled      = true
            include_namespaces = "^elastic"
            exclude_namespaces = "^elastic"
        }

        nvidia {
            taint_gpu_nodes = true
        }

        sysbox = true
    }
}
```

## Example Usage - GCP Provider

```terraform
resource "cpln_mk8s" "gcp-provider" {

    name        = "demo-mk8s-gcp-provider"
    description = "demo-mk8s-gcp-provider"

    tags = {
        terraform_generated = "true"
        acceptance_test     = "true"
    }

    version = "1.28.4"

    firewall {
        source_cidr = "192.168.1.255"
        description = "hello world"
    }

    gcp_provider {
        project_id         = "coke-267310"
        region             = "us-west1"
        tags               = ["tag1", "tag2", "tag3"]
        network            = "mk8s"
        sa_key_link        = "/org/terraform-test-org/secret/gcp"
        pre_install_script = "#! echo hello world"

        labels = {
            hello = "world"
        }

        metadata = {
            drink = "water"
            eat   = "chicken"
            play  = "basketball"
        }

        networking {}

        image {
            recommended = "ubuntu/jammy-22.04"
        }

        node_pool {
            name                   = "my-gcp-node-pool"
            machine_type           = "n1-standard-2"
            assign_public_ip       = true
            zone                   = "us-west1-a"
            boot_disk_size         = 30
            min_size               = 0
            max_size               = 0
            preemptible            = true
            subnet                 = "mk8s"
            local_persistent_disks = 1

            labels = {
                hello = "world"
            }

            taint {
                key    = "hello"
                value  = "world"
                effect = "NoSchedule"
            }

            override_image {
                family = {
                    project = "ubuntu-os-cloud"
                    family  = "ubuntu-2204-lts"
                }
            }
        }

        autoscaler {
            expander              = ["most-pods"]
            unneeded_time         = "10m"
            unready_time          = "20m"
            utilization_threshold = 0.7
        }
    }

    add_ons {
        dashboard = true

        azure_workload_identity {
            tenant_id = "7f43458a-a34e-4bfa-9e56-e2289e49c4ec"
        }

        aws_workload_identity = true
        local_path_storage    = true

        metrics {
            kube_state    = true
            core_dns      = true
            kubelet       = true
            api_server    = true
            node_exporter = true
            cadvisor      = true

            scrape_annotated {
                interval_seconds   = 30
                include_namespaces = "^elastic"
                exclude_namespaces = "^elastic"
                retain_labels      = "^\\w+$"
            }
        }

        logs {
            audit_enabled      = true
            include_namespaces = "^elastic"
            exclude_namespaces = "^elastic"
        }

        registry_mirror {
            mirror {
                registry = "registry.mycompany.com"
                mirrors  = ["https://mirror1.mycompany.com"]
            }

            mirror {
                registry = "docker.io"
                mirrors  = ["https://us-mirror.gcr.io"]
            }
        }

        nvidia {
            taint_gpu_nodes = true
        }

        azure_acr {
            client_id = "4e25b134-160b-4a9d-b392-13b381ced5ef"
        }

        byok = {
            ignore_updates = false
            location       = "/org/terraform-test-org/location/test-byok"

            config = {
                actuator = {
                    min_cpu    = "50m"
                    max_cpu    = "8001m"
                    min_memory = "200Mi"
                    max_memory = "8000Mi"
                    log_level  = "info"
                    env = {
                        CACHE_PERIOD_DATA_SERVICE = "600"
                        LABEL_NODES               = "false"
                    }
                }

                middlebox = {
                    enabled              = false
                    bandwidth_alert_mbps = 650
                }

                common = {
                    deployment_replicas = 1

                    pdb = {
                        max_unavailable = 1
                    }
                }

                longhorn = {
                    replicas = 2
                }

                ingress = {
                    cpu            = "50m"
                    memory         = "200Mi"
                    target_percent = 6000
                }

                istio = {
                    istiod = {
                        replicas   = 2
                        min_cpu    = "50m"
                        max_cpu    = "8001m"
                        min_memory = "100Mi"
                        max_memory = "8000Mi"
                        pdb        = 0
                    }

                    ingress_gateway = {
                        replicas   = 2
                        max_cpu    = "1"
                        max_memory = "1Gi"
                    }

                    sidecar = {
                        min_cpu    = "0m"
                        min_memory = "200Mi"
                    }
                }

                log_splitter = {
                    min_cpu         = "1m"
                    max_cpu         = "1000m"
                    min_memory      = "10Mi"
                    max_memory      = "2000Mi"
                    mem_buffer_size = "128M"
                    per_pod_rate    = 10000
                }

                monitoring = {
                    min_memory = "100Mi"
                    max_memory = "20Gi"

                    kube_state_metrics = {
                        min_memory = "25Mi"
                    }

                    prometheus = {
                        main = {
                            storage = "10Gi"
                        }
                    }
                }

                redis = {
                    min_cpu    = "10m"
                    max_cpu    = "2001m"
                    min_memory = "100Mi"
                    max_memory = "1000Mi"
                    storage    = "8Gi"
                }

                redis_ha = {
                    min_cpu    = "50m"
                    max_cpu    = "2001m"
                    min_memory = "100Mi"
                    max_memory = "1000Mi"
                    storage    = 0
                }

                redis_sentinel = {
                    min_cpu    = "10m"
                    max_cpu    = "500m"
                    min_memory = "10Mi"
                    max_memory = "400Mi"
                    storage    = 0
                }

                tempo_agent = {
                    min_cpu    = "0m"
                    min_memory = "10Mi"
                }

                internal_dns = {
                    min_cpu    = "0m"
                    max_cpu    = "500m"
                    min_memory = "10Mi"
                    max_memory = "400Mi"
                }
            }
        }

        sysbox = true
    }
}
```

# Example Usage - Digital Ocean Provider

```terraform
resource "cpln_mk8s" "digital-ocean-provider" {

    name        = "demo-mk8s-digital-ocean-provider"
    description = "demo-mk8s-digital-ocean-provider"

    tags = {
        terraform_generated = "true"
        acceptance_test     = "true"
    }

    version = "1.28.4"
	
    firewall {
        source_cidr = "192.168.1.255"
        description = "hello world"
    }

    digital_ocean_provider {
        region             = "ams3"
        pre_install_script = "#! echo hello world"
        tokenSecretLink    = "/org/terraform-test-org/secret/digital-ocean"
        vpc_id             = "vpc-1"
        image              = "debian-11"

        digital_ocean_tags = ["tag1"]
        ssh_keys           = ["key1"]
        extra_ssh_keys     = ["extraKey1"]
        reserved_ips       = ["192.0.2.10"]
        
        networking {}

        node_pool {
            name                = "my-triton-node-pool"
            droplet_size        = "s-1vcpu-512mb-10gb"
            override_image      = "ubuntu-22"
            min_size            = 0
            max_size            = 0

            labels = {
                hello = "world"
            }

            taint {
                key    = "hello"
                value  = "world"
                effect = "NoSchedule"
            }
        }
        
        autoscaler {
            expander 	  		  = ["most-pods"]
            unneeded_time         = "10m"
            unready_time  		  = "20m"
            utilization_threshold = 0.7
        }
    }

    add_ons {
        dashboard = true

        azure_workload_identity {
            tenant_id = "7f43458a-a34e-4bfa-9e56-e2289e49c4ec"
        }

        aws_workload_identity = true
        local_path_storage    = true

        metrics {
            kube_state    = true
            core_dns      = true
            kubelet       = true
            api_server    = true
            node_exporter = true
            cadvisor      = true

            scrape_annotated {
                interval_seconds   = 30
                include_namespaces = "^elastic"
                exclude_namespaces = "^elastic"
                retain_labels      = "^\\w+$"
            }
        }

        logs {
            audit_enabled      = true
            include_namespaces = "^elastic"
            exclude_namespaces = "^elastic"
        }

        registry_mirror {
            mirror {
                registry = "registry.mycompany.com"
                mirrors  = ["https://mirror1.mycompany.com"]
            }

            mirror {
                registry = "docker.io"
                mirrors  = ["https://us-mirror.gcr.io"]
            }
        }

        nvidia {
            taint_gpu_nodes = true
        }

        azure_acr {
            client_id = "4e25b134-160b-4a9d-b392-13b381ced5ef"
        }

        byok = {
            ignore_updates = false
            location       = "/org/terraform-test-org/location/test-byok"

            config = {
                actuator = {
                    min_cpu    = "50m"
                    max_cpu    = "8001m"
                    min_memory = "200Mi"
                    max_memory = "8000Mi"
                    log_level  = "info"
                    env = {
                        CACHE_PERIOD_DATA_SERVICE = "600"
                        LABEL_NODES               = "false"
                    }
                }

                middlebox = {
                    enabled              = false
                    bandwidth_alert_mbps = 650
                }

                common = {
                    deployment_replicas = 1

                    pdb = {
                        max_unavailable = 1
                    }
                }

                longhorn = {
                    replicas = 2
                }

                ingress = {
                    cpu            = "50m"
                    memory         = "200Mi"
                    target_percent = 6000
                }

                istio = {
                    istiod = {
                        replicas   = 2
                        min_cpu    = "50m"
                        max_cpu    = "8001m"
                        min_memory = "100Mi"
                        max_memory = "8000Mi"
                        pdb        = 0
                    }

                    ingress_gateway = {
                        replicas   = 2
                        max_cpu    = "1"
                        max_memory = "1Gi"
                    }

                    sidecar = {
                        min_cpu    = "0m"
                        min_memory = "200Mi"
                    }
                }

                log_splitter = {
                    min_cpu         = "1m"
                    max_cpu         = "1000m"
                    min_memory      = "10Mi"
                    max_memory      = "2000Mi"
                    mem_buffer_size = "128M"
                    per_pod_rate    = 10000
                }

                monitoring = {
                    min_memory = "100Mi"
                    max_memory = "20Gi"

                    kube_state_metrics = {
                        min_memory = "25Mi"
                    }

                    prometheus = {
                        main = {
                            storage = "10Gi"
                        }
                    }
                }

                redis = {
                    min_cpu    = "10m"
                    max_cpu    = "2001m"
                    min_memory = "100Mi"
                    max_memory = "1000Mi"
                    storage    = "8Gi"
                }

                redis_ha = {
                    min_cpu    = "50m"
                    max_cpu    = "2001m"
                    min_memory = "100Mi"
                    max_memory = "1000Mi"
                    storage    = 0
                }

                redis_sentinel = {
                    min_cpu    = "10m"
                    max_cpu    = "500m"
                    min_memory = "10Mi"
                    max_memory = "400Mi"
                    storage    = 0
                }

                tempo_agent = {
                    min_cpu    = "0m"
                    min_memory = "10Mi"
                }

                internal_dns = {
                    min_cpu    = "0m"
                    max_cpu    = "500m"
                    min_memory = "10Mi"
                    max_memory = "400Mi"
                }
            }
        }

        sysbox = true
    }
}
```

## Import Syntax

The `terraform import` command is used to bring existing infrastructure resources, created outside of Terraform, into the Terraform state file, enabling their management through Terraform going forward.

To update a statefile with an existing mk8s resource, execute the following import command:

```terraform
terraform import cpln_mk8s.RESOURCE_NAME MK8S_NAME
```

-> 1. Substitute RESOURCE_NAME with the same string that is defined in the HCL file.<br/>2. Substitute MK8S_NAME with the corresponding mk8s defined in the resource.
