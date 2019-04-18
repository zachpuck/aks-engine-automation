const fs = require('fs');

// get api model json
const apiModelRaw = fs.readFileSync('/apimodel.json', 'utf8');
let apiModel = JSON.parse(apiModelRaw);
let existingAgentPools = apiModel.properties.agentPoolProfiles;

// get config json
const configRaw = fs.readFileSync('/cluster.config');
let config = JSON.parse(configRaw);
let configAgentPools = config.properties.agentPoolProfiles;

// find which node pools to remove
let removeNodePools = existingAgentPools.filter(existingPool => {
    let removeNodePool = true;
    for (let i = 0; i < configAgentPools.length; i++) {
        if (existingPool.name === configAgentPools[i].name) {
            removeNodePool = false
        }
    }
    return removeNodePool
});

// find which node pools to remove
let filteredAgentPools = existingAgentPools.filter(existingPool => {
    let removeNodePool = true;
    removeNodePools.forEach(removePool => {
        if (existingPool.name === removePool.name) {
            removeNodePool = false;
        }
    });
    return removeNodePool
});

// update apimodel to remove agent pools
apiModel.properties.agentPoolProfiles = filteredAgentPools;

// return apimodel.json as file with removed node pool(s)
fs.writeFileSync('/updatedApiModel.json', JSON.stringify(apiModel));

// return a array of node pool names as a JSON file
let nodePoolsToDrain = removeNodePools.map(pool => pool.name);
fs.writeFileSync('/nodePoolsToDrain.json', JSON.stringify(nodePoolsToDrain));
