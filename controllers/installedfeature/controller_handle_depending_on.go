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

	status := r.Client.GetInstalledFeaturePatchBase(instance)

	missingDependencies := make([]featuresv1alpha1.InstalledFeatureRef, 0)
	for _, dependency := range instance.Spec.DependsOn {
		locator := types.NamespacedName{
			Namespace: dependency.Namespace,
			Name:      dependency.Name,
		}

		ift, err := r.Client.LoadInstalledFeature(ctx, locator)
		if err != nil || ift.DeletionTimestamp != nil {
			r.markDependencyAsMissing(instance, dependency, reqLogger)
			missingDependencies = append(missingDependencies, dependency)
			continue // next dependency
		}

		reqLogger.Info("working on dependency", "dependency", dependency)

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
		r.Client.WarnEvent(instance, eventReason, NoteMissingDependencies, missingDependencies)
		return changed, fmt.Errorf("missing dependencies: %v", missingDependencies)
	}

	reqLogger.Info("added the dependency to status")
	return changed, nil
}
