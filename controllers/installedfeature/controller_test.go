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
	"context"
	"errors"
	"github.com/golang/mock/gomock"
	. "github.com/klenkes74/k8s-installed-features-catalogue/controllers/installedfeature"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	k8sclient "sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	// +kubebuilder:scaffold:imports
)

const (
	name        = "basic-feature"
	otherName   = "other-feature"
	namespace   = "default"
	version     = "1.0.0-alpha1"
	provider    = "Kaiserpfalz EDV-Service"
	description = "a basic demonstration feature"
	uri         = "https://www.kaiserpfalz-edv.de/k8s/"
)

var (
	ctx = context.Background()

	successResult = ctrl.Result{Requeue: false}
	errorResult   = ctrl.Result{Requeue: false, RequeueAfter: RequeueTime}

	iftLookupKey        = types.NamespacedName{Name: name, Namespace: namespace}
	iftReconcileRequest = reconcile.Request{
		NamespacedName: iftLookupKey,
	}

	otherLookupKey = types.NamespacedName{Name: otherName, Namespace: namespace}
)

var _ = Describe("InstalledFeature controller basics", func() {
	Context("Finalizer Handling", func() {
		It("should add the finalizer when the finalizer is not set", func() {
			By("By creating a new InstalledFeature without finalizer")

			ift := createIFT(name, namespace, version, provider, description, uri, false, false)
			client.EXPECT().LoadInstalledFeature(gomock.Any(), iftLookupKey).Return(ift, nil)

			expected := copyIFT(ift)
			expected.Finalizers = make([]string, 1)
			expected.Finalizers[0] = FinalizerName

			client.EXPECT().SaveInstalledFeature(gomock.Any(), expected).Return(nil)
			client.EXPECT().InfoEvent(ift, "Create", NoteChangedFeature, ift.GetNamespace(), ift.GetName())

			client.EXPECT().GetInstalledFeaturePatchBase(gomock.Any()).Return(k8sclient.MergeFrom(ift))
			client.EXPECT().PatchInstalledFeatureStatus(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)

			result, err := sut.Reconcile(iftReconcileRequest)
			Expect(result).Should(Equal(successResult))
			Expect(err).ToNot(HaveOccurred())
		})

		It("should remove the finalizer when the finalizer is set while being deleted", func() {
			By("By creating a new InstalledFeature without finalizer")

			ift := createIFT(name, namespace, version, provider, description, uri, true, true)
			client.EXPECT().LoadInstalledFeature(gomock.Any(), iftLookupKey).Return(ift, nil)

			expected := copyIFT(ift)
			expected.Finalizers = make([]string, 0)

			client.EXPECT().SaveInstalledFeature(gomock.Any(), expected).Return(nil)

			client.EXPECT().GetInstalledFeaturePatchBase(gomock.Any()).Return(k8sclient.MergeFrom(ift))
			client.EXPECT().PatchInstalledFeatureStatus(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)

			client.EXPECT().InfoEvent(ift, "Delete", NoteChangedFeature, ift.GetNamespace(), ift.GetName())

			result, err := sut.Reconcile(iftReconcileRequest)
			Expect(result).Should(Equal(successResult))
			Expect(err).ToNot(HaveOccurred())
		})
	})

	Context("Basic creation and deletion handling without dependencies, conflicts and groups", func() {
		It("should create the feature when nothing special is given", func() {
			By("By creating a new InstalledFeature")

			ift := createIFT(name, namespace, version, provider, description, uri, true, false)
			client.EXPECT().LoadInstalledFeature(gomock.Any(), iftLookupKey).Return(ift, nil)

			client.EXPECT().GetInstalledFeaturePatchBase(gomock.Any()).Return(k8sclient.MergeFrom(ift))

			client.EXPECT().InfoEvent(ift, "Update", NoteChangedFeature, ift.GetNamespace(), ift.GetName())

			result, err := sut.Reconcile(iftReconcileRequest)

			Expect(result).Should(Equal(successResult))
			Expect(err).ToNot(HaveOccurred())
		})
	})

	Context("Technical Error handling", func() {
		It("should drop the request when the ift can't be loaded due to NotFoundError", func() {
			By("By having a problem loading the ift")

			client.EXPECT().LoadInstalledFeature(gomock.Any(), iftLookupKey).Return(nil, errors.New("some error"))

			result, err := sut.Reconcile(iftReconcileRequest)
			Expect(result).Should(Equal(errorResult))
			Expect(err).To(HaveOccurred())
		})

		It("should requeue request when the ift can't be loaded due to another error but NotFoundError", func() {
			client.EXPECT().LoadInstalledFeature(gomock.Any(), iftLookupKey).Return(nil, k8serrors.NewNotFound(schema.GroupResource{
				Group:    "features.kaiserpfalz-edv.de",
				Resource: "installedfeatures",
			}, iftLookupKey.Name))

			result, err := sut.Reconcile(iftReconcileRequest)
			Expect(result).Should(Equal(successResult))
			Expect(err).To(HaveOccurred())
		})

		It("should requeue request when writing the reconciled object fails", func() {
			By("By getting a failure while saving the data back into the k8s cluster")

			ift := createIFT(name, namespace, version, provider, description, uri, false, false)
			client.EXPECT().LoadInstalledFeature(gomock.Any(), iftLookupKey).Return(ift, nil)

			expected := copyIFT(ift)
			expected.Finalizers = make([]string, 1)
			expected.Finalizers[0] = FinalizerName

			client.EXPECT().SaveInstalledFeature(gomock.Any(), expected).Return(errors.New("some error"))

			By("sending events", func() {
				client.EXPECT().WarnEvent(ift, "Create",
					NoteFeatureSaveFailed,
					ift.GetNamespace(), ift.GetName(),
					gomock.Any(),
				)
			})

			result, err := sut.Reconcile(iftReconcileRequest)
			Expect(result).Should(Equal(errorResult))
			Expect(err).To(HaveOccurred())

		})

		It("should requeue the request when updating the status fails", func() {
			By("By getting an error when updating the status")

			ift := createIFT(name, namespace, version, provider, description, uri, false, false)
			client.EXPECT().LoadInstalledFeature(gomock.Any(), iftLookupKey).Return(ift, nil)

			expected := copyIFT(ift)
			expected.Finalizers = make([]string, 1)
			expected.Finalizers[0] = FinalizerName

			client.EXPECT().SaveInstalledFeature(gomock.Any(), expected).Return(nil)

			client.EXPECT().GetInstalledFeaturePatchBase(gomock.Any()).Return(k8sclient.MergeFrom(expected))
			client.EXPECT().PatchInstalledFeatureStatus(gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.New("patching status failed"))

			By("Sending events", func() {
				client.EXPECT().WarnEvent(ift, "Create", NoteStatusUpdateFailed, ift.GetNamespace(), ift.GetName(), gomock.Any())
			})

			result, err := sut.Reconcile(iftReconcileRequest)

			Expect(result).Should(Equal(errorResult))
			Expect(err).To(HaveOccurred())
		})
	})
})
