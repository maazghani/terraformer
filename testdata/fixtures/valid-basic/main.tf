terraform {
  required_version = ">= 1.0"
}

resource "null_resource" "example" {
  triggers = {
    always_run = "${timestamp()}"
  }
}
