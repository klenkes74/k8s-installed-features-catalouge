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

package installedfeature_test

import (
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	. "github.com/klenkes74/k8s-installed-features-catalogue/api/v1alpha1"
	"github.com/klenkes74/k8s-installed-features-catalogue/controllers/installedfeature"
	"github.com/klenkes74/k8s-installed-features-catalogue/controllers/installedfeaturegroup"
	"github.com/klenkes74/k8s-installed-features-catalogue/generated"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/controller-runtime/pkg/envtest/printer"
	"testing"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
	// +kubebuilder:scaffold:imports
)

func TestInstalledFeatureController(t *testing.T) {
	RegisterFailHandler(Fail)

	RunSpecsWithDefaultAndCustomReporters(t,
		"InstalledFeature Controller Suite",
		[]Reporter{printer.NewlineReporter{}})
}

// These tests use Ginkgo (BDD-style Go testing framework). Refer to
// http://onsi.github.io/ginkgo/ to learn more about Ginkgo.

var (
	ctrlMock *gomock.Controller

	client *generated.MockOcpClient
	sut    installedfeature.Reconciler
)

var _ = BeforeEach(func() {
	logf.SetLogger(zap.New(zap.UseDevMode(true), zap.WriteTo(GinkgoWriter)))

	scheme := runtime.NewScheme()
	Expect(clientgoscheme.AddToScheme(scheme)).Should(Succeed())
	Expect(AddToScheme(scheme)).Should(Succeed())

	ctrlMock = gomock.NewController(GinkgoT())
	client = generated.NewMockOcpClient(ctrlMock)

	sut = installedfeature.Reconciler{
		Client: client,
		Log:    logf.Log,
		Scheme: scheme,
	}
})

var _ = AfterEach(func() {
	By("tearing down the mock controller")
	ctrlMock.Finish()
})

func createIFT(name string, namespace string, version string, provider string, description string, uri string, finalizer bool, deleted bool) *InstalledFeature {
	result := &InstalledFeature{
		TypeMeta: metav1.TypeMeta{
			Kind:       "InstalledFeature",
			APIVersion: GroupVersion.Group + "/" + GroupVersion.Version,
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:              name,
			Namespace:         namespace,
			CreationTimestamp: metav1.Time{Time: time.Now().Add(24 * time.Hour)},
			ResourceVersion:   "1",
			Generation:        0,
			UID:               types.UID(uuid.New().String()),
		},
		Spec: InstalledFeatureSpec{
			Kind:        name,
			Version:     version,
			Provider:    provider,
			Description: description,
			Uri:         uri,
		},
	}

	if finalizer {
		result.Finalizers = make([]string, 1)
		result.Finalizers[0] = installedfeature.FinalizerName
	}

	if deleted {
		deletionGracePeriod := int64(60)
		result.DeletionGracePeriodSeconds = &deletionGracePeriod
		result.DeletionTimestamp = &metav1.Time{Time: time.Now().Add(2 * time.Minute)}
	}

	return result
}

