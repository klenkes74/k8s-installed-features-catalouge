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
	. "github.com/klenkes74/k8s-installed-features-catalogue/api/v1alpha1"
	"github.com/klenkes74/k8s-installed-features-catalogue/controllers/installedfeature"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	k8sclient "sigs.k8s.io/controller-runtime/pkg/client"
)

var _ = Describe("InstalledFeature depending feature handling", func() {
	Context("Handling dependencies", func() {
		It("Should add dependency status when there is a dependency defined that has already other depending features", func() {
			ift := createIFT(name, namespace, version, provider, description, uri, true, false)
			patch := k8sclient.MergeFrom(ift)
			ift.Spec.DependsOn = []InstalledFeatureRef{
				{Namespace: namespace, Name: otherName},
			}

			other := createIFT(otherName, namespace, version, provider, description, uri, true, false)
			other.Status.DependingFeatures = []InstalledFeatureRef{
				{Namespace: "other", Name: "other"},
			}

			By("loading the reconciled feature", func() {
				client.EXPECT().LoadInstalledFeature(gomock.Any(), iftLookupKey).Return(ift, nil)
			})

			By("loading, marking the dependency status, and reconciling the dependency", func() {
				client.EXPECT().LoadInstalledFeature(gomock.Any(), otherLookupKey).Return(other, nil)

				otherPatch := k8sclient.MergeFrom(other)
				client.EXPECT().GetInstalledFeaturePatchBase(other).Return(otherPatch)
				client.EXPECT().PatchInstalledFeatureStatus(ctx, other, otherPatch).Return(nil)

				client.EXPECT().ReconcileFeature(ctx, other).Return(nil)
			})

			By("saving the reconciled feature", func() {
				client.EXPECT().GetInstalledFeaturePatchBase(ift).Return(patch)
				client.EXPECT().InfoEvent(ift, "Update",
					installedfeature.NoteChangedFeature,
					ift.GetNamespace(),
					ift.GetName(),
				)
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
				client.EXPECT().GetInstalledFeaturePatchBase(ift).Return(iftPatch)
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

				client.EXPECT().ReconcileFeature(gomock.Any(), other)
			})

			client.EXPECT().InfoEvent(ift, "Update", installedfeature.NoteChangedFeature, ift.GetNamespace(), ift.GetName())

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
				client.EXPECT().GetInstalledFeaturePatchBase(ift).Return(iftPatch)
				client.EXPECT().PatchInstalledFeatureStatus(gomock.Any(), gomock.Any(), iftPatch).Return(nil)
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

				client.EXPECT().ReconcileFeature(gomock.Any(), other)
			})

			By("Sending events", func() {
				client.EXPECT().InfoEvent(ift, "Delete", installedfeature.NoteChangedFeature, ift.GetNamespace(), ift.GetName())
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

				patch := k8sclient.MergeFrom(ift)
				client.EXPECT().GetInstalledFeaturePatchBase(ift).Return(patch)
			})

			other := createIFT(otherName, namespace, version, provider, description, uri, true, false)
			By("Updating the dependent list in the status of the dependency", func() {
				client.EXPECT().LoadInstalledFeature(gomock.Any(), otherLookupKey).Return(other, nil)

				patch := k8sclient.MergeFrom(other)
				client.EXPECT().GetInstalledFeaturePatchBase(other).Return(patch)
				client.EXPECT().PatchInstalledFeatureStatus(gomock.Any(), gomock.Any(), patch).Return(nil)

				client.EXPECT().ReconcileFeature(gomock.Any(), other)
			})

			By("Sending events", func() {
				client.EXPECT().InfoEvent(ift, "Update", installedfeature.NoteChangedFeature, ift.GetNamespace(), ift.GetName())
			})

			result, err := sut.Reconcile(iftReconcileRequest)

			Expect(result).Should(Equal(successResult))
			Expect(err).ShouldNot(HaveOccurred())
		})
		/* TODO 2020-10-29 klenkes74 can't test via gomock
		It("Should requeue the reconcile when the dependency status can not be changed", func() {
			ift := createIFT(name, namespace, version, provider, description, uri, true, false)
			ift.Spec.DependsOn = []InstalledFeatureRef{
				{Namespace: namespace, Name: otherName},
			}

			By("Loading and saving the feature", func() {
				client.EXPECT().LoadInstalledFeature(gomock.Any(), iftLookupKey).Return(ift, nil)

				ift.Status.MissingDependencies = ift.Spec.DependsOn
			})

			other := createIFT(otherName, namespace, version, provider, description, uri, true, false)
			otherPatch := k8sclient.MergeFrom(other)
			other.Status.DependingFeatures = []InstalledFeatureRef{
				{
					Namespace: "other",
					Name:      "other",
				},
			}

			By("Updating the dependent list in the status of the dependency", func() {
				client.EXPECT().LoadInstalledFeature(gomock.Any(), otherLookupKey).Return(other, nil)

				client.EXPECT().GetInstalledFeaturePatchBase(other).Return(otherPatch)
				client.EXPECT().PatchInstalledFeatureStatus(gomock.Any(), gomock.Any(), otherPatch)

				client.EXPECT().ReconcileFeature(gomock.Any(), other)
			})

			By("sending events", func() {
				client.EXPECT().WarnEvent(ift, "Update",
					installedfeature.NoteStatusUpdateFailed,
					ift.GetNamespace(), ift.GetName(),
					gomock.Any(),
				)
			})

			result, err := sut.Reconcile(iftReconcileRequest)

			Expect(result).Should(Equal(errorResult))
			Expect(err).Should(HaveOccurred())
		})
		*/
		/* TODO 2020-10-29 klenkes74 -- gomock can't handle interface lists as parameter
		It("Should requeue the reconcile when the instance status can not be changed", func() {
			ift := createIFT(name, namespace, version, provider, description, uri, true, false)

			ift.Spec.DependsOn = []InstalledFeatureRef{
				{Namespace: namespace, Name: otherName},
			}
			By("Loading and saving the feature", func() {
				client.EXPECT().LoadInstalledFeature(gomock.Any(), iftLookupKey).Return(ift, nil)

				ift.Status.MissingDependencies = ift.Spec.DependsOn
			})

			other := createIFT(otherName, namespace, version, provider, description, uri, true, false)
			By("Updating the dependent list in the status of the dependency", func() {
				client.EXPECT().LoadInstalledFeature(gomock.Any(), otherLookupKey).Return(other, nil)

				patch := k8sclient.MergeFrom(ift)
				client.EXPECT().GetInstalledFeaturePatchBase(other).Return(patch)
				client.EXPECT().PatchInstalledFeatureStatus(gomock.Any(), gomock.Any(), patch)

				client.EXPECT().ReconcileFeature(gomock.Any(), other).Return(nil)
			})

			By("sending events", func() {
				client.EXPECT().WarnEvent(ift, "Update",
					installedfeature.NoteUpdatingDependenciesFailed,
					types.NamespacedName{Namespace: ift.GetNamespace(), Name: ift.GetName()},
				)
			})
			result, err := sut.Reconcile(iftReconcileRequest)

			Expect(result).Should(Equal(errorResult))
			Expect(err).Should(HaveOccurred())
		})
		*/
		/* TODO 2020-10-29 klenkes74 -- gomock can't handle interface lists as parameter
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

			By("sending events", func() {
				client.EXPECT().WarnEvent(ift, "Update", installedfeature.NoteMissingDependencies, gomock.Any())
				client.EXPECT().WarnEvent(ift, "Update", installedfeature.NoteUpdatingDependenciesFailed, gomock.Any())
			})

			result, err := sut.Reconcile(iftReconcileRequest)

			Expect(result).Should(Equal(errorResult))
			Expect(err).Should(HaveOccurred())
		})
		*/
		/* TODO 2020-10-29 klenkes74 -- gomock can't handle interface lists as parameter
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

				By("sending events", func() {
					client.EXPECT().WarnEvent(ift, "Update", installedfeature.NoteMissingDependencies, gomock.Any())
					client.EXPECT().WarnEvent(ift, "Update", installedfeature.NoteUpdatingDependenciesFailed, gomock.Any())
				})

				result, err := sut.Reconcile(iftReconcileRequest)

				Expect(result).Should(Equal(errorResult))
				Expect(err).Should(HaveOccurred())
			})
		})
		*/
		/* TODO 2020-10-29 klenkes74 -- gomock can't handle interface lists as parameter
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
					client.EXPECT().PatchInstalledFeatureStatus(gomock.Any(), other, k8sclient.MergeFrom(other)).Return(errors.New("patching the dependency failed"))
				})

				By("sending events", func() {
					client.EXPECT().WarnEvent(ift, "Update",
						installedfeature.NoteUpdatingDependenciesFailed,
						types.NamespacedName{Namespace: ift.GetNamespace(), Name: ift.GetName()},
					)
				})

				result, err := sut.Reconcile(iftReconcileRequest)

				Expect(result).Should(Equal(errorResult))
				Expect(err).Should(HaveOccurred())
			})
		*/
		/* TODO 2020-10-29 klenkes74 -- gomock can't handle interface lists as parameter
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

				client.EXPECT().ReconcileFeature(gomock.Any(), other)
			})

			By("sending events", func() {
				client.EXPECT().WarnEvent(ift, "Update",
					installedfeature.NoteUpdatingDependenciesFailed,
					types.NamespacedName{Namespace: ift.GetNamespace(), Name: ift.GetName()},
				)
			})

			result, err := sut.Reconcile(iftReconcileRequest)

			Expect(result).Should(Equal(errorResult))
			Expect(err).Should(HaveOccurred())
		})
		*/
	})
})
