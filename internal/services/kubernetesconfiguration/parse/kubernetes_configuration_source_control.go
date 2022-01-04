package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

type KubernetesConfigurationSourceControlId struct {
	SubscriptionId                 string
	ResourceGroup                  string
	ClusterRp                      string
	ClusterResourceName            string
	ClusterName                    string
	SourceControlConfigurationName string
}

func NewKubernetesConfigurationSourceControlID(subscriptionId, resourceGroup, clusterRp, clusterResourceName, clusterName, sourceControlConfigurationName string) KubernetesConfigurationSourceControlId {
	return KubernetesConfigurationSourceControlId{
		SubscriptionId:                 subscriptionId,
		ResourceGroup:                  resourceGroup,
		ClusterRp:                      clusterRp,
		ClusterResourceName:            clusterResourceName,
		ClusterName:                    clusterName,
		SourceControlConfigurationName: sourceControlConfigurationName,
	}
}

func (id KubernetesConfigurationSourceControlId) String() string {
	segments := []string{
		fmt.Sprintf("Source Control Configuration Name %q", id.SourceControlConfigurationName),
		fmt.Sprintf("Cluster Name %q", id.ClusterName),
		fmt.Sprintf("Cluster Resource Name %q", id.ClusterResourceName),
		fmt.Sprintf("Cluster RP %q", id.ClusterRp),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Kubernetes Configuration Source Control Configuration", segmentsStr)
}

func (id KubernetesConfigurationSourceControlId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/%s/%s/%s/providers/Microsoft.KubernetesConfiguration/sourceControlConfigurations/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, "Microsoft.ContainerService", "managedClusters", id.ClusterName, id.SourceControlConfigurationName)
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

	if resourceId.ClusterName, err = id.PopSegment("clusterResource1"); err != nil {
		return nil, err
	}
	if resourceId.SourceControlConfigurationName, err = id.PopSegment("sourceControlConfigurations"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
