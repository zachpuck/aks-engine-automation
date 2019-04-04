package opctlutil

import (
	azurev1beta1 "github.com/zachpuck/aks-engine-automation/pkg/apis/azure/v1beta1"
)

// ClusterConfig is used to create an aks-engine cluster in azure
type ClusterConfig struct {
	APIVersion string     `json:"apiVersion"`
	Properties Properties `json:"properties"`
}
type OrchestratorProfile struct {
	OrchestratorType    string `json:"orchestratorType"`
	OrchestratorRelease string `json:"orchestratorRelease"`
}
type MasterProfile struct {
	Count     int    `json:"count"`
	DNSPrefix string `json:"dnsPrefix"`
	VMSize    string `json:"vmSize"`
}
type AgentPoolProfiles struct {
	Name   string `json:"name"`
	Count  int    `json:"count"`
	VMSize string `json:"vmSize"`
}
type PublicKeys struct {
	KeyData string `json:"keyData"`
}
type SSH struct {
	PublicKeys []PublicKeys `json:"publicKeys"`
}
type LinuxProfile struct {
	AdminUsername string `json:"adminUsername"`
	SSH           SSH    `json:"ssh"`
}
type ServicePrincipalProfile struct {
	ClientID string `json:"clientId"`
	Secret   string `json:"secret"`
}
type Properties struct {
	OrchestratorProfile     OrchestratorProfile     `json:"orchestratorProfile"`
	MasterProfile           MasterProfile           `json:"masterProfile"`
	AgentPoolProfiles       []AgentPoolProfiles     `json:"agentPoolProfiles"`
	LinuxProfile            LinuxProfile            `json:"linuxProfile"`
	ServicePrincipalProfile ServicePrincipalProfile `json:"servicePrincipalProfile"`
}

type CreateClusterInput struct {
	Credentials azurev1beta1.AzureCredentials
	Location    string
	ClusterName string
	Config      ClusterConfig
}

type CreateClusterOutput struct {
	OpId string
}
