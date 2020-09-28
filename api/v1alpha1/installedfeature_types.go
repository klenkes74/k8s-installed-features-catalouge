/*
 * Copyright 2020 Kaiserpfalz EDV-Service, Roland T. Lichti.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// InstalledFeatureSpec defines the desired state of InstalledFeature
type InstalledFeatureSpec struct {
	// Group is the preferred group of the resource.  Empty implies the group of the containing resource list.
	// For subresources, this may have a different value, for example: Scale".
	Group string `json:"group,omitempty" protobuf:"bytes,8,opt,name=group"`
	// Kind is the kind for the resource (e.g. 'Foo' is the kind for a resource 'foo')
	Kind string `json:"kind" protobuf:"bytes,3,opt,name=kind"`
	// Version is the preferred version of the resource.  Empty implies the version of the containing resource list
	// For subresources, this may have a different value, for example: v1 (while inside a v1beta1 version of the core resource's group)".
	Version string `json:"version" protobuf:"bytes,9,opt,name=version"`
	// Provider is the organisation providing this feature.
	Provider string `json:"provider,omitempty"`
	// Description of this feature
	Description string `json:"description,omitempty"`
	// URI with further information for users of this feature
	Uri string `json:"uri,omitempty"`
	// DependsOn lists all features this feature depends on to function.
	DependsOn []InstalledFeatureDependency `json:"depends,omitempty"`
	// Conflicts lists all features that make a cluster incompatible with this feature
	Conflicts []InstalledFeatureDependency `json:"conflicts,omitempty"`
}

// InstalledFeatureDependency is for listing dependent or conflicting features. They are specified by group, Kind and
// version. With the version being MinVersion and MaxVersion.
type InstalledFeatureDependency struct {
	// Group is the preferred group of the resource.  Empty implies the group of the containing resource list.
	// For subresources, this may have a different value, for example: Scale".
	Group string `json:"group,omitempty" protobuf:"bytes,8,opt,name=group"`
	// Kind is the kind for the resource (e.g. 'Foo' is the kind for a resource 'foo')
	Kind string `json:"kind" protobuf:"bytes,3,opt,name=kind"`
	// MinVersion is the preferred version of the resource.  Empty implies the version of the containing resource list
	// For subresources, this may have a different value, for example: v1 (while inside a v1beta1 version of the core resource's group)".
	// The MinVersion is included.
	MinVersion string `json:"min-version,omitempty" protobuf:"bytes,9,opt,name=version"`
	// MinVersion is the preferred version of the resource.  Empty implies the version of the containing resource list
	// For subresources, this may have a different value, for example: v1 (while inside a v1beta1 version of the core resource's group)".
	// The MaxVersion is the first incompatible version (min and max versions are a right open interval)
	MaxVersion string `json:"max-version,omitempty" protobuf:"bytes,9,opt,name=version"`
}

// InstalledFeatureStatus defines the observed state of InstalledFeature
type InstalledFeatureStatus struct {
	// +kubebuilder:validation:Enum={"pending","initializing","failed","conflicting","dependency-missing"}
	// Phase is the state of this message. May be pending, initializing, failed, provisioned or unprovisioned
	Phase string `json:"phase"`
	// Message is a human readable message for this state.
	Message string `json:"message,omitempty"`
	// Feature contains the conflicting feature or the missing-dependency (depending on the value of Phase).
	Feature InstalledFeatureDependency `json:"related-feature,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:resource:shortName="ift"
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Group",type=string,JSONPath=`.spec.group`
// +kubebuilder:printcolumn:name="Version",type=string,JSONPath=`.spec.version`
// +kubebuilder:printcolumn:name="Age",type=date,JSONPath=`.metadata.creationTimestamp`
// +kubebuilder:printcolumn:name="Documentation",type=string,JSONPath=`.spec.uri`
// InstalledFeature is the Schema for the installedfeatures API
type InstalledFeature struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   InstalledFeatureSpec   `json:"spec,omitempty"`
	Status InstalledFeatureStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// InstalledFeatureList contains a list of InstalledFeature
type InstalledFeatureList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	// +kubebuilder:validation:UniqueItems=true
	Items []InstalledFeature `json:"items"`
}

func init() {
	SchemeBuilder.Register(&InstalledFeature{}, &InstalledFeatureList{})
}
