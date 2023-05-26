#!/bin/bash

FLAG="FALSE"
if [ ! -e /root/minecraft/version.txt ]; then
mkdir /root/minecraft 
FLAG="TRUE"
fi;
if [ -e /root/minecraft/version.txt ]; then
    if [ $(cat /root/minecraft/version.txt) != $(/root/run) ]; then
FLAG="TRUE"
    fi;
fi;
if [ $FLAG = "TRUE" ]; then 
    chmod +x /root/backup.sh /root/run
    /root/backup.sh &
    curl -sLo /root/bedrock_server.zip $(/root/run)
    unzip /root/bedrock_server.zip -d /root/unzip
    if [ ! -e /root/minecraft/server.properties ]; then
        echo -e "\nemit-server-telemetry=true\n" >> /root/unzip/server.properties
        cp /root/unzip/server.properties /root/minecraft
    fi;
    rm /root/unzip/server.properties
    if [ ! -e /root/minecraft/allowlist.json ]; then
        cp /root/unzip/allowlist.json /root/minecraft
    fi;
    rm /root/unzip/allowlist.json
    if [ ! -e /root/minecraft/permissions.json ]; then
        cp /root/unzip/permissions.json /root/minecraft
    fi;
    rm /root/unzip/permissions.json
    echo $(/root/run) > /root/unzip/version.txt
    rm -rf /root/minecraft/behavior_packs /root/minecraft/config /root/minecraft/definitions /root/minecraft/resource_packs /root/minecraft/structures
    mv -f /root/unzip/* /root/minecraft/
    rm -rf /root/unzip
    chmod +x /root/minecraft/bedrock_server
fi;
cd /root/minecraft
./bedrock_server
