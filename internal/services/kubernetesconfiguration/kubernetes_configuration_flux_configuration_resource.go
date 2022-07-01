package kubernetesconfiguration

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"regexp"
	"time"

	"github.com/hashicorp/terraform-provider-azurerm/utils"

	"github.com/hashicorp/terraform-provider-azurerm/internal/services/attestation/validate"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/kubernetesconfiguration/2022-03-01/fluxconfiguration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type KubernetesConfigurationFluxConfigurationModel struct {
	Name                string                         `tfschema:"name"`
	ResourceGroupName   string                         `tfschema:"resource_group_name"`
	ClusterResourceName string                         `tfschema:"cluster_resource_name"`
	ClusterName         string                         `tfschema:"cluster_name"`
	Bucket              []BucketDefinitionModel        `tfschema:"bucket"`
	GitRepository       []GitRepositoryDefinitionModel `tfschema:"git_repository"`
	Kustomizations      []KustomizationDefinitionModel `tfschema:"kustomizations"`
	Namespace           string                         `tfschema:"namespace"`
	Scope               fluxconfiguration.ScopeType    `tfschema:"scope"`
	Suspend             bool                           `tfschema:"suspend"`
}

type BucketDefinitionModel struct {
	AccessKey             string `tfschema:"access_key"`
	BucketSecretKey       string `tfschema:"bucket_secret_key"`
	BucketName            string `tfschema:"bucket_name"`
	Insecure              bool   `tfschema:"insecure"`
	LocalAuthRef          string `tfschema:"local_auth_ref"`
	SyncIntervalInSeconds int64  `tfschema:"sync_interval_in_seconds"`
	TimeoutInSeconds      int64  `tfschema:"timeout_in_seconds"`
	Url                   string `tfschema:"url"`
}

type GitRepositoryDefinitionModel struct {
	HttpsCACert           string                         `tfschema:"https_ca_cert"`
	HttpsUser             string                         `tfschema:"https_user"`
	HttpsKey              string                         `tfschema:"https_key"`
	LocalAuthRef          string                         `tfschema:"local_auth_ref"`
	RepositoryRef         []RepositoryRefDefinitionModel `tfschema:"repository_ref"`
	SshKnownHosts         string                         `tfschema:"ssh_known_hosts"`
	SshPrivateKey         string                         `tfschema:"ssh_private_key"`
	SyncIntervalInSeconds int64                          `tfschema:"sync_interval_in_seconds"`
	TimeoutInSeconds      int64                          `tfschema:"timeout_in_seconds"`
	Url                   string                         `tfschema:"url"`
}

type KustomizationDefinitionModel struct {
	Name                   string   `tfschema:"name"`
	Path                   string   `tfschema:"path"`
	TimeoutInSeconds       int64    `tfschema:"timeout_in_seconds"`
	SyncIntervalInSeconds  int64    `tfschema:"sync_interval_in_seconds"`
	RetryIntervalInSeconds int64    `tfschema:"retry_interval_in_seconds"`
	Force                  bool     `tfschema:"force"`
	Prune                  bool     `tfschema:"prune"`
	DependsOn              []string `tfschema:"depends_on"`
}

type RepositoryRefDefinitionModel struct {
	Branch string `tfschema:"branch"`
	Commit string `tfschema:"commit"`
	Semver string `tfschema:"semver"`
	Tag    string `tfschema:"tag"`
}

type KubernetesConfigurationFluxConfigurationResource struct{}

var _ sdk.ResourceWithUpdate = KubernetesConfigurationFluxConfigurationResource{}

func (r KubernetesConfigurationFluxConfigurationResource) ResourceType() string {
	return "azurerm_kubernetes_configuration_flux_configuration"
}

func (r KubernetesConfigurationFluxConfigurationResource) ModelObject() interface{} {
	return &KubernetesConfigurationFluxConfigurationModel{}
}

func (r KubernetesConfigurationFluxConfigurationResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return fluxconfiguration.ValidateFluxConfigurationID
}

