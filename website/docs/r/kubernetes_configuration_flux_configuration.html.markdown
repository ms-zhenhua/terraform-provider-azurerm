---
subcategory: "Kubernetes Configuration"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_kubernetes_configuration_flux_configuration"
description: |-
  Manages a Kubernetes Configuration Flux Configuration.
---

# azurerm_kubernetes_configuration_flux_configuration

Manages a Kubernetes Configuration Flux Configuration.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_kubernetes_configuration_flux_configuration" "example" {
  name                    = "example-kcfc"
  resource_group_name     = azurerm_resource_group.example.name
  cluster_rp              = ""
  cluster_resource_name   = ""
  cluster_name            = ""
  compliance_state        = ""
  error_message           = ""
  namespace               = ""
  repository_public_key   = ""
  scope                   = ""
  source_kind             = ""
  source_synced_commit_id = ""
  source_updated_at       = ""
  status_updated_at       = ""
  suspend                 = false
  azure_blob {
    account_key              = ""
    container_name           = ""
    local_auth_ref           = ""
    sas_token                = ""
    sync_interval_in_seconds = 0
    timeout_in_seconds       = 0
    url                      = ""
    managed_identity {
      client_id = ""
    }
    service_principal {
      client_certificate            = ""
      client_certificate_password   = ""
      client_certificate_send_chain = false
      client_id                     = ""
      client_secret                 = ""
      tenant_id                     = ""
    }
  }
  bucket {
    access_key               = ""
    bucket_name              = ""
    insecure                 = false
    local_auth_ref           = ""
    sync_interval_in_seconds = 0
    timeout_in_seconds       = 0
    url                      = ""
  }
  git_repository {
    https_ca_cert            = ""
    https_user               = ""
    local_auth_ref           = ""
    ssh_known_hosts          = ""
    sync_interval_in_seconds = 0
    timeout_in_seconds       = 0
    url                      = ""
    repository_ref {
      branch = ""
      commit = ""
      semver = ""
      tag    = ""
    }
  }
  statuses {
    compliance_state = ""
    kind             = ""
    name             = ""
    namespace        = ""
    applied_by {
      name      = ""
      namespace = ""
    }
    helm_release_properties {
      failure_count         = 0
      install_failure_count = 0
      last_revision_applied = 0
      upgrade_failure_count = 0
      helm_chart_ref {
        name      = ""
        namespace = ""
      }
    }
    status_conditions {
      last_transition_time = ""
      message              = ""
      reason               = ""
      status               = ""
      type                 = ""
    }
  }
  kustomizations = jsonencode({
    "key" : {}
  })
  configuration_protected_settings = {
    key = ""
  }
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) Specifies the name which should be used for this Kubernetes Configuration Flux Configuration. Changing this forces a new Kubernetes Configuration Flux Configuration to be created.

* `resource_group_name` - (Required) Specifies the name of the Resource Group where the Kubernetes Configuration Flux Configuration should exist. Changing this forces a new Kubernetes Configuration Flux Configuration to be created.

* `cluster_rp` - (Required) Specifies the Cluster Rp. Changing this forces a new Kubernetes Configuration Flux Configuration to be created.

* `cluster_resource_name` - (Required) Specifies the Cluster Resource Name. Changing this forces a new Kubernetes Configuration Flux Configuration to be created.

* `cluster_name` - (Required) Specifies the Cluster Name. Changing this forces a new Kubernetes Configuration Flux Configuration to be created.

* `azure_blob` - (Optional) An `azure_blob` block as defined below.

* `bucket` - (Optional) A `bucket` block as defined below.

* `configuration_protected_settings` - (Optional) Key-value pairs of protected configuration settings for the configuration.

* `git_repository` - (Optional) A `git_repository` block as defined below.

* `kustomizations` - (Optional) Array of kustomizations used to reconcile the artifact pulled by the source type on the cluster.

* `namespace` - (Optional) Specifies the namespace to which this configuration is installed to. Maximum of 253 lower case alphanumeric characters, hyphen and period only. Changing this forces a new Kubernetes Configuration Flux Configuration to be created.

* `scope` - (Optional) Scope at which the operator will be installed. Changing this forces a new Kubernetes Configuration Flux Configuration to be created.

* `source_kind` - (Optional) Source Kind to pull the configuration data from.

* `statuses` - (Optional) A `statuses` block as defined below.

* `suspend` - (Optional) Whether this configuration should suspend its reconciliation of its kustomizations and sources.

---

An `azure_blob` block supports the following:

* `account_key` - (Optional) Specifies the account key (shared key) to access the storage account.

* `container_name` - (Optional) Specifies the Azure Blob container name to sync from the url endpoint for the flux configuration.

* `local_auth_ref` - (Optional) Name of a local secret on the Kubernetes cluster to use as the authentication secret rather than the managed or user-provided configuration secrets.

* `managed_identity` - (Optional) A `managed_identity` block as defined below.

* `sas_token` - (Optional) Specifies the Shared Access token to access the storage container.

* `service_principal` - (Optional) A `service_principal` block as defined below.

* `sync_interval_in_seconds` - (Optional) Specifies the interval at which to re-reconcile the cluster Azure Blob source with the remote.

* `timeout_in_seconds` - (Optional) Specifies the maximum time to attempt to reconcile the cluster Azure Blob source with the remote.

* `url` - (Optional) Specifies the URL to sync for the flux configuration Azure Blob storage account.

A `managed_identity` block supports the following:

* `client_id` - (Optional) Specifies the client Id for authenticating a Managed Identity.

A `service_principal` block supports the following:

* `client_certificate` - (Optional) Base64-encoded certificate used to authenticate a Service Principal .

