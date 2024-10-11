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

- **aws_tags** (Map of String) Extra tags to attach to all created objects.
- **pre_install_script** (String) Optional shell script that will be run before K8s is installed. Supports SSM.
- **key_pair** (String) Name of keyPair. Supports SSM
- **disk_encryption_key_arn** (String) KMS key used to encrypt volumes. Supports SSM.
- **security_group_ids** (List of String) Security groups to deploy nodes to. Security groups control if the cluster is multi-zone or single-zon.
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
- **private_key_secret_link** (String) Link to a SSH or opaque secret.

Optional:

- **user** (String)

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
- **nvidia** (Block List, Max: 1) ([see below](#nestedblock--add_ons--nvidia))
- **aws_efs** (Block List, Max: 1) ([see below](#nestedblock--add_ons--aws--efs))
- **aws_ecr** (Block List, Max: 1) ([see below](#nestedblock--add_ons--aws--ecr))
- **aws_elb** (Block List, Max: 1) ([see below](#nestedblock--add_ons--aws--elb))
- **azure_acr** (Block List, Max: 1) ([see below](#nestedblock--add_ons--azure_acr))
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
- **scrape_annotated** (Block List, Max: 1) ([see below](#nestedblock--add_ons--metrics--scrape-annotated))

<a id="nestedblock--add_ons--metrics--scrape-annotated"></a>

### `add_ons.metrics.scrape-annotated`

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
                include_namespaces = "^\\d+$"
                exclude_namespaces  = "^[a-z]$"
                retain_labels      = "^\\w+$"
            }
        }

        logs {
            audit_enabled      = true
            include_namespaces = "^\\d+$"
            exclude_namespaces  = "^[a-z]$"
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

            server_type    = "cx11"
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
                include_namespaces = "^\\d+$"
                exclude_namespaces  = "^[a-z]$"
                retain_labels      = "^\\w+$"
            }
        }

        logs {
            audit_enabled      = true
            include_namespaces = "^\\d+$"
            exclude_namespaces  = "^[a-z]$"
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
                include_namespaces = "^\\d+$"
                exclude_namespaces  = "^[a-z]$"
                retain_labels      = "^\\w+$"
            }
        }

        logs {
            audit_enabled      = true
            include_namespaces = "^\\d+$"
            exclude_namespaces  = "^[a-z]$"
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
                include_namespaces = "^\\d+$"
                exclude_namespaces = "^[a-z]$"
                retain_labels      = "^\\w+$"
            }
        }

        logs {
            audit_enabled      = true
            include_namespaces = "^\\d+$"
            exclude_namespaces = "^[a-z]$"
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
                include_namespaces = "^\\d+$"
                exclude_namespaces  = "^[a-z]$"
                retain_labels      = "^\\w+$"
            }
        }

        logs {
            audit_enabled      = true
            include_namespaces = "^\\d+$"
            exclude_namespaces  = "^[a-z]$"
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
                include_namespaces = "^\\d+$"
                exclude_namespaces  = "^[a-z]$"
                retain_labels      = "^\\w+$"
            }
        }

        logs {
            audit_enabled      = true
            include_namespaces = "^\\d+$"
            exclude_namespaces  = "^[a-z]$"
        }

        nvidia {
            taint_gpu_nodes = true
        }

        azure_acr {
            client_id = "4e25b134-160b-4a9d-b392-13b381ced5ef"
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
        "cpln/ignore"       = "true"
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
                include_namespaces = "^\\d+$"
                exclude_namespaces = "^[a-z]$"
                retain_labels      = "^\\w+$"
            }
        }

        logs {
            audit_enabled      = true
            include_namespaces = "^\\d+$"
            exclude_namespaces = "^[a-z]$"
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
                include_namespaces = "^\\d+$"
                exclude_namespaces  = "^[a-z]$"
                retain_labels      = "^\\w+$"
            }
        }

        logs {
            audit_enabled      = true
            include_namespaces = "^\\d+$"
            exclude_namespaces  = "^[a-z]$"
        }

        nvidia {
            taint_gpu_nodes = true
        }

        azure_acr {
            client_id = "4e25b134-160b-4a9d-b392-13b381ced5ef"
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
                exclude_namespaces = "^[a-z]$"
                retain_labels      = "^\\w+$"
            }
        }

        logs {
            audit_enabled      = true
            include_namespaces = "^\\d+$"
            exclude_namespaces = "^[a-z]$"
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

## Import Syntax

The `terraform import` command is used to bring existing infrastructure resources, created outside of Terraform, into the Terraform state file, enabling their management through Terraform going forward.

To update a statefile with an existing mk8s resource, execute the following import command:

```terraform
terraform import cpln_mk8s.RESOURCE_NAME MK8S_NAME
```

-> 1. Substitute RESOURCE_NAME with the same string that is defined in the HCL file.<br/>2. Substitute MK8S_NAME with the corresponding mk8s defined in the resource.
