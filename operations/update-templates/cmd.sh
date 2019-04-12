#!/usr/bin/env bash

echo "downloading aks-engine binary"
curl -o get-akse.sh https://raw.githubusercontent.com/Azure/aks-engine/master/scripts/get-akse.sh
chmod 700 get-akse.sh
./get-akse.sh

# aks-engine looks for the _output folder when generating templates
mkdir -p _output/${clusterName}
cp -r /templates/ /_output/${clusterName}
cp /modifiedApiModel.json /_output/${clusterName}/apimodel.json

echo "update ARM resource templates"
aks-engine generate --api-model _output/${clusterName}/apimodel.json

echo "uploading ARM resource templates to storage account"
az storage blob upload-batch \
    --destination ${clusterName} \
    --source _output/${clusterName} \
    --account-name ${storageAccountName} \
    --account-key ${storageAccountKey} > /dev/null
