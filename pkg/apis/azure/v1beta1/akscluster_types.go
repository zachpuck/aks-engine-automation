/*
Copyright 2019 Zach Puckett.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1beta1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

const AksClusterFinalizer = "akscluster.cnct.io"

// AksClusterSpec defines the desired state of AksCluster
type AksClusterSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Credentials references a secret used to access the azure environment (must have access to blob storageaccount)
	Credentials string `json:"credentials,omitempty"`

	// Region where you are creating the cluster
	Location string `json:"location,omitempty"`

	// Kubernetes version
	KubernetesVersion string `json:"kubernetesVersion"`

	// Master Profile configuration
	MasterProfile MasterProfileConfig `json:"masterProfile"`

	// Agent Pool Profiles
	AgentPoolProfiles []AgentPoolProfileConfig `json:"agentPoolProfiles,omitempty"`
}

// AzureCredentials is used to authenticate to azure
type AzureCredentials struct {
	TenantId       string `json:"tenantId,omitempty"`
	SubscriptionId string `json:"subscriptionId,omitempty"`
	LoginId        string `json:"loginId,omitempty"`
	LoginSecret    string `json:"loginSecret,omitempty"`
}

// MasterProfileConfig defines the resource configuration for the control plane nodes
type MasterProfileConfig struct {
	Count     int    `json:"count,omitempty"`
	DnsPrefix string `json:"dnsPrefix,omitempty"`
	VmSize    string `json:"vmSize,omitempty"`
}

// AgentPoolProfileConfig defines the resource configuration for an agent nodes pool
type AgentPoolProfileConfig struct {
	Name                   string `json:"name,omitempty"`
	Count                  int    `json:"count,omitempty"`
	VmSize                 string `json:"vmSize,omitempty"`
	EnableVMSSNodePublicIP bool   `json:"enableVMSSNodePublicIP,omitempty"`
}

type ClusterStatusPhase string
type ClusterKubernetesVersion string
type ClusterNodePoolCount int

// AksClusterStatus defines the observed state of AksCluster
type AksClusterStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Current number of node pools
	NodePoolCount int `json:"nodePoolCount,omitempty"`
	// Cluster kubernetes version
	KubernetesVersion ClusterKubernetesVersion `json:"kubernetesVersion,omitempty"`
	// Cluster status
	Phase ClusterStatusPhase `json:"phase,omitempty"`
}

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// AksCluster is the Schema for the aksclusters API
// +k8s:openapi-gen=true
type AksCluster struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   AksClusterSpec   `json:"spec,omitempty"`
	Status AksClusterStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// AksClusterList contains a list of AksCluster
type AksClusterList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []AksCluster `json:"items"`
}

func init() {
	SchemeBuilder.Register(&AksCluster{}, &AksClusterList{})
}
