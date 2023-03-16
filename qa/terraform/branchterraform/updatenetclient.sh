#!/bin/bash

# this is a defferent updatenetcleint.sh file for terraform to work with. needs to have go and gcc installed first.

if (($# <2))
then
        echo "you should have two arguments. the netclient branch to pull and the testing version"
        exit
elif (($# > 2))
then
        echo "you should have two arguments. the netclient branch to pull and the testing version"
        exit
else
        echo "you are pulling from branch: $1 and making version: $2"
fi

rm -rf netclient
git clone https://www.github.com/gravitl/netclient
wait
cd netclient
git checkout $1
git pull origin $1
go mod tidy
go build -tags headless -ldflags="-X main.version=$2"
wait


echo update finished