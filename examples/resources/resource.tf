resource "fabric_workspace" "myWorkspace" {
  name        = "Name of the Test-Workspace11"
  description = "Description Text of the Workspace."
}


terraform {
  required_providers {
    fabric = {
      source  = "hashicorp.com/edu/fabric"
      version = ">= 0.0.0"
    }
  }
}

provider "fabric" {
  client_secret   = ""
  client_id       = ""
  tenant_id       = ""
  subscription_id = ""
}