func (r KubernetesConfigurationFluxConfigurationResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
			ValidateFunc: validation.StringMatch(
				regexp.MustCompile("^[a-z\\d]([-a-z\\d]{0,28}[a-z\\d])?$"),
				"`name` must be between 1 and 30 characters. It can contain only lowercase letters, numbers, and hyphens (-). It must start and end with a lowercase letter or number.",
			),
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"cluster_resource_name": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			ForceNew: true,
			ValidateFunc: validation.StringInSlice([]string{
				"managedClusters",
			}, false),
			Default: "managedClusters",
		},

		"cluster_name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"kustomizations": {
			Type:     pluginsdk.TypeSet,
			Required: true,
			MinItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"name": {
						Type:     pluginsdk.TypeString,
						Required: true,
						ValidateFunc: validation.StringMatch(
							regexp.MustCompile("^[a-z\\d]([-a-z\\d]{0,28}[a-z\\d])?$"),
							"`name` of `kustomizations` must be between 1 and 30 characters. It can contain only lowercase letters, numbers, and hyphens (-). It must start and end with a lowercase letter or number.",
						),
					},

					"path": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"timeout_in_seconds": {
						Type:         pluginsdk.TypeInt,
						Optional:     true,
						Default:      600,
						ValidateFunc: validation.IntBetween(1, 35791394),
					},

					"sync_interval_in_seconds": {
						Type:         pluginsdk.TypeInt,
						Optional:     true,
						Default:      600,
						ValidateFunc: validation.IntBetween(1, 35791394),
					},

					"retry_interval_in_seconds": {
						Type:         pluginsdk.TypeInt,
						Optional:     true,
						Default:      600,
						ValidateFunc: validation.IntBetween(1, 35791394),
					},

					"force": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						Default:  false,
					},

					"prune": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						Default:  false,
					},

					"depends_on": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						Elem: &pluginsdk.Schema{
							Type: pluginsdk.TypeString,
						},
					},
				},
			},
			Set: func(v interface{}) int {
				var buf bytes.Buffer
				m := v.(map[string]interface{})
				buf.WriteString(m["name"].(string))
				return pluginsdk.HashString(buf.String())
			},
		},

		"bucket": {
			Type:         pluginsdk.TypeList,
			Optional:     true,
			MaxItems:     1,
			ExactlyOneOf: []string{"bucket", "git_repository"},
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"bucket_name": {
						Type:     pluginsdk.TypeString,
						Required: true,
						ValidateFunc: validation.StringMatch(
							regexp.MustCompile("^[a-z\\d]([-a-z\\d]{0,61}[a-z\\d])?$"),
							"`bucket_name` must be between 1 and 63 characters. It can contain only lowercase letters, numbers, and hyphens (-). It must start and end with a lowercase letter or number.",
						),
					},

					"url": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.IsURLWithHTTPorHTTPS,
					},

					"access_key": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
						ExactlyOneOf: []string{"bucket.0.access_key", "bucket.0.local_auth_ref"},
					},

					"bucket_secret_key": {
						Type:          pluginsdk.TypeString,
						Optional:      true,
						ValidateFunc:  validation.StringIsNotEmpty,
						Sensitive:     true,
						RequiredWith:  []string{"bucket.0.access_key"},
						ConflictsWith: []string{"bucket.0.local_auth_ref"},
					},

					"local_auth_ref": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						ValidateFunc: validation.StringMatch(
							regexp.MustCompile("^[a-z\\d]([-a-z\\d]{0,61}[a-z\\d])?$"),
							"`local_auth_ref` must be between 1 and 63 characters. It can contain only lowercase letters, numbers, and hyphens (-). It must start and end with a lowercase letter or number.",
						),
						ExactlyOneOf: []string{"bucket.0.access_key", "bucket.0.local_auth_ref"},
					},

					"insecure": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						Default:  false,
					},

					"sync_interval_in_seconds": {
						Type:         pluginsdk.TypeInt,
						Optional:     true,
						Default:      600,
						ValidateFunc: validation.IntBetween(1, 35791394),
					},

					"timeout_in_seconds": {
						Type:         pluginsdk.TypeInt,
						Optional:     true,
						Default:      600,
						ValidateFunc: validation.IntBetween(1, 35791394),
					},
				},
			},
		},

		"git_repository": {
			Type:         pluginsdk.TypeList,
			Optional:     true,
			MaxItems:     1,
			ExactlyOneOf: []string{"bucket", "git_repository"},
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"url": {
						Type:     pluginsdk.TypeString,
						Required: true,
						ValidateFunc: validation.Any(
							validation.StringMatch(regexp.MustCompile(`^http://(.)+$`), "The URL must begin with either http://, https://, git@, ssh://"),
							validation.StringMatch(regexp.MustCompile(`^https://(.)+$`), ""),
							validation.StringMatch(regexp.MustCompile(`^git@(.)+$`), ""),
							validation.StringMatch(regexp.MustCompile(`^ssh://(.)+$`), ""),
						),
					},

					"repository_ref": {
						Type:     pluginsdk.TypeList,
						Required: true,
						MaxItems: 1,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"branch": {
									Type:         pluginsdk.TypeString,
									Optional:     true,
									ExactlyOneOf: []string{"git_repository.0.repository_ref.0.branch", "git_repository.0.repository_ref.0.commit", "git_repository.0.repository_ref.0.semver", "git_repository.0.repository_ref.0.tag"},
									ValidateFunc: validation.StringIsNotEmpty,
								},

								"commit": {
									Type:         pluginsdk.TypeString,
									Optional:     true,
									ExactlyOneOf: []string{"git_repository.0.repository_ref.0.branch", "git_repository.0.repository_ref.0.commit", "git_repository.0.repository_ref.0.semver", "git_repository.0.repository_ref.0.tag"},
									ValidateFunc: validation.StringIsNotEmpty,
								},

								"semver": {
									Type:         pluginsdk.TypeString,
									Optional:     true,
									ExactlyOneOf: []string{"git_repository.0.repository_ref.0.branch", "git_repository.0.repository_ref.0.commit", "git_repository.0.repository_ref.0.semver", "git_repository.0.repository_ref.0.tag"},
									ValidateFunc: validation.StringIsNotEmpty,
								},

								"tag": {
									Type:         pluginsdk.TypeString,
									Optional:     true,
									ExactlyOneOf: []string{"git_repository.0.repository_ref.0.branch", "git_repository.0.repository_ref.0.commit", "git_repository.0.repository_ref.0.semver", "git_repository.0.repository_ref.0.tag"},
									ValidateFunc: validation.StringIsNotEmpty,
								},
							},
						},
					},

					"https_ca_cert": {
						Type:          pluginsdk.TypeString,
						Optional:      true,
						ValidateFunc:  validation.StringIsNotEmpty,
						ConflictsWith: []string{"git_repository.0.local_auth_ref", "git_repository.0.ssh_private_key", "git_repository.0.ssh_known_hosts"},
					},

					"https_user": {
						Type:          pluginsdk.TypeString,
						Optional:      true,
						ValidateFunc:  validation.StringIsNotEmpty,
						ConflictsWith: []string{"git_repository.0.local_auth_ref", "git_repository.0.ssh_private_key", "git_repository.0.ssh_known_hosts"},
					},

					"https_key": {
						Type:          pluginsdk.TypeString,
						Optional:      true,
						ValidateFunc:  validation.StringIsNotEmpty,
						Sensitive:     true,
						RequiredWith:  []string{"git_repository.0.https_user"},
						ConflictsWith: []string{"git_repository.0.local_auth_ref", "git_repository.0.ssh_private_key", "git_repository.0.ssh_known_hosts"},
					},

					"local_auth_ref": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						ValidateFunc: validation.StringMatch(
							regexp.MustCompile("^[a-z\\d]([-a-z\\d]{0,61}[a-z\\d])?$"),
							"`local_auth_ref` must be between 1 and 63 characters. It can contain only lowercase letters, numbers, and hyphens (-). It must start and end with a lowercase letter or number.",
						),
						ConflictsWith: []string{"git_repository.0.ssh_known_hosts", "git_repository.0.ssh_private_key", "git_repository.0.https_key", "git_repository.0.https_ca_cert"},
					},

					"ssh_private_key": {
						Type:          pluginsdk.TypeString,
						Optional:      true,
						ValidateFunc:  validate.IsCert,
						Sensitive:     true,
						ConflictsWith: []string{"git_repository.0.local_auth_ref", "git_repository.0.https_user", "git_repository.0.https_key", "git_repository.0.https_ca_cert"},
					},

					"ssh_known_hosts": {
						Type:          pluginsdk.TypeString,
						Optional:      true,
						ValidateFunc:  validation.StringIsNotEmpty,
						ConflictsWith: []string{"git_repository.0.local_auth_ref", "git_repository.0.https_user", "git_repository.0.https_key", "git_repository.0.https_ca_cert"},
					},

					"sync_interval_in_seconds": {
						Type:     pluginsdk.TypeInt,
						Optional: true,
						Default:  600,
					},

					"timeout_in_seconds": {
						Type:     pluginsdk.TypeInt,
						Optional: true,
						Default:  600,
					},
				},
			},
		},

		"namespace": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			ForceNew: true,
			ValidateFunc: validation.StringMatch(
				regexp.MustCompile("^[a-z\\d]([-a-z\\d]{0,61}[a-z\\d])?$"),
				"`name` must be between 1 and 63 characters. It can contain only lowercase letters, numbers, and hyphens (-). It must start and end with a lowercase letter or number.",
			),
			Default: "default",
		},

		"scope": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			ForceNew: true,
			ValidateFunc: validation.StringInSlice([]string{
				string(fluxconfiguration.ScopeTypeNamespace),
				string(fluxconfiguration.ScopeTypeCluster),
			}, false),
			Default: string(fluxconfiguration.ScopeTypeNamespace),
		},

		"suspend": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  false,
		},
	}
}

