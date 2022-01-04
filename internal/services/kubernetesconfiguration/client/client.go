package client

import (
	"github.com/Azure/azure-sdk-for-go/services/kubernetesconfiguration/mgmt/2021-03-01/kubernetesconfiguration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

type Client struct {
	SourceControlClient *kubernetesconfiguration.SourceControlConfigurationsClient
}

func NewClient(o *common.ClientOptions) *Client {
	sourceControlConfigurationClient := kubernetesconfiguration.NewSourceControlConfigurationsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&sourceControlConfigurationClient.Client, o.ResourceManagerAuthorizer)

	return &Client{
		SourceControlClient: &sourceControlConfigurationClient,
	}
}
