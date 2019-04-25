#!/bin/sh

#echo "downloading aks-engine binary"
#curl -o get-akse.sh https://raw.githubusercontent.com/Azure/aks-engine/master/scripts/get-akse.sh
#chmod 700 get-akse.sh
#./get-akse.sh

# aks-engine looks for the _output folder when generating templates
mkdir -p _output/${clusterName}
cp -r /templates/. /_output/${clusterName}

echo "upgrade cluster"
aks-engine upgrade \
  --subscription-id ${subscriptionId} \
  --api-model /_output/${clusterName}/apimodel.json \
  --location ${location} \
  --resource-group ${resourceGroup} \
  --upgrade-version ${upgradeVersion} \
  --auth-method client_secret \
  --client-id ${loginId} \
  --client-secret ${loginSecret}

echo "uploading ARM resource templates to storage account"
az storage blob upload-batch \
    --destination ${clusterName} \
    --source _output/${clusterName} \
    --account-name ${storageAccountName} \
    --account-key ${storageAccountKey} > /dev/null
