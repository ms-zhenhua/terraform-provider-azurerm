package workloads_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-sdk/resource-manager/workloads/2023-04-01/providerinstances"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type WorkloadsProviderInstanceResource struct{}

func TestAccWorkloadsProviderInstance_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_workloads_provider_instance", "test")
	r := WorkloadsProviderInstanceResource{}
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

func TestAccWorkloadsProviderInstance_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_workloads_provider_instance", "test")
	r := WorkloadsProviderInstanceResource{}
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

func TestAccWorkloadsProviderInstance_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_workloads_provider_instance", "test")
	r := WorkloadsProviderInstanceResource{}
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

func TestAccWorkloadsProviderInstance_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_workloads_provider_instance", "test")
	r := WorkloadsProviderInstanceResource{}
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

func (r WorkloadsProviderInstanceResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := providerinstances.ParseProviderInstanceID(state.ID)
	if err != nil {
		return nil, err
	}

	client := clients.Workloads.ProviderInstancesClient
	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return utils.Bool(false), nil
		}
		return nil, fmt.Errorf("retrieving %s: %+v", id, err)
	}
	return utils.Bool(resp.Model != nil), nil
}

func (r WorkloadsProviderInstanceResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctest-rg-%d"
  location = "%s"
}
resource "azurerm_workloads_provider_instance" "test" {
  name                = "acctest-wpi-%d"
  resource_group_name = azurerm_resource_group.test.name
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (r WorkloadsProviderInstanceResource) basic(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
				%s

resource "azurerm_workloads_provider_instance" "test" {
  name                 = "acctest-wpi-%d"
  workloads_monitor_id = azurerm_workloads_monitor.test.id
}
`, template, data.RandomInteger)
}

func (r WorkloadsProviderInstanceResource) requiresImport(data acceptance.TestData) string {
	config := r.basic(data)
	return fmt.Sprintf(`
			%s

resource "azurerm_workloads_provider_instance" "import" {
  name                 = azurerm_workloads_provider_instance.test.name
  workloads_monitor_id = azurerm_workloads_monitor.test.id
}
`, config)
}

func (r WorkloadsProviderInstanceResource) complete(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
			%s

resource "azurerm_workloads_provider_instance" "test" {
  name                 = "acctest-wpi-%d"
  workloads_monitor_id = azurerm_workloads_monitor.test.id
  identity {
    type         = "UserAssigned"
    identity_ids = []
  }
  provider_settings {
    provider_type = ""
  }

}
`, template, data.RandomInteger)
}

func (r WorkloadsProviderInstanceResource) update(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
			%s

resource "azurerm_workloads_provider_instance" "test" {
  name                 = "acctest-wpi-%d"
  workloads_monitor_id = azurerm_workloads_monitor.test.id
  identity {
    type         = "UserAssigned"
    identity_ids = []
  }
  provider_settings {
    provider_type = ""
  }

}
`, template, data.RandomInteger)
}
