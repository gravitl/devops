variable "clients" {
  type= list(string)

}

variable "do_token" {}
variable "pvt_key" {}
variable "extclient" {}
variable "docker" {}
variable "egress" {}
variable "ssh_public_key_path" {
  description = "Local public ssh key"
  default = "/root/.ssh/"
}
