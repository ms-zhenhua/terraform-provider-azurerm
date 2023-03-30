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
	"github.com/hashicorp/go-azure-sdk/resource-manager/workloads/2023-04-01/sapcentralinstances"
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

type WorkloadsSAPCentralInstanceModel struct {
	Name                               string                                       `tfschema:"name"`
	WorkloadsSapVirtualInstanceId      string                                       `tfschema:"workloads_sap_virtual_instance_id"`
	Location                           string                                       `tfschema:"location"`
	Tags                               map[string]string                            `tfschema:"tags"`
	EnqueueReplicationServerProperties []EnqueueReplicationServerPropertiesModel    `tfschema:"enqueue_replication_server_properties"`
	EnqueueServerProperties            []EnqueueServerPropertiesModel               `tfschema:"enqueue_server_properties"`
	Errors                             []SAPVirtualInstanceErrorModel               `tfschema:"errors"`
	GatewayServerProperties            []GatewayServerPropertiesModel               `tfschema:"gateway_server_properties"`
	Health                             sapcentralinstances.SAPHealthState           `tfschema:"health"`
	InstanceNo                         string                                       `tfschema:"instance_no"`
	KernelPatch                        string                                       `tfschema:"kernel_patch"`
	KernelVersion                      string                                       `tfschema:"kernel_version"`
	LoadBalancerDetails                []LoadBalancerDetailsModel                   `tfschema:"load_balancer_details"`
	MessageServerProperties            []MessageServerPropertiesModel               `tfschema:"message_server_properties"`
	Status                             sapcentralinstances.SAPVirtualInstanceStatus `tfschema:"status"`
	Subnet                             string                                       `tfschema:"subnet"`
	VMDetails                          []CentralServerVMDetailsModel                `tfschema:"vm_details"`
}

type EnqueueReplicationServerPropertiesModel struct {
	ErsVersion    sapcentralinstances.EnqueueReplicationServerType `tfschema:"ers_version"`
	Health        sapcentralinstances.SAPHealthState               `tfschema:"health"`
	Hostname      string                                           `tfschema:"hostname"`
	IPAddress     string                                           `tfschema:"ip_address"`
	InstanceNo    string                                           `tfschema:"instance_no"`
	KernelPatch   string                                           `tfschema:"kernel_patch"`
	KernelVersion string                                           `tfschema:"kernel_version"`
}

type EnqueueServerPropertiesModel struct {
	Health    sapcentralinstances.SAPHealthState `tfschema:"health"`
	Hostname  string                             `tfschema:"hostname"`
	IPAddress string                             `tfschema:"ip_address"`
	Port      int64                              `tfschema:"port"`
}

type SAPVirtualInstanceErrorModel struct {
	Properties []ErrorDefinitionModel `tfschema:"properties"`
}

type ErrorDefinitionModel struct {
	Code    string `tfschema:"code"`
	Message string `tfschema:"message"`
}

type GatewayServerPropertiesModel struct {
	Health sapcentralinstances.SAPHealthState `tfschema:"health"`
	Port   int64                              `tfschema:"port"`
}

type LoadBalancerDetailsModel struct {
	Id string `tfschema:"id"`
}

type MessageServerPropertiesModel struct {
	HTTPPort       int64                              `tfschema:"http_port"`
	HTTPSPort      int64                              `tfschema:"https_port"`
	Health         sapcentralinstances.SAPHealthState `tfschema:"health"`
	Hostname       string                             `tfschema:"hostname"`
	IPAddress      string                             `tfschema:"ip_address"`
	InternalMsPort int64                              `tfschema:"internal_ms_port"`
	MsPort         int64                              `tfschema:"ms_port"`
}

type CentralServerVMDetailsModel struct {
	StorageDetails   []StorageInformationModel                           `tfschema:"storage_details"`
	Type             sapcentralinstances.CentralServerVirtualMachineType `tfschema:"type"`
	VirtualMachineId string                                              `tfschema:"virtual_machine_id"`
}

type StorageInformationModel struct {
	Id string `tfschema:"id"`
}

type WorkloadsSAPCentralInstanceResource struct{}

var _ sdk.ResourceWithUpdate = WorkloadsSAPCentralInstanceResource{}

func (r WorkloadsSAPCentralInstanceResource) ResourceType() string {
	return "azurerm_workloads_sap_central_instance"
}

