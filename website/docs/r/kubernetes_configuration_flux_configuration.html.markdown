---
subcategory: "Kubernetes Configuration"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_kubernetes_configuration_flux_configuration"
description: |-
  Manages a Kubernetes Flux Configuration.
---

# azurerm_kubernetes_configuration_flux_configuration

Manages a Kubernetes Flux Configuration.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources2"
  location = "West Europe"
}

resource "azurerm_kubernetes_cluster" "example" {
  name                = "exampleaks"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  dns_prefix          = "exampleaks"

  default_node_pool {
    name       = "default"
    node_count = 1
    vm_size    = "Standard_DS2_v2"
  }

  identity {
    type = "SystemAssigned"
  }
}

resource "azurerm_resource_group_template_deployment" "example" {
  name                = "example-template"
  resource_group_name = azurerm_resource_group.example.name
  deployment_mode     = "Incremental"
  depends_on          = [azurerm_kubernetes_cluster.example]

  template_content = <<TEMPLATE
{
  "$schema": "https://schema.management.azure.com/schemas/2015-01-01/deploymentTemplate.json#",
  "contentVersion": "1.0.0.0",
  "resources": [
     {
      "type": "Microsoft.KubernetesConfiguration/extensions",
      "apiVersion": "2021-09-01",
      "name": "flux",
      "properties": {
        "extensionType": "microsoft.flux",
        "autoUpgradeMinorVersion": true
      },
	  "scope": "Microsoft.ContainerService/managedClusters/exampleaks"
    }
  ]
}
TEMPLATE
}

resource "azurerm_kubernetes_configuration_flux_configuration" "example" {
  name                = "example-fc"
  resource_group_name = azurerm_resource_group.example.name
  cluster_name        = azurerm_kubernetes_cluster.example.name

  git_repository {
    url = "https://github.com/Azure/arc-k8s-demo"
    repository_ref {
      branch = "main"
    }
  }

  kustomizations {
    name = "kustomization-1"
  }

  depends_on = [azurerm_resource_group_template_deployment.example]
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) Specifies the name which should be used for this Kubernetes Flux Configuration. It must be between 1 and 30 characters. It can contain only lowercase letters, numbers, and hyphens (-). It must start and end with a lowercase letter or number. Changing this forces a new Kubernetes Flux Configuration to be created.

* `resource_group_name` - (Required) Specifies the name of the Resource Group where the Kubernetes Flux Configuration should exist. Changing this forces a new Kubernetes Flux Configuration to be created.

* `cluster_name` - (Required) Specifies the Cluster Name. Changing this forces a new Kubernetes Flux Configuration to be created.

* `kustomizations` - (Required) A `kustomizations` block as defined below.

* `cluster_resource_name` - (Optional) Specifies the Cluster Resource Name. The only possible value is `managedClusters`. Defaults to `managedClusters`. Changing this forces a new Kubernetes Flux Configuration to be created.

* `bucket` - (Optional) A `bucket` block as defined below.

* `git_repository` - (Optional) A `git_repository` block as defined below.

* `namespace` - (Optional) Specifies the namespace to which this configuration is installed to. It must be between 1 and 63 characters. It can contain only lowercase letters, numbers, and hyphens (-). It must start and end with a lowercase letter or number. Changing this forces a new Kubernetes Flux Configuration to be created.

* `scope` - (Optional) Scope at which the operator will be installed. Defaults to `namespace`. Changing this forces a new Kubernetes Flux Configuration to be created.

* `suspend` - (Optional) Whether this configuration should suspend its reconciliation of its kustomizations and sources. Defaults to `false`.

---

A `kustomizations` block supports the following:

* `name` - (Required) Specifies the name of the Kustomization. It must be between 1 and 30 characters. It can contain only lowercase letters, numbers, and hyphens (-). It must start and end with a lowercase letter or number.

* `path` - (Optional) Specifies the path in the source reference to reconcile on the cluster.

* `timeout_in_seconds` - (Optional) The maximum time to attempt to reconcile the Kustomization on the cluster. Defaults to `600`.

* `sync_interval_in_seconds` - (Optional) The interval at which to re-reconcile the Kustomization on the cluster. Defaults to `600`.

* `retry_interval_in_seconds` - (Optional) The interval at which to re-reconcile the Kustomization on the cluster in the event of failure on reconciliation. Defaults to `600`.