func (r KubernetesConfigurationFluxConfigurationResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r KubernetesConfigurationFluxConfigurationResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var model KubernetesConfigurationFluxConfigurationModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			client := metadata.Client.KubernetesConfiguration.FluxConfigurationClient
			subscriptionId := metadata.Client.Account.SubscriptionId

			var clusterRp string
			if model.ClusterResourceName == "managedClusters" {
				clusterRp = "Microsoft.ContainerService"
			} else {
				return fmt.Errorf("cluster resource of %s is not supported", model.ClusterResourceName)
			}

			id := fluxconfiguration.NewFluxConfigurationID(subscriptionId, model.ResourceGroupName, clusterRp, model.ClusterResourceName, model.ClusterName, model.Name)
			existing, err := client.Get(ctx, id)
			if err != nil && !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for existing %s: %+v", id, err)
			}

			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			properties := &fluxconfiguration.FluxConfiguration{
				Properties: &fluxconfiguration.FluxConfigurationProperties{
					Scope:   &model.Scope,
					Suspend: &model.Suspend,
				},
			}

			if len(model.GitRepository) > 0 {
				sourceKind := fluxconfiguration.SourceKindTypeGitRepository
				gitRepositoryValue, configurationProtectedSettings, err := expandGitRepositoryDefinitionModel(model.GitRepository)
				if err != nil {
					return err
				}

				properties.Properties.GitRepository = gitRepositoryValue
				properties.Properties.SourceKind = &sourceKind
				properties.Properties.ConfigurationProtectedSettings = configurationProtectedSettings
			} else {
				sourceKind := fluxconfiguration.SourceKindTypeBucket
				properties.Properties.SourceKind = &sourceKind
				bucketValue, configurationProtectedSettings, err := expandBucketDefinitionModel(model.Bucket)
				if err != nil {
					return err
				}

				properties.Properties.Bucket = bucketValue
				properties.Properties.SourceKind = &sourceKind
				properties.Properties.ConfigurationProtectedSettings = configurationProtectedSettings
			}

			kustomizationsValue, err := expandKustomizationDefinitionModel(model.Kustomizations)
			if err != nil {
				return err
			}

			properties.Properties.Kustomizations = kustomizationsValue

			if model.Namespace != "" {
				properties.Properties.Namespace = &model.Namespace
			}

			if err := client.CreateOrUpdateThenPoll(ctx, id, *properties); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func (r KubernetesConfigurationFluxConfigurationResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.KubernetesConfiguration.FluxConfigurationClient

			id, err := fluxconfiguration.ParseFluxConfigurationID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var model KubernetesConfigurationFluxConfigurationModel
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

			if metadata.ResourceData.HasChange("bucket") {
				bucketValue, configurationProtectedSettings, err := expandBucketDefinitionModel(model.Bucket)
				if err != nil {
					return err
				}

				properties.Properties.Bucket = bucketValue
				if properties.Properties.Bucket != nil {
					properties.Properties.ConfigurationProtectedSettings = configurationProtectedSettings
				}
			}

			if metadata.ResourceData.HasChange("git_repository") {
				gitRepositoryValue, configurationProtectedSettings, err := expandGitRepositoryDefinitionModel(model.GitRepository)
				if err != nil {
					return err
				}

				properties.Properties.GitRepository = gitRepositoryValue
				if properties.Properties.GitRepository != nil {
					properties.Properties.ConfigurationProtectedSettings = configurationProtectedSettings
				}
			}

			var sourceKind fluxconfiguration.SourceKindType
			if properties.Properties.Bucket != nil {
				sourceKind = fluxconfiguration.SourceKindTypeBucket
			} else {
				sourceKind = fluxconfiguration.SourceKindTypeGitRepository
			}

			properties.Properties.SourceKind = &sourceKind

			if metadata.ResourceData.HasChange("kustomizations") {
				kustomizationsValue, err := expandKustomizationDefinitionModel(model.Kustomizations)
				if err != nil {
					return err
				}

				properties.Properties.Kustomizations = kustomizationsValue
			}

			if metadata.ResourceData.HasChange("suspend") {
				properties.Properties.Suspend = &model.Suspend
			}

			properties.SystemData = nil

			if err := client.CreateOrUpdateThenPoll(ctx, *id, *properties); err != nil {
				return fmt.Errorf("updating %s: %+v", *id, err)
			}

			return nil
		},
	}
}

