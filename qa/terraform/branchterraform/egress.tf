

resource "digitalocean_droplet" "egress" {
  image = "ubuntu-22-04-x64"
  name = var.egress
  size = "s-2vcpu-2gb"
  ipv6 = true
  ssh_keys = [
    for v in data.digitalocean_ssh_keys.keys.ssh_keys : v.id
  ]
  tags = [var.egress, var.do_tag]

  connection {
    host = self.ipv4_address
    user = "root"
    type = "ssh"
    private_key = var.pvt_key
    timeout = "2m"
  }
  
  provisioner "remote-exec" {
    inline = [
      "export PATH=$PATH:/usr/bin",
      # install netclient
      "DEBIAN_FRONTEND=noninteractive apt-get -y update",
      "DEBIAN_FRONTEND=noninteractive apt-get -y update",
      "snap install go --classic",
      "snap install go --classic",
      "DEBIAN_FRONTEND=noninteractive apt install -y wireguard-tools gcc lxc",
      "DEBIAN_FRONTEND=noninteractive apt install -y wireguard-tools gcc lxc",
      "lxc-create -n container -t download -- -d ubuntu -r jammy -a amd64",
      "lxc-start container",
      "lxc-attach container -- ip a 10.0.3.183 dev eth0",
      "git clone https://www.github.com/gravitl/netclient",
      "cd netclient",
      "git checkout ${var.clientbranch}",
      "git pull origin ${var.clientbranch}",
      "go mod tidy",
      "go build .",
      "./netclient install"

    ]
  }
}

data "digitalocean_droplet" "egressserverip" {
   id = digitalocean_droplet.egress.id
   depends_on = [digitalocean_droplet.egress]
}

resource "local_file" "egressipaddresses" {
   depends_on = [data.digitalocean_droplet.egressserverip]
   content = data.digitalocean_droplet.egressserverip.ipv4_address
   filename = "ipaddress${var.egress}.txt"

}

