# core.file_naming

Enforces that different type of terraform-resources are written
in the correct files

## Why

Consistent file-structure across multiple projects. Keeps scalability and
error analysis simple.

Example "I want to see resources that i load from external (data-resources). I
can just open data.tf"

### Mapping TF-Type to File

| Type      | Filename     |
|-----------|--------------|
| output    | outputs.tf   |
| variable  | variables.tf |
| locals    | locals.tf    |
| provider  | providers.tf |
| terraform | terraform.tf |
| data      | data.tf      |
