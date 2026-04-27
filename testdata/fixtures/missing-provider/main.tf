terraform {
  required_version = ">= 1.0"
  required_providers {
    nonexistent = {
      source  = "example.com/nonexistent/nonexistent"
      version = "1.0.0"
    }
  }
}

resource "nonexistent_resource" "example" {
  name = "test"
}
