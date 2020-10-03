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

package installedfeaturegroup_test

import (
	"errors"
	"github.com/golang/mock/gomock"
	. "github.com/klenkes74/k8s-installed-features-catalogue/api/v1alpha1"
	. "github.com/klenkes74/k8s-installed-features-catalogue/controllers/installedfeaturegroup"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pborman/uuid"
	errors2 "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	k8sclient "sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"time"
	// +kubebuilder:scaffold:imports
)

var _ = Describe("InstalledFeature controller", func() {
	const (
		name        = "basic-feature-group"
		namespace   = "default"
		provider    = "Kaiserpfalz EDV-Service"
		description = "a basic demonstration feature"
		uri         = "https://www.kaiserpfalz-edv.de/k8s/"
	)
	var (
		iftLookupKey        = types.NamespacedName{Name: name, Namespace: namespace}
		iftReconcileRequest = reconcile.Request{
			NamespacedName: iftLookupKey,
		}
	)

	Context("When installing a InstalledFeature CR", func() {
		It("should be created when there are no conflicting features installed and all dependencies met", func() {
			By("By creating a new InstalledFeature")

			ift := createIFTG(name, namespace, provider, description, uri, true, false)
			client.EXPECT().LoadInstalledFeatureGroup(gomock.Any(), iftLookupKey).Return(ift, nil)

			client.EXPECT().GetInstalledFeatureGroupPatchBase(gomock.Any()).Return(k8sclient.MergeFrom(ift))
			client.EXPECT().PatchInstalledFeatureGroupStatus(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)

			result, err := sut.Reconcile(iftReconcileRequest)

			Expect(result).Should(Equal(reconcile.Result{Requeue: false}))
			Expect(err).ShouldNot(HaveOccurred())
		})

		It("should add the finalizer when the finalizer is not set", func() {
			By("By creating a new InstalledFeature without finalizer")

			ift := createIFTG(name, namespace, provider, description, uri, false, false)
			client.EXPECT().LoadInstalledFeatureGroup(gomock.Any(), iftLookupKey).Return(ift, nil)

			expected := copyIFTG(ift)
			expected.Finalizers = make([]string, 1)
			expected.Finalizers[0] = FinalizerName

			client.EXPECT().SaveInstalledFeatureGroup(gomock.Any(), expected).Return(nil)

			client.EXPECT().GetInstalledFeatureGroupPatchBase(gomock.Any()).Return(k8sclient.MergeFrom(ift))
			client.EXPECT().PatchInstalledFeatureGroupStatus(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)

			result, err := sut.Reconcile(iftReconcileRequest)
			Expect(result).Should(Equal(reconcile.Result{Requeue: false}))
			Expect(err).ShouldNot(HaveOccurred())
		})

		It("should remove the finalizer when the instance is deleted", func() {
			By("By creating a new InstalledFeature without finalizer")

			ift := createIFTG(name, namespace, provider, description, uri, true, true)
			client.EXPECT().LoadInstalledFeatureGroup(gomock.Any(), iftLookupKey).Return(ift, nil)

			expected := copyIFTG(ift)
			expected.Finalizers = make([]string, 0)

			client.EXPECT().SaveInstalledFeatureGroup(gomock.Any(), expected).Return(nil)

			client.EXPECT().GetInstalledFeatureGroupPatchBase(gomock.Any()).Return(k8sclient.MergeFrom(ift))
			client.EXPECT().PatchInstalledFeatureGroupStatus(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)

			result, err := sut.Reconcile(iftReconcileRequest)
			Expect(result).Should(Equal(reconcile.Result{Requeue: false}))
			Expect(err).ShouldNot(HaveOccurred())
		})
	})

	Context("Delete an existing InstalledFeature", func() {
		It("should be deleted when there are no dependencies on the removed feature", func() {
			// TODO 2020-09-26 klenkes74 Implement this test
		})

		It("should not be deleted when there are dependencies on the removed feature", func() {
			// TODO 2020-09-26 klenkes74 Implement this test
		})
	})

	Context("Technical handling", func() {
		It("should drop the request when the ift can't be loaded due to NotFoundError", func() {
			By("By having a problem loading the ift")

			client.EXPECT().LoadInstalledFeatureGroup(gomock.Any(), iftLookupKey).Return(nil, errors.New("some error"))

			result, err := sut.Reconcile(iftReconcileRequest)
			Expect(result).Should(Equal(reconcile.Result{RequeueAfter: 60}))
			Expect(err).To(HaveOccurred())
		})

		It("should requeue request when the ift can't be loaded due to another error but NotFoundError", func() {
			client.EXPECT().LoadInstalledFeatureGroup(gomock.Any(), iftLookupKey).Return(nil, errors2.NewNotFound(schema.GroupResource{
				Group:    "features.kaiserpfalz-edv.de",
				Resource: "installedfeatures",
			}, iftLookupKey.Name))

			result, err := sut.Reconcile(iftReconcileRequest)
			Expect(result).Should(Equal(reconcile.Result{Requeue: false}))
			Expect(err).To(HaveOccurred())
		})

		It("should requeue request when writing the reconciled object fails", func() {
			By("By getting a failure while saving the data back into the k8s cluster")

			ift := createIFTG(name, namespace, provider, description, uri, false, false)
			client.EXPECT().LoadInstalledFeatureGroup(gomock.Any(), iftLookupKey).Return(ift, nil)

			expected := copyIFTG(ift)
			expected.Finalizers = make([]string, 1)
			expected.Finalizers[0] = FinalizerName

			client.EXPECT().SaveInstalledFeatureGroup(gomock.Any(), expected).Return(errors.New("some error"))

			result, err := sut.Reconcile(iftReconcileRequest)
			Expect(result).Should(Equal(reconcile.Result{RequeueAfter: 60}))
			Expect(err).To(HaveOccurred())

		})

		It("should requeue request when writing the reconciled object fails", func() {
			By("By getting a failure while saving the data back into the k8s cluster")

			ift := createIFTG(name, namespace, provider, description, uri, false, false)
			client.EXPECT().LoadInstalledFeatureGroup(gomock.Any(), iftLookupKey).Return(ift, nil)

			expected := copyIFTG(ift)
			expected.Finalizers = make([]string, 1)
			expected.Finalizers[0] = FinalizerName

			client.EXPECT().SaveInstalledFeatureGroup(gomock.Any(), expected).Return(errors.New("some error"))

			result, err := sut.Reconcile(iftReconcileRequest)
			Expect(result).Should(Equal(reconcile.Result{RequeueAfter: 60}))
			Expect(err).To(HaveOccurred())

		})

		It("should requeue the request when updating the status fails", func() {
			By("By getting an error when updating the status")

			ift := createIFTG(name, namespace, provider, description, uri, true, false)
			client.EXPECT().LoadInstalledFeatureGroup(gomock.Any(), iftLookupKey).Return(ift, nil)

			client.EXPECT().GetInstalledFeatureGroupPatchBase(gomock.Any()).Return(k8sclient.MergeFrom(ift))
			client.EXPECT().
				PatchInstalledFeatureGroupStatus(gomock.Any(), gomock.Any(), gomock.Any()).
				Return(errors.New("patching failed"))

			result, err := sut.Reconcile(iftReconcileRequest)

			Expect(result).Should(Equal(reconcile.Result{RequeueAfter: 60}))
			Expect(err).To(HaveOccurred())
		})
	})
})

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
			UID:               types.UID(uuid.New()),
		},
		Spec: InstalledFeatureGroupSpec{
			Provider:    provider,
			Description: description,
			Uri:         uri,
		},
	}

	if finalizer {
		result.Finalizers = make([]string, 1)
		result.Finalizers[0] = FinalizerName
	}

	if deleted {
		deletionGracePeriod := int64(60)
		result.DeletionGracePeriodSeconds = &deletionGracePeriod
		result.DeletionTimestamp = &metav1.Time{Time: time.Now().Add(2 * time.Minute)}
	}

	return result
}

func copyIFTG(orig *InstalledFeatureGroup) *InstalledFeatureGroup {
	//goland:noinspection GoDeprecation
	result := &InstalledFeatureGroup{
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
		Spec: InstalledFeatureGroupSpec{
			Provider:    orig.Spec.Provider,
			Description: orig.Spec.Description,
			Uri:         orig.Spec.Uri,
		},
		Status: InstalledFeatureGroupStatus{
			Phase:   orig.Status.Phase,
			Message: orig.Status.Message,
		},
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

	if len(orig.Status.Features) > 0 {
		result.Status.Features = make([]InstalledFeatureGroupListedFeature, len(orig.Status.Features))
		for i, feature := range orig.Status.Features {
			result.Status.Features[i] = InstalledFeatureGroupListedFeature{
				Namespace: feature.Namespace,
				Name:      feature.Name,
			}
		}
	}

	return result
}
