package validate

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import "testing"

func TestKubernetesConfigurationSourceControlID(t *testing.T) {
	cases := []struct {
		Input string
		Valid bool
	}{

		{
			// empty
			Input: "",
			Valid: false,
		},

		{
			// missing SubscriptionId
			Input: "/",
			Valid: false,
		},

		{
			// missing value for SubscriptionId
			Input: "/subscriptions/",
			Valid: false,
		},

		{
			// missing ResourceGroup
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/",
			Valid: false,
		},

		{
			// missing value for ResourceGroup
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/",
			Valid: false,
		},

		{
			// missing ClusterResource1Name
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resourceGroup1/providers/clusterRp1/",
			Valid: false,
		},

		{
			// missing value for ClusterResource1Name
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resourceGroup1/providers/clusterRp1/clusterResource1/",
			Valid: false,
		},

		{
			// missing SourceControlConfigurationName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resourceGroup1/providers/clusterRp1/clusterResource1/cluster1/providers/Microsoft.KubernetesConfiguration/",
			Valid: false,
		},

		{
			// missing value for SourceControlConfigurationName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resourceGroup1/providers/clusterRp1/clusterResource1/cluster1/providers/Microsoft.KubernetesConfiguration/sourceControlConfigurations/",
			Valid: false,
		},

		{
			// valid
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resourceGroup1/providers/clusterRp1/clusterResource1/cluster1/providers/Microsoft.KubernetesConfiguration/sourceControlConfigurations/sourceControlConfiguration1",
			Valid: true,
		},

		{
			// upper-cased
			Input: "/SUBSCRIPTIONS/12345678-1234-9876-4563-123456789012/RESOURCEGROUPS/RESOURCEGROUP1/PROVIDERS/CLUSTERRP1/CLUSTERRESOURCE1/CLUSTER1/PROVIDERS/MICROSOFT.KUBERNETESCONFIGURATION/SOURCECONTROLCONFIGURATIONS/SOURCECONTROLCONFIGURATION1",
			Valid: false,
		},
	}
	for _, tc := range cases {
		t.Logf("[DEBUG] Testing Value %s", tc.Input)
		_, errors := KubernetesConfigurationSourceControlID(tc.Input, "test")
		valid := len(errors) == 0

		if tc.Valid != valid {
			t.Fatalf("Expected %t but got %t", tc.Valid, valid)
		}
	}
}
