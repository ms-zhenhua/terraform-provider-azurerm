---
subcategory: "Kubernetes Configuration"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_kubernetes_configuration_source_control"
description: |-
  Manages a Kubernetes Configuration Source Control.
---

# azurerm_kubernetes_configuration_source_control

Manages a Kubernetes Configuration Source Control.

## Example Usage

```hcl
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

resource "azurerm_kubernetes_cluster" "test" {
  name                = "acctestaks%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  dns_prefix          = "acctestaks%[1]d"

  default_node_pool {
    name                   = "default"
    node_count             = 1
    vm_size                = "Standard_DS2_v2"
  }

  identity {
    type = "SystemAssigned"
  }
}

resource "azurerm_kubernetes_configuration_source_control" "test" {
  name = "acctest-kscc-%d"
  resource_group_name = azurerm_resource_group.test.name
  cluster_name = azurerm_kubernetes_cluster.test.name
  repository_url = "example@github.com:example/flux-get-started"
  configuration_protected_settings = {
    protectedSetting1Key = "cHJvdGVjdGVkU2V0dGluZzFWYWx1ZQ=="
  }
  enable_helm_operator = true
  helm_operator_properties {
    chart_values = "--set git.ssh.secretName=flux-git-deploy --set tillerNamespace=kube-system"
    chart_version = "0.3.0"
  }
  operator_instance_name = "SRSGitHubFluxOp-01"
  operator_namespace = "SRS_Namespace"
  operator_params = "--git-email=xyzgituser@users.srs.github.com"
  operator_scope = "namespace"
  ssh_known_hosts_contents = "c3NoLmRldi5henVyZS5jb20gc3NoLXJzYSBBQUFBQjNOemFDMXljMkVBQUFBREFRQUJBQUFCQVFDN0hyMW9UV3FOcU9sekdKT2ZHSjROYWtWeUl6ZjFyWFlkNGQ3d282akJsa0x2Q0E0b2RCbEwwbURVeVowL1FVZlRUcWV1K3RtMjJnT3N2K1ZyVlRNazZ2d1JVNzVnWS95OXV0NU1iM2JSNUJWNThkS1h5cTlBOVVlQjVDYWtlaG41WmdtNngxbUtvVnlmK0ZGbjI2aVlxWEpSZ3pJWlpjWjVWNmhyRTBRZzM5a1ptNGF6NDhvMEFVYmY2U3A0U0xkdm51TWEyc1ZOd0hCYm9TN0VKa201N1hRUFZVMy9RcHlOTEhiV0Rkend0cmxTK2V6MzBTM0FkWWhMS0VPeEFHOHdlT255cnRMSkFVZW45bVRrb2w4b0lJMWVkZjdtV1diV1ZmMG5CbWx5MjErblpjbUNUSVNRQnRkY3lQYUVubzdmRlFNREQyNi9zMGxmS29iNEt3OEg="
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name which should be used for this Kubernetes Configuration Source Control. Changing this forces a new Kubernetes Configuration Source Control to be created.

* `resource_group_name` - (Required) The name of the Resource Group where the Kubernetes Configuration Source Control should exist. Changing this forces a new Kubernetes Configuration Source Control to be created.

* `cluster_name` - (Required) The name of the kubernetes cluster. Changing this forces a new Kubernetes Configuration Source Control to be created.

* `repository_url` - (Required) Url of the Source Control Repository.

---

* `configuration_protected_settings` - (Optional) Name-value pairs of protected configuration settings for the configuration.

* `enable_helm_operator` - (Optional) Specify if the Helm Operator is enabled. The default value is `false`.

* `helm_operator_properties` - (Optional) A `helm_operator_properties` block as defined below.

* `operator_instance_name` - (Optional) Instance name of the operator. Changing this forces a new Kubernetes Configuration Source Control to be created.

* `operator_namespace` - (Optional) The namespace to which this operator is installed to. Changing this forces a new Kubernetes Configuration Source Control to be created.

* `operator_params` - (Optional) Any Parameters for the Operator instance in string format. The default value is `--git-readonly`.

* `operator_scope` - (Optional) Scope at which the operator will be installed. Possible values are `cluster` and `namespace`. The default value is `cluster`. Changing this forces a new Kubernetes Configuration Source Control to be created.

* `ssh_known_hosts_contents` - (Optional) Base64-encoded known_hosts contents containing public SSH keys required to access private Git instances.

---

An `helm_operator_properties` block exports the following:

* `chart_values` - (Optional) Values override for the operator Helm chart.

* `chart_version` - (Optional) Version of the operator Helm chart.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Kubernetes Configuration Source Control.

* `compliance_status` - Compliance Status of the Configuration. A `compliance_status` block as defined below.

* `repository_public_key` - Public Key associated with this Source Control Configuration (either generated within the cluster or provided by the user).

---

An `compliance_status` block exports the following:

* `compliance_state` - The compliance state of the configuration.

* `last_config_applied` - Datetime the configuration was last applied.

* `message` - Message from when the configuration was applied.

* `message_level` - Level of the message.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Kubernetes Configuration Source Control.
* `read` - (Defaults to 5 minutes) Used when retrieving the Kubernetes Configuration Source Control.
* `update` - (Defaults to 30 minutes) Used when updating the Kubernetes Configuration Source Control.
* `delete` - (Defaults to 30 minutes) Used when deleting the Kubernetes Configuration Source Control.

## Import

Kubernetes Configuration Source Controls can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_kubernetes_configuration_source_control.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.ContainerService/managedClusters/cluster1/providers/Microsoft.KubernetesConfiguration/sourceControlConfigurations/SourceControl1
```