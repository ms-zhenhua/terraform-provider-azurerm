package kubernetesconfiguration_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-sdk/resource-manager/kubernetesconfiguration/2022-07-01/fluxconfiguration"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type KubernetesConfigurationFluxConfigurationResource struct{}

func TestAccKubernetesConfigurationFluxConfiguration_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_configuration_flux_configuration", "test")
	r := KubernetesConfigurationFluxConfigurationResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccKubernetesConfigurationFluxConfiguration_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_configuration_flux_configuration", "test")
	r := KubernetesConfigurationFluxConfigurationResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.RequiresImportErrorStep(r.requiresImport),
	})
}

func TestAccKubernetesConfigurationFluxConfiguration_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_configuration_flux_configuration", "test")
	r := KubernetesConfigurationFluxConfigurationResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccKubernetesConfigurationFluxConfiguration_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_configuration_flux_configuration", "test")
	r := KubernetesConfigurationFluxConfigurationResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.update(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (r KubernetesConfigurationFluxConfigurationResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := fluxconfiguration.ParseFluxConfigurationID(state.ID)
	if err != nil {
		return nil, err
	}

	client := clients.KubernetesConfiguration.FluxConfigurationClient
	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return utils.Bool(false), nil
		}
		return nil, fmt.Errorf("retrieving %s: %+v", id, err)
	}
	return utils.Bool(resp.Model != nil), nil
}

func (r KubernetesConfigurationFluxConfigurationResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctest-rg-%d"
  location = "%s"
}
`, data.RandomInteger, data.Locations.Primary)
}

func (r KubernetesConfigurationFluxConfigurationResource) basic(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
				%s

resource "azurerm_kubernetes_configuration_flux_configuration" "test" {
  name                  = "acctest-kcfc-%d"
  resource_group_name   = azurerm_resource_group.test.name
  cluster_rp            = ""
  cluster_resource_name = ""
  cluster_name          = ""
}
`, template, data.RandomInteger)
}

func (r KubernetesConfigurationFluxConfigurationResource) requiresImport(data acceptance.TestData) string {
	config := r.basic(data)
	return fmt.Sprintf(`
			%s

resource "azurerm_kubernetes_configuration_flux_configuration" "import" {
  name                  = azurerm_kubernetes_configuration_flux_configuration.test.name
  resource_group_name   = azurerm_resource_group.test.name
  cluster_rp            = ""
  cluster_resource_name = ""
  cluster_name          = ""
}
`, config)
}

func (r KubernetesConfigurationFluxConfigurationResource) complete(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
			%s

resource "azurerm_kubernetes_configuration_flux_configuration" "test" {
  name                    = "acctest-kcfc-%d"
  resource_group_name     = azurerm_resource_group.test.name
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
`, template, data.RandomInteger)
}

func (r KubernetesConfigurationFluxConfigurationResource) update(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
			%s

resource "azurerm_kubernetes_configuration_flux_configuration" "test" {
  name                    = "acctest-kcfc-%d"
  resource_group_name     = azurerm_resource_group.test.name
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
`, template, data.RandomInteger)
}
