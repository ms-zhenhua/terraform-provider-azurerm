package kubernetesconfiguration

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
	"github.com/hashicorp/go-azure-sdk/resource-manager/kubernetesconfiguration/2022-07-01/fluxconfiguration"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	azValidate "github.com/hashicorp/terraform-provider-azurerm/helpers/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/kubernetesconfiguration/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type KubernetesConfigurationFluxConfigurationModel struct {
	Name                           string                                `tfschema:"name"`
	ResourceGroupName              string                                `tfschema:"resource_group_name"`
	ClusterRp                      string                                `tfschema:"cluster_rp"`
	ClusterResourceName            string                                `tfschema:"cluster_resource_name"`
	ClusterName                    string                                `tfschema:"cluster_name"`
	AzureBlob                      []AzureBlobDefinitionModel            `tfschema:"azure_blob"`
	Bucket                         []BucketDefinitionModel               `tfschema:"bucket"`
	ConfigurationProtectedSettings map[string]string                     `tfschema:"configuration_protected_settings"`
	GitRepository                  []GitRepositoryDefinitionModel        `tfschema:"git_repository"`
	Kustomizations                 string                                `tfschema:"kustomizations"`
	Namespace                      string                                `tfschema:"namespace"`
	Scope                          fluxconfiguration.ScopeType           `tfschema:"scope"`
	SourceKind                     fluxconfiguration.SourceKindType      `tfschema:"source_kind"`
	Suspend                        bool                                  `tfschema:"suspend"`
	ComplianceState                fluxconfiguration.FluxComplianceState `tfschema:"compliance_state"`
	ErrorMessage                   string                                `tfschema:"error_message"`
	RepositoryPublicKey            string                                `tfschema:"repository_public_key"`
	SourceSyncedCommitId           string                                `tfschema:"source_synced_commit_id"`
	SourceUpdatedAt                string                                `tfschema:"source_updated_at"`
	StatusUpdatedAt                string                                `tfschema:"status_updated_at"`
	Statuses                       []ObjectStatusDefinitionModel         `tfschema:"statuses"`
}

type AzureBlobDefinitionModel struct {
	AccountKey            string                            `tfschema:"account_key"`
	ContainerName         string                            `tfschema:"container_name"`
	LocalAuthRef          string                            `tfschema:"local_auth_ref"`
	ManagedIdentity       []ManagedIdentityDefinitionModel  `tfschema:"managed_identity"`
	SasToken              string                            `tfschema:"sas_token"`
	ServicePrincipal      []ServicePrincipalDefinitionModel `tfschema:"service_principal"`
	SyncIntervalInSeconds int64                             `tfschema:"sync_interval_in_seconds"`
	TimeoutInSeconds      int64                             `tfschema:"timeout_in_seconds"`
	Url                   string                            `tfschema:"url"`
}

type ManagedIdentityDefinitionModel struct {
	ClientId string `tfschema:"client_id"`
}

type ServicePrincipalDefinitionModel struct {
	ClientCertificate          string `tfschema:"client_certificate"`
	ClientCertificatePassword  string `tfschema:"client_certificate_password"`
	ClientCertificateSendChain bool   `tfschema:"client_certificate_send_chain"`
	ClientId                   string `tfschema:"client_id"`
	ClientSecret               string `tfschema:"client_secret"`
	TenantId                   string `tfschema:"tenant_id"`
}

type BucketDefinitionModel struct {
	AccessKey             string `tfschema:"access_key"`
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
	LocalAuthRef          string                         `tfschema:"local_auth_ref"`
	RepositoryRef         []RepositoryRefDefinitionModel `tfschema:"repository_ref"`
	SshKnownHosts         string                         `tfschema:"ssh_known_hosts"`
	SyncIntervalInSeconds int64                          `tfschema:"sync_interval_in_seconds"`
	TimeoutInSeconds      int64                          `tfschema:"timeout_in_seconds"`
	Url                   string                         `tfschema:"url"`
}

type RepositoryRefDefinitionModel struct {
	Branch string `tfschema:"branch"`
	Commit string `tfschema:"commit"`
	Semver string `tfschema:"semver"`
	Tag    string `tfschema:"tag"`
}

