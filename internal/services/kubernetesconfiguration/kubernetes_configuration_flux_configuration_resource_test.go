package kubernetesconfiguration_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/go-azure-sdk/resource-manager/kubernetesconfiguration/2022-03-01/fluxconfiguration"

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
			Config: r.privateGitRepositoryWithHttpKey(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("git_repository.0.https_key"),
	})
}

func TestAccKubernetesConfigurationFluxConfiguration_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_configuration_flux_configuration", "test")
	r := KubernetesConfigurationFluxConfigurationResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.bucket(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("bucket.0.bucket_secret_key"),
		{
			Config: r.privateGitRepositoryWithHttpKey(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("git_repository.0.https_key"),
	})
}

func TestAccKubernetesConfigurationFluxConfiguration_privateRepositoryWithSshKey(t *testing.T) {
	const FluxUrl = "KUBERNETES_FLUX_CONFIGURATION_SSH_URL"
	const PrivateSshKey = "KUBERNETES_FLUX_CONFIGURATION_SSH_KEY"
	const KnownHosts = "KUBERNETES_FLUX_CONFIGURATION_KNOWN_HOSTS"

	if os.Getenv(FluxUrl) == "" || os.Getenv(PrivateSshKey) == "" || os.Getenv(KnownHosts) == "" {
		t.Skip(fmt.Sprintf("Acceptance test skipped unless env `%s`, `%s` and `%s` set", FluxUrl, PrivateSshKey, KnownHosts))
		return
	}

	data := acceptance.BuildTestData(t, "azurerm_kubernetes_configuration_flux_configuration", "test")
	r := KubernetesConfigurationFluxConfigurationResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.privateRepositoryWithSshKey(data, os.Getenv(FluxUrl), os.Getenv(PrivateSshKey), os.Getenv(KnownHosts)),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("git_repository.0.ssh_private_key"),
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
  name     = "acctest-rg-%[1]d"
  location = "%[2]s"
}

resource "azurerm_kubernetes_cluster" "test" {
  name                = "acctestaks%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  dns_prefix          = "acctestaks%[1]d"

  default_node_pool {
    name       = "default"
    node_count = 1
    vm_size    = "Standard_DS2_v2"
  }

  identity {
    type = "SystemAssigned"
  }
}

resource "azurerm_resource_group_template_deployment" "test" {
  name                = "acctesttemplate-%[1]d"
  resource_group_name = azurerm_resource_group.test.name
  depends_on          = [azurerm_kubernetes_cluster.test]
  deployment_mode     = "Incremental"

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
	  "scope": "Microsoft.ContainerService/managedClusters/acctestaks%[1]d"
    }
  ]
}
TEMPLATE
}

`, data.RandomInteger, data.Locations.Primary)
}

func (r KubernetesConfigurationFluxConfigurationResource) basic(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
				%s

resource "azurerm_kubernetes_configuration_flux_configuration" "test" {
  name                = "acctest-fc-%d"
  resource_group_name = azurerm_resource_group.test.name
  cluster_name        = azurerm_kubernetes_cluster.test.name

  git_repository {
    url = "https://github.com/Azure/arc-k8s-demo"
    repository_ref {
      branch = "main"
    }
  }

  kustomizations {
    name = "kustomization-1"
  }

  depends_on = [azurerm_resource_group_template_deployment.test]
}
`, template, data.RandomInteger)
}

func (r KubernetesConfigurationFluxConfigurationResource) requiresImport(data acceptance.TestData) string {
	config := r.basic(data)
	return fmt.Sprintf(`
			%s

resource "azurerm_kubernetes_configuration_flux_configuration" "import" {
  name                = azurerm_kubernetes_configuration_flux_configuration.test.name
  resource_group_name = azurerm_kubernetes_configuration_flux_configuration.test.resource_group_name
  cluster_name        = azurerm_kubernetes_configuration_flux_configuration.test.cluster_name

  git_repository {
    url = "https://github.com/Azure/arc-k8s-demo"
    repository_ref {
      branch = "main"
    }
  }

  kustomizations {
    name = "kustomization-1"
  }
}
`, config)
}

func (r KubernetesConfigurationFluxConfigurationResource) privateGitRepositoryWithHttpKey(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
				%s

resource "azurerm_kubernetes_configuration_flux_configuration" "test" {
  name                = "acctest-fc-%d"
  resource_group_name = azurerm_resource_group.test.name
  cluster_name        = azurerm_kubernetes_cluster.test.name
  namespace           = "example"
  scope               = "cluster"

  git_repository {
    url                      = "https://github.com/Azure/arc-k8s-demo"
    https_user               = "example"
    https_key                = "example"
    https_ca_cert            = "example"
    sync_interval_in_seconds = 6000
    timeout_in_seconds       = 6000

    repository_ref {
      branch = "main"
    }
  }

  kustomizations {
    name                      = "kustomization-1"
    path                      = "./test"
    timeout_in_seconds        = 6000
    sync_interval_in_seconds  = 6000
    retry_interval_in_seconds = 6000
    force                     = true
    prune                     = true
  }

  kustomizations {
    name       = "kustomization-2"
    depends_on = ["kustomization-1"]
  }

  depends_on = [azurerm_resource_group_template_deployment.test]
}
`, template, data.RandomInteger)
}

func (r KubernetesConfigurationFluxConfigurationResource) bucket(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
			%s

resource "azurerm_kubernetes_configuration_flux_configuration" "test" {
  name                = "acctest-fc-%d"
  resource_group_name = azurerm_resource_group.test.name
  cluster_name        = azurerm_kubernetes_cluster.test.name
  namespace           = "example"
  scope               = "cluster"

  bucket {
    access_key               = "example"
    bucket_secret_key        = "example"
    bucket_name              = "flux"
    sync_interval_in_seconds = 6000
    timeout_in_seconds       = 6000
    url                      = "https://fluxminiotest.az.minio.io"
  }

  kustomizations {
    name = "kustomization-1"
  }
}
`, template, data.RandomInteger)
}

func (r KubernetesConfigurationFluxConfigurationResource) privateRepositoryWithSshKey(data acceptance.TestData, url string, sshKey string, knownHosts string) string {
	template := r.template(data)
	return fmt.Sprintf(`
				%s

resource "azurerm_kubernetes_configuration_flux_configuration" "test" {
  name                = "acctest-fc-%d"
  resource_group_name = azurerm_resource_group.test.name
  cluster_name        = azurerm_kubernetes_cluster.test.name

  git_repository {
    url             = "%s"
    ssh_private_key = file("%s")
    ssh_known_hosts = file("%s")

    repository_ref {
      branch = "main"
    }
  }

  kustomizations {
    name = "kustomization-1"
  }

  depends_on = [azurerm_resource_group_template_deployment.test]
}
`, template, data.RandomInteger, url, sshKey, knownHosts)
}
