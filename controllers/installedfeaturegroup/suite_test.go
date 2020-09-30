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
	"github.com/golang/mock/gomock"
	featuresv1alpha1 "github.com/klenkes74/k8s-installed-features-catalogue/api/v1alpha1"
	"github.com/klenkes74/k8s-installed-features-catalogue/controllers/installedfeaturegroup"
	"github.com/klenkes74/k8s-installed-features-catalogue/generated"
	"k8s.io/apimachinery/pkg/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"sigs.k8s.io/controller-runtime/pkg/envtest/printer"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
	// +kubebuilder:scaffold:imports
)

// These tests use Ginkgo (BDD-style Go testing framework). Refer to
// http://onsi.github.io/ginkgo/ to learn more about Ginkgo.

var (
	ctrlMock *gomock.Controller

	client *generated.MockOcpClient
	sut    installedfeaturegroup.Reconciler
)

func TestInstalledFeatureController(t *testing.T) {
	RegisterFailHandler(Fail)

	RunSpecsWithDefaultAndCustomReporters(t,
		"InstalledFeatureGroup Controller Suite",
		[]Reporter{printer.NewlineReporter{}})
}

var _ = BeforeSuite(func() {
	logf.SetLogger(zap.New(zap.UseDevMode(true), zap.WriteTo(GinkgoWriter)))

	scheme := runtime.NewScheme()
	Expect(clientgoscheme.AddToScheme(scheme)).Should(Succeed())
	Expect(featuresv1alpha1.AddToScheme(scheme)).Should(Succeed())

	ctrlMock = gomock.NewController(GinkgoT())
	client = generated.NewMockOcpClient(ctrlMock)

	sut = installedfeaturegroup.Reconciler{
		Client: client,
		Log:    logf.Log,
		Scheme: scheme,
	}
})

var _ = AfterSuite(func() {
	By("tearing down the mock controller")
	ctrlMock.Finish()
})
