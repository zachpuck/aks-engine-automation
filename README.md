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

#### Usage
Once the Operator is installed and running in a kubernetes cluster you can begin creating AksCluster Custom Resources.
An sample resource is located in [config/samples](config/samples).
The sample contains two resources.
1. A kubernetes secret that contains your azure credentials.
1. An `AksCluster` custom resource used to define the details of your cluster.

The Operator has two containers. One container is the kubernetes operator itself and will show the logs 
related managing the `Akscluster` Custom resources. The second container is the `opctl`.
The logs from `opctl` container show the indivdiual results of each ["operation"](operations) 
(the individual steps of managing clusters: create, update, add node, ect..). 
Each of these operations returns the results to standard out. 

The `operations` are created using [opctl](https://opctl.io/docs/)

#### Build
To build new images:
`opctl run build`

This will build

#### Local Debug:

Requirements:
1. [minikube](https://kubernetes.io/docs/tasks/tools/install-minikube/)

 Steps:
1. `minikube start`
1. `make install`
1. `make run`

In as separate terminal: 
1. update sample CR with your azure credentials: ./config/samples/azure_v1beta1_akscluster.yaml
1. `kubectl apply -f ./config/samples/azure_v1beta1_akscluster.yaml`

You can now see the created resource by typing `kubectl get akcluster`

**Features needed**:
1. scaling
1. upgrade kubernetes version
1. custom vnets
1. availability zones
1. Virtual machine scale set for masters