* `client_certificate_password` - (Optional) Specifies the password for the certificate used to authenticate a Service Principal .

* `client_certificate_send_chain` - (Optional) Specifies whether to include x5c header in client claims when acquiring a token to enable subject name / issuer based authentication for the Client Certificate.

* `client_id` - (Optional) Specifies the client Id for authenticating a Service Principal.

* `client_secret` - (Optional) Specifies the client secret for authenticating a Service Principal.

* `tenant_id` - (Optional) Specifies the tenant Id for authenticating a Service Principal.

---

A `bucket` block supports the following:

* `access_key` - (Optional) Plaintext access key used to securely access the S3 bucket.

* `bucket_name` - (Optional) Specifies the bucket name to sync from the url endpoint for the flux configuration.

* `insecure` - (Optional) Specify whether to use insecure communication when puling data from the S3 bucket.

* `local_auth_ref` - (Optional) Name of a local secret on the Kubernetes cluster to use as the authentication secret rather than the managed or user-provided configuration secrets.

* `sync_interval_in_seconds` - (Optional) Specifies the interval at which to re-reconcile the cluster bucket source with the remote.

* `timeout_in_seconds` - (Optional) Specifies the maximum time to attempt to reconcile the cluster bucket source with the remote.

* `url` - (Optional) Specifies the URL to sync for the flux configuration S3 bucket.

---

A `git_repository` block supports the following:

* `https_ca_cert` - (Optional) Base64-encoded HTTPS certificate authority contents used to access git private git repositories over HTTPS.

* `https_user` - (Optional) Plaintext HTTPS username used to access private git repositories over HTTPS.

* `local_auth_ref` - (Optional) Name of a local secret on the Kubernetes cluster to use as the authentication secret rather than the managed or user-provided configuration secrets.

* `repository_ref` - (Optional) A `repository_ref` block as defined below.

* `ssh_known_hosts` - (Optional) Base64-encoded known_hosts value containing public SSH keys required to access private git repositories over SSH.

* `sync_interval_in_seconds` - (Optional) Specifies the interval at which to re-reconcile the cluster git repository source with the remote.

* `timeout_in_seconds` - (Optional) Specifies the maximum time to attempt to reconcile the cluster git repository source with the remote.

* `url` - (Optional) Specifies the URL to sync for the flux configuration git repository.

A `repository_ref` block supports the following:

* `branch` - (Optional) Specifies the git repository branch name to checkout.

* `commit` - (Optional) Specifies the commit SHA to checkout. This value must be combined with the branch name to be valid. This takes precedence over semver.

* `semver` - (Optional) Specifies the semver range used to match against git repository tags. This takes precedence over tag.

* `tag` - (Optional) Specifies the git repository tag name to checkout. This takes precedence over branch.

---

A `statuses` block supports the following:

* `applied_by` - (Optional) An `applied_by` block as defined below.

* `compliance_state` - (Optional) Compliance state of the applied object showing whether the applied object has come into a ready state on the cluster.

* `helm_release_properties` - (Optional) A `helm_release_properties` block as defined below.

* `kind` - (Optional) Kind of the applied object.

* `name` - (Optional) Name of the applied object.

* `namespace` - (Optional) Namespace of the applied object.

* `status_conditions` - (Optional) A `status_conditions` block as defined below.

An `applied_by` block supports the following:

* `name` - (Optional) Name of the object.

* `namespace` - (Optional) Namespace of the object.

A `helm_release_properties` block supports the following:

* `failure_count` - (Optional) Total number of times that the HelmRelease failed to install or upgrade.

* `helm_chart_ref` - (Optional) A `helm_chart_ref` block as defined below.

* `install_failure_count` - (Optional) Number of times that the HelmRelease failed to install.

* `last_revision_applied` - (Optional) Specifies the revision number of the last released object change.

* `upgrade_failure_count` - (Optional) Number of times that the HelmRelease failed to upgrade.

A `helm_chart_ref` block supports the following:

* `name` - (Optional) Name of the object.

* `namespace` - (Optional) Namespace of the object.

A `status_conditions` block supports the following:

* `last_transition_time` - (Optional) Last time this status condition has changed.

* `message` - (Optional) A more verbose description of the object status condition.

* `reason` - (Optional) Reason for the specified status condition type status.

* `status` - (Optional) Status of the Kubernetes object condition type.

* `type` - (Optional) Object status condition type for this object.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Kubernetes Configuration Flux Configuration.

* `compliance_state` - Combined status of the Flux Kubernetes resources created by the fluxConfiguration or created by the managed objects.

* `error_message` - Error message returned to the user in the case of provisioning failure.

* `repository_public_key` - Public Key associated with this fluxConfiguration (either generated within the cluster or provided by the user).

* `source_synced_commit_id` - Branch and/or SHA of the source commit synced with the cluster.

* `source_updated_at` - Datetime the fluxConfiguration synced its source on the cluster.

* `status_updated_at` - Datetime the fluxConfiguration synced its status on the cluster with Azure.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Kubernetes Configuration Flux Configuration.
* `read` - (Defaults to 5 minutes) Used when retrieving the Kubernetes Configuration Flux Configuration.
* `update` - (Defaults to 30 minutes) Used when updating the Kubernetes Configuration Flux Configuration.
* `delete` - (Defaults to 30 minutes) Used when deleting the Kubernetes Configuration Flux Configuration.

## Import

Kubernetes Configuration Flux Configuration can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_kubernetes_configuration_flux_configuration.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resourceGroup1/providers/{clusterRp}/clusterResource1/cluster1/providers/Microsoft.KubernetesConfiguration/fluxConfigurations/fluxConfiguration1
```
