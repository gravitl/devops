resource "digitalocean_droplet" "dockerclient" {
  image = "ubuntu-22-10-x64"
  name = var.docker
  region = "nyc3"
  size = "s-1vcpu-1gb"
  ipv6 = true
  ssh_keys = [
    for v in data.digitalocean_ssh_keys.keys.ssh_keys : v.id
  ]
  tags = [var.docker ,var.branch != "develop" ? var.branch : var.clientbranch]  
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
      # install netclient
      "pwd",
      "apt-get -y update",
      "apt-get -y update",
      "apt install -y wireguard-tools docker.io",
      "apt install -y wireguard-tools docker.io",
      "git clone https://www.github.com/gravitl/netclient",
      "cd netclient",
      "git checkout ${var.clientbranch}",
      "git pull origin ${var.clientbranch}",
      "docker build --build-arg version=${var.clientbranch} -t terraform/test . "
    ]
  }
}

data "digitalocean_droplet" "dockerserverip" {
   id = digitalocean_droplet.dockerclient.id
   depends_on = [digitalocean_droplet.dockerclient]
}

resource "local_file" "dockeripaddresses" {
   depends_on = [data.digitalocean_droplet.dockerserverip]
   content = data.digitalocean_droplet.dockerserverip.ipv4_address
   filename = "ipaddress${var.docker}.txt"

}