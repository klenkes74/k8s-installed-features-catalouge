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
	. "github.com/onsi/ginkgo"
	// +kubebuilder:scaffold:imports
)

var _ = Describe("InstalledFeature controller handling featuregroups", func() {
	/*
		Context("Handle Library Groups", func() {
			It("should add the status entry on the IFTG when the IFTG has no features yet", func() {
				ift := createIFT(name, namespace, version, provider, description, uri, true, false)
				setGroupToIFT(ift, group, namespace)
				iftg := createIFTG(group, namespace, provider, description, uri, true, false)

				client.EXPECT().LoadInstalledFeature(gomock.Any(), iftLookupKey).Return(ift, nil)

				client.EXPECT().LoadInstalledFeatureGroup(gomock.Any(), iftgLookupKey).Return(iftg, nil)
				client.EXPECT().GetInstalledFeatureGroupPatchBase(gomock.Any()).Return(k8sclient.MergeFrom(iftg))
				client.EXPECT().PatchInstalledFeatureGroupStatus(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)

				client.EXPECT().GetInstalledFeaturePatchBase(gomock.Any()).Return(k8sclient.MergeFrom(ift))
				client.EXPECT().PatchInstalledFeatureStatus(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)

				client.EXPECT().InfoEvent(ift, "Update", "Changed feature %s/%s", ift.GetNamespace(), ift.GetName())

				result, err := sut.Reconcile(iftReconcileRequest)

				Expect(result).Should(Equal(successResult))
				Expect(err).ToNot(HaveOccurred())
			})

			It("should add the status entry on the IFTG when the IFTG has already features", func() {
				ift := createIFT(name, namespace, version, provider, description, uri, true, false)
				setGroupToIFT(ift, group, namespace)
				iftg := createIFTG(group, namespace, provider, description, uri, true, false)
				iftg.Status.Features = make([]InstalledFeatureGroupListedFeature, 1)
				iftg.Status.Features[0] = InstalledFeatureGroupListedFeature{
					Namespace: namespace,
					Name:      "other-feature",
				}

				client.EXPECT().LoadInstalledFeature(gomock.Any(), iftLookupKey).Return(ift, nil)

				client.EXPECT().LoadInstalledFeatureGroup(gomock.Any(), iftgLookupKey).Return(iftg, nil)
				client.EXPECT().GetInstalledFeatureGroupPatchBase(gomock.Any()).Return(k8sclient.MergeFrom(iftg))
				client.EXPECT().PatchInstalledFeatureGroupStatus(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)

				client.EXPECT().GetInstalledFeaturePatchBase(gomock.Any()).Return(k8sclient.MergeFrom(ift))
				client.EXPECT().PatchInstalledFeatureStatus(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)

				client.EXPECT().InfoEvent(ift, "Update", installedfeature.NoteChangedFeature, ift.GetNamespace(), ift.GetName())

				result, err := sut.Reconcile(iftReconcileRequest)

				Expect(result).Should(Equal(successResult))
				Expect(err).ToNot(HaveOccurred())
			})

			It("should not add the feature to the status entry on the IFTG when the IFTG already lists this feature", func() {
				ift := createIFT(name, namespace, version, provider, description, uri, true, false)
				setGroupToIFT(ift, group, namespace)
				iftg := createIFTG(group, namespace, provider, description, uri, true, false)
				iftg.Status.Features = make([]InstalledFeatureGroupListedFeature, 1)
				iftg.Status.Features[0] = InstalledFeatureGroupListedFeature{
					Namespace: namespace,
					Name:      name,
				}

				client.EXPECT().LoadInstalledFeature(gomock.Any(), iftLookupKey).Return(ift, nil)

				client.EXPECT().LoadInstalledFeatureGroup(gomock.Any(), iftgLookupKey).Return(iftg, nil)
				client.EXPECT().GetInstalledFeatureGroupPatchBase(gomock.Any()).Return(k8sclient.MergeFrom(iftg))

				client.EXPECT().GetInstalledFeaturePatchBase(gomock.Any()).Return(k8sclient.MergeFrom(ift))
				client.EXPECT().PatchInstalledFeatureStatus(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)

				client.EXPECT().InfoEvent(ift, "Update", installedfeature.NoteChangedFeature, ift.GetNamespace(), ift.GetName())

				result, err := sut.Reconcile(iftReconcileRequest)

				Expect(result).Should(Equal(successResult))
				Expect(err).ToNot(HaveOccurred())
			})

			It("should remove the status entry on the IFTG when IFT is deleted and is listed in IFTG", func() {
				ift := createIFT(name, namespace, version, provider, description, uri, false, true)
				setGroupToIFT(ift, group, namespace)
				iftg := createIFTG(group, namespace, provider, description, uri, true, false)
				iftg.Status.Features = make([]InstalledFeatureGroupListedFeature, 1)
				iftg.Status.Features[0] = InstalledFeatureGroupListedFeature{
					Namespace: namespace,
					Name:      name,
				}

				client.EXPECT().LoadInstalledFeature(gomock.Any(), iftLookupKey).Return(ift, nil)

				client.EXPECT().LoadInstalledFeatureGroup(gomock.Any(), iftgLookupKey).Return(iftg, nil)
				client.EXPECT().GetInstalledFeatureGroupPatchBase(gomock.Any()).Return(k8sclient.MergeFrom(iftg))
				client.EXPECT().PatchInstalledFeatureGroupStatus(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)

				client.EXPECT().GetInstalledFeaturePatchBase(gomock.Any()).Return(k8sclient.MergeFrom(ift))
				client.EXPECT().PatchInstalledFeatureStatus(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)

				client.EXPECT().InfoEvent(ift, "Delete", installedfeature.NoteChangedFeature, ift.GetNamespace(), ift.GetName())

				result, err := sut.Reconcile(iftReconcileRequest)

				Expect(result).Should(Equal(successResult))
				Expect(err).ToNot(HaveOccurred())
			})

			It("should requeue the request when IFTG can't be loaded", func() {
				ift := createIFT(name, namespace, version, provider, description, uri, false, true)
				setGroupToIFT(ift, group, namespace)

				client.EXPECT().LoadInstalledFeature(gomock.Any(), iftLookupKey).Return(ift, nil)

				client.EXPECT().LoadInstalledFeatureGroup(gomock.Any(), iftgLookupKey).Return(nil, errors.New("can not load IFTG"))

				By("Sending events", func() {

					client.EXPECT().WarnEvent(ift, "Delete", installedfeature.NoteUpdatingGroupFailed, types.NamespacedName{Namespace: ift.GetNamespace(), Name: ift.GetName()})
				})

				result, err := sut.Reconcile(iftReconcileRequest)

				Expect(result).Should(Equal(errorResult))
				Expect(err).Should(HaveOccurred())
			})
		})
	*/
})
