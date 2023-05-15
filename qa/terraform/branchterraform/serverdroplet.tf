# Create a new tag for the server. the branch tag will be
# a variable that is passed in at plan and apply
resource "digitalocean_tag" "server_tag" {
  name = "server"
}

# Create the droplet. This will create the droplet and
# setup netmaker locally on the droplet
resource "digitalocean_droplet" "terraformnetmakerserver" {
  image = "ubuntu-22-10-x64"
  name = "server"
  region = "nyc3"
  size = "s-2vcpu-2gb-intel"
  ssh_keys = [    for v in data.digitalocean_ssh_keys.keys.ssh_keys : v.id ] 
  tags   = [ digitalocean_tag.server_tag.id, var.do_tag]
  
  #get a connection to the created droplet
  connection {
    host = self.ipv4_address
    user = "root"
    type = "ssh"
    private_key = var.pvt_key
    timeout = "2m"
  }
  
  #use remote-exec to install netmaker onto the server
  provisioner "remote-exec" {
    inline = [
      "export PATH=$PATH:/usr/bin",
      # install netmaker
      "wget https://raw.githubusercontent.com/gravitl/netmaker/develop/scripts/nm-quick.sh",
      "apt-get -y update",
      "apt-get -y update",
      "apt install -y docker-compose docker.io",
      "apt install -y docker-compose docker.io",
      "bash nm-quick.sh -b local -t ${var.branch} -a"
      
    ]
  }
}

#this will get a reference to the ip of the droplet
data "digitalocean_droplet" "serverip" {
   id = digitalocean_droplet.terraformnetmakerserver.id
   depends_on = [digitalocean_droplet.terraformnetmakerserver]
}

# This null resource will scp the docker-compose over to the terraform server to gather required information
resource "null_resource" "getdockercompose" {

  depends_on = [data.digitalocean_droplet.serverip, digitalocean_droplet.terraformnetmakerserver]

  provisioner "local-exec" {
     command = "scp -o StrictHostKeyChecking=no -o UserKnownHostsFile=/dev/null root@${digitalocean_droplet.terraformnetmakerserver.ipv4_address}:/root/docker-compose.yml ."
  }
}

# This null_resource will run a shell script that will extract the information from the docker-compose and populate a txt file
resource "null_resource" "getserverinfo" {
  
  depends_on = [data.digitalocean_droplet.serverip, digitalocean_droplet.terraformnetmakerserver, null_resource.getdockercompose, local_file.ipaddresses, local_file.extipaddresses, local_file.dockeripaddresses, local_file.egressipaddresses]
  provisioner "local-exec" {
    interpreter = ["/bin/bash" ,"-c"]
    command = "sudo bash getserverinfo.sh ${var.do_tag} ${var.clientbranch}"
  }
}