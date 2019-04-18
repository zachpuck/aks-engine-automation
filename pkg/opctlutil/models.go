package opctlutil

import (
	azurev1beta1 "github.com/zachpuck/aks-engine-automation/pkg/apis/azure/v1beta1"
	"time"
)

// ClusterConfig is used to create an aks-engine cluster in azure
type ClusterConfig struct {
	APIVersion string     `json:"apiVersion"`
	Properties Properties `json:"properties"`
}
type OrchestratorProfile struct {
	OrchestratorType    string `json:"orchestratorType"`
	OrchestratorRelease string `json:"orchestratorRelease"`
	OrchestratorVersion string `json:"orchestratorVersion"`
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

type GetOpEventsInput struct {
	OpId      string
	StartTime time.Time
}

type GetOpEventsOutput struct {
	Outcome string
}

type DeleteClusterInput struct {
	Credentials azurev1beta1.AzureCredentials
	Location    string
	ClusterName string
}

type DeleteClusterOutput struct {
	OpId string
}

type AddNodePoolGroupInput struct {
	Credentials azurev1beta1.AzureCredentials
	ClusterName string
	Config      ClusterConfig
}

type AddNodePoolGroupOutput struct {
	OpId string
}

type RemoveNodePoolGroupInput struct {
	Credentials azurev1beta1.AzureCredentials
	Location    string
	ClusterName string
	Config      ClusterConfig
}

type RemoveNodePoolGroupOutput struct {
	OpId string
}
