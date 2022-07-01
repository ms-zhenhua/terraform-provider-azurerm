package client

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/kubernetesconfiguration/2022-03-01/fluxconfiguration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

type Client struct {
	FluxConfigurationClient *fluxconfiguration.FluxConfigurationClient
}

func NewClient(o *common.ClientOptions) *Client {

	fluxConfigurationClient := fluxconfiguration.NewFluxConfigurationClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&fluxConfigurationClient.Client, o.ResourceManagerAuthorizer)

	return &Client{
		FluxConfigurationClient: &fluxConfigurationClient,
	}
}
