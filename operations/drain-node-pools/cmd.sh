#!/usr/bin/env bash

export KUBECONFIG=/kubeConfig.json

for pool in $(echo "${nodePools}" | jq -r '.[]'); do
    echo "Draining nodes in node pool: ${pool}"
    kubectl drain \
        $(kubectl get nodes --selector=agentpool=${pool} --no-headers -o custom-columns=:metadata.name) \
        --ignore-daemonsets
done
