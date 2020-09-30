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
	"github.com/klenkes74/k8s-installed-features-catalogue/controllers"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	"github.com/go-logr/logr"
	featuresv1alpha1 "github.com/klenkes74/k8s-installed-features-catalogue/api/v1alpha1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
)

const (
	FinalizerName = "features.kaiserpfalz-edv.de/installedfeature-controller"
)

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
			return ctrl.Result{Requeue: false}, err
		}

		return ctrl.Result{Requeue: true, RequeueAfter: 10}, err
	}

	changed = r.handleGroupEntry(ctx, instance, reqLogger, changed)
	changed = r.handleFinalizer(ctx, instance, reqLogger, changed)

	return r.handleUpdate(ctx, instance, reqLogger, changed)
}

func (r *Reconciler) handleGroupEntry(ctx context.Context, instance *featuresv1alpha1.InstalledFeature, reqLogger logr.Logger, changed bool) bool {
	if instance.Spec.Group == nil {
		return changed
	}

	log := reqLogger.WithValues("group", instance.Spec.Group)

	group, err := r.Client.LoadInstalledFeatureGroup(ctx, types.NamespacedName{
		Namespace: instance.Spec.Group.Namespace,
		Name:      instance.Spec.Group.Name,
	})
	if err != nil {
		log.Info("could not load group - will not update the group information")

		return changed
	}

	patch := r.Client.GetInstalledFeatureGroupPatchBase(group)
	if len(group.Status.Features) > 0 {
		for i, feature := range group.Status.Features {
			if feature.Name == instance.Name && feature.Namespace == instance.Namespace {
				if instance.DeletionTimestamp == nil {
					log.Info("feature already listed in feature group")
					return changed
				} else {
					log.Info("removing feature from feature group")

					r.removeFromGroup(group.Status.Features, i)
					break
				}
			}
		}

		if instance.DeletionTimestamp == nil {
			group.Status.Features = append(group.Status.Features, featuresv1alpha1.InstalledFeatureGroupListedFeature{
				Namespace: instance.Namespace,
				Name:      instance.Name,
			})
			log.Info("added feature to feature group")
		}
	} else {
		group.Status.Features = make([]featuresv1alpha1.InstalledFeatureGroupListedFeature, 1)
		group.Status.Features[0] = featuresv1alpha1.InstalledFeatureGroupListedFeature{
			Namespace: instance.Namespace,
			Name:      instance.Name,
		}
	}
	r.Client.PatchInstalledFeatureGroupStatus(ctx, group, patch)

	return changed
}

func (r *Reconciler) removeFromGroup(features []featuresv1alpha1.InstalledFeatureGroupListedFeature, pos int) []featuresv1alpha1.InstalledFeatureGroupListedFeature {
	features[len(features)-1], features[pos] = features[pos], features[len(features)-1]
	return features[:len(features)-1]
}

func (r *Reconciler) handleFinalizer(_ context.Context, instance *featuresv1alpha1.InstalledFeature, reqLogger logr.Logger, changed bool) bool {
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

func (r *Reconciler) handleUpdate(ctx context.Context, instance *featuresv1alpha1.InstalledFeature, reqLogger logr.Logger, changed bool) (ctrl.Result, error) {
	if changed {
		reqLogger.Info("rewriting the installedfeature")

		err := r.Client.SaveInstalledFeature(ctx, instance)
		if err != nil {
			reqLogger.Error(err, "could not rewrite the installedfeature")

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

func (r *Reconciler) modifyStatus(ctx context.Context, instance *featuresv1alpha1.InstalledFeature, phase string, message string) error {
	status := r.Client.GetInstalledFeaturePatchBase(instance)
	instance.Status.Phase = phase
	instance.Status.Message = message
	return r.Client.PatchInstalledFeatureStatus(ctx, instance, status)
}
