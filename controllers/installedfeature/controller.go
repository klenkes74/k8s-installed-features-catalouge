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

package installedfeature

import (
	"context"
	"github.com/go-logr/logr"
	featuresv1alpha1 "github.com/klenkes74/k8s-installed-features-catalogue/api/v1alpha1"
	"github.com/klenkes74/k8s-installed-features-catalogue/controllers"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"time"
)

const (
	// FinalizerName is the name added to the finalizer of the managed objects.
	FinalizerName = "features.kaiserpfalz-edv.de/installedfeature-controller"

	// RequeueTime is the default requeuing time when the operator is running in problems
	RequeueTime = 60 * time.Second
)

// The default requeue for error handling
var errorRequeue = ctrl.Result{RequeueAfter: RequeueTime}

// Reconciler reconciles a InstalledFeature object
type Reconciler struct {
	Client controllers.OcpClient

	Log    logr.Logger
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=features.kaiserpfalz-edv.de,resources=installedfeatures,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=features.kaiserpfalz-edv.de,resources=installedfeatures/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=features.kaiserpfalz-edv.de,resources=installedfeaturegroups/status,verbs=get;update;patch

func (r *Reconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&featuresv1alpha1.InstalledFeature{}).
		Complete(r)
}

func (r *Reconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	ctx := context.Background()
	reqLogger := r.Log.WithValues("installedfeature", req.NamespacedName)
	reqLogger.Info("working on", "ctx", ctx)

	changed := false

	instance, err := r.Client.LoadInstalledFeature(ctx, req.NamespacedName)
	if err != nil {
		if errors.IsNotFound(err) {
			return ctrl.Result{}, err
		}

		return errorRequeue, err
	}

	changed, err = r.handleDependingOn(ctx, instance, reqLogger, changed)
	if err != nil {
		return errorRequeue, err
	}

	changed, err = r.handleDependent(ctx, instance, reqLogger, changed)
	if err != nil {
		return errorRequeue, err
	}

	changed, err = r.handleGroupEntry(ctx, instance, reqLogger, changed)
	if err != nil {
		return errorRequeue, err
	}

	changed = r.handleFinalizer(ctx, instance, reqLogger, changed)

	return r.handleUpdate(ctx, instance, reqLogger, changed)
}

func (r *Reconciler) markDependencyAsMissing(instance *featuresv1alpha1.InstalledFeature, dependency featuresv1alpha1.InstalledFeatureRef, reqLogger logr.Logger) {
	reqLogger.Info("mark the missing feature", "feature", dependency)

	if instance.Status.MissingDependencies == nil {
		instance.Status.MissingDependencies = make([]featuresv1alpha1.InstalledFeatureRef, 0)
	}

	instance.Status.MissingDependencies = append(instance.Status.MissingDependencies, dependency)
}

func (r *Reconciler) removeMissingDependencyStatus(instance *featuresv1alpha1.InstalledFeature, dependency featuresv1alpha1.InstalledFeatureRef, reqLogger logr.Logger) bool {
	i := r.indexOfMissingDependency(instance, dependency)
	if i == -1 {
		return false
	}

	reqLogger.Info("remove the marked missing feature", "feature", dependency)

	instance.Status.MissingDependencies[i] = instance.Status.MissingDependencies[len(instance.Status.MissingDependencies)-1]
	instance.Status.MissingDependencies = instance.Status.MissingDependencies[:len(instance.Status.MissingDependencies)-1]

	return true
}

func (r *Reconciler) indexOfMissingDependency(instance *featuresv1alpha1.InstalledFeature, dependency featuresv1alpha1.InstalledFeatureRef) int {
	if instance.Status.MissingDependencies == nil {
		return -1
	}

	for i, d := range instance.Status.MissingDependencies {
		if d.Name == dependency.Name && d.Namespace == dependency.Namespace {
			return i
		}
	}

	return -1
}
