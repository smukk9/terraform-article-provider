# Provider Learnings


## Setting up provider for local testing

1. create/update ~/.terraformrc file with dev_overrides

```
provider_installation {

  dev_overrides {
    "article" = "/Users/{home}/github/terrafrom-article-provider/terraform-provider-article/"
  }
  direct {}
}

```

2. Ensure the name of the provider used like below, terraform will take care of matching. 

```
terraform {
  required_providers {
    article = {
      source = "article" # This should match the dev_overrides in .terraformrc
    }
  }
}

```

3. Ensure you build the executable of the provider code with terraform-provider-<name> in our case terraform-provider-article, use the taskfile



4. `plugin.Serve` starts the provider server and listens for requests from Terraform. This architecture allows Terraform to continue functioning even if your provider fails.

```go
func main() {
// This tells Terraform how to use your provider.
    plugin.Serve(&plugin.ServeOpts{
    ProviderFunc: provider,
    })
}

```

5. The `meta` variable in the `resourceArticleCreate` function contains the provider configuration that was passed when Terraform initialized the provider. This configuration can be accessed as needed. If you don’t set this correctly, the default values will result in an error.

6. A best practice in Terraform is to ensure your `createObject` (i.e., `resourceArticleCreate`) calls `resourceArticleRead` at the end and returns the error from the read function. This ensures that Terraform considers the entire apply process done when it can read back and update the "known after apply" values. This will help keep the state in sync with the object created by the `createFunc`.

7. You can use `ConfigureFunc` to define your provider's configuration. In this case, check the `configure` function for this setup.

8. After successfully creating a resource, you **must** set the `d.SetId(fmt.Sprintf("%v", createdArticle["id"])))` so that Terraform can identify it in the state file. If the id is not set correctly, Terraform won't know the resource exists, and it will consider it as not created.

9. During a `terraform plan`, Terraform checks the state file `(terraform.tfstate)` to see if the resource already exists. It uses the ID and other resource-specific attributes to decide whether the resource needs to be created, updated, or deleted. The read function is invoked to retrieve the current state of the resource so that Terraform can compare it with the desired state defined in the configuration.

10. When you write a provider you job is ensure you have all CRUD function along with the `Schema`, `ResourcesMap`, `DatasourcesMap` and ConfigureFunc are defined with in the `provider` function