#!/bin/bash

COUNTRY="$1"
BASEURL="$2"
EARLYDATE=$(date --date="-3 days" --iso-8601 2>&- || gdate --date="-3 days" --iso-8601)

mkdir -p "data/${COUNTRY}/hourly" || true

for DAY in $(curl -s "${BASEURL}diagnosis-keys/country/${COUNTRY}/date" | jq -r '.[]');
do 
    echo "Checking if ${DAY} is after ${EARLYDATE}:	";
    if [ "${DAY}" \> "${EARLYDATE}" ];
    then
        echo "Downloading ${DAY}:	";
        wget "${BASEURL}diagnosis-keys/country/${COUNTRY}/date/${DAY}" -O "data/${COUNTRY}/${DAY}.zip";

        for HOUR in $(curl -s "${BASEURL}diagnosis-keys/country/${COUNTRY}/date/${DAY}/hour" | jq -r '.[]');
        do
            echo "Downloading hour ${HOUR}:	";
            wget "${BASEURL}diagnosis-keys/country/${COUNTRY}/date/${DAY}/hour/${HOUR}" -O "data/${COUNTRY}/hourly/${DAY}-${HOUR}.zip";
        done
    fi;
done

wget "${BASEURL}configuration/country/${COUNTRY}/app_config" -O "data/${COUNTRY}/app_config.zip"