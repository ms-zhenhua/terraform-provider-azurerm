---
subcategory: "Kubernetesconfiguration"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_kubernetesconfiguration_source_control_configuration"
description: |-
  Manages a kubernetesconfiguration SourceControlConfiguration.
---

# azurerm_kubernetesconfiguration_source_control_configuration

Manages a kubernetesconfiguration SourceControlConfiguration.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-kubernetesconfiguration"
  location = "West Europe"
}

resource "azurerm_kubernetesconfiguration_source_control_configuration" "example" {
  name = "example-sourcecontrolconfiguration"
  resource_group_name = azurerm_resource_group.example.name
  cluster_name = "clusterName1"
  cluster_resource_name = "connectedClusters"
  cluster_rp = "Microsoft.Kubernetes"
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name which should be used for this kubernetesconfiguration SourceControlConfiguration. Changing this forces a new kubernetesconfiguration SourceControlConfiguration to be created.

* `resource_group_name` - (Required) The name of the Resource Group where the kubernetesconfiguration SourceControlConfiguration should exist. Changing this forces a new kubernetesconfiguration SourceControlConfiguration to be created.

* `cluster_name` - (Required) The name of the kubernetes cluster. Changing this forces a new kubernetesconfiguration SourceControlConfiguration to be created.

* `cluster_resource_name` - (Required) The Kubernetes cluster resource name - either managedClusters (for AKS clusters) or connectedClusters (for OnPrem K8S clusters). Possible values are "managedClusters" and "connectedClusters" is allowed. Changing this forces a new kubernetesconfiguration SourceControlConfiguration to be created.

* `cluster_rp` - (Required) The Kubernetes cluster RP - either Microsoft.ContainerService (for AKS clusters) or Microsoft.Kubernetes (for OnPrem K8S clusters). Possible values are "Microsoft.ContainerService" and "Microsoft.Kubernetes" is allowed. Changing this forces a new kubernetesconfiguration SourceControlConfiguration to be created.

---

* `configuration_protected_settings` - (Optional) Name-value pairs of protected configuration settings for the configuration.

* `enable_helm_operator` - (Optional) Option to enable Helm Operator for this git configuration.

* `helm_operator_properties` - (Optional) A `helm_operator_properties` block as defined below.

* `operator_instance_name` - (Optional) Instance name of the operator - identifying the specific configuration.

* `operator_namespace` - (Optional) The namespace to which this operator is installed to. Maximum of 253 lower case alphanumeric characters, hyphen and period only.

* `operator_params` - (Optional) Any Parameters for the Operator instance in string format.

* `operator_scope` - (Optional) Scope at which the operator will be installed. Possible values are "cluster" and "namespace" is allowed.

* `operator_type` - (Optional) Type of the operator. Possible value is &#34;Flux&#34;is allowed.

* `repository_url` - (Optional) Url of the SourceControl Repository.

* `ssh_known_hosts_contents` - (Optional) Base64-encoded known_hosts contents containing public SSH keys required to access private Git instances.

---

An `helm_operator_properties` block exports the following:

* `chart_values` - (Optional) Values override for the operator Helm chart.

* `chart_version` - (Optional) Version of the operator Helm chart.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the kubernetesconfiguration SourceControlConfiguration.

* `compliance_status` - Compliance Status of the Configuration. A `compliance_status` block as defined below.

* `repository_public_key` - Public Key associated with this SourceControl configuration (either generated within the cluster or provided by the user).

* `type` - The type of the resource. E.g. "Microsoft.Compute/virtualMachines" or "Microsoft.Storage/storageAccounts".

---

An `compliance_status` block exports the following:

* `compliance_state` - The compliance state of the configuration.

* `last_config_applied` - Datetime the configuration was last applied.

* `message` - Message from when the configuration was applied.

* `message_level` - Level of the message.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the kubernetesconfiguration SourceControlConfiguration.
* `read` - (Defaults to 5 minutes) Used when retrieving the kubernetesconfiguration SourceControlConfiguration.
* `delete` - (Defaults to 30 minutes) Used when deleting the kubernetesconfiguration SourceControlConfiguration.

## Import

kubernetesconfiguration SourceControlConfigurations can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_kubernetesconfiguration_source_control_configuration.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/{clusterRp}1/clusterResource1/cluster1/providers/Microsoft.KubernetesConfiguration/sourceControlConfigurations/sourceControlConfiguration1
```