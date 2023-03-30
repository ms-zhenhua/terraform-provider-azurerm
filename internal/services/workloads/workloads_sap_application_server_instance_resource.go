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
	"github.com/hashicorp/go-azure-sdk/resource-manager/workloads/2023-04-01/sapapplicationserverinstances"
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

type WorkloadsSAPApplicationServerInstanceModel struct {
	Name                          string                                                 `tfschema:"name"`
	WorkloadsSapVirtualInstanceId string                                                 `tfschema:"workloads_sap_virtual_instance_id"`
	Location                      string                                                 `tfschema:"location"`
	Tags                          map[string]string                                      `tfschema:"tags"`
	Errors                        []SAPVirtualInstanceErrorModel                         `tfschema:"errors"`
	GatewayPort                   int64                                                  `tfschema:"gateway_port"`
	Health                        sapapplicationserverinstances.SAPHealthState           `tfschema:"health"`
	Hostname                      string                                                 `tfschema:"hostname"`
	IPAddress                     string                                                 `tfschema:"ip_address"`
	IcmHTTPPort                   int64                                                  `tfschema:"icm_http_port"`
	IcmHTTPSPort                  int64                                                  `tfschema:"icm_https_port"`
	InstanceNo                    string                                                 `tfschema:"instance_no"`
	KernelPatch                   string                                                 `tfschema:"kernel_patch"`
	KernelVersion                 string                                                 `tfschema:"kernel_version"`
	LoadBalancerDetails           []LoadBalancerDetailsModel                             `tfschema:"load_balancer_details"`
	Status                        sapapplicationserverinstances.SAPVirtualInstanceStatus `tfschema:"status"`
	Subnet                        string                                                 `tfschema:"subnet"`
	VMDetails                     []ApplicationServerVMDetailsModel                      `tfschema:"vm_details"`
}

type SAPVirtualInstanceErrorModel struct {
	Properties []ErrorDefinitionModel `tfschema:"properties"`
}

type ErrorDefinitionModel struct {
	Code    string `tfschema:"code"`
	Message string `tfschema:"message"`
}

type LoadBalancerDetailsModel struct {
	Id string `tfschema:"id"`
}

type ApplicationServerVMDetailsModel struct {
	StorageDetails   []StorageInformationModel                                         `tfschema:"storage_details"`
	Type             sapapplicationserverinstances.ApplicationServerVirtualMachineType `tfschema:"type"`
	VirtualMachineId string                                                            `tfschema:"virtual_machine_id"`
}

type StorageInformationModel struct {
	Id string `tfschema:"id"`
}

type WorkloadsSAPApplicationServerInstanceResource struct{}

var _ sdk.ResourceWithUpdate = WorkloadsSAPApplicationServerInstanceResource{}

func (r WorkloadsSAPApplicationServerInstanceResource) ResourceType() string {
	return "azurerm_workloads_sap_application_server_instance"
}

func (r WorkloadsSAPApplicationServerInstanceResource) ModelObject() interface{} {
	return &WorkloadsSAPApplicationServerInstanceModel{}
}

func (r WorkloadsSAPApplicationServerInstanceResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return sapapplicationserverinstances.ValidateSAPApplicationServerInstanceID
}

func (r WorkloadsSAPApplicationServerInstanceResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"workloads_sap_virtual_instance_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: sapvirtualinstances.ValidateSapVirtualInstanceID,
		},

		"location": commonschema.Location(),

		"tags": commonschema.Tags(),
	}
}

func (r WorkloadsSAPApplicationServerInstanceResource) Attributes() map[string]*pluginsdk.Schema {
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

		"gateway_port": {
			Type:     pluginsdk.TypeInt,
			Computed: true,
		},

		"health": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"hostname": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"ip_address": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"icm_http_port": {
			Type:     pluginsdk.TypeInt,
			Computed: true,
		},

		"icm_https_port": {
			Type:     pluginsdk.TypeInt,
			Computed: true,
		},

		"instance_no": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"kernel_patch": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"kernel_version": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"load_balancer_details": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"id": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},
				},
			},
		},

		"status": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"subnet": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"vm_details": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"storage_details": {
						Type:     pluginsdk.TypeList,
						Computed: true,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"id": {
									Type:     pluginsdk.TypeString,
									Computed: true,
								},
							},
						},
					},

					"type": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},

					"virtual_machine_id": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},
				},
			},
		},
	}
}

func (r WorkloadsSAPApplicationServerInstanceResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var model WorkloadsSAPApplicationServerInstanceModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			client := metadata.Client.Workloads.SAPApplicationServerInstancesClient
			sapVirtualInstanceId, err := sapvirtualinstances.ParseSapVirtualInstanceID(model.WorkloadsSapVirtualInstanceId)
			if err != nil {
				return err
			}

			id := sapapplicationserverinstances.NewApplicationInstanceID(sapVirtualInstanceId.SubscriptionId, sapVirtualInstanceId.ResourceGroupName, sapVirtualInstanceId.SapVirtualInstanceName, model.Name)
			existing, err := client.Get(ctx, id)
			if err != nil && !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for existing %s: %+v", id, err)
			}

			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			properties := &sapapplicationserverinstances.SAPApplicationServerInstance{
				Location:   location.Normalize(model.Location),
				Properties: &sapapplicationserverinstances.SAPApplicationServerProperties{},
				Tags:       &model.Tags,
			}

			if err := client.CreateThenPoll(ctx, id, *properties); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func (r WorkloadsSAPApplicationServerInstanceResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Workloads.SAPApplicationServerInstancesClient

			id, err := sapapplicationserverinstances.ParseApplicationInstanceID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var model WorkloadsSAPApplicationServerInstanceModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			properties := &sapapplicationserverinstances.UpdateSAPApplicationInstanceRequest{}

			if metadata.ResourceData.HasChange("tags") {
				properties.Tags = &model.Tags
			}

			if err := client.UpdateThenPoll(ctx, *id, *properties); err != nil {
				return fmt.Errorf("updating %s: %+v", *id, err)
			}

			return nil
		},
	}
}

