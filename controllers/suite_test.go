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

package controllers_test

import (
	"path/filepath"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/envtest"
	"sigs.k8s.io/controller-runtime/pkg/envtest/printer"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"

	featuresv1alpha1 "github.com/klenkes74/k8s-installed-features-catalogue/api/v1alpha1"
	// +kubebuilder:scaffold:imports
)

// These tests use Ginkgo (BDD-style Go testing framework). Refer to
// http://onsi.github.io/ginkgo/ to learn more about Ginkgo.

var cfg *rest.Config
var k8sClient client.Client
var testEnv *envtest.Environment

func TestAPIs(t *testing.T) {
	RegisterFailHandler(Fail)

	RunSpecsWithDefaultAndCustomReporters(t,
		"Controller Suite",
		[]Reporter{printer.NewlineReporter{}})
}

var _ = BeforeSuite(func(done Done) {
	logf.SetLogger(zap.New(zap.UseDevMode(true), zap.WriteTo(GinkgoWriter)))

	By("bootstrapping test environment")
	testEnv = &envtest.Environment{
		CRDDirectoryPaths: []string{filepath.Join("..", "config", "crd", "bases")},
	}

	var err error
	cfg, err = testEnv.Start()
	Expect(err).ToNot(HaveOccurred())
	Expect(cfg).ToNot(BeNil())

	err = featuresv1alpha1.AddToScheme(scheme.Scheme)
	Expect(err).NotTo(HaveOccurred())

	// +kubebuilder:scaffold:scheme

	k8sClient, err = client.New(cfg, client.Options{Scheme: scheme.Scheme})
	Expect(err).ToNot(HaveOccurred())
	Expect(k8sClient).ToNot(BeNil())

	close(done)
}, 60)

var _ = AfterSuite(func() {
	By("tearing down the test environment")
	err := testEnv.Stop()
	Expect(err).ToNot(HaveOccurred())
})

var _ = Describe("Creating a new InstalledFeature", func() {
	It("should be created when there are no conflicting features installed and all dependencies met", func() {
		// TODO 2020-09-26 klenkes74 Implement this test
	})

	It("should be created with failing state when there are conflicting features installed", func() {
		// TODO 2020-09-26 klenkes74 Implement this test
	})

	It("should be created with failing state when there are dependencies missing completely", func() {
		// TODO 2020-09-26 klenkes74 Implement this test
	})

	It("should be created with failing state when the dependency version is too low", func() {
		// TODO 2020-09-26 klenkes74 Implement this test
	})

	It("should be created with failing state when the dependency version is too high", func() {
		// TODO 2020-09-26 klenkes74 Implement this test
	})

	It("should not be created when the same feature is already installed", func() {
		// TODO 2020-09-26 klenkes74 Implement this test
	})
})

var _ = Describe("Delete an existing InstalledFeature", func() {
	It("should be deleted when there are no dependencies on the removed feature", func() {
		// TODO 2020-09-26 klenkes74 Implement this test
	})

	It("should not be deleted when there are dependencies on the removed feature", func() {
		// TODO 2020-09-26 klenkes74 Implement this test
	})
})

var _ = Describe("Technical handling", func() {
	It("should ignore failures to load a resource when the reason is 'NotFound'", func() {
		// TODO 2020-09-26 klenkes74 Implement this test
	})

	It("should requeue the request when the resource can not be loaded for any reason but 'NotFound'", func() {
		// TODO 2020-09-26 klenkes74 Implement this test
	})

	It("should requeue the request when a changed resource fails to update", func() {
		// TODO 2020-09-26 klenkes74 Implement this test
	})

	It("should add the finalizer when the finalizer is not set", func() {
		// TODO 2020-09-26 klenkes74 Implement this test
	})
})
