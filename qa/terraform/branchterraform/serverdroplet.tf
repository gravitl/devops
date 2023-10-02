
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
      "rm -rf netclient",
      "rm -rf netclient",
      "bash nm-quick.sh -a -b local -t ${var.branch} -d ${var.server}.clustercat.com",
      "snap install go --classic",
      "snap install go --classic",
      "DEBIAN_FRONTEND=noninteractive apt install -y wireguard-tools gcc",
      "DEBIAN_FRONTEND=noninteractive apt install -y wireguard-tools gcc",
      #remove the netclient binary fetched from install script. running twice to ensure removal.
      "rm netclient",
      "rm netclient",
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



# This null resource will scp the docker-compose over to the terraform server to gather required information
resource "null_resource" "getdockercompose" {

  depends_on = [null_resource.terraformnetmakerserver]

  provisioner "local-exec" {
     command = "scp -o StrictHostKeyChecking=no -o UserKnownHostsFile=/dev/null root@server.${var.server}.clustercat.com:/root/netmaker.env ."
  }
}

# This null_resource will run a shell script that will extract the information from the docker-compose and populate a txt file
resource "null_resource" "getserverinfo" {
  
  depends_on = [ null_resource.terraformnetmakerserver, null_resource.getdockercompose, local_file.ipaddresses, local_file.extipaddresses, local_file.dockeripaddresses, local_file.egressipaddresses]
  provisioner "local-exec" {
    interpreter = ["/bin/bash" ,"-c"]
    command = "sudo bash getserverinfo.sh ${var.do_tag} ${var.clientbranch}"
  }
}