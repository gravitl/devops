# this will creat n clients depending on the clients array in terraform.tfvars.
# you can add morre just by adding to the array.
resource "digitalocean_droplet" "clients" {
  count  = length(var.clients)
  image  = "ubuntu-22-10-x64"
  name   = var.clients[count.index]
  region = "nyc3"
  size   = "s-1vcpu-1gb"
  ipv6   = true
  ssh_keys = [
    data.digitalocean_ssh_key.terraform.id
  ]
  tags = [var.clients[count.index], var.branch]
  # creates ssh connection
  connection {
    host        = self.ipv4_address
    user        = "root"
    type        = "ssh"
    private_key = file(var.pvt_key)
    timeout     = "2m"
  }
  #installs tools needed for netclient installation
  provisioner "remote-exec" {
    inline = [
      "export PATH=$PATH:/usr/bin",
      # install netmaker
      "pwd",
      "apt-get -y update",
      "apt-get -y update",
      "snap install go --classic",
      "snap install go --classic",
      "apt install -y wireguard-tools gcc",
      "apt install -y wireguard-tools gcc"
    ]
  }
}

# gets the data of each client
data "digitalocean_droplet" "serverips" {
  count      = length(var.clients)
  name       = var.clients[count.index]
  depends_on = [digitalocean_droplet.clients]
}

# ouputs their ipaddres to a local file to be proccessed by getserverinfo.sh
resource "local_file" "ipaddresses" {
  depends_on = [data.digitalocean_droplet.serverips]
  count      = length(var.clients)
  content    = data.digitalocean_droplet.serverips[count.index].ipv4_address
  filename   = "ipaddress${var.clients[count.index]}.txt"

}