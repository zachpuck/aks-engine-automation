#!/usr/bin/env bash

set -e

### begin login
loginCmd='az login -u "$loginId" -p "$loginSecret"'

# handle opts
if [ "$loginTenantId" != " " ]; then
    loginCmd=$(printf "%s --tenant %s" "$loginCmd" "$loginTenantId")
fi

case "$loginType" in
    "user")
        echo "logging in as user"
        ;;
    "sp")
        echo "logging in as service principal"
        loginCmd=$(printf "%s --service-principal" "$loginCmd")
        ;;
esac
eval "$loginCmd" >/dev/null

echo "setting default subscription"
az account set --subscription "$subscriptionId"
### end login

# Find the azure virtual machine scale set names for each node pool.
agentPoolsArray=()
for pool in $(echo "${nodePools}" | jq -r '.[]'); do
    agentPoolsArray+=$(az vmss list -g ${clusterName}-group  --subscription ${subscriptionId} --query "[].{name:name}[?contains(name,'${pool}')]")
done

# Delete virtual machine scale sets in azure.
for vmss in $(echo "${agentPoolsArray}" | jq -r '.[].name'); do
    echo "Deleting Azure Virtual Machine Scale Set: ${vmss}"
    az vmss delete --name ${vmss} --resource-group ${clusterName}-group --no-wait
done
