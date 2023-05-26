#!/bin/bash

cd /root/minecraft;
while true;
    do sleep 21600;
    if [ ! -d './backup' ]; then
        mkdir ./backup;
    fi;
    cp -R ./worlds ./backup/$(date \"+"+"%s"+"\").backup;
    cd ./backup;
        find ./ -mtime +2 -name \"*.backup\" -type d | xargs rm -rf;
    cd ..;
done
