
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
      "git clone https://github.com/gravitl/netmaker.git",
      "cd netmaker",
      "git checkout ${var.branch}",
      "git pull origin ${var.branch}",
      "docker build --build-arg \"tags=ee\" -t gravitl/netmaker:${var.branch} .",

      "cd ~",

      "docker run -d --name nanomq --network testnet -p 1883:1883 -p 8083:8083 -p 8883:8883 emqx/nanomq:latest",

      "export SERVER_HOST=$(dig -4 myip.opendns.com @resolver1.opendns.com +short || curl -s ifconfig.me)",
      "export NETMAKER_BASE_DOMAIN=${var.server}.clustercat.com",
      "export MASTER_KEY=$(tr -dc A-Za-z0-9 </dev/urandom | head -c 30)",

      "touch netmaker.env",

      "echo \"SERVER_HOST=api.${var.server}.clustercat.com\" >> netmaker.env",
      "echo \"MASTER_KEY=$MASTER_KEY\" >> netmaker.env",
      "echo \"NETMAKER_BASE_DOMAIN=$NETMAKER_BASE_DOMAIN\" >> netmaker.env",
      "echo \"NM_DOMAIN=${var.server}.clustercat.com\" >> netmaker.env"

      "docker run -d --name netmaker -e SERVER_BROKER_ENDPOINT=\"nanomq:1883\" --env-file netmaker.env --network testnet -p 8081:8081 gravitl/netmaker:${var.branch}",
      "docker restart caddy",
        
      "nmctl context set default --endpoint=\"https://api.$NETMAKER_BASE_DOMAIN\" --master_key=\"$MASTER_KEY\"",
      "nmctl context use default",
      "nmctl network create --name netmaker --ipv4_addr 10.101.0.0/16",
      "export TOKEN=$(nmctl enrollment_key create --tags netmaker --unlimited --networks netmaker | jq -r .token)",

      # "export PATH=$PATH:/usr/bin",
      # "wget https://raw.githubusercontent.com/gravitl/devops/${var.devopsbranch}/qa/nm-quick.sh",
      # "chmod +x nm-quick.sh",
      # "chmod +x nm-quick.sh",
      # "rm -rf netclient",
      # "rm -rf netclient",
      # "bash nm-quick.sh -a -b local -t ${var.branch} -d ${var.server}.clustercat.com",
      # "snap install go --classic",
      # "snap install go --classic",
      "DEBIAN_FRONTEND=noninteractive apt install -y wireguard-tools gcc",
      "DEBIAN_FRONTEND=noninteractive apt install -y wireguard-tools gcc",
      # remove the netclient binary fetched from install script. running twice to ensure removal.
      "rm netclient",
      "rm netclient",
      "git clone https://www.github.com/gravitl/netclient",
      "cd netclient",
      "git checkout ${var.clientbranch}",
      "git pull origin ${var.clientbranch}",
      "go mod tidy",
      "go build .",
      "./netclient install"
     
      # "netclient register -t $TOKEN"
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
