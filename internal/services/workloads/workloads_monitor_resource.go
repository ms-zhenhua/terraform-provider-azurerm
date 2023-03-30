package workloads

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/workloads/2023-04-01/monitors"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type WorkloadsMonitorModel struct {
	Name                              string                               `tfschema:"name"`
	ResourceGroupName                 string                               `tfschema:"resource_group_name"`
	AppLocation                       string                               `tfschema:"app_location"`
	Location                          string                               `tfschema:"location"`
	LogAnalyticsWorkspaceArmId        string                               `tfschema:"log_analytics_workspace_arm_id"`
	ManagedResourceGroupConfiguration []MonitorManagedRGConfigurationModel `tfschema:"managed_resource_group_configuration"`
	MonitorSubnet                     string                               `tfschema:"monitor_subnet"`
	RoutingPreference                 monitors.RoutingPreference           `tfschema:"routing_preference"`
	Tags                              map[string]string                    `tfschema:"tags"`
	ZoneRedundancyPreference          string                               `tfschema:"zone_redundancy_preference"`
	Errors                            []ErrorModel                         `tfschema:"errors"`
	MsiArmId                          string                               `tfschema:"msi_arm_id"`
	StorageAccountArmId               string                               `tfschema:"storage_account_arm_id"`
}

type MonitorManagedRGConfigurationModel struct {
	Name string `tfschema:"name"`
}

type ErrorModel struct {
	Code       string                 `tfschema:"code"`
	InnerError []ErrorInnerErrorModel `tfschema:"inner_error"`
	Message    string                 `tfschema:"message"`
	Target     string                 `tfschema:"target"`
}

type WorkloadsMonitorResource struct{}

var _ sdk.ResourceWithUpdate = WorkloadsMonitorResource{}

func (r WorkloadsMonitorResource) ResourceType() string {
	return "azurerm_workloads_monitor"
}

func (r WorkloadsMonitorResource) ModelObject() interface{} {
	return &WorkloadsMonitorModel{}
}

func (r WorkloadsMonitorResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return monitors.ValidateMonitorID
}

func (r WorkloadsMonitorResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"app_location": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"identity": commonschema.UserAssignedIdentityOptional(),

		"location": commonschema.Location(),

		"log_analytics_workspace_arm_id": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"managed_resource_group_configuration": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			ForceNew: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"name": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ForceNew:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
				},
			},
		},

		"monitor_subnet": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"routing_preference": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			ForceNew: true,
			ValidateFunc: validation.StringInSlice([]string{
				string(monitors.RoutingPreferenceDefault),
				string(monitors.RoutingPreferenceRouteAll),
			}, false),
		},

		"tags": commonschema.Tags(),

		"zone_redundancy_preference": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},
	}
}

func (r WorkloadsMonitorResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"errors": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"code": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},

					"inner_error": {
						Type:     pluginsdk.TypeList,
						Computed: true,
					},

					"message": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},

					"target": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},
				},
			},
		},

		"msi_arm_id": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"storage_account_arm_id": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
	}
}

func (r WorkloadsMonitorResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var model WorkloadsMonitorModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			client := metadata.Client.Workloads.MonitorsClient
			subscriptionId := metadata.Client.Account.SubscriptionId
			id := monitors.NewMonitorID(subscriptionId, model.ResourceGroupName, model.Name)
			existing, err := client.Get(ctx, id)
			if err != nil && !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for existing %s: %+v", id, err)
			}

			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			identityValue, err := identity.ExpandUserAssignedMap(metadata.ResourceData.Get("identity").([]interface{}))
			if err != nil {
				return fmt.Errorf("expanding `identity`: %+v", err)
			}
			properties := &monitors.Monitor{
				Identity: identityValue,
				Location: location.Normalize(model.Location),
				Properties: &monitors.MonitorProperties{
					ManagedResourceGroupConfiguration: expandManagedRGConfigurationModel(model.ManagedResourceGroupConfiguration),
					RoutingPreference:                 &model.RoutingPreference,
				},
				Tags: &model.Tags,
			}

			if model.AppLocation != "" {
				properties.Properties.AppLocation = &model.AppLocation
			}

			if model.LogAnalyticsWorkspaceArmId != "" {
				properties.Properties.LogAnalyticsWorkspaceArmId = &model.LogAnalyticsWorkspaceArmId
			}

			if model.MonitorSubnet != "" {
				properties.Properties.MonitorSubnet = &model.MonitorSubnet
			}

			if model.ZoneRedundancyPreference != "" {
				properties.Properties.ZoneRedundancyPreference = &model.ZoneRedundancyPreference
			}

			if err := client.CreateThenPoll(ctx, id, *properties); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func (r WorkloadsMonitorResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Workloads.MonitorsClient

			id, err := monitors.ParseMonitorID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var model WorkloadsMonitorModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			properties := &monitors.UpdateMonitorRequest{}

			if metadata.ResourceData.HasChange("identity") {
				identityValue, err := identity.ExpandUserAssignedMap(metadata.ResourceData.Get("identity").([]interface{}))
				if err != nil {
					return fmt.Errorf("expanding `identity`: %+v", err)
				}
				properties.Identity = identityValue
			}

			if metadata.ResourceData.HasChange("tags") {
				properties.Tags = &model.Tags
			}

			if _, err := client.Update(ctx, *id, *properties); err != nil {
				return fmt.Errorf("updating %s: %+v", *id, err)
			}

			return nil
		},
	}
}

