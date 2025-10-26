# Rules
## Core
| Rule | Summary |
|--------|---------|
| [Avoid using hashicorp/null provider](core/avoid_null_provider.md) | With newer Terraform version, use locals and terraform_data as native replacement for hashicorp/null |
| [Enforce Variable Description](core/enforce_variable_description.md) | To understand what that variable does (even if it seems trivial), always add a description |
| [File Naming](core/file_naming.md) | File naming should follow a strict convention. |
| [Naming Convention](core/naming_convention.md) | Terraform names should only contain lowercase alphanumeric characters and underscores. |
| [Required Provider Must Be Declared](core/required_provider_must_be_declared.md) | All providers used in resources or data sources are declared in the terraform.required_providers block. |
