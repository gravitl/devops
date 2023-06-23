resource "digitalocean_droplet_snapshot" "egress_snapshot" {
  droplet_id = "347216123"
  name = "egresssnapshot${var.do_tag}"
}

resource "digitalocean_droplet" "egress" {
  depends_on = [
    digitalocean_droplet_snapshot.egress_snapshot
  ]
  image = digitalocean_droplet_snapshot.egress_snapshot.id
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
      "apt-get -y update",
      "apt-get -y update",
      "snap install go --classic",
      "snap install go --classic",
      "apt install -y wireguard-tools gcc",
      "apt install -y wireguard-tools gcc",
      "git clone https://www.github.com/gravitl/netclient",
      "cd netclient",
      "git checkout ${var.clientbranch}",
      "git pull origin ${var.clientbranch}",
      "go mod tidy",
      "go build -tags headless",
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

resource "null_resource" "remove_snapshot" {
  depends_on = [digitalocean_droplet.egress]
  provisioner "local-exec" {
      command = <<EOT
        curl -X DELETE \
        -H 'Content-Type: application/json' \
        -H "Authorization: Bearer ${var.do_token}" \
        "https://api.digitalocean.com/v2/snapshots/egresssnapshot/${var.do_tag}"
EOT
  }
}