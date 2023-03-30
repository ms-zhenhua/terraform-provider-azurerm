package workloads

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	tagsHelper "github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/workloads/2023-04-01/monitors"
	"github.com/hashicorp/go-azure-sdk/resource-manager/workloads/2023-04-01/providerinstances"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	azValidate "github.com/hashicorp/terraform-provider-azurerm/helpers/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type WorkloadsProviderInstanceModel struct {
	Name               string                            `tfschema:"name"`
	WorkloadsMonitorId string                            `tfschema:"workloads_monitor_id"`
	ProviderSettings   []ProviderSpecificPropertiesModel `tfschema:"provider_settings"`
	Errors             []ErrorModel                      `tfschema:"errors"`
}

type ProviderSpecificPropertiesModel struct {
	ProviderType string `tfschema:"provider_type"`
}

type ErrorModel struct {
	Code       string                 `tfschema:"code"`
	InnerError []ErrorInnerErrorModel `tfschema:"inner_error"`
	Message    string                 `tfschema:"message"`
	Target     string                 `tfschema:"target"`
}

type WorkloadsProviderInstanceResource struct{}

var _ sdk.ResourceWithUpdate = WorkloadsProviderInstanceResource{}

func (r WorkloadsProviderInstanceResource) ResourceType() string {
	return "azurerm_workloads_provider_instance"
}

func (r WorkloadsProviderInstanceResource) ModelObject() interface{} {
	return &WorkloadsProviderInstanceModel{}
}

func (r WorkloadsProviderInstanceResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return providerinstances.ValidateProviderInstanceID
}

func (r WorkloadsProviderInstanceResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"workloads_monitor_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: monitors.ValidateMonitorID,
		},

		"identity": commonschema.UserAssignedIdentityOptional(),

		"provider_settings": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"provider_type": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
				},
			},
		},
	}
}

func (r WorkloadsProviderInstanceResource) Attributes() map[string]*pluginsdk.Schema {
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
	}
}

func (r WorkloadsProviderInstanceResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var model WorkloadsProviderInstanceModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			client := metadata.Client.Workloads.ProviderInstancesClient
			monitorId, err := monitors.ParseMonitorID(model.WorkloadsMonitorId)
			if err != nil {
				return err
			}

			id := providerinstances.NewProviderInstanceID(monitorId.SubscriptionId, monitorId.ResourceGroupName, monitorId.MonitorName, model.Name)
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
			properties := &providerinstances.ProviderInstance{
				Identity: identityValue,
				Properties: &providerinstances.ProviderInstanceProperties{
					ProviderSettings: expandProviderSpecificPropertiesModel(model.ProviderSettings),
				},
			}

			if err := client.CreateThenPoll(ctx, id, *properties); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func (r WorkloadsProviderInstanceResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Workloads.ProviderInstancesClient

			id, err := providerinstances.ParseProviderInstanceID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var model WorkloadsProviderInstanceModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			resp, err := client.Get(ctx, *id)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			properties := resp.Model
			if properties == nil {
				return fmt.Errorf("retrieving %s: properties was nil", id)
			}

			if metadata.ResourceData.HasChange("identity") {
				identityValue, err := identity.ExpandUserAssignedMap(metadata.ResourceData.Get("identity").([]interface{}))
				if err != nil {
					return fmt.Errorf("expanding `identity`: %+v", err)
				}
				properties.Identity = identityValue
			}

			if metadata.ResourceData.HasChange("provider_settings") {
				properties.Properties.ProviderSettings = expandProviderSpecificPropertiesModel(model.ProviderSettings)
			}

			properties.SystemData = nil

			if err := client.CreateThenPoll(ctx, *id, *properties); err != nil {
				return fmt.Errorf("updating %s: %+v", *id, err)
			}

			return nil
		},
	}
}

func (r WorkloadsProviderInstanceResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Workloads.ProviderInstancesClient

			id, err := providerinstances.ParseProviderInstanceID(metadata.ResourceData.Id())
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

			state := WorkloadsProviderInstanceModel{
				Name:               id.ProviderInstanceName,
				WorkloadsMonitorId: monitors.NewMonitorID(id.SubscriptionId, id.ResourceGroupName, id.MonitorName).ID(),
			}

			identityValue, err := identity.FlattenUserAssignedMap(model.Identity)
			if err != nil {
				return fmt.Errorf("flattening `identity`: %+v", err)
			}

			if err := metadata.ResourceData.Set("identity", identityValue); err != nil {
				return fmt.Errorf("setting `identity`: %+v", err)
			}

			if properties := model.Properties; properties != nil {
				state.Errors = flattenErrorModel(properties.Errors)

				state.ProviderSettings = flattenProviderSpecificPropertiesModel(properties.ProviderSettings)
			}

			return metadata.Encode(&state)
		},
	}
}

func (r WorkloadsProviderInstanceResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Workloads.ProviderInstancesClient

			id, err := providerinstances.ParseProviderInstanceID(metadata.ResourceData.Id())
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

func expandProviderSpecificPropertiesModel(inputList []ProviderSpecificPropertiesModel) *providerinstances.ProviderSpecificProperties {
	if len(inputList) == 0 {
		return nil
	}
	input := &inputList[0]
	output := providerinstances.ProviderSpecificProperties{
		ProviderType: input.ProviderType,
	}

	return &output
}

func flattenErrorModel(input *providerinstances.Error) []ErrorModel {
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

func flattenErrorInnerErrorModel(input *providerinstances.ErrorInnerError) []ErrorInnerErrorModel {
	var outputList []ErrorInnerErrorModel
	if input == nil {
		return outputList
	}
	output := ErrorInnerErrorModel{}

	return append(outputList, output)
}

func flattenProviderSpecificPropertiesModel(input *providerinstances.ProviderSpecificProperties) []ProviderSpecificPropertiesModel {
	var outputList []ProviderSpecificPropertiesModel
	if input == nil {
		return outputList
	}
	output := ProviderSpecificPropertiesModel{
		ProviderType: input.ProviderType,
	}

	return append(outputList, output)
}
