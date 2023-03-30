package client

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/workloads/2023-04-01/monitors"
	"github.com/hashicorp/go-azure-sdk/resource-manager/workloads/2023-04-01/providerinstances"
	"github.com/hashicorp/go-azure-sdk/resource-manager/workloads/2023-04-01/sapapplicationserverinstances"
	"github.com/hashicorp/go-azure-sdk/resource-manager/workloads/2023-04-01/sapcentralinstances"
	"github.com/hashicorp/go-azure-sdk/resource-manager/workloads/2023-04-01/sapdatabaseinstances"
	"github.com/hashicorp/go-azure-sdk/resource-manager/workloads/2023-04-01/saplandscapemonitor"
	"github.com/hashicorp/go-azure-sdk/resource-manager/workloads/2023-04-01/sapvirtualinstances"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

type Client struct {
	SAPApplicationServerInstancesClient *sapapplicationserverinstances.SAPApplicationServerInstancesClient
	SAPDatabaseInstancesClient          *sapdatabaseinstances.SAPDatabaseInstancesClient
	ProviderInstancesClient             *providerinstances.ProviderInstancesClient
	SAPVirtualInstancesClient           *sapvirtualinstances.SAPVirtualInstancesClient
	SapLandscapeMonitorClient           *saplandscapemonitor.SapLandscapeMonitorClient
	SAPCentralInstancesClient           *sapcentralinstances.SAPCentralInstancesClient
	MonitorsClient                      *monitors.MonitorsClient
}

func NewClient(o *common.ClientOptions) (*Client, error) {

	sAPApplicationServerInstancesClient, err := sapapplicationserverinstances.NewSAPApplicationServerInstancesClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building SAPApplicationServerInstances client: %+v", err)
	}

	o.Configure(sAPApplicationServerInstancesClient.Client, o.Authorizers.ResourceManager)

	sAPDatabaseInstancesClient, err := sapdatabaseinstances.NewSAPDatabaseInstancesClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building SAPDatabaseInstances client: %+v", err)
	}

	o.Configure(sAPDatabaseInstancesClient.Client, o.Authorizers.ResourceManager)

	providerInstancesClient, err := providerinstances.NewProviderInstancesClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building ProviderInstances client: %+v", err)
	}

	o.Configure(providerInstancesClient.Client, o.Authorizers.ResourceManager)

	sAPVirtualInstancesClient, err := sapvirtualinstances.NewSAPVirtualInstancesClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building SAPVirtualInstances client: %+v", err)
	}

	o.Configure(sAPVirtualInstancesClient.Client, o.Authorizers.ResourceManager)

	sapLandscapeMonitorClient, err := saplandscapemonitor.NewSapLandscapeMonitorClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building SapLandscapeMonitor client: %+v", err)
	}

	o.Configure(sapLandscapeMonitorClient.Client, o.Authorizers.ResourceManager)

	sAPCentralInstancesClient, err := sapcentralinstances.NewSAPCentralInstancesClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building SAPCentralInstances client: %+v", err)
	}

	o.Configure(sAPCentralInstancesClient.Client, o.Authorizers.ResourceManager)

	monitorsClient, err := monitors.NewMonitorsClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Monitors client: %+v", err)
	}

	o.Configure(monitorsClient.Client, o.Authorizers.ResourceManager)

	return &Client{
		SAPApplicationServerInstancesClient: sAPApplicationServerInstancesClient,
		SAPDatabaseInstancesClient:          sAPDatabaseInstancesClient,
		ProviderInstancesClient:             providerInstancesClient,
		SAPVirtualInstancesClient:           sAPVirtualInstancesClient,
		SapLandscapeMonitorClient:           sapLandscapeMonitorClient,
		SAPCentralInstancesClient:           sAPCentralInstancesClient,
		MonitorsClient:                      monitorsClient,
	}, nil
}
