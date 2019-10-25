#
# Provider Configuration
#

provider "aws" {
  region  = "us-west-2"
  version = "~> 2.0"
}

terraform {
  backend "local" {}
}