func (r WorkloadsSAPApplicationServerInstanceResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Workloads.SAPApplicationServerInstancesClient

			id, err := sapapplicationserverinstances.ParseApplicationInstanceID(metadata.ResourceData.Id())
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

			state := WorkloadsSAPApplicationServerInstanceModel{
				Name:                          id.ApplicationInstanceName,
				WorkloadsSapVirtualInstanceId: sapvirtualinstances.NewSapVirtualInstanceID(id.SubscriptionId, id.ResourceGroupName, id.SapVirtualInstanceName).ID(),
				Location:                      location.Normalize(model.Location),
			}

			if properties := model.Properties; properties != nil {
				state.Errors = flattenSAPVirtualInstanceErrorModel(properties.Errors)

				if properties.GatewayPort != nil {
					state.GatewayPort = *properties.GatewayPort
				}

				if properties.Health != nil {
					state.Health = *properties.Health
				}

				if properties.Hostname != nil {
					state.Hostname = *properties.Hostname
				}

				if properties.IPAddress != nil {
					state.IPAddress = *properties.IPAddress
				}

				if properties.IcmHTTPPort != nil {
					state.IcmHTTPPort = *properties.IcmHTTPPort
				}

				if properties.IcmHTTPSPort != nil {
					state.IcmHTTPSPort = *properties.IcmHTTPSPort
				}

				if properties.InstanceNo != nil {
					state.InstanceNo = *properties.InstanceNo
				}

				if properties.KernelPatch != nil {
					state.KernelPatch = *properties.KernelPatch
				}

				if properties.KernelVersion != nil {
					state.KernelVersion = *properties.KernelVersion
				}

				state.LoadBalancerDetails = flattenLoadBalancerDetailsModel(properties.LoadBalancerDetails)

				if properties.Status != nil {
					state.Status = *properties.Status
				}

				if properties.Subnet != nil {
					state.Subnet = *properties.Subnet
				}

				state.VMDetails = flattenApplicationServerVMDetailsModelArray(properties.VMDetails)
			}
			if model.Tags != nil {
				state.Tags = *model.Tags
			}

			return metadata.Encode(&state)
		},
	}
}

func (r WorkloadsSAPApplicationServerInstanceResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Workloads.SAPApplicationServerInstancesClient

			id, err := sapapplicationserverinstances.ParseApplicationInstanceID(metadata.ResourceData.Id())
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

func flattenSAPVirtualInstanceErrorModel(input *sapapplicationserverinstances.SAPVirtualInstanceError) []SAPVirtualInstanceErrorModel {
	var outputList []SAPVirtualInstanceErrorModel
	if input == nil {
		return outputList
	}
	output := SAPVirtualInstanceErrorModel{
		Properties: flattenErrorDefinitionModel(input.Properties),
	}

	return append(outputList, output)
}

func flattenErrorDefinitionModel(input *sapapplicationserverinstances.ErrorDefinition) []ErrorDefinitionModel {
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

func flattenLoadBalancerDetailsModel(input *sapapplicationserverinstances.LoadBalancerDetails) []LoadBalancerDetailsModel {
	var outputList []LoadBalancerDetailsModel
	if input == nil {
		return outputList
	}
	output := LoadBalancerDetailsModel{}
	if input.Id != nil {
		output.Id = *input.Id
	}

	return append(outputList, output)
}

func flattenApplicationServerVMDetailsModelArray(inputList *[]sapapplicationserverinstances.ApplicationServerVMDetails) []ApplicationServerVMDetailsModel {
	var outputList []ApplicationServerVMDetailsModel
	if inputList == nil {
		return outputList
	}
	for _, input := range *inputList {
		output := ApplicationServerVMDetailsModel{
			StorageDetails: flattenStorageInformationModelArray(input.StorageDetails),
		}

		if input.Type != nil {
			output.Type = *input.Type
		}

		if input.VirtualMachineId != nil {
			output.VirtualMachineId = *input.VirtualMachineId
		}
		outputList = append(outputList, output)
	}
	return outputList
}

func flattenStorageInformationModelArray(inputList *[]sapapplicationserverinstances.StorageInformation) []StorageInformationModel {
	var outputList []StorageInformationModel
	if inputList == nil {
		return outputList
	}
	for _, input := range *inputList {
		output := StorageInformationModel{}

		if input.Id != nil {
			output.Id = *input.Id
		}
		outputList = append(outputList, output)
	}
	return outputList
}
