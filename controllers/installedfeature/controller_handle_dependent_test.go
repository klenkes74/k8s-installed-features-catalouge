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
	. "github.com/klenkes74/k8s-installed-features-catalogue/api/v1alpha1"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	k8sclient "sigs.k8s.io/controller-runtime/pkg/client"
)

var _ = Describe("InstalledFeature dependent feature handling", func() {
	Context("When the dependent feature cant be loaded", func() {
		It("should reque the request when the dependent feature is not 'not found'", func() {
			var other *InstalledFeature

			By("reconciling the dependent feature", func() {
				other = createIFT(otherName, namespace, version, provider, description, uri, true, false)
				other.Status.DependingFeatures = []InstalledFeatureRef{
					{Namespace: namespace, Name: name},
				}

				client.EXPECT().LoadInstalledFeature(ctx, otherLookupKey).Return(other, nil)
			})

			By("working on the depending feature", func() {
				client.EXPECT().LoadInstalledFeature(ctx, iftLookupKey).Return(nil, errors.New("dependent feature not found"))
			})

			result, err := sut.Reconcile(otherReconcileRequest)

			Expect(result).Should(Equal(errorResult))
			Expect(err).Should(HaveOccurred())
		})

		It("should not requeue when dependent features is 'not found'", func() {
			var other *InstalledFeature

			By("reconciling the dependent feature", func() {
				other = createIFT(otherName, namespace, version, provider, description, uri, true, false)
				other.Status.DependingFeatures = []InstalledFeatureRef{
					{Namespace: namespace, Name: name},
				}

				client.EXPECT().LoadInstalledFeature(ctx, otherLookupKey).Return(other, nil)

				client.EXPECT().GetInstalledFeaturePatchBase(other).Return(k8sclient.MergeFrom(other))
				client.EXPECT().PatchInstalledFeatureStatus(ctx, other, k8sclient.MergeFrom(other)).Return(nil)
			})

			By("working on the depending feature", func() {
				client.EXPECT().LoadInstalledFeature(ctx, iftLookupKey).Return(nil, createNotFound("InstalledFeature", name))
			})

			result, err := sut.Reconcile(otherReconcileRequest)

			Expect(result).Should(Equal(successResult))
			Expect(err).ShouldNot(HaveOccurred())
		})
	})

	Context("When removing a dependent InstalledFeature", func() {
		It("", func() {

		})
	})
})
