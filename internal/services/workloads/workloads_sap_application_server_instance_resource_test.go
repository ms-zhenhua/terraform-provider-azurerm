package workloads_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-sdk/resource-manager/workloads/2023-04-01/sapapplicationserverinstances"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type WorkloadsSAPApplicationServerInstanceResource struct{}

func TestAccWorkloadsSAPApplicationServerInstance_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_workloads_sap_application_server_instance", "test")
	r := WorkloadsSAPApplicationServerInstanceResource{}
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

func TestAccWorkloadsSAPApplicationServerInstance_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_workloads_sap_application_server_instance", "test")
	r := WorkloadsSAPApplicationServerInstanceResource{}
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

func TestAccWorkloadsSAPApplicationServerInstance_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_workloads_sap_application_server_instance", "test")
	r := WorkloadsSAPApplicationServerInstanceResource{}
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

func TestAccWorkloadsSAPApplicationServerInstance_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_workloads_sap_application_server_instance", "test")
	r := WorkloadsSAPApplicationServerInstanceResource{}
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

func (r WorkloadsSAPApplicationServerInstanceResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := sapapplicationserverinstances.ParseApplicationInstanceID(state.ID)
	if err != nil {
		return nil, err
	}

	client := clients.Workloads.SAPApplicationServerInstancesClient
	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return utils.Bool(false), nil
		}
		return nil, fmt.Errorf("retrieving %s: %+v", id, err)
	}
	return utils.Bool(resp.Model != nil), nil
}

func (r WorkloadsSAPApplicationServerInstanceResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctest-rg-%d"
  location = "%s"
}
resource "azurerm_workloads_sap_application_server_instance" "test" {
  name                = "acctest-wsasi-%d"
  resource_group_name = azurerm_resource_group.test.name
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (r WorkloadsSAPApplicationServerInstanceResource) basic(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
				%s

resource "azurerm_workloads_sap_application_server_instance" "test" {
  name                              = "acctest-wsasi-%d"
  workloads_sap_virtual_instance_id = azurerm_workloads_sap_virtual_instance.test.id
  location                          = "%s"
}
`, template, data.RandomInteger, data.Locations.Primary)
}

func (r WorkloadsSAPApplicationServerInstanceResource) requiresImport(data acceptance.TestData) string {
	config := r.basic(data)
	return fmt.Sprintf(`
			%s

resource "azurerm_workloads_sap_application_server_instance" "import" {
  name                              = azurerm_workloads_sap_application_server_instance.test.name
  workloads_sap_virtual_instance_id = azurerm_workloads_sap_virtual_instance.test.id
  location                          = "%s"
}
`, config, data.Locations.Primary)
}

func (r WorkloadsSAPApplicationServerInstanceResource) complete(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
			%s

resource "azurerm_workloads_sap_application_server_instance" "test" {
  name                              = "acctest-wsasi-%d"
  workloads_sap_virtual_instance_id = azurerm_workloads_sap_virtual_instance.test.id
  location                          = "%s"
  tags = {
    key = "value"
  }

}
`, template, data.RandomInteger, data.Locations.Primary)
}

func (r WorkloadsSAPApplicationServerInstanceResource) update(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
			%s

resource "azurerm_workloads_sap_application_server_instance" "test" {
  name                              = "acctest-wsasi-%d"
  workloads_sap_virtual_instance_id = azurerm_workloads_sap_virtual_instance.test.id
  location                          = "%s"
  tags = {
    key = "value"
  }

}
`, template, data.RandomInteger, data.Locations.Primary)
}
