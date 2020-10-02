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
	"fmt"
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

		return ctrl.Result{RequeueAfter: 60}, err
	}

	changed, err = r.handleDependingOn(ctx, instance, reqLogger, changed)
	if err != nil {
		return ctrl.Result{RequeueAfter: 60}, err
	}

	changed, err = r.handleGroupEntry(ctx, instance, reqLogger, changed)
	if err != nil {
		return ctrl.Result{RequeueAfter: 60}, err
	}

	changed = r.handleFinalizer(ctx, instance, reqLogger, changed)

	return r.handleUpdate(ctx, instance, reqLogger, changed)
}

func (r *Reconciler) handleDependingOn(ctx context.Context, instance *featuresv1alpha1.InstalledFeature, reqLogger logr.Logger, changed bool) (bool, error) {
	if instance.Spec.DependsOn == nil || len(instance.Spec.DependsOn) == 0 {
		return changed, nil
	}

	reqLogger.Info("handling dependencies")

	status := r.Client.GetInstalledFeaturePatchBase(instance)

	missingDependencies := make([]featuresv1alpha1.InstalledFeatureRef, 0)
	for _, dependency := range instance.Spec.DependsOn {
		locator := types.NamespacedName{
			Namespace: dependency.Feature.Namespace,
			Name:      dependency.Feature.Name,
		}

		ift, err := r.Client.LoadInstalledFeature(ctx, locator)
		if err != nil || ift.DeletionTimestamp != nil {
			r.markDependencyAsMissing(instance, dependency, reqLogger)
			missingDependencies = append(missingDependencies, dependency.Feature)
			continue // next dependency
		}

		reqLogger.Info("working on dependency", "dependency", dependency.Feature)

		if ift.Status.DependingFeatures == nil && instance.DeletionTimestamp == nil {
			ift.Status.DependingFeatures = make([]featuresv1alpha1.InstalledFeatureRef, 0)
		}
		iftStatus := r.Client.GetInstalledFeaturePatchBase(ift)

		dependencyChanged := true
		alreadyRegistered := false
		for i, ft := range ift.Status.DependingFeatures {
			if ft.Namespace == instance.Namespace && ft.Name == instance.Name {
				alreadyRegistered = true

				if instance.DeletionTimestamp != nil {
					reqLogger.Info("instance is deleted - remove registered depending feature", "feature", ift.Name)

					ift.Status.DependingFeatures[i] = ift.Status.DependingFeatures[len(ift.Status.DependingFeatures)-1]
					// We do not need to put s[i] at the end, as it will be discarded anyway
					ift.Status.DependingFeatures = ift.Status.DependingFeatures[:len(ift.Status.DependingFeatures)-1]
				} else {
					reqLogger.Info("feature already registered as depending feature", "feature", ift.Name)
					dependencyChanged = false
				}
				break
			}
		}

		if !alreadyRegistered && instance.DeletionTimestamp == nil {
			ift.Status.DependingFeatures = append(ift.Status.DependingFeatures, featuresv1alpha1.InstalledFeatureRef{
				Namespace: instance.Namespace,
				Name:      instance.Name,
			})
		}

		if dependencyChanged {
			err = r.Client.PatchInstalledFeatureStatus(ctx, ift, iftStatus)
			if err != nil {
				reqLogger.Info("can not update entry with dependency information", "feature", ift)

				return changed, err
			}
		}
	}

	err := r.Client.PatchInstalledFeatureStatus(ctx, instance, status)
	if err != nil {
		reqLogger.Info("dependency status could not be set.")

		return changed, err
	}

	if len(missingDependencies) > 0 {
		return changed, fmt.Errorf("missing dependencies: %v", missingDependencies)
	}

	reqLogger.Info("added the dependency to status")
	return changed, nil
}

func (r *Reconciler) markDependencyAsMissing(instance *featuresv1alpha1.InstalledFeature, dependency featuresv1alpha1.InstalledFeatureDependency, reqLogger logr.Logger) {
	reqLogger.Info("can not load feature we depend on", "feature", dependency.Feature)

	if instance.Status.MissingDependencies == nil {
		instance.Status.MissingDependencies = make([]featuresv1alpha1.InstalledFeatureDependency, 0)
	}

	instance.Status.MissingDependencies = append(instance.Status.MissingDependencies, dependency)
}

func (r *Reconciler) handleGroupEntry(ctx context.Context, instance *featuresv1alpha1.InstalledFeature, reqLogger logr.Logger, changed bool) (bool, error) {
	if instance.Spec.Group == nil {
		return changed, nil
	}

	log := reqLogger.WithValues("group", instance.Spec.Group)

	log.Info("handling group entry")

	group, err := r.Client.LoadInstalledFeatureGroup(ctx, types.NamespacedName{
		Namespace: instance.Spec.Group.Namespace,
		Name:      instance.Spec.Group.Name,
	})
	if err != nil {
		log.Info("could not load group - will not update the group information")

		return changed, err
	}

	patch := r.Client.GetInstalledFeatureGroupPatchBase(group)
	if len(group.Status.Features) > 0 {
		for i, feature := range group.Status.Features {
			if feature.Name == instance.Name && feature.Namespace == instance.Namespace {
				if instance.DeletionTimestamp == nil {
					log.Info("feature already listed in feature group")
					return changed, nil
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

	return changed, r.Client.PatchInstalledFeatureGroupStatus(ctx, group, patch)
}

func (r *Reconciler) removeFromGroup(features []featuresv1alpha1.InstalledFeatureGroupListedFeature, pos int) []featuresv1alpha1.InstalledFeatureGroupListedFeature {
	features[len(features)-1], features[pos] = features[pos], features[len(features)-1]
	return features[:len(features)-1]
}

func (r *Reconciler) handleFinalizer(_ context.Context, instance *featuresv1alpha1.InstalledFeature, reqLogger logr.Logger, changed bool) bool {
	reqLogger.Info("handling finalizer")

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
		err := r.Client.SaveInstalledFeature(ctx, instance)
		if err != nil {
			reqLogger.Error(err, "could not rewrite the installedfeature")

			return ctrl.Result{RequeueAfter: 60}, err
		}
	}

	statusChanged := false
	status := r.Client.GetInstalledFeaturePatchBase(instance)

	if len(instance.Status.MissingDependencies) > 0 {
		instance.Status.Phase = "pending"
		instance.Status.Message = "dependencies are missing"
		statusChanged = true
	} else if len(instance.Status.ConflictingFeatures) > 0 {
		instance.Status.Phase = "pending"
		instance.Status.Message = "there are conflicting features"
		statusChanged = true
	} else if instance.Status.Phase != "provisioned" {
		instance.Status.Phase = "provisioned"
		instance.Status.Message = ""
		statusChanged = true
	}

	if statusChanged {
		err := r.Client.PatchInstalledFeatureStatus(ctx, instance, status)
		if err != nil {
			reqLogger.Error(err, "could not set the status to the installedfeature")

			return ctrl.Result{RequeueAfter: 60}, err
		}
	}

	return ctrl.Result{}, nil
}