func (r KubernetesConfigurationFluxConfigurationResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.KubernetesConfiguration.FluxConfigurationClient

			id, err := fluxconfiguration.ParseFluxConfigurationID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var configModel KubernetesConfigurationFluxConfigurationModel
			if err := metadata.Decode(&configModel); err != nil {
				return fmt.Errorf("decoding: %+v", err)
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
				return fmt.Errorf("retrieving %s: model was nil", id)
			}

			state := KubernetesConfigurationFluxConfigurationModel{
				Name:                id.FluxConfigurationName,
				ResourceGroupName:   id.ResourceGroupName,
				ClusterResourceName: id.ClusterResourceName,
				ClusterName:         id.ClusterName,
			}

			if properties := model.Properties; properties != nil {
				bucketValue, err := flattenBucketDefinitionModel(properties.Bucket, &configModel)
				if err != nil {
					return err
				}

				state.Bucket = bucketValue

				gitRepositoryValue, err := flattenGitRepositoryDefinitionModel(properties.GitRepository, &configModel)
				if err != nil {
					return err
				}

				state.GitRepository = gitRepositoryValue

				kustomizationsValue, err := flattenKustomizationDefinitionModel(properties.Kustomizations)
				if err != nil {
					return err
				}

				state.Kustomizations = kustomizationsValue

				if properties.Namespace != nil {
					state.Namespace = *properties.Namespace
				}

				if properties.Scope != nil {
					state.Scope = *properties.Scope
				}

				if properties.Suspend != nil {
					state.Suspend = *properties.Suspend
				}
			}

			return metadata.Encode(&state)
		},
	}
}

