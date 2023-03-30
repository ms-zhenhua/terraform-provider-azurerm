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
	"github.com/hashicorp/go-azure-sdk/resource-manager/workloads/2023-04-01/sapvirtualinstances"
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

type WorkloadsSAPVirtualInstanceModel struct {
	Name                              string                                       `tfschema:"name"`
	ResourceGroupName                 string                                       `tfschema:"resource_group_name"`
	Configuration                     []SAPConfigurationModel                      `tfschema:"configuration"`
	Environment                       sapvirtualinstances.SAPEnvironmentType       `tfschema:"environment"`
	Location                          string                                       `tfschema:"location"`
	ManagedResourceGroupConfiguration []ManagedRGConfigurationModel                `tfschema:"managed_resource_group_configuration"`
	SapProduct                        sapvirtualinstances.SAPProductType           `tfschema:"sap_product"`
	Tags                              map[string]string                            `tfschema:"tags"`
	Errors                            []SAPVirtualInstanceErrorModel               `tfschema:"errors"`
	Health                            sapvirtualinstances.SAPHealthState           `tfschema:"health"`
	State                             sapvirtualinstances.SAPVirtualInstanceState  `tfschema:"state"`
	Status                            sapvirtualinstances.SAPVirtualInstanceStatus `tfschema:"status"`
}

type SAPConfigurationModel struct {
	ConfigurationType sapvirtualinstances.SAPConfigurationType `tfschema:"configuration_type"`
}

type ManagedRGConfigurationModel struct {
	Name string `tfschema:"name"`
}

type SAPVirtualInstanceErrorModel struct {
	Properties []ErrorDefinitionModel `tfschema:"properties"`
}

type ErrorDefinitionModel struct {
	Code    string `tfschema:"code"`
	Message string `tfschema:"message"`
}

type WorkloadsSAPVirtualInstanceResource struct{}

var _ sdk.ResourceWithUpdate = WorkloadsSAPVirtualInstanceResource{}

func (r WorkloadsSAPVirtualInstanceResource) ResourceType() string {
	return "azurerm_workloads_sap_virtual_instance"
}

func (r WorkloadsSAPVirtualInstanceResource) ModelObject() interface{} {
	return &WorkloadsSAPVirtualInstanceModel{}
}

func (r WorkloadsSAPVirtualInstanceResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return sapvirtualinstances.ValidateSAPVirtualInstanceID
}

func (r WorkloadsSAPVirtualInstanceResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"configuration": {
			Type:     pluginsdk.TypeList,
			Required: true,
			ForceNew: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"configuration_type": {
						Type:     pluginsdk.TypeString,
						Required: true,
						ForceNew: true,
						ValidateFunc: validation.StringInSlice([]string{
							string(sapvirtualinstances.SAPConfigurationTypeDiscovery),
							string(sapvirtualinstances.SAPConfigurationTypeDeploymentWithOSConfig),
							string(sapvirtualinstances.SAPConfigurationTypeDeployment),
						}, false),
					},
				},
			},
		},

		"environment": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
			ValidateFunc: validation.StringInSlice([]string{
				string(sapvirtualinstances.SAPEnvironmentTypeNonProd),
				string(sapvirtualinstances.SAPEnvironmentTypeProd),
			}, false),
		},

		"identity": commonschema.UserAssignedIdentityOptional(),

		"location": commonschema.Location(),

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

		"sap_product": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
			ValidateFunc: validation.StringInSlice([]string{
				string(sapvirtualinstances.SAPProductTypeECC),
				string(sapvirtualinstances.SAPProductTypeSFourHANA),
				string(sapvirtualinstances.SAPProductTypeOther),
			}, false),
		},

		"tags": commonschema.Tags(),
	}
}

func (r WorkloadsSAPVirtualInstanceResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"errors": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"properties": {
						Type:     pluginsdk.TypeList,
						Computed: true,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"code": {
									Type:     pluginsdk.TypeString,
									Computed: true,
								},

								"message": {
									Type:     pluginsdk.TypeString,
									Computed: true,
								},
							},
						},
					},
				},
			},
		},

		"health": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"state": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"status": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
	}
}

func (r WorkloadsSAPVirtualInstanceResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var model WorkloadsSAPVirtualInstanceModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			client := metadata.Client.Workloads.SAPVirtualInstancesClient
			subscriptionId := metadata.Client.Account.SubscriptionId
			id := sapvirtualinstances.NewSapVirtualInstanceID(subscriptionId, model.ResourceGroupName, model.Name)
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
			properties := &sapvirtualinstances.SAPVirtualInstance{
				Identity: identityValue,
				Location: location.Normalize(model.Location),
				Properties: sapvirtualinstances.SAPVirtualInstanceProperties{
					Environment:                       model.Environment,
					ManagedResourceGroupConfiguration: expandManagedRGConfigurationModel(model.ManagedResourceGroupConfiguration),
					SapProduct:                        model.SapProduct,
				},
				Tags: &model.Tags,
			}

			configurationValue := expandSAPConfigurationModel(model.Configuration)
			if configurationValue != nil {
				properties.Properties.Configuration = *configurationValue
			}

			if err := client.CreateThenPoll(ctx, id, *properties); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func (r WorkloadsSAPVirtualInstanceResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Workloads.SAPVirtualInstancesClient

			id, err := sapvirtualinstances.ParseSapVirtualInstanceID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var model WorkloadsSAPVirtualInstanceModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			properties := &sapvirtualinstances.UpdateSAPVirtualInstanceRequest{}

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

func (r WorkloadsSAPVirtualInstanceResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Workloads.SAPVirtualInstancesClient

			id, err := sapvirtualinstances.ParseSapVirtualInstanceID(metadata.ResourceData.Id())
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

			state := WorkloadsSAPVirtualInstanceModel{
				Name:              id.SapVirtualInstanceName,
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

			properties := &model.Properties
			state.Configuration = flattenSAPConfigurationModel(&properties.Configuration)

			state.Environment = properties.Environment

			state.Errors = flattenSAPVirtualInstanceErrorModel(properties.Errors)

			if properties.Health != nil {
				state.Health = *properties.Health
			}

			state.ManagedResourceGroupConfiguration = flattenManagedRGConfigurationModel(properties.ManagedResourceGroupConfiguration)

			state.SapProduct = properties.SapProduct

			if properties.State != nil {
				state.State = *properties.State
			}

			if properties.Status != nil {
				state.Status = *properties.Status
			}
			if model.Tags != nil {
				state.Tags = *model.Tags
			}

			return metadata.Encode(&state)
		},
	}
}

func (r WorkloadsSAPVirtualInstanceResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Workloads.SAPVirtualInstancesClient

			id, err := sapvirtualinstances.ParseSapVirtualInstanceID(metadata.ResourceData.Id())
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

func expandSAPConfigurationModel(inputList []SAPConfigurationModel) *sapvirtualinstances.SAPConfiguration {
	if len(inputList) == 0 {
		return nil
	}
	input := &inputList[0]
	output := sapvirtualinstances.SAPConfiguration{
		ConfigurationType: input.ConfigurationType,
	}

	return &output
}

func expandManagedRGConfigurationModel(inputList []ManagedRGConfigurationModel) *sapvirtualinstances.ManagedRGConfiguration {
	if len(inputList) == 0 {
		return nil
	}
	input := &inputList[0]
	output := sapvirtualinstances.ManagedRGConfiguration{}
	if input.Name != "" {
		output.Name = &input.Name
	}

	return &output
}

func flattenSAPConfigurationModel(input *sapvirtualinstances.SAPConfiguration) []SAPConfigurationModel {
	var outputList []SAPConfigurationModel
	if input == nil {
		return outputList
	}
	output := SAPConfigurationModel{
		ConfigurationType: input.ConfigurationType,
	}

	return append(outputList, output)
}

func flattenSAPVirtualInstanceErrorModel(input *sapvirtualinstances.SAPVirtualInstanceError) []SAPVirtualInstanceErrorModel {
	var outputList []SAPVirtualInstanceErrorModel
	if input == nil {
		return outputList
	}
	output := SAPVirtualInstanceErrorModel{
		Properties: flattenErrorDefinitionModel(input.Properties),
	}

	return append(outputList, output)
}

func flattenErrorDefinitionModel(input *sapvirtualinstances.ErrorDefinition) []ErrorDefinitionModel {
	var outputList []ErrorDefinitionModel
	if input == nil {
		return outputList
	}
	output := ErrorDefinitionModel{}
	if input.Code != nil {
		output.Code = *input.Code
	}

	if input.Message != nil {
		output.Message = *input.Message
	}

	return append(outputList, output)
}

func flattenManagedRGConfigurationModel(input *sapvirtualinstances.ManagedRGConfiguration) []ManagedRGConfigurationModel {
	var outputList []ManagedRGConfigurationModel
	if input == nil {
		return outputList
	}
	output := ManagedRGConfigurationModel{}
	if input.Name != nil {
		output.Name = *input.Name
	}

	return append(outputList, output)
}
