variable "clients" {
  type= list(string)
}

variable "branch" {
   default = "develop"
}

variable "clientbranch" {
   default = "develop"
}

variable "do_token" {}
variable "do_tag"{}
variable "extclient" {}
variable "docker" {}
variable "egress" {}
variable "pvt_key" {}