func (r KubernetesConfigurationFluxConfigurationResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.KubernetesConfiguration.FluxConfigurationClient

			id, err := fluxconfiguration.ParseFluxConfigurationID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			if err := client.DeleteThenPoll(ctx, *id, fluxconfiguration.DeleteOperationOptions{}); err != nil {
				return fmt.Errorf("deleting %s: %+v", id, err)
			}

			return nil
		},
	}
}

func expandKustomizationDefinitionModel(inputList []KustomizationDefinitionModel) (*map[string]fluxconfiguration.KustomizationDefinition, error) {
	if len(inputList) == 0 {
		return nil, nil
	}

	outputList := make(map[string]fluxconfiguration.KustomizationDefinition)
	for _, v := range inputList {
		input := v
		output := fluxconfiguration.KustomizationDefinition{
			DependsOn:              &input.DependsOn,
			Force:                  &input.Force,
			Name:                   &input.Name,
			Prune:                  &input.Prune,
			RetryIntervalInSeconds: &input.RetryIntervalInSeconds,
			SyncIntervalInSeconds:  &input.SyncIntervalInSeconds,
			TimeoutInSeconds:       &input.TimeoutInSeconds,
		}

		if input.Path != "" {
			output.Path = utils.String(input.Path)
		}

		outputList[input.Name] = output
	}

	return &outputList, nil
}

