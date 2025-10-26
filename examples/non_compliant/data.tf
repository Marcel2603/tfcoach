data "null_data_source" "values" {
  inputs = {
    all_colors = concat(
      "green",
      "blue",
    )
  }
}