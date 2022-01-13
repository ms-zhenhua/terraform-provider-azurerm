package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/resourceid"
)

var _ resourceid.Formatter = KubernetesConfigurationSourceControlId{}

func TestKubernetesConfigurationSourceControlIDFormatter(t *testing.T) {
	actual := NewKubernetesConfigurationSourceControlID("12345678-1234-9876-4563-123456789012", "resourceGroup1", "cluster1", "SourceControl1").ID()
	expected := "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resourceGroup1/providers/Microsoft.ContainerService/managedClusters/cluster1/providers/Microsoft.KubernetesConfiguration/SourceControls/SourceControl1"
	if actual != expected {
		t.Fatalf("Expected %q but got %q", expected, actual)
	}
}

func TestKubernetesConfigurationSourceControlID(t *testing.T) {
	testData := []struct {
		Input    string
		Error    bool
		Expected *KubernetesConfigurationSourceControlId
	}{

		{
			// empty
			Input: "",
			Error: true,
		},

		{
			// missing SubscriptionId
			Input: "/",
			Error: true,
		},

		{
			// missing value for SubscriptionId
			Input: "/subscriptions/",
			Error: true,
		},

		{
			// missing ResourceGroup
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/",
			Error: true,
		},

		{
			// missing value for ResourceGroup
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/",
			Error: true,
		},

		{
			// missing ManagedClusterName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resourceGroup1/providers/Microsoft.ContainerService/",
			Error: true,
		},

		{
			// missing value for ManagedClusterName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resourceGroup1/providers/Microsoft.ContainerService/managedClusters/",
			Error: true,
		},

		{
			// missing SourceControlName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resourceGroup1/providers/Microsoft.ContainerService/managedClusters/cluster1/providers/Microsoft.KubernetesConfiguration/",
			Error: true,
		},

		{
			// missing value for SourceControlName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resourceGroup1/providers/Microsoft.ContainerService/managedClusters/cluster1/providers/Microsoft.KubernetesConfiguration/SourceControls/",
			Error: true,
		},

		{
			// valid
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resourceGroup1/providers/Microsoft.ContainerService/managedClusters/cluster1/providers/Microsoft.KubernetesConfiguration/SourceControls/SourceControl1",
			Expected: &KubernetesConfigurationSourceControlId{
				SubscriptionId:    "12345678-1234-9876-4563-123456789012",
				ResourceGroup:     "resourceGroup1",
				ClusterName:       "cluster1",
				SourceControlName: "SourceControl1",
			},
		},

		{
			// upper-cased
			Input: "/SUBSCRIPTIONS/12345678-1234-9876-4563-123456789012/RESOURCEGROUPS/RESOURCEGROUP1/PROVIDERS/MICROSOFT.CONTAINERSERVICE/MANAGEDCLUSTERS/CLUSTER1/PROVIDERS/MICROSOFT.KUBERNETESCONFIGURATION/SourceControlS/SourceControl1",
			Error: true,
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q", v.Input)

		actual, err := KubernetesConfigurationSourceControlID(v.Input)
		if err != nil {
			if v.Error {
				continue
			}

			t.Fatalf("Expect a value but got an error: %s", err)
		}
		if v.Error {
			t.Fatal("Expect an error but didn't get one")
		}

		if actual.SubscriptionId != v.Expected.SubscriptionId {
			t.Fatalf("Expected %q but got %q for SubscriptionId", v.Expected.SubscriptionId, actual.SubscriptionId)
		}
		if actual.ResourceGroup != v.Expected.ResourceGroup {
			t.Fatalf("Expected %q but got %q for ResourceGroup", v.Expected.ResourceGroup, actual.ResourceGroup)
		}
		if actual.ClusterName != v.Expected.ClusterName {
			t.Fatalf("Expected %q but got %q for ManagedClusterName", v.Expected.ClusterName, actual.ClusterName)
		}
		if actual.SourceControlName != v.Expected.SourceControlName {
			t.Fatalf("Expected %q but got %q for SourceControlName", v.Expected.SourceControlName, actual.SourceControlName)
		}
	}
}