func expandBucketDefinitionModel(inputList []BucketDefinitionModel) (*fluxconfiguration.BucketDefinition, *map[string]string, error) {
	if len(inputList) == 0 {
		return nil, nil, nil
	}

	input := &inputList[0]
	output := fluxconfiguration.BucketDefinition{
		Insecure:              &input.Insecure,
		SyncIntervalInSeconds: &input.SyncIntervalInSeconds,
		TimeoutInSeconds:      &input.TimeoutInSeconds,
	}

	if input.AccessKey != "" {
		output.AccessKey = &input.AccessKey
	}

	if input.BucketName != "" {
		output.BucketName = &input.BucketName
	}

	if input.LocalAuthRef != "" {
		output.LocalAuthRef = &input.LocalAuthRef
	}

	if input.Url != "" {
		output.Url = &input.Url
	}

	var configSettings = make(map[string]string)
	if input.BucketSecretKey != "" {
		configSettings["bucketSecretKey"] = base64.StdEncoding.EncodeToString([]byte(input.BucketSecretKey))
	}

	var outputConfigSettings *map[string]string = nil
	if len(configSettings) > 0 {
		outputConfigSettings = &configSettings
	}

	return &output, outputConfigSettings, nil
}

func expandGitRepositoryDefinitionModel(inputList []GitRepositoryDefinitionModel) (*fluxconfiguration.GitRepositoryDefinition, *map[string]string, error) {
	if len(inputList) == 0 {
		return nil, nil, nil
	}

	input := &inputList[0]
	output := fluxconfiguration.GitRepositoryDefinition{
		SyncIntervalInSeconds: &input.SyncIntervalInSeconds,
		TimeoutInSeconds:      &input.TimeoutInSeconds,
	}

	if input.HttpsCACert != "" {
		encodedValue := base64.StdEncoding.EncodeToString([]byte(input.HttpsCACert))
		output.HttpsCACert = &encodedValue
	}

	if input.HttpsUser != "" {
		output.HttpsUser = &input.HttpsUser
	}

	if input.LocalAuthRef != "" {
		output.LocalAuthRef = &input.LocalAuthRef
	}

	repositoryRefValue, err := expandRepositoryRefDefinitionModel(input.RepositoryRef)
	if err != nil {
		return nil, nil, err
	}

	output.RepositoryRef = repositoryRefValue

	if input.SshKnownHosts != "" {
		encodedValue := base64.StdEncoding.EncodeToString([]byte(input.SshKnownHosts))
		output.SshKnownHosts = &encodedValue
	}

	if input.Url != "" {
		output.Url = &input.Url
	}

	var configSettings = make(map[string]string)
	if input.HttpsKey != "" {
		configSettings["httpsKey"] = base64.StdEncoding.EncodeToString([]byte(input.HttpsKey))
	}

	if input.SshPrivateKey != "" {
		configSettings["sshPrivateKey"] = base64.StdEncoding.EncodeToString([]byte(input.SshPrivateKey))
	}

	return &output, &configSettings, nil
}

func expandRepositoryRefDefinitionModel(inputList []RepositoryRefDefinitionModel) (*fluxconfiguration.RepositoryRefDefinition, error) {
	if len(inputList) == 0 {
		return nil, nil
	}

	input := &inputList[0]
	output := fluxconfiguration.RepositoryRefDefinition{}

	if input.Branch != "" {
		output.Branch = &input.Branch
	}

	if input.Commit != "" {
		output.Commit = &input.Commit
	}

	if input.Semver != "" {
		output.Semver = &input.Semver
	}

	if input.Tag != "" {
		output.Tag = &input.Tag
	}

	return &output, nil
}

