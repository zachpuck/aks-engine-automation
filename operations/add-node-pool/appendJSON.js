const fs = require('fs');

// get api model json
const apiModelRaw = fs.readFileSync('/apimodel.json', 'utf8');
let apiModel = JSON.parse(apiModelRaw);
let existingAgentPools = apiModel.properties.agentPoolProfiles;

// get config json
const configRaw = fs.readFileSync('/cluster.config');
let config = JSON.parse(configRaw);
let configAgentPools = config.properties.agentPoolProfiles;

// find which node pools to add
let newNodePools = configAgentPools.filter(function (configAgentPool) {
    let addNodePool = true
    for (let i = 0; i < existingAgentPools.length; i++) {
        if (configAgentPool.name === existingAgentPools[i].name) {
            addNodePool = false
        }
    }
    return addNodePool
});

// create copy of existing apimodel agentPoolProfile
function jsonCopy(src) {
    return JSON.parse(JSON.stringify(src));
}

// create new api model agentPoolProfile for each new node pool
newNodePools.forEach(function (pool) {
    // create copy of an existing agent pool profile
    let newProfile = jsonCopy(apiModel.properties.agentPoolProfiles[0]);

    // update with new pool properties
    newProfile.name = pool.name;
    newProfile.count = pool.count;
    newProfile.vmSize = pool.vmSize;

    // TODO: only certain vm types support acceleratedNetworking so disabling by default (needs improvement)
    newProfile.acceleratedNetworkingEnabled = false;

    // add new node pool to existing apimodel.json
    apiModel.properties.agentPoolProfiles.push(newProfile);
});

// return apimodel.json as file with new node pools
fs.writeFileSync('/newApiModel.json', JSON.stringify(apiModel));