func (r WorkloadsSAPCentralInstanceResource) ModelObject() interface{} {
	return &WorkloadsSAPCentralInstanceModel{}
}

func (r WorkloadsSAPCentralInstanceResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return sapcentralinstances.ValidateSAPCentralInstanceID
}

func (r WorkloadsSAPCentralInstanceResource) Arguments() map[string]*pluginsdk.Schema {
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

func (r WorkloadsSAPCentralInstanceResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"enqueue_replication_server_properties": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"ers_version": {
						Type:     pluginsdk.TypeString,
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
				},
			},
		},

		"enqueue_server_properties": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
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

					"port": {
						Type:     pluginsdk.TypeInt,
						Computed: true,
					},
				},
			},
		},

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

		"gateway_server_properties": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"health": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},

					"port": {
						Type:     pluginsdk.TypeInt,
						Computed: true,
					},
				},
			},
		},

		"health": {
			Type:     pluginsdk.TypeString,
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

		"message_server_properties": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"http_port": {
						Type:     pluginsdk.TypeInt,
						Computed: true,
					},

					"https_port": {
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

					"internal_ms_port": {
						Type:     pluginsdk.TypeInt,
						Computed: true,
					},

					"ms_port": {
						Type:     pluginsdk.TypeInt,
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

func (r WorkloadsSAPCentralInstanceResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var model WorkloadsSAPCentralInstanceModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			client := metadata.Client.Workloads.SAPCentralInstancesClient
			sapVirtualInstanceId, err := sapvirtualinstances.ParseSapVirtualInstanceID(model.WorkloadsSapVirtualInstanceId)
			if err != nil {
				return err
			}

			id := sapcentralinstances.NewCentralInstanceID(sapVirtualInstanceId.SubscriptionId, sapVirtualInstanceId.ResourceGroupName, sapVirtualInstanceId.SapVirtualInstanceName, model.Name)
			existing, err := client.Get(ctx, id)
			if err != nil && !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for existing %s: %+v", id, err)
			}

			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			properties := &sapcentralinstances.SAPCentralServerInstance{
				Location:   location.Normalize(model.Location),
				Properties: &sapcentralinstances.SAPCentralServerProperties{},
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

func (r WorkloadsSAPCentralInstanceResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Workloads.SAPCentralInstancesClient

			id, err := sapcentralinstances.ParseCentralInstanceID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var model WorkloadsSAPCentralInstanceModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			properties := &sapcentralinstances.UpdateSAPCentralInstanceRequest{}

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

func (r WorkloadsSAPCentralInstanceResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Workloads.SAPCentralInstancesClient

			id, err := sapcentralinstances.ParseCentralInstanceID(metadata.ResourceData.Id())
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

			state := WorkloadsSAPCentralInstanceModel{
				Name:                          id.CentralInstanceName,
				WorkloadsSapVirtualInstanceId: sapvirtualinstances.NewSapVirtualInstanceID(id.SubscriptionId, id.ResourceGroupName, id.SapVirtualInstanceName).ID(),
				Location:                      location.Normalize(model.Location),
			}

			if properties := model.Properties; properties != nil {
				state.EnqueueReplicationServerProperties = flattenEnqueueReplicationServerPropertiesModel(properties.EnqueueReplicationServerProperties)

				state.EnqueueServerProperties = flattenEnqueueServerPropertiesModel(properties.EnqueueServerProperties)

				state.Errors = flattenSAPVirtualInstanceErrorModel(properties.Errors)

				state.GatewayServerProperties = flattenGatewayServerPropertiesModel(properties.GatewayServerProperties)

				if properties.Health != nil {
					state.Health = *properties.Health
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

				state.MessageServerProperties = flattenMessageServerPropertiesModel(properties.MessageServerProperties)

				if properties.Status != nil {
					state.Status = *properties.Status
				}

				if properties.Subnet != nil {
					state.Subnet = *properties.Subnet
				}

				state.VMDetails = flattenCentralServerVMDetailsModelArray(properties.VMDetails)
			}
			if model.Tags != nil {
				state.Tags = *model.Tags
			}

			return metadata.Encode(&state)
		},
	}
}

func (r WorkloadsSAPCentralInstanceResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Workloads.SAPCentralInstancesClient

			id, err := sapcentralinstances.ParseCentralInstanceID(metadata.ResourceData.Id())
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

func flattenEnqueueReplicationServerPropertiesModel(input *sapcentralinstances.EnqueueReplicationServerProperties) []EnqueueReplicationServerPropertiesModel {
	var outputList []EnqueueReplicationServerPropertiesModel
	if input == nil {
		return outputList
	}
	output := EnqueueReplicationServerPropertiesModel{}
	if input.ErsVersion != nil {
		output.ErsVersion = *input.ErsVersion
	}

	if input.Health != nil {
		output.Health = *input.Health
	}

	if input.Hostname != nil {
		output.Hostname = *input.Hostname
	}

	if input.IPAddress != nil {
		output.IPAddress = *input.IPAddress
	}

	if input.InstanceNo != nil {
		output.InstanceNo = *input.InstanceNo
	}

	if input.KernelPatch != nil {
		output.KernelPatch = *input.KernelPatch
	}

	if input.KernelVersion != nil {
		output.KernelVersion = *input.KernelVersion
	}

	return append(outputList, output)
}

func flattenEnqueueServerPropertiesModel(input *sapcentralinstances.EnqueueServerProperties) []EnqueueServerPropertiesModel {
	var outputList []EnqueueServerPropertiesModel
	if input == nil {
		return outputList
	}
	output := EnqueueServerPropertiesModel{}
	if input.Health != nil {
		output.Health = *input.Health
	}

	if input.Hostname != nil {
		output.Hostname = *input.Hostname
	}

	if input.IPAddress != nil {
		output.IPAddress = *input.IPAddress
	}

	if input.Port != nil {
		output.Port = *input.Port
	}

	return append(outputList, output)
}

func flattenSAPVirtualInstanceErrorModel(input *sapcentralinstances.SAPVirtualInstanceError) []SAPVirtualInstanceErrorModel {
	var outputList []SAPVirtualInstanceErrorModel
	if input == nil {
		return outputList
	}
	output := SAPVirtualInstanceErrorModel{
		Properties: flattenErrorDefinitionModel(input.Properties),
	}

	return append(outputList, output)
}

func flattenErrorDefinitionModel(input *sapcentralinstances.ErrorDefinition) []ErrorDefinitionModel {
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

func flattenGatewayServerPropertiesModel(input *sapcentralinstances.GatewayServerProperties) []GatewayServerPropertiesModel {
	var outputList []GatewayServerPropertiesModel
	if input == nil {
		return outputList
	}
	output := GatewayServerPropertiesModel{}
	if input.Health != nil {
		output.Health = *input.Health
	}

	if input.Port != nil {
		output.Port = *input.Port
	}

	return append(outputList, output)
}

func flattenLoadBalancerDetailsModel(input *sapcentralinstances.LoadBalancerDetails) []LoadBalancerDetailsModel {
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

func flattenMessageServerPropertiesModel(input *sapcentralinstances.MessageServerProperties) []MessageServerPropertiesModel {
	var outputList []MessageServerPropertiesModel
	if input == nil {
		return outputList
	}
	output := MessageServerPropertiesModel{}
	if input.HTTPPort != nil {
		output.HTTPPort = *input.HTTPPort
	}

	if input.HTTPSPort != nil {
		output.HTTPSPort = *input.HTTPSPort
	}

	if input.Health != nil {
		output.Health = *input.Health
	}

	if input.Hostname != nil {
		output.Hostname = *input.Hostname
	}

	if input.IPAddress != nil {
		output.IPAddress = *input.IPAddress
	}

	if input.InternalMsPort != nil {
		output.InternalMsPort = *input.InternalMsPort
	}

	if input.MsPort != nil {
		output.MsPort = *input.MsPort
	}

	return append(outputList, output)
}

func flattenCentralServerVMDetailsModelArray(inputList *[]sapcentralinstances.CentralServerVMDetails) []CentralServerVMDetailsModel {
	var outputList []CentralServerVMDetailsModel
	if inputList == nil {
		return outputList
	}
	for _, input := range *inputList {
		output := CentralServerVMDetailsModel{
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

func flattenStorageInformationModelArray(inputList *[]sapcentralinstances.StorageInformation) []StorageInformationModel {
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
