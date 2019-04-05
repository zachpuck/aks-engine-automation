// +build !ignore_autogenerated

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
// Code generated by main. DO NOT EDIT.

package v1beta1

import (
	runtime "k8s.io/apimachinery/pkg/runtime"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *AgentPoolProfileConfig) DeepCopyInto(out *AgentPoolProfileConfig) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new AgentPoolProfileConfig.
func (in *AgentPoolProfileConfig) DeepCopy() *AgentPoolProfileConfig {
	if in == nil {
		return nil
	}
	out := new(AgentPoolProfileConfig)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *AksCluster) DeepCopyInto(out *AksCluster) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	out.Status = in.Status
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new AksCluster.
func (in *AksCluster) DeepCopy() *AksCluster {
	if in == nil {
		return nil
	}
	out := new(AksCluster)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *AksCluster) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *AksClusterList) DeepCopyInto(out *AksClusterList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	out.ListMeta = in.ListMeta
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]AksCluster, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new AksClusterList.
func (in *AksClusterList) DeepCopy() *AksClusterList {
	if in == nil {
		return nil
	}
	out := new(AksClusterList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *AksClusterList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *AksClusterSpec) DeepCopyInto(out *AksClusterSpec) {
	*out = *in
	out.MasterProfile = in.MasterProfile
	if in.AgentPoolProfiles != nil {
		in, out := &in.AgentPoolProfiles, &out.AgentPoolProfiles
		*out = make([]AgentPoolProfileConfig, len(*in))
		copy(*out, *in)
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new AksClusterSpec.
func (in *AksClusterSpec) DeepCopy() *AksClusterSpec {
	if in == nil {
		return nil
	}
	out := new(AksClusterSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *AksClusterStatus) DeepCopyInto(out *AksClusterStatus) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new AksClusterStatus.
func (in *AksClusterStatus) DeepCopy() *AksClusterStatus {
	if in == nil {
		return nil
	}
	out := new(AksClusterStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *AzureCredentials) DeepCopyInto(out *AzureCredentials) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new AzureCredentials.
func (in *AzureCredentials) DeepCopy() *AzureCredentials {
	if in == nil {
		return nil
	}
	out := new(AzureCredentials)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *MasterProfileConfig) DeepCopyInto(out *MasterProfileConfig) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new MasterProfileConfig.
func (in *MasterProfileConfig) DeepCopy() *MasterProfileConfig {
	if in == nil {
		return nil
	}
	out := new(MasterProfileConfig)
	in.DeepCopyInto(out)
	return out
}
