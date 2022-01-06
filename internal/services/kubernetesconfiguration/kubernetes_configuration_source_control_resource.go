package kubernetesconfiguration

import (
    "fmt"
    "github.com/Azure/azure-sdk-for-go/services/kubernetesconfiguration/mgmt/2021-03-01/kubernetesconfiguration"
    "github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
    "github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
    "github.com/hashicorp/terraform-provider-azurerm/internal/clients"
    "github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
    "github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
    "github.com/hashicorp/terraform-provider-azurerm/utils"
    "github.com/hashicorp/terraform-provider-azurerm/internal/services/kubernetesconfiguration/parse"
    "github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
    "log"
    "time"
)

func resourceKubernetesConfigurationSourceControl() *pluginsdk.Resource {
    return &pluginsdk.Resource{
        Create: resourceKubernetesconfigurationSourceControlConfigurationCreateUpdate,
        Read:   resourceKubernetesconfigurationSourceControlConfigurationRead,
        Update: resourceKubernetesconfigurationSourceControlConfigurationCreateUpdate,
        Delete: resourceKubernetesconfigurationSourceControlConfigurationDelete,

        Timeouts: &pluginsdk.ResourceTimeout{
            Create: pluginsdk.DefaultTimeout(30 * time.Minute),
            Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
            Update: pluginsdk.DefaultTimeout(30 * time.Minute),
            Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
        },

        Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
            _, err := parse.KubernetesConfigurationSourceControlID(id)
            return err
        }),


        Schema: map[string]*pluginsdk.Schema{
            "name": {
                Type:     pluginsdk.TypeString,
                Required: true,
                ForceNew: true,
                ValidateFunc:,
            },

            "resource_group_name": azure.SchemaResourceGroupName(),

            "cluster_name": {
                Type:     pluginsdk.TypeString,
                Required: true,
                ForceNew: true,
                ValidateFunc:,
            },

            "repository_url": {
                Type:     pluginsdk.TypeString,
                Required: true,
                ValidateFunc:,
            },

            "configuration_protected_settings": {
                Type:     pluginsdk.TypeMap,
                Optional: true,
                Elem: &pluginsdk.Schema{
                    Type:     pluginsdk.TypeString,
                },
            },

            "enable_helm_operator": {
                Type:     pluginsdk.TypeBool,
                Optional: true,
            },

            "helm_operator_properties": {
                Type:     pluginsdk.TypeList,
                Optional: true,
                MaxItems: 1,
                Elem: &pluginsdk.Resource{
                    Schema: map[string]*pluginsdk.Schema{
                        "chart_values": {
                            Type:     pluginsdk.TypeString,
                            Optional: true,
                        },

                        "chart_version": {
                            Type:     pluginsdk.TypeString,
                            Optional: true,
                        },
                    },
                },
            },

            "operator_instance_name": {
                Type:     pluginsdk.TypeString,
                Optional: true,
            },

            "operator_namespace": {
                Type:     pluginsdk.TypeString,
                Optional: true,
                ForceNew: true,
                Default: "default",
                ValidateFunc:,
            },

            "operator_params": {
                Type:     pluginsdk.TypeString,
                Optional: true,
            },

            "operator_scope": {
                Type:     pluginsdk.TypeString,
                Optional: true,
                ValidateFunc: validation.StringInSlice([]string{
                    string(kubernetesconfiguration.Cluster),
                    string(kubernetesconfiguration.Namespace),
                }, false),
                Default: kubernetesconfiguration.Cluster,
            },

            "operator_type": {
                Type:     pluginsdk.TypeString,
                Optional: true,
                ValidateFunc: validation.StringInSlice([]string{
                    string(kubernetesconfiguration.Flux),
                }, false),
            },

            "ssh_known_hosts_contents": {
                Type:     pluginsdk.TypeString,
                Optional: true,
            },

            "compliance_status": {
                Type:     pluginsdk.TypeList,
                Computed: true,
                MaxItems: 1,
                Elem: &pluginsdk.Resource{
                    Schema: map[string]*pluginsdk.Schema{
                        "compliance_state": {
                            Type:     pluginsdk.TypeString,
                            Computed: true,
                        },

                        "last_config_applied": {
                            Type:     pluginsdk.TypeString,
                            Computed: true,
                        },

                        "message": {
                            Type:     pluginsdk.TypeString,
                            Computed: true,
                        },

                        "message_level": {
                            Type:     pluginsdk.TypeString,
                            Computed: true,
                        },
                    },
                },
            },

            "repository_public_key": {
                Type:     pluginsdk.TypeString,
                Computed: true,
            },
        },
    }
}
func resourceKubernetesconfigurationSourceControlConfigurationCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
    subscriptionId := meta.(*clients.Client).Account.SubscriptionId
    client := meta.(*clients.Client).KubernetesConfiguration.SourceControlClient
    ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
    defer cancel()

    id := parse.NewKubernetesConfigurationSourceControlID(subscriptionId, d.Get("resource_group_name").(string), d.Get("cluster_rp").(string), d.Get("cluster_resource_name").(string), d.Get("cluster_name").(string), d.Get("name").(string))

    if d.IsNewResource() {
        existing, err :=client.Get(ctx, id.ResourceGroup, kubernetesconfiguration.Componentsintqrsparametersclusterrpparameterschema(id.ProviderName), kubernetesconfiguration.Components1c8uz7parametersclusterresourcenameparameterschema(id.{clusterRp}Name), id.{clusterResourceName}Name, id.SourceControlConfigurationName)
        if err != nil {
            if !utils.ResponseWasNotFound(existing.Response) {
                return fmt.Errorf("checking for existing %s: %+v", id, err)
            }
        }
        if existing.ID != nil && *existing.ID != "" {
            return tf.ImportAsExistsError("azurerm_kubernetesconfiguration_source_control_configuration", *existing.ID)
        }
    }


    props := kubernetesconfiguration.SourceControlConfiguration{
        SourceControlConfigurationProperties: &kubernetesconfiguration.SourceControlConfigurationProperties{
            ConfigurationProtectedSettings: utils.ExpandMapStringPtrString(d.Get("configuration_protected_settings").(map[string]interface{})),
            EnableHelmOperator: utils.Bool(d.Get("enable_helm_operator").(bool)),
            HelmOperatorProperties: expandSourceControlConfigurationHelmOperatorProperties(d.Get("helm_operator_properties").([]interface{})),
            OperatorInstanceName: utils.String(d.Get("operator_instance_name").(string)),
            OperatorNamespace: utils.String(d.Get("operator_namespace").(string)),
            OperatorParams: utils.String(d.Get("operator_params").(string)),
            OperatorScope: kubernetesconfiguration.OperatorScopeType(d.Get("operator_scope").(string)),
            OperatorType: kubernetesconfiguration.OperatorType(d.Get("operator_type").(string)),
            RepositoryURL: utils.String(d.Get("repository_url").(string)),
            SSHKnownHostsContents: utils.String(d.Get("ssh_known_hosts_contents").(string)),
        },
    }
    if _, err :=client.CreateOrUpdate(ctx, id.ResourceGroup, kubernetesconfiguration.Componentsintqrsparametersclusterrpparameterschema(id.ProviderName), kubernetesconfiguration.Components1c8uz7parametersclusterresourcenameparameterschema(id.{clusterRp}Name), id.{clusterResourceName}Name, id.SourceControlConfigurationName, props); err != nil {
        return fmt.Errorf("creating/updating %s: %+v", id, err)
    }

    d.SetId(id.ID())
    return resourceKubernetesconfigurationSourceControlConfigurationRead(d, meta)
}

