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
rm docker-compose.yml

# sets some variables
masterkey=$(cat serverinfo.txt | grep master_key | awk '{print $2;}' | tr -d '"')
#echo $masterkey
ipv6addr=7b65:4206:9653:2021::/64
ipv4addr=10.22.145.0/24
apiref=$(cat serverinfo.txt | grep api_addr | awk '{print$2;}' | tr -d '"')
#echo "api reference: ${apiref}"
netid='terranet'



# uses api calls to the server to setup a network and access key for docker and registration key for clients, then records those keys to serverinfo.txt
curl -d '{"addressrange":"10.22.145.0/24","addressrange6":"7b65:4206:9653:2021::/64","netid":"terranet"}' -H "Authorization: Bearer ${masterkey}" -H 'Content-Type: application/json' https://$apiref/api/networks
curl -X POST -H "Authorization: Bearer $masterkey" -d '{"expiration":0,"uses_remaining":10,"networks":["terranet"],"unlimited":false,"tags":[]}' https://$apiref/api/v1/enrollment-keys
curl -H "Authorization: Bearer $masterkey" https://$apiref/api/v1/enrollment-keys | jq | grep token >> serverinfo.txt

#grabs ip addresses from all created clients
tail -n +1 ipaddress*.txt | tr -d "=<>"  >> serverinfo.txt

# ssh into each client and registers with the server
regtoken=$(cat serverinfo.txt | grep token | awk '{print$2;}' | tr -d '",')
#echo "regtoken: ${regtoken}"
bash updatenetclient.sh $1 v0.18.5

cd /home/runner/

host1key=$(cat /root/ipaddresshost1.txt)
scp -o StrictHostKeyChecking=no -o UserKnownHostsFile=/dev/null /root/netclient/netclient  root@$host1key:~ 
ssh -o StrictHostKeyChecking=no -o UserKnownHostsFile=/dev/null root@$host1key "./netclient install && netclient register -t ${regtoken}"

ingresskey=$(cat /root/ipaddressingress.txt)
scp -o StrictHostKeyChecking=no -o UserKnownHostsFile=/dev/null /root/netclient/netclient  root@$ingresskey:~ 
ssh -o StrictHostKeyChecking=no -o UserKnownHostsFile=/dev/null root@$ingresskey "./netclient install && netclient register -t ${regtoken}"

egresskey=$(cat /root/ipaddressegress.txt)
scp -o StrictHostKeyChecking=no -o UserKnownHostsFile=/dev/null /root/netclient/netclient  root@$egresskey:~ 
ssh -o StrictHostKeyChecking=no -o UserKnownHostsFile=/dev/null root@$egresskey "./netclient install && netclient register -t ${regtoken}"

relaykey=$(cat /root/ipaddressrelay.txt)
scp -o StrictHostKeyChecking=no -o UserKnownHostsFile=/dev/null /root/netclient/netclient  root@$relaykey:~ 
ssh -o StrictHostKeyChecking=no -o UserKnownHostsFile=/dev/null root@$relaykey "./netclient install && netclient register -t ${regtoken}"

dockerkey=$(cat /root/ipaddressdocker.txt)
ssh -o StrictHostKeyChecking=no -o UserKnownHostsFile=/dev/null root@$dockerkey "docker run -d --network host  --privileged -e TOKEN=${regtoken} -v /etc/netclient:/etc/netclient --name netclient2 terraform/test"
