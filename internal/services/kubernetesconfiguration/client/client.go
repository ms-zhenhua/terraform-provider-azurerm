package client

import (
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/kubernetesconfiguration/sdk/2022-07-01/fluxconfiguration"
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
