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
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"time"
)

const (
	// ControllerName is the name of the controller used internally
	ControllerName = "installedfeature-controller"

	// FinalizerName is the name added to the finalizer of the managed objects.
	FinalizerName = "features.kaiserpfalz-edv.de/" + ControllerName

	// RequeueTime is the default requeuing time when the operator is running in problems
	RequeueTime = 60 * time.Second

	// NoteChangedFeature is the event text for a successful update/creation/deletion
	NoteChangedFeature = "Changed feature %s/%s"
	// NoteUpdatingDepenciesFailed is the event text when dependencies could not be linked to this feature
	NoteUpdatingDependenciesFailed = "Could not update the dependencies of %s"
	// NoteUpdatingDependentFeaturesFailed is the event text when updating all dependent features failed
	NoteUpdatingDependentFeaturesFailed = "Could not handle the dependent features of %s"
	// NoteUpdatingGroupFailed is the event text when updating the group relation failed
	NoteUpdatingGroupFailed = "Could not handle the group relation of %s"
	// NoteMissingDependencies is the event text listing the missing dependencies
	NoteMissingDependencies = "Feature has missing dependencies: %v"
	// NoteStatusUpdateFailed is the event text when the status could not be updated
	NoteStatusUpdateFailed = "Could not save the status of feature %s/%s: %s"
	// NoteFeatureSaveFailed is the event text when saving a feature is failing
	NoteFeatureSaveFailed = "Could not save the feature %s/%s: %s"
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
// +kubebuilder:rbac:groups="",resources=events,verbs=create;patch

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

		reqLogger.Error(err, "Could not load installed feature")

		return errorRequeue, err
	}

	eventReason := r.calculateEventReason(instance)

	changed, err = r.handleDependingOn(ctx, instance, eventReason, reqLogger, changed)
	if err != nil {
		r.Client.WarnEvent(instance, eventReason, NoteUpdatingDependenciesFailed, req.NamespacedName)
		return errorRequeue, err
	}

	changed, err = r.handleDependent(ctx, instance, reqLogger, changed)
	if err != nil {
		r.Client.WarnEvent(instance, eventReason, NoteUpdatingDependentFeaturesFailed, req.NamespacedName)
		return errorRequeue, err
	}

	changed, err = r.handleGroupEntry(ctx, instance, reqLogger, changed)
	if err != nil {
		r.Client.WarnEvent(instance, eventReason, NoteUpdatingGroupFailed, req.NamespacedName)

		return errorRequeue, err
	}

	changed = r.handleFinalizer(ctx, instance, reqLogger, changed)

	return r.handleUpdate(ctx, instance, eventReason, reqLogger, changed)
}

func (r *Reconciler) calculateEventReason(instance *featuresv1alpha1.InstalledFeature) string {
	if instance.DeletionTimestamp != nil {
		return "Delete"
	} else if !controllerutil.ContainsFinalizer(instance, FinalizerName) {
		return "Create"
	}

	return "Update"
}

// addRef adds the reference and returns the new array. The second return value is true, if the array has been changed
// and false if the array is unchanged.
func (r *Reconciler) addRef(
	refs []featuresv1alpha1.InstalledFeatureRef,
	dependency featuresv1alpha1.InstalledFeatureRef,
) ([]featuresv1alpha1.InstalledFeatureRef, bool) {
	i := r.indexOfRef(refs, dependency)
	if i != -1 {
		return refs, false
	}

	if refs == nil {
		refs = []featuresv1alpha1.InstalledFeatureRef{}
	}

	return append(refs, dependency), true
}

// removeRef removes the reference and returns the new array. The second return value is true, if the array has been
// changed and false if the array is unchanged.
func (r *Reconciler) removeRef(
	refs []featuresv1alpha1.InstalledFeatureRef,
	dependency featuresv1alpha1.InstalledFeatureRef,
) ([]featuresv1alpha1.InstalledFeatureRef, bool) {
	i := r.indexOfRef(refs, dependency)
	if i == -1 {
		return refs, false
	}

	if len(refs) > 1 {
		refs[i] = refs[len(refs)-1]
		refs = refs[:len(refs)-1]
	} else {
		refs = nil
	}

	return refs, true
}

// indexOfRef returns the index of the reference in the array or -1 if the array is empty, nil or the reference is not
// listed in the array
func (r *Reconciler) indexOfRef(
	refs []featuresv1alpha1.InstalledFeatureRef,
	dependency featuresv1alpha1.InstalledFeatureRef,
) int {
	if refs == nil {
		return -1
	}

	for i, d := range refs {
		if d.Name == dependency.Name && d.Namespace == dependency.Namespace {
			return i
		}
	}

	return -1
}

// generateRef generates a reference pointing to the instance
func (r *Reconciler) generateRef(instance *featuresv1alpha1.InstalledFeature) featuresv1alpha1.InstalledFeatureRef {
	return featuresv1alpha1.InstalledFeatureRef{
		Namespace: instance.Namespace,
		Name:      instance.Name,
	}
}
