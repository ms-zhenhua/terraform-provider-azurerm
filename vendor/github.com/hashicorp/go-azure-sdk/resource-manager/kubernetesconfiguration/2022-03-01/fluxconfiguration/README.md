
## `github.com/hashicorp/go-azure-sdk/resource-manager/kubernetesconfiguration/2022-03-01/fluxconfiguration` Documentation

The `fluxconfiguration` SDK allows for interaction with the Azure Resource Manager Service `kubernetesconfiguration` (API Version `2022-03-01`).

This readme covers example usages, but further information on [using this SDK can be found in the project root](https://github.com/hashicorp/go-azure-sdk/tree/main/docs).

### Import Path

```go
import "github.com/hashicorp/go-azure-sdk/resource-manager/kubernetesconfiguration/2022-03-01/fluxconfiguration"
```


### Client Initialization

```go
client := fluxconfiguration.NewFluxConfigurationClientWithBaseURI("https://management.azure.com")
client.Client.Authorizer = authorizer
if err != nil {
	// handle the error
}
```


### Example Usage: `FluxConfigurationClient.CreateOrUpdate`

```go
ctx := context.TODO()
id := fluxconfiguration.NewFluxConfigurationID("12345678-1234-9876-4563-123456789012", "example-resource-group", "clusterRpValue", "clusterResourceValue", "clusterValue", "fluxConfigurationValue")

payload := fluxconfiguration.FluxConfiguration{
	// ...
}

future, err := client.CreateOrUpdate(ctx, id, payload)
if err != nil {
	// handle the error
}
if err := future.Poller.PollUntilDone(); err != nil {
	// handle the error
}
```


### Example Usage: `FluxConfigurationClient.Delete`

```go
ctx := context.TODO()
id := fluxconfiguration.NewFluxConfigurationID("12345678-1234-9876-4563-123456789012", "example-resource-group", "clusterRpValue", "clusterResourceValue", "clusterValue", "fluxConfigurationValue")
future, err := client.Delete(ctx, id, fluxconfiguration.DefaultDeleteOperationOptions())
if err != nil {
	// handle the error
}
if err := future.Poller.PollUntilDone(); err != nil {
	// handle the error
}
```


### Example Usage: `FluxConfigurationClient.Get`

```go
ctx := context.TODO()
id := fluxconfiguration.NewFluxConfigurationID("12345678-1234-9876-4563-123456789012", "example-resource-group", "clusterRpValue", "clusterResourceValue", "clusterValue", "fluxConfigurationValue")
read, err := client.Get(ctx, id)
if err != nil {
	// handle the error
}
if model := read.Model; model != nil {
	// do something with the model/response object
}
```


### Example Usage: `FluxConfigurationClient.List`

```go
ctx := context.TODO()
id := fluxconfiguration.NewProviderID("12345678-1234-9876-4563-123456789012", "example-resource-group", "clusterRpValue", "clusterResourceValue", "clusterValue")
// alternatively `client.List(ctx, id)` can be used to do batched pagination
items, err := client.ListComplete(ctx, id)
if err != nil {
	// handle the error
}
for _, item := range items {
	// do something
}
```


### Example Usage: `FluxConfigurationClient.Update`

```go
ctx := context.TODO()
id := fluxconfiguration.NewFluxConfigurationID("12345678-1234-9876-4563-123456789012", "example-resource-group", "clusterRpValue", "clusterResourceValue", "clusterValue", "fluxConfigurationValue")

payload := fluxconfiguration.FluxConfigurationPatch{
	// ...
}

future, err := client.Update(ctx, id, payload)
if err != nil {
	// handle the error
}
if err := future.Poller.PollUntilDone(); err != nil {
	// handle the error
}
```
