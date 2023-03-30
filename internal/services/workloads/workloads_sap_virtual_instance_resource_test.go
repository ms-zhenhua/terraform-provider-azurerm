package workloads_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-sdk/resource-manager/workloads/2023-04-01/sapvirtualinstances"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type WorkloadsSAPVirtualInstanceResource struct{}

func TestAccWorkloadsSAPVirtualInstance_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_workloads_sap_virtual_instance", "test")
	r := WorkloadsSAPVirtualInstanceResource{}
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

func TestAccWorkloadsSAPVirtualInstance_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_workloads_sap_virtual_instance", "test")
	r := WorkloadsSAPVirtualInstanceResource{}
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

func TestAccWorkloadsSAPVirtualInstance_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_workloads_sap_virtual_instance", "test")
	r := WorkloadsSAPVirtualInstanceResource{}
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

func TestAccWorkloadsSAPVirtualInstance_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_workloads_sap_virtual_instance", "test")
	r := WorkloadsSAPVirtualInstanceResource{}
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

func (r WorkloadsSAPVirtualInstanceResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := sapvirtualinstances.ParseSapVirtualInstanceID(state.ID)
	if err != nil {
		return nil, err
	}

	client := clients.Workloads.SAPVirtualInstancesClient
	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return utils.Bool(false), nil
		}
		return nil, fmt.Errorf("retrieving %s: %+v", id, err)
	}
	return utils.Bool(resp.Model != nil), nil
}

func (r WorkloadsSAPVirtualInstanceResource) template(data acceptance.TestData) string {
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

func (r WorkloadsSAPVirtualInstanceResource) basic(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
				%s

resource "azurerm_workloads_sap_virtual_instance" "test" {
  name                = "acctest-wsvi-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = "%s"
  environment         = ""
  sap_product         = ""
  configuration {
    configuration_type = ""
  }
}
`, template, data.RandomInteger, data.Locations.Primary)
}

func (r WorkloadsSAPVirtualInstanceResource) requiresImport(data acceptance.TestData) string {
	config := r.basic(data)
	return fmt.Sprintf(`
			%s

resource "azurerm_workloads_sap_virtual_instance" "import" {
  name                = azurerm_workloads_sap_virtual_instance.test.name
  resource_group_name = azurerm_resource_group.test.name
  location            = "%s"
  environment         = ""
  sap_product         = ""
  configuration {
    configuration_type = ""
  }
}
`, config, data.Locations.Primary)
}

func (r WorkloadsSAPVirtualInstanceResource) complete(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
			%s

resource "azurerm_workloads_sap_virtual_instance" "test" {
  name                = "acctest-wsvi-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = "%s"
  environment         = ""
  sap_product         = ""
  identity {
    type         = "UserAssigned"
    identity_ids = []
  }
  configuration {
    configuration_type = ""
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

func (r WorkloadsSAPVirtualInstanceResource) update(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
			%s

resource "azurerm_workloads_sap_virtual_instance" "test" {
  name                = "acctest-wsvi-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = "%s"
  environment         = ""
  sap_product         = ""
  identity {
    type         = "UserAssigned"
    identity_ids = []
  }
  configuration {
    configuration_type = ""
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
