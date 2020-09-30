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
	reqLogger := r.Log.WithValues("installedfeaturegroup", req.NamespacedName)
	reqLogger.Info("working on", "ctx", ctx)

	changed := false

	instance, err := r.Client.LoadInstalledFeatureGroup(ctx, req.NamespacedName)
	if err != nil {
		if errors.IsNotFound(err) {
			return ctrl.Result{Requeue: false}, err
		}

		return ctrl.Result{Requeue: true, RequeueAfter: 10}, err
	}

	changed = r.handleFinalizer(instance, changed)
	return r.handleUpdate(changed, reqLogger, ctx, instance)
}

func (r *Reconciler) handleFinalizer(instance *featuresv1alpha1.InstalledFeatureGroup, changed bool) bool {
	if !controllerutil.ContainsFinalizer(instance, FinalizerName) && instance.DeletionTimestamp == nil {
		controllerutil.AddFinalizer(instance, FinalizerName)

		changed = true
	} else if controllerutil.ContainsFinalizer(instance, FinalizerName) && instance.DeletionTimestamp != nil {
		controllerutil.RemoveFinalizer(instance, FinalizerName)

		changed = true
	}
	return changed
}

func (r *Reconciler) handleUpdate(changed bool, reqLogger logr.Logger, ctx context.Context, instance *featuresv1alpha1.InstalledFeatureGroup) (ctrl.Result, error) {
	if changed {
		reqLogger.Info("rewriting the installedfeaturegroup")

		err := r.Client.SaveInstalledFeatureGroup(ctx, instance)
		if err != nil {
			r.Log.Error(err, "could not rewrite the installedfeaturegroup")

			return ctrl.Result{Requeue: true, RequeueAfter: 10}, err
		}
	}

	if instance.Status.Phase == "" {
		err := r.modifyStatus(ctx, instance, "provisioned", "ok")
		if err != nil {
			reqLogger.Error(err, "could not set the status to the installedfeature")

			return ctrl.Result{RequeueAfter: 10, Requeue: true}, err
		}
	}

	return ctrl.Result{}, nil
}

func (r *Reconciler) modifyStatus(ctx context.Context, instance *featuresv1alpha1.InstalledFeatureGroup, phase string, message string) error {
	status := r.Client.GetInstalledFeatureGroupPatchBase(instance)
	instance.Status.Phase = phase
	instance.Status.Message = message
	return r.Client.PatchInstalledFeatureGroupStatus(ctx, instance, status)
}