func resourceKubernetesconfigurationSourceControlConfigurationRead(d *pluginsdk.ResourceData, meta interface{}) error {
    client := meta.(*clients.Client).Kubernetesconfiguration.SourceControlConfigurationClient
    ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
    defer cancel()

    id, err := parse.KubernetesconfigurationSourceControlConfigurationID(d.Id())
    if err != nil {
        return err
    }

    resp, err :=client.Get(ctx, id.ResourceGroup, kubernetesconfiguration.Componentsintqrsparametersclusterrpparameterschema(id.ProviderName), kubernetesconfiguration.Components1c8uz7parametersclusterresourcenameparameterschema(id.{clusterRp}Name), id.{clusterResourceName}Name, id.SourceControlConfigurationName)
    if err != nil {
        if utils.ResponseWasNotFound(resp.Response) {
            log.Printf("[INFO] kubernetesconfiguration %q does not exist - removing from state", d.Id())
            d.SetId("")
            return nil
        }
        return fmt.Errorf("retrieving %s: %+v", id, err)
    }
    d.Set("name", id.SourceControlConfigurationName)
    d.Set("resource_group_name", id.ResourceGroup)
    d.Set("cluster_name", id.{clusterResourceName}Name)
    d.Set("cluster_resource_name", id.{clusterRp}Name)
    d.Set("cluster_rp", id.ProviderName)
    if props := resp.SourceControlConfigurationProperties; props != nil {
        d.Set("configuration_protected_settings", utils.FlattenMapStringPtrString(props.ConfigurationProtectedSettings))
        d.Set("enable_helm_operator", props.EnableHelmOperator)
        if err := d.Set("helm_operator_properties", flattenSourceControlConfigurationHelmOperatorProperties(props.HelmOperatorProperties)); err != nil {
            return fmt.Errorf("setting `helm_operator_properties`: %+v", err)
        }
        d.Set("operator_instance_name", props.OperatorInstanceName)
        d.Set("operator_namespace", props.OperatorNamespace)
        d.Set("operator_params", props.OperatorParams)
        d.Set("operator_scope", props.OperatorScope)
        d.Set("operator_type", props.OperatorType)
        d.Set("repository_url", props.RepositoryURL)
        d.Set("ssh_known_hosts_contents", props.SSHKnownHostsContents)
        if err := d.Set("compliance_status", flattenSourceControlConfigurationComplianceStatus(props.ComplianceStatus)); err != nil {
            return fmt.Errorf("setting `compliance_status`: %+v", err)
        }
        d.Set("repository_public_key", props.RepositoryPublicKey)
    }
    return nil
}

