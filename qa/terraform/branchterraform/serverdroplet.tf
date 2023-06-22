
# ssh into passed in server and reset netmaker on it with the branch
resource "null_resource" "terraformnetmakerserver" {
  
  #get a connection to the passed in droplet
  connection {
    host = "server.${var.server}.clustercat.com"
    user = "root"
    type = "ssh"
    private_key = var.pvt_key
    timeout = "2m"
  }
  
  #use remote-exec to install netmaker onto the server
  provisioner "remote-exec" {
    inline = [
      "export PATH=$PATH:/usr/bin",
      "wget https://raw.githubusercontent.com/gravitl/netmaker/${var.branch}/scripts/nm-quick.sh",
      "chmod +x nm-quick.sh",
      "chmod +x nm-quick.sh",
      "apt-get -y update",
      "apt-get -y update",
      "bash nm-quick.sh -a -b local -t ${var.branch} -d ${var.server}.clustercat.com"
      
    ]
  }
}



# This null resource will scp the docker-compose over to the terraform server to gather required information
resource "null_resource" "getdockercompose" {

  depends_on = [digitalocean_droplet.terraformnetmakerserver]

  provisioner "local-exec" {
     command = "scp -o StrictHostKeyChecking=no -o UserKnownHostsFile=/dev/null root@server.${var.server}.clustercat.com:/root/netmaker.env ."
  }
}

# This null_resource will run a shell script that will extract the information from the docker-compose and populate a txt file
resource "null_resource" "getserverinfo" {
  
  depends_on = [ digitalocean_droplet.terraformnetmakerserver, null_resource.getdockercompose, local_file.ipaddresses, local_file.extipaddresses, local_file.dockeripaddresses, local_file.egressipaddresses]
  provisioner "local-exec" {
    interpreter = ["/bin/bash" ,"-c"]
    command = "sudo bash getserverinfo.sh ${var.do_tag} ${var.clientbranch}"
  }
}