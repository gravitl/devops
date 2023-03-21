resource "digitalocean_droplet" "clients" {
  count = length(var.clients)
  image = "ubuntu-22-10-x64"
  name = var.clients[count.index]
  region = "nyc3"
  size = "s-1vcpu-1gb"
  ipv6 = true
  ssh_keys = [
    for v in data.digitalocean_ssh_keys.keys.ssh_keys : v.id
  ]
  tags = [var.clients[count.index] ,var.clientbranch]  
  connection {
    host = self.ipv4_address
    user = "root"
    type = "ssh"
    private_key = file(var.pvt_key)
    timeout = "2m"
  }
  
  provisioner "remote-exec" {
    inline = [
      "export PATH=$PATH:/usr/bin",
      # install netmaker
      "apt-get -y update",
      "apt-get -y update",
      "snap install go --classic",
      "snap install go --classic",
      "apt install -y wireguard-tools gcc",
      "apt install -y wireguard-tools gcc"
    ]
  }
}

data "digitalocean_droplet" "serverips" {
   count = length(var.clients)
   id = digitalocean_droplet.clients[count.index].id
   depends_on = [digitalocean_droplet.clients]
}

resource "local_file" "ipaddresses" {
   depends_on = [data.digitalocean_droplet.serverips]
   count = length(var.clients)
   content = data.digitalocean_droplet.serverips[count.index].ipv4_address
   filename = "ipaddress${var.clients[count.index]}.txt"
   
}