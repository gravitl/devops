# sets up terraform with the digitalocean provider
terraform {
  required_providers {
    digitalocean = {
      source  = "digitalocean/digitalocean"
      version = "~> 2.0"
    }
  }
}


# sets the token aquired from digital ocean that will be exported.
provider "digitalocean" {
  token = var.do_token
}

data "digitalocean_ssh_key" "terraform" {
  # if you stored your ssh key under a different name in digital ocean, change this to that name.
  name = "terraform"
}