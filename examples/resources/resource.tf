resource "fabric_workspace" "examplews" {
  name = "testtf"
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

}
