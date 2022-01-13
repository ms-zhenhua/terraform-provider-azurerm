package client

import (
	"github.com/Azure/azure-sdk-for-go/services/kubernetesconfiguration/mgmt/2021-03-01/kubernetesconfiguration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

type Client struct {
	SourceControlClient *kubernetesconfiguration.SourceControlConfigurationsClient
}

func NewClient(o *common.ClientOptions) *Client {
	SourceControlClient := kubernetesconfiguration.NewSourceControlConfigurationsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&SourceControlClient.Client, o.ResourceManagerAuthorizer)

	return &Client{
		SourceControlClient: &SourceControlClient,
	}
}
