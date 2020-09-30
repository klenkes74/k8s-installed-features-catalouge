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

// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.
// Important: Run "make" to regenerate code after modifying this file

// InstalledFeatureGroupSpec defines the desired state of InstalledFeatureGroup
type InstalledFeatureGroupSpec struct {
	// Provider is the organisation providing this feature
	Provider string `json:"provider,omitempty"`
	// Description of this feature
	Description string `json:"description,omitempty"`
	// URI with further information for users of this feature
	Uri string `json:"uri,omitempty"`
}

// InstaledFeatureGroupListedFeature defines subfeatures by namespace and name
type InstalledFeatureGroupListedFeature struct {
	// Namespace is the namespace of the feature listed
	Namespace string `json:"namespace,omitempty"`
	// Name is the name of the feature listed
	Name string `json:"name"`
}

// InstalledFeatureGroupStatus defines the observed state of InstalledFeatureGroup
type InstalledFeatureGroupStatus struct {
	// +kubebuilder:validation:Enum={"pending","initializing","failed","provisioned"}
	// Phase is the state of this message. May be pending, initializing, failed, provisioned
	Phase string `json:"phase"`
	// Message is a human readable message for this state
	Message string `json:"message,omitempty"`
	// Features contain all features of this feature group
	Features []InstalledFeatureGroupListedFeature `json:"features,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:shortName="iftg"
// +kubebuilder:printcolumn:name="Group",type=string,JSONPath=`.metadata.name`
// +kubebuilder:printcolumn:name="Age",type=date,JSONPath=`.metadata.creationTimestamp`
// +kubebuilder:printcolumn:name="Documentation",type=string,JSONPath=`.spec.uri`

// InstalledFeatureGroup is the Schema for the installedfeaturegroups API
type InstalledFeatureGroup struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   InstalledFeatureGroupSpec   `json:"spec,omitempty"`
	Status InstalledFeatureGroupStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// InstalledFeatureGroupList contains a list of InstalledFeatureGroup
type InstalledFeatureGroupList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []InstalledFeatureGroup `json:"items"`
}

func init() {
	SchemeBuilder.Register(&InstalledFeatureGroup{}, &InstalledFeatureGroupList{})
}
