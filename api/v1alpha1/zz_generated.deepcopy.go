// +build !ignore_autogenerated

/*
Copyright 2020 Kaiserpfalz EDV-Service, Roland T. Lichti.

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

// Code generated by controller-gen. DO NOT EDIT.

package v1alpha1

import (
	runtime "k8s.io/apimachinery/pkg/runtime"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *InstalledFeature) DeepCopyInto(out *InstalledFeature) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	out.Status = in.Status
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new InstalledFeature.
func (in *InstalledFeature) DeepCopy() *InstalledFeature {
	if in == nil {
		return nil
	}
	out := new(InstalledFeature)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *InstalledFeature) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *InstalledFeatureDependency) DeepCopyInto(out *InstalledFeatureDependency) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new InstalledFeatureDependency.
func (in *InstalledFeatureDependency) DeepCopy() *InstalledFeatureDependency {
	if in == nil {
		return nil
	}
	out := new(InstalledFeatureDependency)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *InstalledFeatureList) DeepCopyInto(out *InstalledFeatureList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]InstalledFeature, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new InstalledFeatureList.
func (in *InstalledFeatureList) DeepCopy() *InstalledFeatureList {
	if in == nil {
		return nil
	}
	out := new(InstalledFeatureList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *InstalledFeatureList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *InstalledFeatureSpec) DeepCopyInto(out *InstalledFeatureSpec) {
	*out = *in
	if in.DependsOn != nil {
		in, out := &in.DependsOn, &out.DependsOn
		*out = make([]InstalledFeatureDependency, len(*in))
		copy(*out, *in)
	}
	if in.Conflicts != nil {
		in, out := &in.Conflicts, &out.Conflicts
		*out = make([]InstalledFeatureDependency, len(*in))
		copy(*out, *in)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new InstalledFeatureSpec.
func (in *InstalledFeatureSpec) DeepCopy() *InstalledFeatureSpec {
	if in == nil {
		return nil
	}
	out := new(InstalledFeatureSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *InstalledFeatureStatus) DeepCopyInto(out *InstalledFeatureStatus) {
	*out = *in
	out.Feature = in.Feature
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new InstalledFeatureStatus.
func (in *InstalledFeatureStatus) DeepCopy() *InstalledFeatureStatus {
	if in == nil {
		return nil
	}
	out := new(InstalledFeatureStatus)
	in.DeepCopyInto(out)
	return out
}
