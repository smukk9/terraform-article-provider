## Provider Learnings

1. `plugin.Serve` starts the provider server and listens for requests from Terraform. This architecture allows Terraform to continue functioning even if your provider fails.

```go
func main() {
// This tells Terraform how to use your provider.
    plugin.Serve(&plugin.ServeOpts{
    ProviderFunc: provider,
    })
}

```

2. The `meta` variable in the `resourceArticleCreate` function contains the provider configuration that was passed when Terraform initialized the provider. This configuration can be accessed as needed. If you don’t set this correctly, the default values will result in an error.

3. A best practice in Terraform is to ensure your `createObject` (i.e., `resourceArticleCreate`) calls `resourceArticleRead` at the end and returns the error from the read function. This ensures that Terraform considers the entire apply process done when it can read back and update the "known after apply" values. This will help keep the state in sync with the object created by the `createFunc`.

4. You can use `ConfigureFunc` to define your provider's configuration. In this case, check the `configure` function for this setup.

5. After successfully creating a resource, you **must** set the `d.SetId(fmt.Sprintf("%v", createdArticle["id"])))` so that Terraform can identify it in the state file. If the id is not set correctly, Terraform won't know the resource exists, and it will consider it as not created.

6. During a `terraform plan`, Terraform checks the state file `(terraform.tfstate)` to see if the resource already exists. It uses the ID and other resource-specific attributes to decide whether the resource needs to be created, updated, or deleted. The read function is invoked to retrieve the current state of the resource so that Terraform can compare it with the desired state defined in the configuration.

7. When you write a provider you job is ensure you have all CRUD function along with the `Schema`, `ResourcesMap`, `DatasourcesMap` and ConfigureFunc are defined with in the `provider` function