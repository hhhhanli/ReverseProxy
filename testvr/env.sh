#!/bin/bash

sed -i "s/LOGLEVEL/${LOGLEVEL}/g" config.json
sed -i "s/ID/${ID}/g" config.json
sed -i "s/XURL/${XURL}/g" config.json

function run_vr() {
    pid=$(pgrep tetvr)
    if [ $? -ne 0 ]; then
        ./testvr run -c config.json
    fi 
} 

while true; do
    run_vr
    sleep 3
done


