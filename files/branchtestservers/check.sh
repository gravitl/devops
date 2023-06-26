#!/bin/bash
#
if [ -e /tmp/branchtest ]
then
    echo fail
    exit 0
else
    touch /tmp/branchtest
    docker-compose down
    docker system prune -a -f --volumes
    rm ~/* 
    echo pass
    exit 0
fi