type ObjectStatusDefinitionModel struct {
	AppliedBy             []ObjectReferenceDefinitionModel       `tfschema:"applied_by"`
	ComplianceState       fluxconfiguration.FluxComplianceState  `tfschema:"compliance_state"`
	HelmReleaseProperties []HelmReleasePropertiesDefinitionModel `tfschema:"helm_release_properties"`
	Kind                  string                                 `tfschema:"kind"`
	Name                  string                                 `tfschema:"name"`
	Namespace             string                                 `tfschema:"namespace"`
	StatusConditions      []ObjectStatusConditionDefinitionModel `tfschema:"status_conditions"`
}

type ObjectReferenceDefinitionModel struct {
	Name      string `tfschema:"name"`
	Namespace string `tfschema:"namespace"`
}

type HelmReleasePropertiesDefinitionModel struct {
	FailureCount        int64                            `tfschema:"failure_count"`
	HelmChartRef        []ObjectReferenceDefinitionModel `tfschema:"helm_chart_ref"`
	InstallFailureCount int64                            `tfschema:"install_failure_count"`
	LastRevisionApplied int64                            `tfschema:"last_revision_applied"`
	UpgradeFailureCount int64                            `tfschema:"upgrade_failure_count"`
}

type ObjectStatusConditionDefinitionModel struct {
	LastTransitionTime string `tfschema:"last_transition_time"`
	Message            string `tfschema:"message"`
	Reason             string `tfschema:"reason"`
	Status             string `tfschema:"status"`
	Type               string `tfschema:"type"`
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
	return kubernetesconfiguration.ValidateFluxConfigurationID
}

func (r KubernetesConfigurationFluxConfigurationResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"cluster_rp": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"cluster_resource_name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"cluster_name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"azure_blob": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"account_key": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"container_name": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"local_auth_ref": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"managed_identity": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						MaxItems: 1,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"client_id": {
									Type:         pluginsdk.TypeString,
									Optional:     true,
									ValidateFunc: validation.StringIsNotEmpty,
								},
							},
						},
					},

					"sas_token": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"service_principal": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						MaxItems: 1,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"client_certificate": {
									Type:         pluginsdk.TypeString,
									Optional:     true,
									ValidateFunc: validation.StringIsNotEmpty,
								},

								"client_certificate_password": {
									Type:         pluginsdk.TypeString,
									Optional:     true,
									ValidateFunc: validation.StringIsNotEmpty,
								},

								"client_certificate_send_chain": {
									Type:     pluginsdk.TypeBool,
									Optional: true,
								},

								"client_id": {
									Type:         pluginsdk.TypeString,
									Optional:     true,
									ValidateFunc: validation.StringIsNotEmpty,
								},

								"client_secret": {
									Type:         pluginsdk.TypeString,
									Optional:     true,
									ValidateFunc: validation.StringIsNotEmpty,
								},

								"tenant_id": {
									Type:         pluginsdk.TypeString,
									Optional:     true,
									ValidateFunc: validation.StringIsNotEmpty,
								},
							},
						},
					},

					"sync_interval_in_seconds": {
						Type:     pluginsdk.TypeInt,
						Optional: true,
					},

					"timeout_in_seconds": {
						Type:     pluginsdk.TypeInt,
						Optional: true,
					},

					"url": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
				},
			},
		},

		"bucket": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"access_key": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"bucket_name": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"insecure": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
					},

					"local_auth_ref": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"sync_interval_in_seconds": {
						Type:     pluginsdk.TypeInt,
						Optional: true,
					},

					"timeout_in_seconds": {
						Type:     pluginsdk.TypeInt,
						Optional: true,
					},

					"url": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
				},
			},
		},

		"configuration_protected_settings": {
			Type:     pluginsdk.TypeMap,
			Optional: true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},

		"git_repository": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"https_ca_cert": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"https_user": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"local_auth_ref": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"repository_ref": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						MaxItems: 1,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"branch": {
									Type:         pluginsdk.TypeString,
									Optional:     true,
									ValidateFunc: validation.StringIsNotEmpty,
								},

								"commit": {
									Type:         pluginsdk.TypeString,
									Optional:     true,
									ValidateFunc: validation.StringIsNotEmpty,
								},

								"semver": {
									Type:         pluginsdk.TypeString,
									Optional:     true,
									ValidateFunc: validation.StringIsNotEmpty,
								},

								"tag": {
									Type:         pluginsdk.TypeString,
									Optional:     true,
									ValidateFunc: validation.StringIsNotEmpty,
								},
							},
						},
					},

					"ssh_known_hosts": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"sync_interval_in_seconds": {
						Type:     pluginsdk.TypeInt,
						Optional: true,
					},

					"timeout_in_seconds": {
						Type:     pluginsdk.TypeInt,
						Optional: true,
					},

					"url": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
				},
			},
		},

		"kustomizations": {
			Type:             pluginsdk.TypeString,
			Optional:         true,
			ValidateFunc:     validation.StringIsJSON,
			DiffSuppressFunc: pluginsdk.SuppressJsonDiff,
		},

		"namespace": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"scope": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			ForceNew: true,
			ValidateFunc: validation.StringInSlice([]string{
				string(fluxconfiguration.ScopeTypeCluster),
				string(fluxconfiguration.ScopeTypeNamespace),
			}, false),
		},

		"source_kind": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			ValidateFunc: validation.StringInSlice([]string{
				string(fluxconfiguration.SourceKindTypeGitRepository),
				string(fluxconfiguration.SourceKindTypeBucket),
				string(fluxconfiguration.SourceKindTypeAzureBlob),
			}, false),
		},

		"suspend": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
		},
	}
}

