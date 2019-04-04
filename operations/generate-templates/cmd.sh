#!/usr/bin/env bash

echo "check if cluster container exists"
exists=$(az storage container exists -n ${clusterName} \
            --account-name ${storageAccountName} \
            --account-key ${storageAccountKey})

if [ "$(echo $exists | jq '.exists')" == "true" ]; then
    echo "cluster exists in storage"
else
    echo "downloading aks-engine binary"
    curl -o get-akse.sh https://raw.githubusercontent.com/Azure/aks-engine/master/scripts/get-akse.sh
    chmod 700 get-akse.sh
    ./get-akse.sh

    echo "generating ARM resource templates"
    aks-engine generate /cluster.config

    echo "creating storage account container for cluster"
    az storage container create -n ${clusterName} \
            --account-name ${storageAccountName} \
            --account-key ${storageAccountKey}

    echo "uploading ARM resource templates to storage account"
    az storage blob upload-batch \
        --destination ${clusterName} \
        --source /_output/${clusterName} \
        --account-name ${storageAccountName} \
        --account-key ${storageAccountKey} > /dev/null
fi