func (r WorkloadsMonitorResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Workloads.MonitorsClient

			id, err := monitors.ParseMonitorID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(id)
				}

				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			model := resp.Model
			if model == nil {
				return fmt.Errorf("retrieving %s: model was nil", *id)
			}

			state := WorkloadsMonitorModel{
				Name:              id.MonitorName,
				ResourceGroupName: id.ResourceGroupName,
				Location:          location.Normalize(model.Location),
			}

			identityValue, err := identity.FlattenUserAssignedMap(model.Identity)
			if err != nil {
				return fmt.Errorf("flattening `identity`: %+v", err)
			}

			if err := metadata.ResourceData.Set("identity", identityValue); err != nil {
				return fmt.Errorf("setting `identity`: %+v", err)
			}

			if properties := model.Properties; properties != nil {
				if properties.AppLocation != nil {
					state.AppLocation = *properties.AppLocation
				}

				state.Errors = flattenErrorModel(properties.Errors)

				if properties.LogAnalyticsWorkspaceArmId != nil {
					state.LogAnalyticsWorkspaceArmId = *properties.LogAnalyticsWorkspaceArmId
				}

				state.ManagedResourceGroupConfiguration = flattenManagedRGConfigurationModel(properties.ManagedResourceGroupConfiguration)

				if properties.MonitorSubnet != nil {
					state.MonitorSubnet = *properties.MonitorSubnet
				}

				if properties.MsiArmId != nil {
					state.MsiArmId = *properties.MsiArmId
				}

				if properties.RoutingPreference != nil {
					state.RoutingPreference = *properties.RoutingPreference
				}

				if properties.StorageAccountArmId != nil {
					state.StorageAccountArmId = *properties.StorageAccountArmId
				}

				if properties.ZoneRedundancyPreference != nil {
					state.ZoneRedundancyPreference = *properties.ZoneRedundancyPreference
				}
			}
			if model.Tags != nil {
				state.Tags = *model.Tags
			}

			return metadata.Encode(&state)
		},
	}
}

func (r WorkloadsMonitorResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Workloads.MonitorsClient

			id, err := monitors.ParseMonitorID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			if err := client.DeleteThenPoll(ctx, *id); err != nil {
				return fmt.Errorf("deleting %s: %+v", id, err)
			}

			return nil
		},
	}
}

func expandManagedRGConfigurationModel(inputList []MonitorManagedRGConfigurationModel) *monitors.ManagedRGConfiguration {
	if len(inputList) == 0 {
		return nil
	}
	input := &inputList[0]
	output := monitors.ManagedRGConfiguration{}
	if input.Name != "" {
		output.Name = &input.Name
	}

	return &output
}

func flattenErrorModel(input *monitors.Error) []ErrorModel {
	var outputList []ErrorModel
	if input == nil {
		return outputList
	}
	output := ErrorModel{
		InnerError: flattenErrorInnerErrorModel(input.InnerError),
	}
	if input.Code != nil {
		output.Code = *input.Code
	}

	if input.Message != nil {
		output.Message = *input.Message
	}

	if input.Target != nil {
		output.Target = *input.Target
	}

	return append(outputList, output)
}

func flattenErrorInnerErrorModel(input *monitors.ErrorInnerError) []ErrorInnerErrorModel {
	var outputList []ErrorInnerErrorModel
	if input == nil {
		return outputList
	}
	output := ErrorInnerErrorModel{}

	return append(outputList, output)
}

func flattenManagedRGConfigurationModel(input *monitors.ManagedRGConfiguration) []MonitorManagedRGConfigurationModel {
	var outputList []MonitorManagedRGConfigurationModel
	if input == nil {
		return outputList
	}
	output := MonitorManagedRGConfigurationModel{}
	if input.Name != nil {
		output.Name = *input.Name
	}

	return append(outputList, output)
}
