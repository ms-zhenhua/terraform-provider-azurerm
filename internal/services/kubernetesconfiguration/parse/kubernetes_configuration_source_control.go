package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

type KubernetesConfigurationSourceControlId struct {
	SubscriptionId    string
	ResourceGroup     string
	ClusterName       string
	SourceControlName string
}

func NewKubernetesConfigurationSourceControlID(subscriptionId, resourceGroup, managedClusterName, SourceControlName string) KubernetesConfigurationSourceControlId {
	return KubernetesConfigurationSourceControlId{
		SubscriptionId:    subscriptionId,
		ResourceGroup:     resourceGroup,
		ClusterName:       managedClusterName,
		SourceControlName: SourceControlName,
	}
}

func (id KubernetesConfigurationSourceControlId) String() string {
	segments := []string{
		fmt.Sprintf("Source Control Configuration Name %q", id.SourceControlName),
		fmt.Sprintf("Cluster Name %q", id.ClusterName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Kubernetes Configuration Source Control", segmentsStr)
}

func (id KubernetesConfigurationSourceControlId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.ContainerService/managedClusters/%s/providers/Microsoft.KubernetesConfiguration/SourceControls/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.ClusterName, id.SourceControlName)
}

// KubernetesConfigurationSourceControlID parses a KubernetesConfigurationSourceControl ID into an KubernetesConfigurationSourceControlId struct
func KubernetesConfigurationSourceControlID(input string) (*KubernetesConfigurationSourceControlId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := KubernetesConfigurationSourceControlId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
	}

	if resourceId.ClusterName, err = id.PopSegment("managedClusters"); err != nil {
		return nil, err
	}
	if resourceId.SourceControlName, err = id.PopSegment("SourceControls"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
