

terraform {
  required_providers {
    fabric = {
      source  = "hashicorp.com/edu/fabric"
      version = ">= 0.0.0"
    }
  }
}

provider "fabric" {
  base_url        = "api.fabric.microsoft.com"
  client_secret   = "5Og8Q~schfyvm2UfxZJIesL23ynPdvm~HmCg5bE5"
  client_id       = "227f31e9-97ff-4fef-933d-c86c671d766a"
  tenant_id       = "08e2c9c5-f65a-4cd8-8d35-a02272e142d6"
  subscription_id = "8c62f590-4f16-4727-a2a7-f7a0304f308d"
}

data "fabric_workspace" "example" {
  id = "c45f1037-a6ed-4316-9fe9-49b2ea54bc87"
}

output "workspace_id" {
  value = data.fabric_workspace.example.id
}
