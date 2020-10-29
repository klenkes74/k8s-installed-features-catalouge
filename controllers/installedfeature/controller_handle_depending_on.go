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
	"github.com/go-logr/logr"
	featuresv1alpha1 "github.com/klenkes74/k8s-installed-features-catalogue/api/v1alpha1"
	"k8s.io/apimachinery/pkg/types"
)

func (r *Reconciler) handleDependingOn(ctx context.Context, instance *featuresv1alpha1.InstalledFeature, eventReason string, reqLogger logr.Logger, changed bool) (bool, error) {
	if instance.Spec.DependsOn == nil || len(instance.Spec.DependsOn) == 0 {
		return changed, nil
	}

	reqLogger.Info("handling dependencies")

	for _, dependency := range instance.Spec.DependsOn {
		locator := types.NamespacedName{
			Namespace: dependency.Namespace,
			Name:      dependency.Name,
		}

		ift, err := r.Client.LoadInstalledFeature(ctx, locator)
		if err != nil || ift.DeletionTimestamp != nil {
			r.markDependencyAsMissing(instance, dependency, reqLogger)

			changed = true
			continue // next dependency
		}

		if instance.DeletionTimestamp == nil {
			err = r.registerInstanceAsDependingOn(ctx, instance, ift, reqLogger)
			if err != nil {
				reqLogger.Info("can not change state of dependency: ", "error", err.Error())

				return changed, err
			}
		}

		err = r.Client.ReconcileFeature(ctx, ift)
		if err != nil {
			reqLogger.Info("can not reconcile feature",
				"feature", ift,
			)

			return changed, fmt.Errorf("can not start reconcilation of dependency: %s", err.Error())
		}
	}

	if len(instance.Status.MissingDependencies) > 0 {
		r.Client.WarnEvent(instance, eventReason, NoteMissingDependencies, instance.Status.MissingDependencies)
		return changed, fmt.Errorf("missing dependencies: %v", instance.Status.MissingDependencies)
	}

	return changed, nil
}

func (r *Reconciler) markDependencyAsMissing(
	instance *featuresv1alpha1.InstalledFeature,
	dependency featuresv1alpha1.InstalledFeatureRef,
	reqLogger logr.Logger,
) bool {
	result := false

	reqLogger.Info("adding the missing feature to missing feature list", "feature", dependency)
	instance.Status.MissingDependencies, result = r.addRef(instance.Status.MissingDependencies, dependency)

	return result
}

func (r *Reconciler) registerInstanceAsDependingOn(
	ctx context.Context,
	instance, dependency *featuresv1alpha1.InstalledFeature,
	reqLogger logr.Logger,
) error {
	if i := r.indexOfRef(dependency.Status.DependingFeatures, r.generateRef(instance)); i != -1 {
		reqLogger.Info("dependency already listed in feature",
			"feature", dependency,
			"dependency", instance,
		)
		return nil
	}

	reqLogger.Info("adding dependency to feature",
		"feature", dependency,
		"dependency", instance,
	)

	patch := r.Client.GetInstalledFeaturePatchBase(dependency)
	dependency.Status.DependingFeatures, _ = r.addRef(dependency.Status.DependingFeatures, r.generateRef(instance))
	return r.Client.PatchInstalledFeatureStatus(ctx, dependency, patch)
}
