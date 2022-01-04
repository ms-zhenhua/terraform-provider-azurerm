package kubernetesconfiguration_test

import (
    "testing"
    "context"
)

type KubernetesconfigurationSourceControlConfigurationResource struct{}

func TestAccKubernetesconfigurationSourceControlConfiguration_basic(t *testing.T) {
    data := acceptance.BuildTestData(t, "azurerm_kubernetesconfiguration_source_control_configuration", "test")
    r := KubernetesconfigurationSourceControlConfigurationResource{}
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

func TestAccKubernetesconfigurationSourceControlConfiguration_requiresImport(t *testing.T) {
    data := acceptance.BuildTestData(t, "azurerm_kubernetesconfiguration_source_control_configuration", "test")
    r := KubernetesconfigurationSourceControlConfigurationResource{}
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

func TestAccKubernetesconfigurationSourceControlConfiguration_complete(t *testing.T) {
    data := acceptance.BuildTestData(t, "azurerm_kubernetesconfiguration_source_control_configuration", "test")
    r := KubernetesconfigurationSourceControlConfigurationResource{}
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

func TestAccKubernetesconfigurationSourceControlConfiguration_update(t *testing.T) {
    data := acceptance.BuildTestData(t, "azurerm_kubernetesconfiguration_source_control_configuration", "test")
    r := KubernetesconfigurationSourceControlConfigurationResource{}
    data.ResourceTest(t, r, []acceptance.TestStep{
        {
            Config: r.basic(data),
            Check: acceptance.ComposeTestCheckFunc(
                check.That(data.ResourceName).ExistsInAzure(r),
            ),
        },
        data.ImportStep(),
        {
            Config: r.complete(data),
            Check: acceptance.ComposeTestCheckFunc(
                check.That(data.ResourceName).ExistsInAzure(r),
            ),
        },
        data.ImportStep(),
        {
            Config: r.basic(data),
            Check: acceptance.ComposeTestCheckFunc(
                check.That(data.ResourceName).ExistsInAzure(r),
            ),
        },
        data.ImportStep(),
    })
}

func TestAccKubernetesconfigurationSourceControlConfiguration_updateHelmOperatorProperties(t *testing.T) {
    data := acceptance.BuildTestData(t, "azurerm_kubernetesconfiguration_source_control_configuration", "test")
    r := KubernetesconfigurationSourceControlConfigurationResource{}
    data.ResourceTest(t, r, []acceptance.TestStep{
        {
            Config: r.complete(data),
            Check: acceptance.ComposeTestCheckFunc(
                check.That(data.ResourceName).ExistsInAzure(r),
            ),
        },
        data.ImportStep(),
        {
            Config: r.updateHelmOperatorProperties(data),
            Check: acceptance.ComposeTestCheckFunc(
                check.That(data.ResourceName).ExistsInAzure(r),
            ),
        },
        data.ImportStep(),
    })
}

func (r KubernetesconfigurationSourceControlConfigurationResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
    id, err := KubernetesconfigurationSourceControlConfigurationID(state.ID)
    if err != nil {
        return nil, err
    }
    resp, err := clients.Kubernetesconfiguration.SourceControlConfigurationClient.Get(ctx, id.ResourceGroup, id.ClusterRp, id.ClusterResourceName, id.ClusterName, id.SourceControlConfigurationName)
    if err != nil {
        if response.WasNotFound(resp.HttpResponse) {
            return utils.Bool(false), nil
        }
        return nil, fmt.Errorf("retrieving %s: %+v", id, err)
    }
    return utils.Bool(true), nil
}

func (r KubernetesconfigurationSourceControlConfigurationResource) template(data acceptance.TestData) string {
    return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctest-kubernetesconfiguration-%d"
  location = "%s"
}
`, data.RandomInteger, data.Locations.Primary)
}

func (r KubernetesconfigurationSourceControlConfigurationResource) basic(data acceptance.TestData) string {
    template := r.template(data)
    return fmt.Sprintf(`
%s

resource "azurerm_kubernetesconfiguration_source_control_configuration" "test" {
  name = "acctest-kscc-%d"
  resource_group_name = azurerm_resource_group.test.name
  cluster_name = "clusterName1"
  cluster_resource_name = "connectedClusters"
  cluster_rp = "Microsoft.Kubernetes"
}
`, template, data.RandomInteger)
}

func (r KubernetesconfigurationSourceControlConfigurationResource) requiresImport(data acceptance.TestData) string {
    config := r.basic(data)
    return fmt.Sprintf(`
%s

resource "azurerm_kubernetesconfiguration_source_control_configuration" "import" {
  name = azurerm_kubernetesconfiguration_source_control_configuration.test.name
  resource_group_name = azurerm_kubernetesconfiguration_source_control_configuration.test.resource_group_name
  cluster_name = azurerm_kubernetesconfiguration_source_control_configuration.test.cluster_name
  cluster_resource_name = azurerm_kubernetesconfiguration_source_control_configuration.test.cluster_resource_name
  cluster_rp = azurerm_kubernetesconfiguration_source_control_configuration.test.cluster_rp
}
`, config)
}

func (r KubernetesconfigurationSourceControlConfigurationResource) complete(data acceptance.TestData) string {
    template := r.template(data)
    return fmt.Sprintf(`
%s

resource "azurerm_kubernetesconfiguration_source_control_configuration" "test" {
  name = "acctest-kscc-%d"
  resource_group_name = azurerm_resource_group.test.name
  cluster_name = "clusterName1"
  cluster_resource_name = "connectedClusters"
  cluster_rp = "Microsoft.Kubernetes"
  configuration_protected_settings = {
    protectedSetting1Key = "protectedSetting1Value"
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
  operator_type = "Flux"
  repository_url = "git@github.com:k8sdeveloper425/flux-get-started"
  ssh_known_hosts_contents = "c3NoLmRldi5henVyZS5jb20gc3NoLXJzYSBBQUFBQjNOemFDMXljMkVBQUFBREFRQUJBQUFCQVFDN0hyMW9UV3FOcU9sekdKT2ZHSjROYWtWeUl6ZjFyWFlkNGQ3d282akJsa0x2Q0E0b2RCbEwwbURVeVowL1FVZlRUcWV1K3RtMjJnT3N2K1ZyVlRNazZ2d1JVNzVnWS95OXV0NU1iM2JSNUJWNThkS1h5cTlBOVVlQjVDYWtlaG41WmdtNngxbUtvVnlmK0ZGbjI2aVlxWEpSZ3pJWlpjWjVWNmhyRTBRZzM5a1ptNGF6NDhvMEFVYmY2U3A0U0xkdm51TWEyc1ZOd0hCYm9TN0VKa201N1hRUFZVMy9RcHlOTEhiV0Rkend0cmxTK2V6MzBTM0FkWWhMS0VPeEFHOHdlT255cnRMSkFVZW45bVRrb2w4b0lJMWVkZjdtV1diV1ZmMG5CbWx5MjErblpjbUNUSVNRQnRkY3lQYUVubzdmRlFNREQyNi9zMGxmS29iNEt3OEg="
}
`, template, data.RandomInteger)
}

func (r KubernetesconfigurationSourceControlConfigurationResource) updateHelmOperatorProperties(data acceptance.TestData) string {
    template := r.template(data)
    return fmt.Sprintf(`
%s

resource "azurerm_kubernetesconfiguration_source_control_configuration" "test" {
  name = "acctest-kscc-%d"
  resource_group_name = azurerm_resource_group.test.name
  cluster_name = "clusterName1"
  cluster_resource_name = "connectedClusters"
  cluster_rp = "Microsoft.Kubernetes"
  configuration_protected_settings = {
    protectedSetting1Key = "protectedSetting1Value"
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
  operator_type = "Flux"
  repository_url = "git@github.com:k8sdeveloper425/flux-get-started"
  ssh_known_hosts_contents = "c3NoLmRldi5henVyZS5jb20gc3NoLXJzYSBBQUFBQjNOemFDMXljMkVBQUFBREFRQUJBQUFCQVFDN0hyMW9UV3FOcU9sekdKT2ZHSjROYWtWeUl6ZjFyWFlkNGQ3d282akJsa0x2Q0E0b2RCbEwwbURVeVowL1FVZlRUcWV1K3RtMjJnT3N2K1ZyVlRNazZ2d1JVNzVnWS95OXV0NU1iM2JSNUJWNThkS1h5cTlBOVVlQjVDYWtlaG41WmdtNngxbUtvVnlmK0ZGbjI2aVlxWEpSZ3pJWlpjWjVWNmhyRTBRZzM5a1ptNGF6NDhvMEFVYmY2U3A0U0xkdm51TWEyc1ZOd0hCYm9TN0VKa201N1hRUFZVMy9RcHlOTEhiV0Rkend0cmxTK2V6MzBTM0FkWWhMS0VPeEFHOHdlT255cnRMSkFVZW45bVRrb2w4b0lJMWVkZjdtV1diV1ZmMG5CbWx5MjErblpjbUNUSVNRQnRkY3lQYUVubzdmRlFNREQyNi9zMGxmS29iNEt3OEg="
}
`, template, data.RandomInteger)
}
