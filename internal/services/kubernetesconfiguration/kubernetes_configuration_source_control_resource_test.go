package kubernetesconfiguration_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/services/kubernetesconfiguration/parse"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

const (
	ClusterRp           = "Microsoft.ContainerService"
	ClusterResourceName = "managedClusters"
)

type KubernetesConfigurationSourceControlResource struct{}

func TestAccKubernetesConfigurationSourceControl_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_configuration_source_control", "test")
	r := KubernetesConfigurationSourceControlResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("configuration_protected_settings"),
	})
}

func TestAccKubernetesConfigurationSourceControl_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_configuration_source_control", "test")
	r := KubernetesConfigurationSourceControlResource{}
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

func TestAccKubernetesConfigurationSourceControl_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_configuration_source_control", "test")
	r := KubernetesConfigurationSourceControlResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("configuration_protected_settings"),
	})
}

func TestAccKubernetesConfigurationSourceControl_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_kubernetes_configuration_source_control", "test")
	r := KubernetesConfigurationSourceControlResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("configuration_protected_settings"),
		{
			Config: r.update(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("configuration_protected_settings"),
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("configuration_protected_settings"),
	})
}

func (r KubernetesConfigurationSourceControlResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := parse.KubernetesConfigurationSourceControlID(state.ID)
	if err != nil {
		return nil, err
	}
	resp, err := clients.KubernetesConfiguration.SourceControlClient.Get(ctx, id.ResourceGroup, ClusterRp, ClusterResourceName, id.ClusterName, id.SourceControlName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return utils.Bool(false), nil
		}
		return nil, fmt.Errorf("retrieving %s: %+v", id, err)
	}
	return utils.Bool(true), nil
}

func (r KubernetesConfigurationSourceControlResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

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
`, data.RandomInteger, data.Locations.Primary)
}

func (r KubernetesConfigurationSourceControlResource) basic(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
	%s
	
	resource "azurerm_kubernetes_configuration_source_control" "test" {
	  name = "acctest-kscc-%d"
	  resource_group_name = azurerm_resource_group.test.name
	  cluster_name = azurerm_kubernetes_cluster.test.name
	  repository_url = "basic@github.com:example/flux-get-started"
	}
	`, template, data.RandomInteger)
}

func (r KubernetesConfigurationSourceControlResource) requiresImport(data acceptance.TestData) string {
	config := r.basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_kubernetes_configuration_source_control" "import" {
  name = azurerm_kubernetes_configuration_source_control.test.name
  resource_group_name = azurerm_kubernetes_configuration_source_control.test.resource_group_name
  cluster_name = azurerm_kubernetes_configuration_source_control.test.cluster_name
  repository_url = azurerm_kubernetes_configuration_source_control.test.repository_url
}
`, config)
}

func (r KubernetesConfigurationSourceControlResource) complete(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

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
`, template, data.RandomInteger)
}

func (r KubernetesConfigurationSourceControlResource) update(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_kubernetes_configuration_source_control" "test" {
  name = "acctest-kscc-%d"
  resource_group_name = azurerm_resource_group.test.name
  cluster_name = azurerm_kubernetes_cluster.test.name
  repository_url = "update@github.com:example/flux-get-started"
  configuration_protected_settings = {
    protectedSetting1Key = "cHJvdGVjdGVkU2V0dGluZzJWYWx1ZQ=="
  }
  enable_helm_operator = true
  helm_operator_properties {
    chart_values = "--set git.ssh.secretName=flux-git-deploy --set tillerNamespace=kube-system"
    chart_version = "0.3.0"
  }
  operator_params = "--git-email=xyzgituser@users.srs.github.com"
  ssh_known_hosts_contents = "c3NoLmRldi5henVyZS5jb20gc3NoLXJzYSBBQUFBQjNOemFDMXljMkVBQUFBREFRQUJBQUFCQVFDN0hyMW9UV3FOcU9sekdKT2ZHSjROYWtWeUl6ZjFyWFlkNGQ3d282akJsa0x2Q0E0b2RCbEwwbURVeVowL1FVZlRUcWV1K3RtMjJnT3N2K1ZyVlRNazZ2d1JVNzVnWS95OXV0NU1iM2JSNUJWNThkS1h5cTlBOVVlQjVDYWtlaG41WmdtNngxbUtvVnlmK0ZGbjI2aVlxWEpSZ3pJWlpjWjVWNmhyRTBRZzM5a1ptNGF6NDhvMEFVYmY2U3A0U0xkdm51TWEyc1ZOd0hCYm9TN0VKa201N1hRUFZVMy9RcHlOTEhiV0Rkend0cmxTK2V6MzBTM0FkWWhMS0VPeEFHOHdlT255cnRMSkFVZW45bVRrb2w4b0lJMWVkZjdtV1diV1ZmMG5CbWx5MjErblpjbUNUSVNRQnRkY3lQYUVubzdmRlFNREQyNi9zMGxmS29iNEt3OEg="
}
`, template, data.RandomInteger)
}
