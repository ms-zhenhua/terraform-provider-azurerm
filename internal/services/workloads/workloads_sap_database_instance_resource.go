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
	"github.com/hashicorp/go-azure-sdk/resource-manager/workloads/2023-04-01/sapdatabaseinstances"
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

type WorkloadsSAPDatabaseInstanceModel struct {
	Name                          string                                        `tfschema:"name"`
	WorkloadsSapVirtualInstanceId string                                        `tfschema:"workloads_sap_virtual_instance_id"`
	Location                      string                                        `tfschema:"location"`
	Tags                          map[string]string                             `tfschema:"tags"`
	DatabaseSid                   string                                        `tfschema:"database_sid"`
	DatabaseType                  string                                        `tfschema:"database_type"`
	Errors                        []SAPVirtualInstanceErrorModel                `tfschema:"errors"`
	IPAddress                     string                                        `tfschema:"ip_address"`
	LoadBalancerDetails           []LoadBalancerDetailsModel                    `tfschema:"load_balancer_details"`
	Status                        sapdatabaseinstances.SAPVirtualInstanceStatus `tfschema:"status"`
	Subnet                        string                                        `tfschema:"subnet"`
	VMDetails                     []DatabaseVMDetailsModel                      `tfschema:"vm_details"`
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

type DatabaseVMDetailsModel struct {
	Status           sapdatabaseinstances.SAPVirtualInstanceStatus `tfschema:"status"`
	StorageDetails   []StorageInformationModel                     `tfschema:"storage_details"`
	VirtualMachineId string                                        `tfschema:"virtual_machine_id"`
}

type StorageInformationModel struct {
	Id string `tfschema:"id"`
}

type WorkloadsSAPDatabaseInstanceResource struct{}

var _ sdk.ResourceWithUpdate = WorkloadsSAPDatabaseInstanceResource{}

func (r WorkloadsSAPDatabaseInstanceResource) ResourceType() string {
	return "azurerm_workloads_sap_database_instance"
}

func (r WorkloadsSAPDatabaseInstanceResource) ModelObject() interface{} {
	return &WorkloadsSAPDatabaseInstanceModel{}
}

func (r WorkloadsSAPDatabaseInstanceResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return sapdatabaseinstances.ValidateSAPDatabaseInstanceID
}

func (r WorkloadsSAPDatabaseInstanceResource) Arguments() map[string]*pluginsdk.Schema {
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

func (r WorkloadsSAPDatabaseInstanceResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"database_sid": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"database_type": {
			Type:     pluginsdk.TypeString,
			Computed: true,
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

		"ip_address": {
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
					"status": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},

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

					"virtual_machine_id": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},
				},
			},
		},
	}
}

func (r WorkloadsSAPDatabaseInstanceResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var model WorkloadsSAPDatabaseInstanceModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			client := metadata.Client.Workloads.SAPDatabaseInstancesClient
			sapVirtualInstanceId, err := sapvirtualinstances.ParseSapVirtualInstanceID(model.WorkloadsSapVirtualInstanceId)
			if err != nil {
				return err
			}

			id := sapdatabaseinstances.NewDatabaseInstanceID(sapVirtualInstanceId.SubscriptionId, sapVirtualInstanceId.ResourceGroupName, sapVirtualInstanceId.SapVirtualInstanceName, model.Name)
			existing, err := client.Get(ctx, id)
			if err != nil && !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for existing %s: %+v", id, err)
			}

			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			properties := &sapdatabaseinstances.SAPDatabaseInstance{
				Location:   location.Normalize(model.Location),
				Properties: &sapdatabaseinstances.SAPDatabaseProperties{},
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

func (r WorkloadsSAPDatabaseInstanceResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Workloads.SAPDatabaseInstancesClient

			id, err := sapdatabaseinstances.ParseDatabaseInstanceID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var model WorkloadsSAPDatabaseInstanceModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			properties := &sapdatabaseinstances.UpdateSAPDatabaseInstanceRequest{}

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

func (r WorkloadsSAPDatabaseInstanceResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Workloads.SAPDatabaseInstancesClient

			id, err := sapdatabaseinstances.ParseDatabaseInstanceID(metadata.ResourceData.Id())
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

			state := WorkloadsSAPDatabaseInstanceModel{
				Name:                          id.DatabaseInstanceName,
				WorkloadsSapVirtualInstanceId: sapvirtualinstances.NewSapVirtualInstanceID(id.SubscriptionId, id.ResourceGroupName, id.SapVirtualInstanceName).ID(),
				Location:                      location.Normalize(model.Location),
			}

			if properties := model.Properties; properties != nil {
				if properties.DatabaseSid != nil {
					state.DatabaseSid = *properties.DatabaseSid
				}

				if properties.DatabaseType != nil {
					state.DatabaseType = *properties.DatabaseType
				}

				state.Errors = flattenSAPVirtualInstanceErrorModel(properties.Errors)

				if properties.IPAddress != nil {
					state.IPAddress = *properties.IPAddress
				}

				state.LoadBalancerDetails = flattenLoadBalancerDetailsModel(properties.LoadBalancerDetails)

				if properties.Status != nil {
					state.Status = *properties.Status
				}

				if properties.Subnet != nil {
					state.Subnet = *properties.Subnet
				}

				state.VMDetails = flattenDatabaseVMDetailsModelArray(properties.VMDetails)
			}
			if model.Tags != nil {
				state.Tags = *model.Tags
			}

			return metadata.Encode(&state)
		},
	}
}

func (r WorkloadsSAPDatabaseInstanceResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Workloads.SAPDatabaseInstancesClient

			id, err := sapdatabaseinstances.ParseDatabaseInstanceID(metadata.ResourceData.Id())
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

func flattenSAPVirtualInstanceErrorModel(input *sapdatabaseinstances.SAPVirtualInstanceError) []SAPVirtualInstanceErrorModel {
	var outputList []SAPVirtualInstanceErrorModel
	if input == nil {
		return outputList
	}
	output := SAPVirtualInstanceErrorModel{
		Properties: flattenErrorDefinitionModel(input.Properties),
	}

	return append(outputList, output)
}

func flattenErrorDefinitionModel(input *sapdatabaseinstances.ErrorDefinition) []ErrorDefinitionModel {
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

func flattenLoadBalancerDetailsModel(input *sapdatabaseinstances.LoadBalancerDetails) []LoadBalancerDetailsModel {
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

func flattenDatabaseVMDetailsModelArray(inputList *[]sapdatabaseinstances.DatabaseVMDetails) []DatabaseVMDetailsModel {
	var outputList []DatabaseVMDetailsModel
	if inputList == nil {
		return outputList
	}
	for _, input := range *inputList {
		output := DatabaseVMDetailsModel{
			StorageDetails: flattenStorageInformationModelArray(input.StorageDetails),
		}

		if input.Status != nil {
			output.Status = *input.Status
		}

		if input.VirtualMachineId != nil {
			output.VirtualMachineId = *input.VirtualMachineId
		}
		outputList = append(outputList, output)
	}
	return outputList
}

func flattenStorageInformationModelArray(inputList *[]sapdatabaseinstances.StorageInformation) []StorageInformationModel {
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
