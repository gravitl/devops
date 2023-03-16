#! /bin/bash

#starts the serverinfo.txt file
echo "Server-info:" >> serverinfo.txt

#grabs info from docker-compose.
cat docker-compose.yml | grep SERVER_HOST >> serverinfo.txt
cat docker-compose.yml | grep MASTER_KEY >> serverinfo.txt
cat docker-compose.yml | grep SERVER_HTTP_HOST >> serverinfo.txt
cat docker-compose.yml | grep BACKEND_URL >> serverinfo.txt

# renames info from docker-compose
sed -i 's/SERVER_HOST/ip_address/g' serverinfo.txt
sed -i 's/MASTER_KEY/master_key/g' serverinfo.txt

sed -i 's/SERVER_HTTP_HOST/api_addr/g' serverinfo.txt
sed -i 's/BACKEND_URL/dashboard_addr/g' serverinfo.txt
sed -i 's-https://api-dashboard-g' serverinfo.txt
echo '      Role-tag: "server"' >> serverinfo.txt
echo "      branch-tag: $1" >>serverinfo.txt
#rm docker-compose.yml

# sets some variables
masterkey=$(cat serverinfo.txt | grep master_key | awk '{print $2;}' | tr -d '"')
echo $masterkey
ipv6addr=7b65:4206:9653:2021::/64
ipv4addr=10.22.145.0/24
apiref=$(cat serverinfo.txt | grep api_addr | awk '{print$2;}' | tr -d '"')
echo "api reference: ${apiref}"
netid='terranet'



# uses api calls to the server to setup a network and access key, then records that access key to serverinfo.txt
curl -d '{"addressrange":"10.22.145.0/24","addressrange6":"7b65:4206:9653:2021::/64","netid":"terranet"}' -H "Authorization: Bearer ${masterkey}" -H 'Content-Type: application/json' https://$apiref/api/networks
curl -d '{"uses":9999,"name":"mykey"}' -H "Authorization: Bearer $masterkey" -H 'Content-Type: application/json' https://$apiref/api/networks/$netid/keys
curl -H "Authorization: Bearer $masterkey" https://$apiref/api/networks/$netid/keys | jq | grep accessstring >> serverinfo.txt

#grabs ip addresses from all created clients
tail -n +1 ipaddress*.txt | tr -d "=<>"  >> serverinfo.txt

# ssh into each client and joins the network
accesstoken=$(cat serverinfo.txt | grep accessstring | awk '{print$2;}' | tr -d '",')
echo "theaccesstoken is: ${accesstoken}"

# install netclient locally and send the binary to the different clients.
bash updatenetclient.sh $1 $1

host1key=$(cat ipaddresshost1.txt)
scp -o StrictHostKeyChecking=no -o UserKnownHostsFile=/dev/null netclient/netclient  root@$host1key:~ 
ssh -o StrictHostKeyChecking=no -o UserKnownHostsFile=/dev/null root@$host1key "./netclient install && netclient join -t ${accesstoken}"

ingresskey=$(cat ipaddressingress.txt)
scp -o StrictHostKeyChecking=no -o UserKnownHostsFile=/dev/null netclient/netclient  root@$ingresskey:~ 
ssh -o StrictHostKeyChecking=no -o UserKnownHostsFile=/dev/null root@$ingresskey "./netclient install && netclient join -t ${accesstoken}"

egresskey=$(cat ipaddressegress.txt)
scp -o StrictHostKeyChecking=no -o UserKnownHostsFile=/dev/null netclient/netclient  root@$egresskey:~ 
ssh -o StrictHostKeyChecking=no -o UserKnownHostsFile=/dev/null root@$egresskey "./netclient install && netclient join -t ${accesstoken}"

relaykey=$(cat ipaddressrelay.txt)
scp -o StrictHostKeyChecking=no -o UserKnownHostsFile=/dev/null netclient/netclient  root@$relaykey:~ 
ssh -o StrictHostKeyChecking=no -o UserKnownHostsFile=/dev/null root@$relaykey "./netclient install && netclient join -t ${accesstoken}"

dockerkey=$(cat ipaddressdocker.txt)
ssh -o StrictHostKeyChecking=no -o UserKnownHostsFile=/dev/null root@$dockerkey "docker run -d --network host  --privileged -e TOKEN=${accesstoken} -v /etc/netclient:/etc/netclient --name netclient2 terraform/test"