func flattenKustomizationDefinitionModel(inputList *map[string]fluxconfiguration.KustomizationDefinition) ([]KustomizationDefinitionModel, error) {
	var outputList []KustomizationDefinitionModel
	if inputList == nil {
		return outputList, nil
	}

	for _, input := range *inputList {
		output := KustomizationDefinitionModel{}

		if input.DependsOn != nil {
			output.DependsOn = *input.DependsOn
		}

		if input.Force != nil {
			output.Force = *input.Force
		}

		if input.Name != nil {
			output.Name = *input.Name
		}

		if input.Path != nil {
			output.Path = *input.Path
		}

		if input.Prune != nil {
			output.Prune = *input.Prune
		}

		if input.RetryIntervalInSeconds != nil {
			output.RetryIntervalInSeconds = *input.RetryIntervalInSeconds
		}

		if input.SyncIntervalInSeconds != nil {
			output.SyncIntervalInSeconds = *input.SyncIntervalInSeconds
		}

		if input.TimeoutInSeconds != nil {
			output.TimeoutInSeconds = *input.TimeoutInSeconds
		}

		outputList = append(outputList, output)
	}

	return outputList, nil
}

func flattenBucketDefinitionModel(input *fluxconfiguration.BucketDefinition, model *KubernetesConfigurationFluxConfigurationModel) ([]BucketDefinitionModel, error) {
	var outputList []BucketDefinitionModel
	if input == nil {
		return outputList, nil
	}

	output := BucketDefinitionModel{}

	if input.AccessKey != nil {
		output.AccessKey = *input.AccessKey
	}

	if input.BucketName != nil {
		output.BucketName = *input.BucketName
	}

	if input.Insecure != nil {
		output.Insecure = *input.Insecure
	}

	if input.LocalAuthRef != nil {
		output.LocalAuthRef = *input.LocalAuthRef
	}

	if input.SyncIntervalInSeconds != nil {
		output.SyncIntervalInSeconds = *input.SyncIntervalInSeconds
	}

	if input.TimeoutInSeconds != nil {
		output.TimeoutInSeconds = *input.TimeoutInSeconds
	}

	if input.Url != nil {
		output.Url = *input.Url
	}

	if model != nil && len(model.Bucket) > 0 {
		output.BucketSecretKey = model.Bucket[0].BucketSecretKey
	}

	return append(outputList, output), nil
}

func flattenGitRepositoryDefinitionModel(input *fluxconfiguration.GitRepositoryDefinition, model *KubernetesConfigurationFluxConfigurationModel) ([]GitRepositoryDefinitionModel, error) {
	var outputList []GitRepositoryDefinitionModel
	if input == nil {
		return outputList, nil
	}

	output := GitRepositoryDefinitionModel{}

	if input.HttpsCACert != nil {
		decodedValue, err := base64.StdEncoding.DecodeString(*input.HttpsCACert)
		if err != nil {
			return nil, err
		}

		output.HttpsCACert = string(decodedValue)
	}

	if input.HttpsUser != nil {
		output.HttpsUser = *input.HttpsUser
	}

	if input.LocalAuthRef != nil {
		output.LocalAuthRef = *input.LocalAuthRef
	}

	repositoryRefValue, err := flattenRepositoryRefDefinitionModel(input.RepositoryRef)
	if err != nil {
		return nil, err
	}

	output.RepositoryRef = repositoryRefValue

	if input.SshKnownHosts != nil {
		decodedValue, err := base64.StdEncoding.DecodeString(*input.SshKnownHosts)
		if err != nil {
			return nil, err
		}

		output.SshKnownHosts = string(decodedValue)
	}

	if input.SyncIntervalInSeconds != nil {
		output.SyncIntervalInSeconds = *input.SyncIntervalInSeconds
	}

	if input.TimeoutInSeconds != nil {
		output.TimeoutInSeconds = *input.TimeoutInSeconds
	}

	if input.Url != nil {
		output.Url = *input.Url
	}

	if model != nil && len(model.GitRepository) > 0 {
		output.HttpsKey = model.GitRepository[0].HttpsKey
		output.SshPrivateKey = model.GitRepository[0].SshPrivateKey
	}

	return append(outputList, output), nil
}

func flattenRepositoryRefDefinitionModel(input *fluxconfiguration.RepositoryRefDefinition) ([]RepositoryRefDefinitionModel, error) {
	var outputList []RepositoryRefDefinitionModel
	if input == nil {
		return outputList, nil
	}

	output := RepositoryRefDefinitionModel{}

	if input.Branch != nil {
		output.Branch = *input.Branch
	}

	if input.Commit != nil {
		output.Commit = *input.Commit
	}

	if input.Semver != nil {
		output.Semver = *input.Semver
	}

	if input.Tag != nil {
		output.Tag = *input.Tag
	}

	return append(outputList, output), nil
}
