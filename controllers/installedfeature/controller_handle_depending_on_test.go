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
	"errors"
	"github.com/golang/mock/gomock"
	. "github.com/klenkes74/k8s-installed-features-catalogue/api/v1alpha1"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	k8sclient "sigs.k8s.io/controller-runtime/pkg/client"
)

var _ = Describe("InstalledFeature depending feature handling", func() {
	Context("Handling dependencies", func() {
		It("Should add dependency status when there is a dependency defined that has already other depending features", func() {
			ift := createIFT(name, namespace, version, provider, description, uri, true, false)
			ift.Spec.DependsOn = []InstalledFeatureRef{
				{Namespace: namespace, Name: otherName},
			}

			By("Loading and saving the feature", func() {
				client.EXPECT().LoadInstalledFeature(gomock.Any(), iftLookupKey).Return(ift, nil)

				ift.Status.MissingDependencies = ift.Spec.DependsOn

				iftPatch := k8sclient.MergeFrom(ift)
				client.EXPECT().GetInstalledFeaturePatchBase(ift).Return(iftPatch)
				client.EXPECT().GetInstalledFeaturePatchBase(ift).Return(iftPatch)

				client.EXPECT().PatchInstalledFeatureStatus(gomock.Any(), gomock.Any(), iftPatch).Return(nil).Times(2)
			})

			other := createIFT(otherName, namespace, version, provider, description, uri, true, false)
			other.Status.DependingFeatures = []InstalledFeatureRef{
				{Namespace: "other", Name: "other"},
			}
			By("Updating the dependent list in the status of the dependency", func() {
				client.EXPECT().LoadInstalledFeature(gomock.Any(), otherLookupKey).Return(other, nil)

				iftPatch := k8sclient.MergeFrom(other)
				client.EXPECT().GetInstalledFeaturePatchBase(other).Return(iftPatch)
				client.EXPECT().PatchInstalledFeatureStatus(gomock.Any(), gomock.Any(), iftPatch).Return(nil)
			})

			result, err := sut.Reconcile(iftReconcileRequest)

			Expect(result).Should(Equal(successResult))
			Expect(err).ShouldNot(HaveOccurred())
		})

		It("Should add dependency status when the dependency is already listed in dependency status", func() {
			ift := createIFT(name, namespace, version, provider, description, uri, true, false)
			ift.Spec.DependsOn = []InstalledFeatureRef{
				{Namespace: namespace, Name: otherName},
			}

			By("Loading and saving the feature", func() {
				client.EXPECT().LoadInstalledFeature(gomock.Any(), iftLookupKey).Return(ift, nil)

				iftPatch := k8sclient.MergeFrom(ift)
				client.EXPECT().GetInstalledFeaturePatchBase(ift).Return(iftPatch).Times(2)
				client.EXPECT().PatchInstalledFeatureStatus(gomock.Any(), gomock.Any(), iftPatch).Return(nil).Times(2)
			})

			other := createIFT(otherName, namespace, version, provider, description, uri, true, false)
			other.Status.DependingFeatures = []InstalledFeatureRef{
				{
					Namespace: namespace,
					Name:      name,
				},
			}
			By("Updating the dependent list in the status of the dependency", func() {
				client.EXPECT().LoadInstalledFeature(gomock.Any(), otherLookupKey).Return(other, nil)

				iftPatch := k8sclient.MergeFrom(other)
				client.EXPECT().GetInstalledFeaturePatchBase(other).Return(iftPatch)
			})

			result, err := sut.Reconcile(iftReconcileRequest)

			Expect(result).Should(Equal(successResult))
			Expect(err).ShouldNot(HaveOccurred())
		})

		It("Should remove dependency status when the instance is deleted and already listed in the status of dependency", func() {
			ift := createIFT(name, namespace, version, provider, description, uri, true, true)
			ift.Spec.DependsOn = []InstalledFeatureRef{
				{Namespace: namespace, Name: otherName},
			}

			By("Loading and saving the feature", func() {
				client.EXPECT().LoadInstalledFeature(gomock.Any(), iftLookupKey).Return(ift, nil)

				client.EXPECT().SaveInstalledFeature(gomock.Any(), ift).Return(nil)

				iftPatch := k8sclient.MergeFrom(ift)
				client.EXPECT().GetInstalledFeaturePatchBase(ift).Return(iftPatch).Times(2)
				client.EXPECT().PatchInstalledFeatureStatus(gomock.Any(), gomock.Any(), iftPatch).Return(nil).Times(2)
			})

			other := createIFT(otherName, namespace, version, provider, description, uri, true, false)
			other.Status.DependingFeatures = []InstalledFeatureRef{
				{
					Namespace: namespace,
					Name:      name,
				},
			}
			By("Updating the dependent list in the status of the dependency", func() {
				client.EXPECT().LoadInstalledFeature(gomock.Any(), otherLookupKey).Return(other, nil)

				iftPatch := k8sclient.MergeFrom(other)
				client.EXPECT().GetInstalledFeaturePatchBase(other).Return(iftPatch)
				client.EXPECT().PatchInstalledFeatureStatus(gomock.Any(), gomock.Any(), iftPatch).Return(nil)
			})

			result, err := sut.Reconcile(iftReconcileRequest)

			Expect(result).Should(Equal(successResult))
			Expect(err).ShouldNot(HaveOccurred())
		})

		It("Should add dependency status when there is a dependency defined that has no other dependencies", func() {
			ift := createIFT(name, namespace, version, provider, description, uri, true, false)

			ift.Spec.DependsOn = []InstalledFeatureRef{
				{Namespace: namespace, Name: otherName},
			}
			By("Loading and saving the feature", func() {
				client.EXPECT().LoadInstalledFeature(gomock.Any(), iftLookupKey).Return(ift, nil)

				iftPatch := k8sclient.MergeFrom(ift)
				client.EXPECT().GetInstalledFeaturePatchBase(ift).Return(iftPatch).Times(2)
				client.EXPECT().PatchInstalledFeatureStatus(gomock.Any(), gomock.Any(), iftPatch).Return(nil).Times(2)
			})

			other := createIFT(otherName, namespace, version, provider, description, uri, true, false)
			By("Updating the dependent list in the status of the dependency", func() {
				client.EXPECT().LoadInstalledFeature(gomock.Any(), otherLookupKey).Return(other, nil)

				iftPatch := k8sclient.MergeFrom(ift)
				client.EXPECT().GetInstalledFeaturePatchBase(other).Return(iftPatch)
				client.EXPECT().PatchInstalledFeatureStatus(gomock.Any(), gomock.Any(), iftPatch).Return(nil)
			})

			result, err := sut.Reconcile(iftReconcileRequest)

			Expect(result).Should(Equal(successResult))
			Expect(err).ShouldNot(HaveOccurred())
		})

		It("Should requeue the reconcile when the dependency status can not be changed", func() {
			ift := createIFT(name, namespace, version, provider, description, uri, true, false)
			ift.Spec.DependsOn = []InstalledFeatureRef{
				{Namespace: namespace, Name: otherName},
			}

			By("Loading and saving the feature", func() {
				client.EXPECT().LoadInstalledFeature(gomock.Any(), iftLookupKey).Return(ift, nil)

				ift.Status.MissingDependencies = ift.Spec.DependsOn

				iftPatch := k8sclient.MergeFrom(ift)
				client.EXPECT().GetInstalledFeaturePatchBase(ift).Return(iftPatch)
				client.EXPECT().GetInstalledFeaturePatchBase(ift).Return(iftPatch)

				client.EXPECT().PatchInstalledFeatureStatus(gomock.Any(), gomock.Any(), iftPatch).Return(nil)
				client.EXPECT().PatchInstalledFeatureStatus(gomock.Any(), gomock.Any(), iftPatch).Return(errors.New("could not update status"))
			})

			other := createIFT(otherName, namespace, version, provider, description, uri, true, false)
			other.Status.DependingFeatures = []InstalledFeatureRef{
				{
					Namespace: "other",
					Name:      "other",
				},
			}
			By("Updating the dependent list in the status of the dependency", func() {
				client.EXPECT().LoadInstalledFeature(gomock.Any(), otherLookupKey).Return(other, nil)

				iftPatch := k8sclient.MergeFrom(other)
				client.EXPECT().GetInstalledFeaturePatchBase(other).Return(iftPatch)
				client.EXPECT().PatchInstalledFeatureStatus(gomock.Any(), gomock.Any(), iftPatch).Return(nil)
			})

			result, err := sut.Reconcile(iftReconcileRequest)

			Expect(result).Should(Equal(errorResult))
			Expect(err).Should(HaveOccurred())
		})

		It("Should requeue the reconcile when the instance status can not be changed", func() {
			ift := createIFT(name, namespace, version, provider, description, uri, true, false)

			ift.Spec.DependsOn = []InstalledFeatureRef{
				{Namespace: namespace, Name: otherName},
			}
			By("Loading and saving the feature", func() {
				client.EXPECT().LoadInstalledFeature(gomock.Any(), iftLookupKey).Return(ift, nil)

				ift.Status.MissingDependencies = ift.Spec.DependsOn

				iftPatch := k8sclient.MergeFrom(ift)
				client.EXPECT().GetInstalledFeaturePatchBase(ift).Return(iftPatch)
				client.EXPECT().PatchInstalledFeatureStatus(gomock.Any(), gomock.Any(), iftPatch).Return(errors.New("dependency status can not be patched"))
			})

			other := createIFT(otherName, namespace, version, provider, description, uri, true, false)
			By("Updating the dependent list in the status of the dependency", func() {
				client.EXPECT().LoadInstalledFeature(gomock.Any(), otherLookupKey).Return(other, nil)

				iftPatch := k8sclient.MergeFrom(ift)
				client.EXPECT().GetInstalledFeaturePatchBase(other).Return(iftPatch)
			})

			result, err := sut.Reconcile(iftReconcileRequest)

			Expect(result).Should(Equal(errorResult))
			Expect(err).Should(HaveOccurred())
		})

		It("Should mark missing dependency when dependency is marked as deleted", func() {
			ift := createIFT(name, namespace, version, provider, description, uri, true, false)

			ift.Spec.DependsOn = []InstalledFeatureRef{
				{Namespace: namespace, Name: otherName},
			}
			By("Loading and saving the feature", func() {
				client.EXPECT().LoadInstalledFeature(gomock.Any(), iftLookupKey).Return(ift, nil)

				iftPatch := k8sclient.MergeFrom(ift)

				ift.Status.MissingDependencies = []InstalledFeatureRef{
					ift.Spec.DependsOn[0],
				}

				client.EXPECT().GetInstalledFeaturePatchBase(ift).Return(iftPatch)
				client.EXPECT().PatchInstalledFeatureStatus(gomock.Any(), gomock.Any(), iftPatch).Return(nil)
			})

			other := createIFT(otherName, namespace, version, provider, description, uri, true, true)
			By("Updating the dependent list in the status of the dependency", func() {
				client.EXPECT().LoadInstalledFeature(gomock.Any(), otherLookupKey).Return(other, nil)

			})

			result, err := sut.Reconcile(iftReconcileRequest)

			Expect(result).Should(Equal(errorResult))
			Expect(err).Should(HaveOccurred())
		})

		It("Should mark missing dependency when dependency can not be loaded", func() {
			ift := createIFT(name, namespace, version, provider, description, uri, true, false)

			ift.Spec.DependsOn = []InstalledFeatureRef{
				{Namespace: namespace, Name: otherName},
			}
			By("Loading and saving the feature", func() {
				client.EXPECT().LoadInstalledFeature(gomock.Any(), iftLookupKey).Return(ift, nil)

				iftPatch := k8sclient.MergeFrom(ift)
				client.EXPECT().GetInstalledFeaturePatchBase(ift).Return(iftPatch)
				client.EXPECT().PatchInstalledFeatureStatus(gomock.Any(), gomock.Any(), iftPatch).Return(nil)
			})

			By("Updating the dependent list in the status of the dependency", func() {
				client.EXPECT().LoadInstalledFeature(gomock.Any(), otherLookupKey).Return(nil, errors.New("other feature not found"))
			})

			result, err := sut.Reconcile(iftReconcileRequest)

			Expect(result).Should(Equal(errorResult))
			Expect(err).Should(HaveOccurred())
		})
	})

	Context("Handling technical failures", func() {
		It("Should requeue the request when patching the dependency fails", func() {
			ift := createIFT(name, namespace, version, provider, description, uri, true, false)

			ift.Spec.DependsOn = []InstalledFeatureRef{
				{Namespace: namespace, Name: otherName},
			}
			By("Loading and saving the feature", func() {
				client.EXPECT().LoadInstalledFeature(gomock.Any(), iftLookupKey).Return(ift, nil)

				client.EXPECT().GetInstalledFeaturePatchBase(ift).Return(k8sclient.MergeFrom(ift))
			})

			other := createIFT(otherName, namespace, version, provider, description, uri, true, false)
			other.Status.DependingFeatures = []InstalledFeatureRef{
				{Namespace: "other", Name: "other"},
			}
			By("Updating the dependent list in the status of the dependency", func() {
				client.EXPECT().LoadInstalledFeature(gomock.Any(), otherLookupKey).Return(other, nil)

				client.EXPECT().GetInstalledFeaturePatchBase(other).Return(k8sclient.MergeFrom(other))
				client.EXPECT().PatchInstalledFeatureStatus(gomock.Any(), other, k8sclient.MergeFrom(other)).Return(errors.New("patching the dependeny failed"))
			})

			result, err := sut.Reconcile(iftReconcileRequest)

			Expect(result).Should(Equal(errorResult))
			Expect(err).Should(HaveOccurred())
		})

		It("Should requeue the request when patching the feature fails", func() {
			ift := createIFT(name, namespace, version, provider, description, uri, true, false)

			ift.Spec.DependsOn = []InstalledFeatureRef{
				{Namespace: namespace, Name: otherName},
			}
			By("Loading and saving the feature", func() {
				client.EXPECT().LoadInstalledFeature(gomock.Any(), iftLookupKey).Return(ift, nil)

				client.EXPECT().GetInstalledFeaturePatchBase(ift).Return(k8sclient.MergeFrom(ift))
				client.EXPECT().PatchInstalledFeatureStatus(ctx, ift, k8sclient.MergeFrom(ift)).Return(errors.New("patching the feature failed"))
			})

			other := createIFT(otherName, namespace, version, provider, description, uri, true, false)
			other.Status.DependingFeatures = []InstalledFeatureRef{
				{Namespace: "other", Name: "other"},
			}
			By("Updating the dependent list in the status of the dependency", func() {
				client.EXPECT().LoadInstalledFeature(gomock.Any(), otherLookupKey).Return(other, nil)

				client.EXPECT().GetInstalledFeaturePatchBase(other).Return(k8sclient.MergeFrom(other))
				client.EXPECT().PatchInstalledFeatureStatus(gomock.Any(), other, k8sclient.MergeFrom(other)).Return(nil)
			})

			result, err := sut.Reconcile(iftReconcileRequest)

			Expect(result).Should(Equal(errorResult))
			Expect(err).Should(HaveOccurred())
		})
	})
})
