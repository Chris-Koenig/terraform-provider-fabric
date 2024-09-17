

terraform {
  required_providers {
    fabric = {
      source  = "hashicorp.com/edu/fabric"
      version = ">= 0.0.0"
    }
  }
}

provider "fabric" {

}

data "fabric_workspace" "example" {
  id = "5de8160d-8474-4c75-8faf-fef3bbe7f83d"
}

output "workspace_id" {
  value = data.fabric_workspace.example.id
}
