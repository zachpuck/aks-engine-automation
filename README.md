# AKS engine automation

A Kubernetes operator used to provision Kubernetes clusters in azure utilizing [aks-engine](https://github.com/Azure/aks-engine)

## Overview
AKS engine automation provides a Kubernetes native process to creating aks-engine clusters in azure. 
Utilizing custom resource definitions of Kind `AksCluster` to create and update clusters. 
Additional secrets are created for azure credentials(`<cluster-name>-secret`) and vm ssh private key(`<cluster-name>-privatekey`).

## Getting started

**Prerequisites**:
Aks engine automation requires an azure storage account used to store the aks-engine manifests. 
This storage account must be in the same subscription as the Kubernetes clusters.
1. Create an azure resource group: 
`az group create --name "aks-operator-group" --location "westus"`
1. Create an azure storage account: 
`az storage account create --name "aksoperatorstorage01" --resource-group "aks-operator-group"`

Installing the Operator:
```yaml
helm install deployments/helm/aks-engine-automation 
  --name AksOperator 
  --set storageAccount.name="aksoperatorstorage01" 
  --set storageAccount.group="aks-operator-group"
```

#### Build
To build new images:
`opctl run build`

#### Local Debug (using Opctl):

Requirements:
1. docker
2. [opctl](https://opctl.io/docs/getting-started/opctl.html#installation)

 Steps:
1. set dockerSocket in your environment, typically: `export dockerSocket=/var/run/docker.sock`
1. start debug op: `opctl run debug`

This will start a a container running [kind](https://github.com/kubernetes-sigs/kind) and one running `aks-engine-automation` that will use the kubeconfig for the kind cluster
The kind-config-debug.yaml will be located in the root of this repo and can be used to access the kind cluster locally by setting `export KUBECONFIG=<repo-path>/kind-config-debug.yaml`

**Cleanup after debugging**:

Once you have completed locally debugging run `opctl run cleanup` to stop the kind cluster and remove its kubeconfig from the repo.


**Features needed**:
1. scaling
1. add additional node pool
1. upgrade kubernetes version
1. custom vnets
1. availability zones
1. Virtual machine scale set for masters
