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
	"github.com/hashicorp/go-azure-sdk/resource-manager/workloads/2023-04-01/saplandscapemonitor"
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

type WorkloadsSapLandscapeMonitorModel struct {
	Name                 string                                       `tfschema:"name"`
	WorkloadsMonitorId   string                                       `tfschema:"workloads_monitor_id"`
	Grouping             []SapLandscapeMonitorPropertiesGroupingModel `tfschema:"grouping"`
	TopMetricsThresholds []SapLandscapeMonitorMetricThresholdsModel   `tfschema:"top_metrics_thresholds"`
}

type SapLandscapeMonitorPropertiesGroupingModel struct {
	Landscape      []SapLandscapeMonitorSidMappingModel `tfschema:"landscape"`
	SapApplication []SapLandscapeMonitorSidMappingModel `tfschema:"sap_application"`
}

type SapLandscapeMonitorSidMappingModel struct {
	Name   string   `tfschema:"name"`
	TopSid []string `tfschema:"top_sid"`
}

type SapLandscapeMonitorMetricThresholdsModel struct {
	Green  float64 `tfschema:"green"`
	Name   string  `tfschema:"name"`
	Red    float64 `tfschema:"red"`
	Yellow float64 `tfschema:"yellow"`
}

type WorkloadsSapLandscapeMonitorResource struct{}

var _ sdk.ResourceWithUpdate = WorkloadsSapLandscapeMonitorResource{}

func (r WorkloadsSapLandscapeMonitorResource) ResourceType() string {
	return "azurerm_workloads_sap_landscape_monitor"
}

func (r WorkloadsSapLandscapeMonitorResource) ModelObject() interface{} {
	return &WorkloadsSapLandscapeMonitorModel{}
}

func (r WorkloadsSapLandscapeMonitorResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return saplandscapemonitor.ValidateSapLandscapeMonitorID
}

func (r WorkloadsSapLandscapeMonitorResource) Arguments() map[string]*pluginsdk.Schema {
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

		"grouping": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"landscape": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"name": {
									Type:         pluginsdk.TypeString,
									Optional:     true,
									ValidateFunc: validation.StringIsNotEmpty,
								},

								"top_sid": {
									Type:     pluginsdk.TypeList,
									Optional: true,
									Elem: &pluginsdk.Schema{
										Type: pluginsdk.TypeString, ValidateFunc: validation.StringIsNotEmpty,
									},
								},
							},
						},
					},

					"sap_application": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"name": {
									Type:         pluginsdk.TypeString,
									Optional:     true,
									ValidateFunc: validation.StringIsNotEmpty,
								},

								"top_sid": {
									Type:     pluginsdk.TypeList,
									Optional: true,
									Elem: &pluginsdk.Schema{
										Type: pluginsdk.TypeString, ValidateFunc: validation.StringIsNotEmpty,
									},
								},
							},
						},
					},
				},
			},
		},

		"top_metrics_thresholds": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"green": {
						Type:     pluginsdk.TypeFloat,
						Optional: true,
					},

					"name": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"red": {
						Type:     pluginsdk.TypeFloat,
						Optional: true,
					},

					"yellow": {
						Type:     pluginsdk.TypeFloat,
						Optional: true,
					},
				},
			},
		},
	}
}

func (r WorkloadsSapLandscapeMonitorResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r WorkloadsSapLandscapeMonitorResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var model WorkloadsSapLandscapeMonitorModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			client := metadata.Client.Workloads.SapLandscapeMonitorClient
			monitorId, err := monitors.ParseMonitorID(model.WorkloadsMonitorId)
			if err != nil {
				return err
			}

			id := saplandscapemonitor.NewMonitorID(monitorId.SubscriptionId, monitorId.ResourceGroupName, monitorId.MonitorName, model.Name)
			existing, err := client.Get(ctx, id)
			if err != nil && !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for existing %s: %+v", id, err)
			}

			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			properties := &saplandscapemonitor.SapLandscapeMonitor{
				Properties: &saplandscapemonitor.SapLandscapeMonitorProperties{
					Grouping:             expandSapLandscapeMonitorPropertiesGroupingModel(model.Grouping),
					TopMetricsThresholds: expandSapLandscapeMonitorMetricThresholdsModelArray(model.TopMetricsThresholds),
				},
			}

			if _, err := client.Create(ctx, id, *properties); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func (r WorkloadsSapLandscapeMonitorResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Workloads.SapLandscapeMonitorClient

			id, err := saplandscapemonitor.ParseMonitorID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var model WorkloadsSapLandscapeMonitorModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			properties := &saplandscapemonitor.SapLandscapeMonitor{
				Properties: &saplandscapemonitor.SapLandscapeMonitorProperties{},
			}

			if metadata.ResourceData.HasChange("grouping") {
				properties.Properties.Grouping = expandSapLandscapeMonitorPropertiesGroupingModel(model.Grouping)
			}

			if metadata.ResourceData.HasChange("top_metrics_thresholds") {
				properties.Properties.TopMetricsThresholds = expandSapLandscapeMonitorMetricThresholdsModelArray(model.TopMetricsThresholds)
			}

			if _, err := client.Update(ctx, *id, *properties); err != nil {
				return fmt.Errorf("updating %s: %+v", *id, err)
			}

			return nil
		},
	}
}

func (r WorkloadsSapLandscapeMonitorResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Workloads.SapLandscapeMonitorClient

			id, err := saplandscapemonitor.ParseMonitorID(metadata.ResourceData.Id())
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

			state := WorkloadsSapLandscapeMonitorModel{
				Name:               id.MonitorName,
				WorkloadsMonitorId: monitors.NewMonitorID(id.SubscriptionId, id.ResourceGroupName).ID(),
			}

			if properties := model.Properties; properties != nil {
				state.Grouping = flattenSapLandscapeMonitorPropertiesGroupingModel(properties.Grouping)

				state.TopMetricsThresholds = flattenSapLandscapeMonitorMetricThresholdsModelArray(properties.TopMetricsThresholds)
			}

			return metadata.Encode(&state)
		},
	}
}

func (r WorkloadsSapLandscapeMonitorResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Workloads.SapLandscapeMonitorClient

			id, err := saplandscapemonitor.ParseMonitorID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			if _, err := client.Delete(ctx, *id); err != nil {
				return fmt.Errorf("deleting %s: %+v", id, err)
			}

			return nil
		},
	}
}

func expandSapLandscapeMonitorPropertiesGroupingModel(inputList []SapLandscapeMonitorPropertiesGroupingModel) *saplandscapemonitor.SapLandscapeMonitorPropertiesGrouping {
	if len(inputList) == 0 {
		return nil
	}
	input := &inputList[0]
	output := saplandscapemonitor.SapLandscapeMonitorPropertiesGrouping{
		Landscape:      expandSapLandscapeMonitorSidMappingModelArray(input.Landscape),
		SapApplication: expandSapLandscapeMonitorSidMappingModelArray(input.SapApplication),
	}

	return &output
}

func expandSapLandscapeMonitorSidMappingModelArray(inputList []SapLandscapeMonitorSidMappingModel) *[]saplandscapemonitor.SapLandscapeMonitorSidMapping {
	var outputList []saplandscapemonitor.SapLandscapeMonitorSidMapping
	for _, v := range inputList {
		input := v
		output := saplandscapemonitor.SapLandscapeMonitorSidMapping{
			TopSid: &input.TopSid,
		}

		if input.Name != "" {
			output.Name = &input.Name
		}
		outputList = append(outputList, output)
	}
	return &outputList
}

func expandSapLandscapeMonitorMetricThresholdsModelArray(inputList []SapLandscapeMonitorMetricThresholdsModel) *[]saplandscapemonitor.SapLandscapeMonitorMetricThresholds {
	var outputList []saplandscapemonitor.SapLandscapeMonitorMetricThresholds
	for _, v := range inputList {
		input := v
		output := saplandscapemonitor.SapLandscapeMonitorMetricThresholds{
			Green:  &input.Green,
			Red:    &input.Red,
			Yellow: &input.Yellow,
		}

		if input.Name != "" {
			output.Name = &input.Name
		}
		outputList = append(outputList, output)
	}
	return &outputList
}

func flattenSapLandscapeMonitorPropertiesGroupingModel(input *saplandscapemonitor.SapLandscapeMonitorPropertiesGrouping) []SapLandscapeMonitorPropertiesGroupingModel {
	var outputList []SapLandscapeMonitorPropertiesGroupingModel
	if input == nil {
		return outputList
	}
	output := SapLandscapeMonitorPropertiesGroupingModel{
		Landscape:      flattenSapLandscapeMonitorSidMappingModelArray(input.Landscape),
		SapApplication: flattenSapLandscapeMonitorSidMappingModelArray(input.SapApplication),
	}

	return append(outputList, output)
}

func flattenSapLandscapeMonitorSidMappingModelArray(inputList *[]saplandscapemonitor.SapLandscapeMonitorSidMapping) []SapLandscapeMonitorSidMappingModel {
	var outputList []SapLandscapeMonitorSidMappingModel
	if inputList == nil {
		return outputList
	}
	for _, input := range *inputList {
		output := SapLandscapeMonitorSidMappingModel{}

		if input.Name != nil {
			output.Name = *input.Name
		}

		if input.TopSid != nil {
			output.TopSid = *input.TopSid
		}
		outputList = append(outputList, output)
	}
	return outputList
}

func flattenSapLandscapeMonitorMetricThresholdsModelArray(inputList *[]saplandscapemonitor.SapLandscapeMonitorMetricThresholds) []SapLandscapeMonitorMetricThresholdsModel {
	var outputList []SapLandscapeMonitorMetricThresholdsModel
	if inputList == nil {
		return outputList
	}
	for _, input := range *inputList {
		output := SapLandscapeMonitorMetricThresholdsModel{}

		if input.Green != nil {
			output.Green = *input.Green
		}

		if input.Name != nil {
			output.Name = *input.Name
		}

		if input.Red != nil {
			output.Red = *input.Red
		}

		if input.Yellow != nil {
			output.Yellow = *input.Yellow
		}
		outputList = append(outputList, output)
	}
	return outputList
}
