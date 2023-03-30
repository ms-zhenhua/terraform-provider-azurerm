package workloads_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-sdk/resource-manager/workloads/2023-04-01/saplandscapemonitor"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type WorkloadsSapLandscapeMonitorResource struct{}

func TestAccWorkloadsSapLandscapeMonitor_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_workloads_sap_landscape_monitor", "test")
	r := WorkloadsSapLandscapeMonitorResource{}
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

func TestAccWorkloadsSapLandscapeMonitor_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_workloads_sap_landscape_monitor", "test")
	r := WorkloadsSapLandscapeMonitorResource{}
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

func TestAccWorkloadsSapLandscapeMonitor_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_workloads_sap_landscape_monitor", "test")
	r := WorkloadsSapLandscapeMonitorResource{}
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

func TestAccWorkloadsSapLandscapeMonitor_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_workloads_sap_landscape_monitor", "test")
	r := WorkloadsSapLandscapeMonitorResource{}
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

func (r WorkloadsSapLandscapeMonitorResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := saplandscapemonitor.ParseMonitorID(state.ID)
	if err != nil {
		return nil, err
	}

	client := clients.Workloads.SapLandscapeMonitorClient
	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return utils.Bool(false), nil
		}
		return nil, fmt.Errorf("retrieving %s: %+v", id, err)
	}
	return utils.Bool(resp.Model != nil), nil
}

func (r WorkloadsSapLandscapeMonitorResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctest-rg-%d"
  location = "%s"
}
resource "azurerm_workloads_sap_landscape_monitor" "test" {
  name                = "acctest-wslm-%d"
  resource_group_name = azurerm_resource_group.test.name
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (r WorkloadsSapLandscapeMonitorResource) basic(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
				%s

resource "azurerm_workloads_sap_landscape_monitor" "test" {
  name                 = "acctest-wslm-%d"
  workloads_monitor_id = azurerm_workloads_monitor.test.id
}
`, template, data.RandomInteger)
}

func (r WorkloadsSapLandscapeMonitorResource) requiresImport(data acceptance.TestData) string {
	config := r.basic(data)
	return fmt.Sprintf(`
			%s

resource "azurerm_workloads_sap_landscape_monitor" "import" {
  name                 = azurerm_workloads_sap_landscape_monitor.test.name
  workloads_monitor_id = azurerm_workloads_monitor.test.id
}
`, config)
}

func (r WorkloadsSapLandscapeMonitorResource) complete(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
			%s

resource "azurerm_workloads_sap_landscape_monitor" "test" {
  name                 = "acctest-wslm-%d"
  workloads_monitor_id = azurerm_workloads_monitor.test.id
  grouping {
    landscape {
      name    = ""
      top_sid = []
    }
    sap_application {
      name    = ""
      top_sid = []
    }
  }
  top_metrics_thresholds {
    green  = 0.0
    name   = ""
    red    = 0.0
    yellow = 0.0
  }

}
`, template, data.RandomInteger)
}

func (r WorkloadsSapLandscapeMonitorResource) update(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
			%s

resource "azurerm_workloads_sap_landscape_monitor" "test" {
  name                 = "acctest-wslm-%d"
  workloads_monitor_id = azurerm_workloads_monitor.test.id
  grouping {
    landscape {
      name    = ""
      top_sid = []
    }
    sap_application {
      name    = ""
      top_sid = []
    }
  }
  top_metrics_thresholds {
    green  = 0.0
    name   = ""
    red    = 0.0
    yellow = 0.0
  }

}
`, template, data.RandomInteger)
}
