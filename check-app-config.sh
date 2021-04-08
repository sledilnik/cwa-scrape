#!/bin/bash
set -e

BASEURL="$1"

wget "${BASEURL}app_config_android" -O "data/app_config_android.zip"
wget "${BASEURL}app_config_ios"  -O "data/app_config_ios.zip"

cd proto
unzip -p "../data/app_config_android.zip" export.bin | protoc --decode=ApplicationConfigurationAndroid app_config_android.proto > "../data/app_config_android.txt"
unzip -p "../data/app_config_ios.zip"     export.bin | protoc --decode=ApplicationConfigurationIOS     app_config_ios.proto     > "../data/app_config_ios.txt"
cd ..