func (r KubernetesConfigurationFluxConfigurationResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"compliance_state": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"error_message": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"repository_public_key": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"source_synced_commit_id": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"source_updated_at": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"status_updated_at": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"statuses": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"applied_by": {
						Type:     pluginsdk.TypeList,
						Computed: true,
						MaxItems: 1,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"name": {
									Type:     pluginsdk.TypeString,
									Computed: true,
								},

								"namespace": {
									Type:     pluginsdk.TypeString,
									Computed: true,
								},
							},
						},
					},

					"compliance_state": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},

					"helm_release_properties": {
						Type:     pluginsdk.TypeList,
						Computed: true,
						MaxItems: 1,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"failure_count": {
									Type:     pluginsdk.TypeInt,
									Computed: true,
								},

								"helm_chart_ref": {
									Type:     pluginsdk.TypeList,
									Computed: true,
									MaxItems: 1,
									Elem: &pluginsdk.Resource{
										Schema: map[string]*pluginsdk.Schema{
											"name": {
												Type:     pluginsdk.TypeString,
												Computed: true,
											},

											"namespace": {
												Type:     pluginsdk.TypeString,
												Computed: true,
											},
										},
									},
								},

								"install_failure_count": {
									Type:     pluginsdk.TypeInt,
									Computed: true,
								},

								"last_revision_applied": {
									Type:     pluginsdk.TypeInt,
									Computed: true,
								},

								"upgrade_failure_count": {
									Type:     pluginsdk.TypeInt,
									Computed: true,
								},
							},
						},
					},

					"kind": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},

					"name": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},

					"namespace": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},

					"status_conditions": {
						Type:     pluginsdk.TypeList,
						Computed: true,
						MaxItems: 1,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"last_transition_time": {
									Type:     pluginsdk.TypeString,
									Computed: true,
								},

								"message": {
									Type:     pluginsdk.TypeString,
									Computed: true,
								},

								"reason": {
									Type:     pluginsdk.TypeString,
									Computed: true,
								},

								"status": {
									Type:     pluginsdk.TypeString,
									Computed: true,
								},

								"type": {
									Type:     pluginsdk.TypeString,
									Computed: true,
								},
							},
						},
					},
				},
			},
		},
	}
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
			id := fluxconfiguration.NewFluxConfigurationID(subscriptionId, model.ResourceGroupName, model.ClusterRp, model.ClusterResourceName, model.ClusterName, model.Name)
			existing, err := client.Get(ctx, id)
			if err != nil && !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for existing %s: %+v", id, err)
			}

			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			properties := &fluxconfiguration.FluxConfiguration{
				Properties: &fluxconfiguration.FluxConfigurationProperties{
					ConfigurationProtectedSettings: &model.ConfigurationProtectedSettings,
					Scope:                          &model.Scope,
					SourceKind:                     &model.SourceKind,
					Suspend:                        &model.Suspend,
				},
			}

			azureBlobValue, err := expandAzureBlobDefinitionModel(model.AzureBlob)
			if err != nil {
				return err
			}

			properties.Properties.AzureBlob = azureBlobValue

			bucketValue, err := expandBucketDefinitionModel(model.Bucket)
			if err != nil {
				return err
			}

			properties.Properties.Bucket = bucketValue

			gitRepositoryValue, err := expandGitRepositoryDefinitionModel(model.GitRepository)
			if err != nil {
				return err
			}

			properties.Properties.GitRepository = gitRepositoryValue

			if model.Kustomizations != "" {
				var kustomizationsValue map[string]kubernetesconfiguration.KustomizationDefinition
				err = json.Unmarshal([]byte(model.Kustomizations), &kustomizationsValue)
				if err != nil {
					return err
				}
				properties.Properties.Kustomizations = &kustomizationsValue
			}

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

			if metadata.ResourceData.HasChange("azure_blob") {
				azureBlobValue, err := expandAzureBlobDefinitionModel(model.AzureBlob)
				if err != nil {
					return err
				}

				properties.Properties.AzureBlob = azureBlobValue
			}

			if metadata.ResourceData.HasChange("bucket") {
				bucketValue, err := expandBucketDefinitionModel(model.Bucket)
				if err != nil {
					return err
				}

				properties.Properties.Bucket = bucketValue
			}

			if metadata.ResourceData.HasChange("configuration_protected_settings") {
				properties.Properties.ConfigurationProtectedSettings = &model.ConfigurationProtectedSettings
			}

			if metadata.ResourceData.HasChange("git_repository") {
				gitRepositoryValue, err := expandGitRepositoryDefinitionModel(model.GitRepository)
				if err != nil {
					return err
				}

				properties.Properties.GitRepository = gitRepositoryValue
			}

			if metadata.ResourceData.HasChange("kustomizations") {
				var kustomizationsValue map[string]kubernetesconfiguration.KustomizationDefinition
				err := json.Unmarshal([]byte(model.Kustomizations), &kustomizationsValue)
				if err != nil {
					return err
				}

				properties.Properties.Kustomizations = &kustomizationsValue
			}

			if metadata.ResourceData.HasChange("namespace") {
				properties.Properties.Namespace = &model.Namespace
			}

			if metadata.ResourceData.HasChange("scope") {
				properties.Properties.Scope = &model.Scope
			}

			if metadata.ResourceData.HasChange("source_kind") {
				properties.Properties.SourceKind = &model.SourceKind
			}

			if metadata.ResourceData.HasChange("suspend") {
				properties.Properties.Suspend = &model.Suspend
			}

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
				ClusterRp:           id.ClusterRp,
				ClusterResourceName: id.ClusterResourceName,
				ClusterName:         id.ClusterName,
			}

			if properties := model.Properties; properties != nil {
				azureBlobValue, err := flattenAzureBlobDefinitionModel(properties.AzureBlob)
				if err != nil {
					return err
				}

				state.AzureBlob = azureBlobValue

				bucketValue, err := flattenBucketDefinitionModel(properties.Bucket)
				if err != nil {
					return err
				}

				state.Bucket = bucketValue

				if properties.ComplianceState != nil {
					state.ComplianceState = *properties.ComplianceState
				}

				if properties.ConfigurationProtectedSettings != nil {
					state.ConfigurationProtectedSettings = *properties.ConfigurationProtectedSettings
				}

				if properties.ErrorMessage != nil {
					state.ErrorMessage = *properties.ErrorMessage
				}

				gitRepositoryValue, err := flattenGitRepositoryDefinitionModel(properties.GitRepository)
				if err != nil {
					return err
				}

				state.GitRepository = gitRepositoryValue

				if properties.Kustomizations != nil && *properties.Kustomizations != nil {

					kustomizationsValue, err := json.Marshal(*properties.Kustomizations)
					if err != nil {
						return err
					}

					state.Kustomizations = string(kustomizationsValue)
				}

				if properties.Namespace != nil {
					state.Namespace = *properties.Namespace
				}

				if properties.RepositoryPublicKey != nil {
					state.RepositoryPublicKey = *properties.RepositoryPublicKey
				}

				if properties.Scope != nil {
					state.Scope = *properties.Scope
				}

				if properties.SourceKind != nil {
					state.SourceKind = *properties.SourceKind
				}

				if properties.SourceSyncedCommitId != nil {
					state.SourceSyncedCommitId = *properties.SourceSyncedCommitId
				}

				if properties.SourceUpdatedAt != nil {
					state.SourceUpdatedAt = *properties.SourceUpdatedAt
				}

				if properties.StatusUpdatedAt != nil {
					state.StatusUpdatedAt = *properties.StatusUpdatedAt
				}

				statusesValue, err := flattenObjectStatusDefinitionModel(properties.Statuses)
				if err != nil {
					return err
				}

				state.Statuses = statusesValue

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

func expandAzureBlobDefinitionModel(inputList []AzureBlobDefinitionModel) (*fluxconfiguration.AzureBlobDefinition, error) {
	if len(inputList) == 0 {
		return nil, nil
	}

	input := &inputList[0]
	output := fluxconfiguration.AzureBlobDefinition{
		AccountKey:            &input.AccountKey,
		ContainerName:         &input.ContainerName,
		LocalAuthRef:          &input.LocalAuthRef,
		SasToken:              &input.SasToken,
		SyncIntervalInSeconds: &input.SyncIntervalInSeconds,
		TimeoutInSeconds:      &input.TimeoutInSeconds,
		Url:                   &input.Url,
	}

	managedIdentityValue, err := expandManagedIdentityDefinitionModel(input.ManagedIdentity)
	if err != nil {
		return nil, err
	}

	output.ManagedIdentity = managedIdentityValue

	servicePrincipalValue, err := expandServicePrincipalDefinitionModel(input.ServicePrincipal)
	if err != nil {
		return nil, err
	}

	output.ServicePrincipal = servicePrincipalValue

	return &output, nil
}

func expandManagedIdentityDefinitionModel(inputList []ManagedIdentityDefinitionModel) (*fluxconfiguration.ManagedIdentityDefinition, error) {
	if len(inputList) == 0 {
		return nil, nil
	}

	input := &inputList[0]
	output := fluxconfiguration.ManagedIdentityDefinition{
		ClientId: &input.ClientId,
	}

	return &output, nil
}

func expandServicePrincipalDefinitionModel(inputList []ServicePrincipalDefinitionModel) (*fluxconfiguration.ServicePrincipalDefinition, error) {
	if len(inputList) == 0 {
		return nil, nil
	}

	input := &inputList[0]
	output := fluxconfiguration.ServicePrincipalDefinition{
		ClientCertificate:          &input.ClientCertificate,
		ClientCertificatePassword:  &input.ClientCertificatePassword,
		ClientCertificateSendChain: &input.ClientCertificateSendChain,
		ClientId:                   &input.ClientId,
		ClientSecret:               &input.ClientSecret,
		TenantId:                   &input.TenantId,
	}

	return &output, nil
}

func expandBucketDefinitionModel(inputList []BucketDefinitionModel) (*fluxconfiguration.BucketDefinition, error) {
	if len(inputList) == 0 {
		return nil, nil
	}

	input := &inputList[0]
	output := fluxconfiguration.BucketDefinition{
		AccessKey:             &input.AccessKey,
		BucketName:            &input.BucketName,
		Insecure:              &input.Insecure,
		LocalAuthRef:          &input.LocalAuthRef,
		SyncIntervalInSeconds: &input.SyncIntervalInSeconds,
		TimeoutInSeconds:      &input.TimeoutInSeconds,
		Url:                   &input.Url,
	}

	return &output, nil
}

func expandGitRepositoryDefinitionModel(inputList []GitRepositoryDefinitionModel) (*fluxconfiguration.GitRepositoryDefinition, error) {
	if len(inputList) == 0 {
		return nil, nil
	}

	input := &inputList[0]
	output := fluxconfiguration.GitRepositoryDefinition{
		HttpsCACert:           &input.HttpsCACert,
		HttpsUser:             &input.HttpsUser,
		LocalAuthRef:          &input.LocalAuthRef,
		SshKnownHosts:         &input.SshKnownHosts,
		SyncIntervalInSeconds: &input.SyncIntervalInSeconds,
		TimeoutInSeconds:      &input.TimeoutInSeconds,
		Url:                   &input.Url,
	}

	repositoryRefValue, err := expandRepositoryRefDefinitionModel(input.RepositoryRef)
	if err != nil {
		return nil, err
	}

	output.RepositoryRef = repositoryRefValue

	return &output, nil
}

func expandRepositoryRefDefinitionModel(inputList []RepositoryRefDefinitionModel) (*fluxconfiguration.RepositoryRefDefinition, error) {
	if len(inputList) == 0 {
		return nil, nil
	}

	input := &inputList[0]
	output := fluxconfiguration.RepositoryRefDefinition{
		Branch: &input.Branch,
		Commit: &input.Commit,
		Semver: &input.Semver,
		Tag:    &input.Tag,
	}

	return &output, nil
}

func flattenAzureBlobDefinitionModel(input *fluxconfiguration.AzureBlobDefinition) ([]AzureBlobDefinitionModel, error) {
	var outputList []AzureBlobDefinitionModel
	if input == nil {
		return outputList, nil
	}

	output := AzureBlobDefinitionModel{}

	if input.AccountKey != nil {
		output.AccountKey = *input.AccountKey
	}

	if input.ContainerName != nil {
		output.ContainerName = *input.ContainerName
	}

	if input.LocalAuthRef != nil {
		output.LocalAuthRef = *input.LocalAuthRef
	}

	managedIdentityValue, err := flattenManagedIdentityDefinitionModel(input.ManagedIdentity)
	if err != nil {
		return nil, err
	}

	output.ManagedIdentity = managedIdentityValue

	if input.SasToken != nil {
		output.SasToken = *input.SasToken
	}

	servicePrincipalValue, err := flattenServicePrincipalDefinitionModel(input.ServicePrincipal)
	if err != nil {
		return nil, err
	}

	output.ServicePrincipal = servicePrincipalValue

	if input.SyncIntervalInSeconds != nil {
		output.SyncIntervalInSeconds = *input.SyncIntervalInSeconds
	}

	if input.TimeoutInSeconds != nil {
		output.TimeoutInSeconds = *input.TimeoutInSeconds
	}

	if input.Url != nil {
		output.Url = *input.Url
	}

	return append(outputList, output), nil
}

func flattenManagedIdentityDefinitionModel(input *fluxconfiguration.ManagedIdentityDefinition) ([]ManagedIdentityDefinitionModel, error) {
	var outputList []ManagedIdentityDefinitionModel
	if input == nil {
		return outputList, nil
	}

	output := ManagedIdentityDefinitionModel{}

	if input.ClientId != nil {
		output.ClientId = *input.ClientId
	}

	return append(outputList, output), nil
}

func flattenServicePrincipalDefinitionModel(input *fluxconfiguration.ServicePrincipalDefinition) ([]ServicePrincipalDefinitionModel, error) {
	var outputList []ServicePrincipalDefinitionModel
	if input == nil {
		return outputList, nil
	}

	output := ServicePrincipalDefinitionModel{}

	if input.ClientCertificate != nil {
		output.ClientCertificate = *input.ClientCertificate
	}

	if input.ClientCertificatePassword != nil {
		output.ClientCertificatePassword = *input.ClientCertificatePassword
	}

	if input.ClientCertificateSendChain != nil {
		output.ClientCertificateSendChain = *input.ClientCertificateSendChain
	}

	if input.ClientId != nil {
		output.ClientId = *input.ClientId
	}

	if input.ClientSecret != nil {
		output.ClientSecret = *input.ClientSecret
	}

	if input.TenantId != nil {
		output.TenantId = *input.TenantId
	}

	return append(outputList, output), nil
}

func flattenBucketDefinitionModel(input *fluxconfiguration.BucketDefinition) ([]BucketDefinitionModel, error) {
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

	return append(outputList, output), nil
}

func flattenGitRepositoryDefinitionModel(input *fluxconfiguration.GitRepositoryDefinition) ([]GitRepositoryDefinitionModel, error) {
	var outputList []GitRepositoryDefinitionModel
	if input == nil {
		return outputList, nil
	}

	output := GitRepositoryDefinitionModel{}

	if input.HttpsCACert != nil {
		output.HttpsCACert = *input.HttpsCACert
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
		output.SshKnownHosts = *input.SshKnownHosts
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

func flattenObjectStatusDefinitionModel(inputList *[]fluxconfiguration.ObjectStatusDefinition) ([]ObjectStatusDefinitionModel, error) {
	var outputList []ObjectStatusDefinitionModel
	if inputList == nil {
		return outputList, nil
	}

	for _, input := range *inputList {
		output := ObjectStatusDefinitionModel{}

		appliedByValue, err := flattenObjectReferenceDefinitionModel(input.AppliedBy)
		if err != nil {
			return nil, err
		}

		output.AppliedBy = appliedByValue

		if input.ComplianceState != nil {
			output.ComplianceState = *input.ComplianceState
		}

		helmReleasePropertiesValue, err := flattenHelmReleasePropertiesDefinitionModel(input.HelmReleaseProperties)
		if err != nil {
			return nil, err
		}

		output.HelmReleaseProperties = helmReleasePropertiesValue

		if input.Kind != nil {
			output.Kind = *input.Kind
		}

		if input.Name != nil {
			output.Name = *input.Name
		}

		if input.Namespace != nil {
			output.Namespace = *input.Namespace
		}

		statusConditionsValue, err := flattenObjectStatusConditionDefinitionModel(input.StatusConditions)
		if err != nil {
			return nil, err
		}

		output.StatusConditions = statusConditionsValue

		outputList = append(outputList, output)
	}

	return outputList, nil
}

func flattenObjectReferenceDefinitionModel(input *fluxconfiguration.ObjectReferenceDefinition) ([]ObjectReferenceDefinitionModel, error) {
	var outputList []ObjectReferenceDefinitionModel
	if input == nil {
		return outputList, nil
	}

	output := ObjectReferenceDefinitionModel{}

	if input.Name != nil {
		output.Name = *input.Name
	}

	if input.Namespace != nil {
		output.Namespace = *input.Namespace
	}

	return append(outputList, output), nil
}

func flattenHelmReleasePropertiesDefinitionModel(input *fluxconfiguration.HelmReleasePropertiesDefinition) ([]HelmReleasePropertiesDefinitionModel, error) {
	var outputList []HelmReleasePropertiesDefinitionModel
	if input == nil {
		return outputList, nil
	}

	output := HelmReleasePropertiesDefinitionModel{}

	if input.FailureCount != nil {
		output.FailureCount = *input.FailureCount
	}

	helmChartRefValue, err := flattenObjectReferenceDefinitionModel(input.HelmChartRef)
	if err != nil {
		return nil, err
	}

	output.HelmChartRef = helmChartRefValue

	if input.InstallFailureCount != nil {
		output.InstallFailureCount = *input.InstallFailureCount
	}

	if input.LastRevisionApplied != nil {
		output.LastRevisionApplied = *input.LastRevisionApplied
	}

	if input.UpgradeFailureCount != nil {
		output.UpgradeFailureCount = *input.UpgradeFailureCount
	}

	return append(outputList, output), nil
}

func flattenObjectStatusConditionDefinitionModel(inputList *[]fluxconfiguration.ObjectStatusConditionDefinition) ([]ObjectStatusConditionDefinitionModel, error) {
	var outputList []ObjectStatusConditionDefinitionModel
	if inputList == nil {
		return outputList, nil
	}

	for _, input := range *inputList {
		output := ObjectStatusConditionDefinitionModel{}

		if input.LastTransitionTime != nil {
			output.LastTransitionTime = *input.LastTransitionTime
		}

		if input.Message != nil {
			output.Message = *input.Message
		}

		if input.Reason != nil {
			output.Reason = *input.Reason
		}

		if input.Status != nil {
			output.Status = *input.Status
		}

		if input.Type != nil {
			output.Type = *input.Type
		}

		outputList = append(outputList, output)
	}

	return outputList, nil
}
