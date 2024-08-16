#! /bin/bash

#starts the serverinfo.txt file
echo "Server-info:" >> serverinfo.txt

#grabs info from docker-compose.
cat netmaker.env | grep SERVER_HOST= >> serverinfo.txt
cat netmaker.env | grep MASTER_KEY= >> serverinfo.txt
cat netmaker.env | grep NM_DOMAIN= >> serverinfo.txt


# renames info from docker-compose
sed -i 's/SERVER_HOST=/ip_address /g' serverinfo.txt
sed -i 's/MASTER_KEY=/master_key /g' serverinfo.txt
sed -i 's/NM_DOMAIN=/api_addr /g' serverinfo.txt
echo '      Role-tag: "server"' >> serverinfo.txt
echo "      branch-tag: $1" >>serverinfo.txt
rm netmaker.env

cat serverinfo.txt

# sets some variables
masterkey=$(cat serverinfo.txt | grep master_key | awk '{print $2;}' | tr -d '"')
apiref="$(cat serverinfo.txt | grep api_addr | awk '{print$2;}' | tr -d '"')"
echo "API REFERENCE IS: api.$apiref"




# uses an api call to get the netmaker enrollment key that was made during the nm-quick script and records that key to serverinfo.txt
curl -H "Authorization: Bearer $masterkey" https://api.$apiref/api/v1/enrollment-keys | jq | grep token >> serverinfo.txt

#grabs ip addresses from all created clients
tail -n +1 ipaddress*.txt | tr -d "=<>"  >> serverinfo.txt

# ssh into each client and registers with the server
regtoken=$(cat serverinfo.txt | grep token | awk '{print$2;}' | tr -d '",')
regtokennonewline=$(cat serverinfo.txt | grep token | awk '{print$2;}' | tr -d '",' | tr -d '\n')
#get ip addresses
relayedkey=$(cat ipaddressrelayed.txt)
ingresskey=$(cat ipaddressingress.txt)
egresskey=$(cat ipaddressegress.txt)
relaykey=$(cat ipaddressrelay.txt)
dockerkey=$(cat ipaddressdocker.txt)
serverkey=$(cat serverinfo.txt | grep ip_address | awk '{print$2;}' | tr -d '",')


#register with server.
echo "logging into each client and registering with the network."
echo "logging into server"
ssh -i /home/runner/.ssh/deploy.key -o StrictHostKeyChecking=no -o UserKnownHostsFile=/dev/null root@$serverkey "netclient register -t ${regtoken}"
echo "logging into relayed"
echo "logging into docker"
ssh -i /home/runner/.ssh/deploy.key -o StrictHostKeyChecking=no -o UserKnownHostsFile=/dev/null root@$dockerkey "docker run -d --network host --privileged -e TOKEN=${regtokennonewline} -v /etc/netclient:/etc/netclient --name netclient2 terraform/test "
ssh -i /home/runner/.ssh/deploy.key -o StrictHostKeyChecking=no -o UserKnownHostsFile=/dev/null root@$relayedkey "netclient register -t ${regtoken}"
echo "done with relayed. logging into ingress"
ssh -i /home/runner/.ssh/deploy.key -o StrictHostKeyChecking=no -o UserKnownHostsFile=/dev/null root@$ingresskey "netclient register -t ${regtoken}"
echo "done with ingress. logging into egress"
ssh -i /home/runner/.ssh/deploy.key -o StrictHostKeyChecking=no -o UserKnownHostsFile=/dev/null root@$egresskey "netclient register -t ${regtoken}"
echo "done with egress. logging into relay"
ssh -i /home/runner/.ssh/deploy.key -o StrictHostKeyChecking=no -o UserKnownHostsFile=/dev/null root@$relaykey "netclient register -t ${regtoken}"
echo "done with relay. completed"