func resourceKubernetesconfigurationSourceControlConfigurationDelete(d *pluginsdk.ResourceData, meta interface{}) error {
    client := meta.(*clients.Client).Kubernetesconfiguration.SourceControlConfigurationClient
    ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
    defer cancel()

    id, err := parse.KubernetesconfigurationSourceControlConfigurationID(d.Id())
    if err != nil {
        return err
    }

    future, err :=client.Delete(ctx, id.ResourceGroup, kubernetesconfiguration.Componentsintqrsparametersclusterrpparameterschema(id.ProviderName), kubernetesconfiguration.Components1c8uz7parametersclusterresourcenameparameterschema(id.{clusterRp}Name), id.{clusterResourceName}Name, id.SourceControlConfigurationName)
    if err != nil {
        return fmt.Errorf("deleting %s: %+v", id, err)
    }

    if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
        return fmt.Errorf("waiting for deletion of the %s: %+v", id, err)
    }
    return nil
}

func expandSourceControlConfigurationHelmOperatorProperties(input []interface{}) *kubernetesconfiguration.HelmOperatorProperties {
    if len(input) == 0 || input[0] == nil {
        return nil
    }
    v := input[0].(map[string]interface{})
    return &kubernetesconfiguration.HelmOperatorProperties{
        ChartVersion: utils.String(v["chart_version"].(string)),
        ChartValues: utils.String(v["chart_values"].(string)),
    }
}

func flattenSourceControlConfigurationHelmOperatorProperties(input *kubernetesconfiguration.HelmOperatorProperties) []interface{} {
    if input == nil {
        return make([]interface{}, 0)
    }

    var chartValues string
    if input.ChartValues != nil {
        chartValues = *input.ChartValues
    }
    var chartVersion string
    if input.ChartVersion != nil {
        chartVersion = *input.ChartVersion
    }
    return []interface{}{
        map[string]interface{}{
            "chart_values": chartValues,
            "chart_version": chartVersion,
        },
    }
}

func flattenSourceControlConfigurationComplianceStatus(input *kubernetesconfiguration.ComplianceStatus) []interface{} {
    if input == nil {
        return make([]interface{}, 0)
    }

    var complianceState kubernetesconfiguration.ComplianceStateType
    if input.ComplianceState != "" {
        complianceState = input.ComplianceState
    }
    var lastConfigApplied string
    if input.LastConfigApplied != nil {
        lastConfigApplied = input.LastConfigApplied.Format(time.RFC3339)
    }
    var message string
    if input.Message != nil {
        message = *input.Message
    }
    var messageLevel kubernetesconfiguration.MessageLevelType
    if input.MessageLevel != "" {
        messageLevel = input.MessageLevel
    }
    return []interface{}{
        map[string]interface{}{
            "compliance_state": complianceState,
            "last_config_applied": lastConfigApplied,
            "message": message,
            "message_level": messageLevel,
        },
    }
}
