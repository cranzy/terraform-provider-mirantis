
terraform {
  required_version = ">= 1.0.0"
  required_providers {
    mcc = {
      version = "= 0.9.0"
      source  = "mirantis.com/providers/mcc"
    }
  }
}
