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

package installedfeaturegroup

import (
	"context"
	"github.com/klenkes74/k8s-installed-features-catalogue/controllers"
	"k8s.io/apimachinery/pkg/api/errors"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	"github.com/go-logr/logr"
	featuresv1alpha1 "github.com/klenkes74/k8s-installed-features-catalogue/api/v1alpha1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
)

const (
	FinalizerName = "features.kaiserpfalz-edv.de/installedfeature-controller"
)

// Reconciler reconciles a InstalledFeatureGroup object
type Reconciler struct {
	Client controllers.OcpClient

	Log    logr.Logger
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=features.kaiserpfalz-edv.de,resources=installedfeaturegroups,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=features.kaiserpfalz-edv.de,resources=installedfeaturegroups/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=features.kaiserpfalz-edv.de,resources=installedfeatures/status,verbs=get;update;patch

func (r *Reconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&featuresv1alpha1.InstalledFeatureGroup{}).
		Complete(r)
}

func (r *Reconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	ctx := context.Background()
	reqLogger := r.Log.WithValues("installed-feature-group", req.NamespacedName)
	reqLogger.Info("working on", "ctx", ctx)

	changed := false

	instance, err := r.Client.LoadInstalledFeatureGroup(ctx, req.NamespacedName)
	if err != nil {
		if errors.IsNotFound(err) {
			return ctrl.Result{Requeue: false}, err
		}

		return ctrl.Result{RequeueAfter: 60}, err
	}

	changed = r.handleFinalizer(ctx, instance, reqLogger, changed)

	return r.handleUpdate(ctx, instance, reqLogger, changed)
}

func (r *Reconciler) handleFinalizer(_ context.Context, instance *featuresv1alpha1.InstalledFeatureGroup, reqLogger logr.Logger, changed bool) bool {
	if !controllerutil.ContainsFinalizer(instance, FinalizerName) && instance.DeletionTimestamp == nil {
		reqLogger.Info("adding finalizer")

		controllerutil.AddFinalizer(instance, FinalizerName)

		changed = true
	} else if controllerutil.ContainsFinalizer(instance, FinalizerName) && instance.DeletionTimestamp != nil {
		reqLogger.Info("removing finalizer")

		controllerutil.RemoveFinalizer(instance, FinalizerName)

		changed = true
	}
	return changed
}

func (r *Reconciler) handleUpdate(ctx context.Context, instance *featuresv1alpha1.InstalledFeatureGroup, reqLogger logr.Logger, changed bool) (ctrl.Result, error) {
	if changed {
		reqLogger.Info("rewriting the InstalledFeatureGroup")

		err := r.Client.SaveInstalledFeatureGroup(ctx, instance)
		if err != nil {
			r.Log.Error(err, "could not rewrite the InstalledFeatureGroup")

			return ctrl.Result{RequeueAfter: 60}, err
		}
	}

	statusChanged := false
	status := r.Client.GetInstalledFeatureGroupPatchBase(instance)
	if instance.Status.Phase == "" {
		instance.Status.Phase = "provisioned"
		instance.Status.Message = ""
		statusChanged = true
	}

	if statusChanged {
		err := r.Client.PatchInstalledFeatureGroupStatus(ctx, instance, status)
		if err != nil {
			return ctrl.Result{RequeueAfter: 60}, err
		}
	}

	return ctrl.Result{}, nil
}
