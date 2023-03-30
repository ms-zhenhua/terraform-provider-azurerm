package workloads_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-sdk/resource-manager/workloads/2023-04-01/monitors"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type WorkloadsMonitorResource struct{}

func TestAccWorkloadsMonitor_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_workloads_monitor", "test")
	r := WorkloadsMonitorResource{}
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

func TestAccWorkloadsMonitor_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_workloads_monitor", "test")
	r := WorkloadsMonitorResource{}
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

func TestAccWorkloadsMonitor_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_workloads_monitor", "test")
	r := WorkloadsMonitorResource{}
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

func TestAccWorkloadsMonitor_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_workloads_monitor", "test")
	r := WorkloadsMonitorResource{}
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

func (r WorkloadsMonitorResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := monitors.ParseMonitorID(state.ID)
	if err != nil {
		return nil, err
	}

	client := clients.Workloads.MonitorsClient
	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return utils.Bool(false), nil
		}
		return nil, fmt.Errorf("retrieving %s: %+v", id, err)
	}
	return utils.Bool(resp.Model != nil), nil
}

func (r WorkloadsMonitorResource) template(data acceptance.TestData) string {
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

func (r WorkloadsMonitorResource) basic(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
				%s

resource "azurerm_workloads_monitor" "test" {
  name                = "acctest-wm-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = "%s"
}
`, template, data.RandomInteger, data.Locations.Primary)
}

func (r WorkloadsMonitorResource) requiresImport(data acceptance.TestData) string {
	config := r.basic(data)
	return fmt.Sprintf(`
			%s

resource "azurerm_workloads_monitor" "import" {
  name                = azurerm_workloads_monitor.test.name
  resource_group_name = azurerm_resource_group.test.name
  location            = "%s"
}
`, config, data.Locations.Primary)
}

func (r WorkloadsMonitorResource) complete(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
			%s

resource "azurerm_workloads_monitor" "test" {
  name                           = "acctest-wm-%d"
  resource_group_name            = azurerm_resource_group.test.name
  location                       = "%s"
  app_location                   = ""
  log_analytics_workspace_arm_id = ""
  monitor_subnet                 = ""
  routing_preference             = ""
  zone_redundancy_preference     = ""
  identity {
    type         = "UserAssigned"
    identity_ids = []
  }
  managed_resource_group_configuration {
    name = ""
  }
  tags = {
    key = "value"
  }

}
`, template, data.RandomInteger, data.Locations.Primary)
}

func (r WorkloadsMonitorResource) update(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
			%s

resource "azurerm_workloads_monitor" "test" {
  name                           = "acctest-wm-%d"
  resource_group_name            = azurerm_resource_group.test.name
  location                       = "%s"
  app_location                   = ""
  log_analytics_workspace_arm_id = ""
  monitor_subnet                 = ""
  routing_preference             = ""
  zone_redundancy_preference     = ""
  identity {
    type         = "UserAssigned"
    identity_ids = []
  }
  managed_resource_group_configuration {
    name = ""
  }
  tags = {
    key = "value"
  }

}
`, template, data.RandomInteger, data.Locations.Primary)
}
