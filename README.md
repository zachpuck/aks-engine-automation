# AKS engine automation

### Generating the templates.

1. Custom vnets
1. Availability zones
1. Virtual Machine Scale Sets

### Creating the resource group and service principal.

### Deploying to Azure.

### Adding a node (scaling).

### Adding a second node pool.

### Upgrading kubernetes version.



# store in storage account

  keyVaultName:
  KeyVaultRegion:
  StorageAccountName:
  StorageAccountRegion:



you provide a config manifest and in it, the path to the secrets.

i need to, generate templates if they do not exist. and apply?, update? upgrade 

#### Local Debug:

 **Using Opctl**:

 requirements:
1. docker
2. [opctl](https://opctl.io/docs/getting-started/opctl.html#installation)

 steps:
1. set dockerSocket in your environment, typically: `export dockerSocket=/var/run/docker.sock`
1. start debug op: `opctl run debug`

  This will start a a container running [kind](https://github.com/kubernetes-sigs/kind) and one running `aks-engine-automation` that will use the kubeconfig for the kind cluster
 The kind-config-debug.yaml will be located in the root of this repo and can be used to access the kind cluster locally by setting `export KUBECONFIG=<repo-path>/kind-config-debug.yaml`

 **Cleanup after debugging**:

 Once you have completed locally debugging run `opctl run cleanup` to stop the kind cluster and remove its kubeconfig from the repo.