* `force` - (Optional) Whether enable re-creating Kubernetes resources on the cluster when patching fails due to an immutable field change. Defaults to `false`.

* `prune` - (Optional) Whether enable garbage collections of Kubernetes objects created by this Kustomization. Defaults to `false`.

* `depends_on` - (Optional) Specifies other Kustomizations that this Kustomization depends on. This Kustomization will not reconcile until all dependencies have completed their reconciliation.

---

A `bucket` block supports the following:

* `bucket_name` - (Required) Specifies the bucket name to sync from the url endpoint for the flux configuration. It must be between 1 and 63 characters. It can contain only lowercase letters, numbers, and hyphens (-). It must start and end with a lowercase letter or number.

* `url` - (Required) Specifies the URL to sync for the flux configuration S3 bucket. It must start with `http://` or `https://`.

* `access_key` - (Optional) Plaintext access key used to securely access the S3 bucket.

* `bucket_secret_key` - (Optional) Secret key used to authenticate with the bucket source.

* `insecure` - (Optional) Specify whether to use insecure communication when puling data from the S3 bucket. Defaults to `false`.

* `local_auth_ref` - (Optional) Name of a local secret on the Kubernetes cluster to use as the authentication secret rather than the managed or user-provided configuration secrets. It must be between 1 and 63 characters. It can contain only lowercase letters, numbers, and hyphens (-). It must start and end with a lowercase letter or number.

* `sync_interval_in_seconds` - (Optional) Specifies the interval at which to re-reconcile the cluster git repository source with the remote. Defaults to `600`.

* `timeout_in_seconds` - (Optional) Specifies the maximum time to attempt to reconcile the cluster git repository source with the remote. Defaults to `600`.

---

A `git_repository` block supports the following:

* `url` - (Required) Specifies the URL to sync for the flux configuration git repository. It must start with `http://`, `https://`, `git@`, or `ssh://`.

* `repository_ref` - (Required) A `repository_ref` block as defined below.

* `https_ca_cert` - (Optional) Plaintext HTTPS certificate authority contents used to access git private git repositories over HTTPS.

* `https_user` - (Optional) Plaintext HTTPS username used to access private git repositories over HTTPS.

* `https_key` - (Optional) Plaintext HTTPS personal access token or password that will be used to access the repository.

* `local_auth_ref` - (Optional) Name of a local secret on the Kubernetes cluster to use as the authentication secret rather than the managed or user-provided configuration secrets. It must be between 1 and 63 characters. It can contain only lowercase letters, numbers, and hyphens (-). It must start and end with a lowercase letter or number.

* `ssh_private_key` - (Optional) Plaintext SSH private key in PEM format.

* `ssh_known_hosts` - (Optional) Plaintext known_hosts value containing public SSH keys required to access private git repositories over SSH.

* `sync_interval_in_seconds` - (Optional) Specifies the interval at which to re-reconcile the cluster git repository source with the remote. Defaults to `600`.

* `timeout_in_seconds` - (Optional) Specifies the maximum time to attempt to reconcile the cluster git repository source with the remote. Defaults to `600`.

---

A `repository_ref` block supports the following:

* `branch` - (Optional) Specifies the git repository branch name to checkout.

* `commit` - (Optional) Specifies the commit SHA to checkout. This value must be combined with the branch name to be valid. This takes precedence over semver.

* `semver` - (Optional) Specifies the semver range used to match against git repository tags. This takes precedence over tag.

* `tag` - (Optional) Specifies the git repository tag name to checkout. This takes precedence over branch.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Kubernetes Flux Configuration.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Kubernetes Flux Configuration.
* `read` - (Defaults to 5 minutes) Used when retrieving the Kubernetes Flux Configuration.
* `update` - (Defaults to 30 minutes) Used when updating the Kubernetes Flux Configuration.
* `delete` - (Defaults to 30 minutes) Used when deleting the Kubernetes Flux Configuration.

## Import

Kubernetes Flux Configuration can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_kubernetes_configuration_flux_configuration.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resourceGroup1/providers/{clusterRp}/{clusterResourceName}/cluster1/providers/Microsoft.KubernetesConfiguration/fluxConfigurations/fluxConfiguration1
```
