# tfcoach-ignore-file:core.naming_convention
resource "aws_s3_bucket" "This" {}
# tfcoach-ignore: core.required_provider_must_be_declared,core.file_naming
# tfcoach-ignore: core.required_provider_must_be_declared




# tfcoach-ignore:core.avoid_type_in_name
data "aws_s3_bucket" "ignored_s3" {}

# tfcoach-ignore: core.required_provider_must_be_declared
# tfcoach-ignore:core.file_naming





data "aws_s3_bucket" "ignored" {}
data "aws_s3_bucket" "NotIgnored" {}