func copyIFT(orig *InstalledFeature) *InstalledFeature {
	//goland:noinspection GoDeprecation
	result := &InstalledFeature{
		TypeMeta: metav1.TypeMeta{
			Kind:       orig.TypeMeta.Kind,
			APIVersion: orig.TypeMeta.APIVersion,
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:                       orig.ObjectMeta.Name,
			GenerateName:               orig.ObjectMeta.GenerateName,
			Namespace:                  orig.ObjectMeta.Namespace,
			SelfLink:                   orig.ObjectMeta.SelfLink,
			UID:                        orig.ObjectMeta.UID,
			ResourceVersion:            orig.ObjectMeta.ResourceVersion,
			Generation:                 orig.ObjectMeta.Generation,
			CreationTimestamp:          orig.ObjectMeta.CreationTimestamp,
			DeletionTimestamp:          orig.ObjectMeta.DeletionTimestamp,
			DeletionGracePeriodSeconds: orig.ObjectMeta.DeletionGracePeriodSeconds,
			ClusterName:                orig.ObjectMeta.ClusterName,
		},
		Spec: InstalledFeatureSpec{
			Kind:        orig.Spec.Kind,
			Version:     orig.Spec.Version,
			Provider:    orig.Spec.Provider,
			Description: orig.Spec.Description,
			Uri:         orig.Spec.Uri,
		},
	}

	if orig.Spec.Group != nil {
		result = setGroupToIFT(result, orig.Spec.Group.Name, orig.Spec.Group.Namespace)
	}

	if len(orig.ObjectMeta.Labels) > 0 {
		result.ObjectMeta.Labels = make(map[string]string)
		for key, value := range orig.ObjectMeta.Labels {
			result.ObjectMeta.Labels[key] = value
		}
	}

	if len(orig.ObjectMeta.Annotations) > 0 {
		result.ObjectMeta.Annotations = make(map[string]string)
		for key, value := range orig.ObjectMeta.Annotations {
			result.ObjectMeta.Annotations[key] = value
		}
	}

	if len(orig.ObjectMeta.Finalizers) > 0 {
		result.ObjectMeta.Finalizers = make([]string, len(orig.ObjectMeta.Finalizers))
		for i, value := range orig.ObjectMeta.Finalizers {
			result.ObjectMeta.Finalizers[i] = value
		}
	}

	if len(orig.ObjectMeta.OwnerReferences) > 0 {
		result.ObjectMeta.OwnerReferences = make([]metav1.OwnerReference, len(orig.ObjectMeta.OwnerReferences))
		for i, r := range orig.ObjectMeta.OwnerReferences {
			result.ObjectMeta.OwnerReferences[i] = metav1.OwnerReference{
				APIVersion:         r.APIVersion,
				Kind:               r.Kind,
				Name:               r.Name,
				UID:                r.UID,
				Controller:         r.Controller,
				BlockOwnerDeletion: r.BlockOwnerDeletion,
			}
		}
	}

	if len(orig.ObjectMeta.ManagedFields) > 0 {
		result.ObjectMeta.ManagedFields = make([]metav1.ManagedFieldsEntry, len(orig.ObjectMeta.ManagedFields))
		for i, r := range orig.ObjectMeta.ManagedFields {
			result.ObjectMeta.ManagedFields[i] = metav1.ManagedFieldsEntry{
				Manager:    r.Manager,
				Operation:  r.Operation,
				APIVersion: r.APIVersion,
				Time:       r.Time,
				FieldsType: r.FieldsType,
				FieldsV1: &metav1.FieldsV1{
					Raw: r.FieldsV1.Raw,
				},
			}
		}
	}

	if len(orig.Spec.DependsOn) > 0 {
		result.Spec.DependsOn = make([]InstalledFeatureRef, len(orig.Spec.DependsOn))
		for i, r := range orig.Spec.DependsOn {
			result.Spec.DependsOn[i] = createFeatureRef(r)
		}
	}

	if len(orig.Status.DependingFeatures) > 0 {
		result.Status.DependingFeatures = make([]InstalledFeatureRef, len(orig.Status.DependingFeatures))
		for i, r := range orig.Status.DependingFeatures {
			result.Status.DependingFeatures[i] = createFeatureRef(r)
		}
	}

	if len(orig.Status.MissingDependencies) > 0 {
		result.Status.MissingDependencies = make([]InstalledFeatureRef, len(orig.Status.MissingDependencies))
		for i, r := range orig.Status.MissingDependencies {
			result.Status.MissingDependencies[i] = createFeatureRef(r)
		}
	}

	if len(orig.Status.ConflictingFeatures) > 0 {
		result.Status.ConflictingFeatures = make([]InstalledFeatureRef, len(orig.Status.ConflictingFeatures))
		for i, r := range orig.Status.ConflictingFeatures {
			result.Status.ConflictingFeatures[i] = createFeatureRef(r)
		}
	}

	return result
}

func createFeatureRef(orig InstalledFeatureRef) InstalledFeatureRef {
	return InstalledFeatureRef{
		Namespace: orig.Namespace,
		Name:      orig.Name,
	}
}

func setGroupToIFT(instance *InstalledFeature, name string, namespace string) *InstalledFeature {
	instance.Spec.Group = &InstalledFeatureRef{
		Namespace: namespace,
		Name:      name,
	}

	return instance
}

func createIFTG(name string, namespace string, provider string, description string, uri string, finalizer bool, deleted bool) *InstalledFeatureGroup {
	result := &InstalledFeatureGroup{
		TypeMeta: metav1.TypeMeta{
			Kind:       "InstalledFeature",
			APIVersion: GroupVersion.Group + "/" + GroupVersion.Version,
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:              name,
			Namespace:         namespace,
			CreationTimestamp: metav1.Time{Time: time.Now().Add(24 * time.Hour)},
			ResourceVersion:   "1",
			Generation:        0,
			UID:               types.UID(uuid.New().String()),
		},
		Spec: InstalledFeatureGroupSpec{
			Provider:    provider,
			Description: description,
			Uri:         uri,
		},
	}

	if finalizer {
		result.Finalizers = make([]string, 1)
		result.Finalizers[0] = installedfeaturegroup.FinalizerName
	}

	if deleted {
		deletionGracePeriod := int64(60)
		result.DeletionGracePeriodSeconds = &deletionGracePeriod
		result.DeletionTimestamp = &metav1.Time{Time: time.Now().Add(2 * time.Minute)}
	}

	return result
}

func createNotFound(resourceType string, name string) errors.APIStatus {
	return errors.NewNotFound(
		schema.GroupResource{
			Group:    GroupVersion.Group,
			Resource: resourceType,
		},
		name,
	)
}
