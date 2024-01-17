#! /bin/bash

# this file is needed on the docs server (143.198.165.134) for the workflow to run successfully. 
# It calls this script to change the image to inputs.version that is entered on the workflow.

image=$(cat docker-compose.yml | grep image | awk '{print$2; exit}')

sed -i "s+$image+gravitl/netmaker-docs:$1+g" docker-compose.yml

docker pull gravitl/netmaker-docs:$1

docker-compose up -d