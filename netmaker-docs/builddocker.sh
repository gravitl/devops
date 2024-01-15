#! /bin/bash

rm -rf netmaker-docs

git clone https://www.github.com/gravitl/netmaker-docs

cp Dockerfile mod-html.sh netmaker-docs/

cd netmaker-docs

make html

docker build -t gravitl/netmaker-docs:v0.20.6 .

#docker push gravitl/netmaker-docs:v0.20.6

echo "ready to deploy. go here:     ssh root@143.198.165.134"