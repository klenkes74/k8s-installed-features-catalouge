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
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
)

func (r *Reconciler) handleDependent(ctx context.Context, instance *featuresv1alpha1.InstalledFeature, reqLogger logr.Logger, changed bool) (bool, error) {
	if instance.Status.DependingFeatures == nil || len(instance.Status.DependingFeatures) == 0 {
		return changed, nil
	}

	reqLogger.Info("handling dependent features")

	missingDependent := make([]featuresv1alpha1.InstalledFeatureRef, 0)
	for _, feature := range instance.Status.DependingFeatures {
		ift, err := r.Client.LoadInstalledFeature(ctx, types.NamespacedName{Namespace: feature.Namespace, Name: feature.Name})
		if err != nil {
			if errors.IsNotFound(err) {
				reqLogger.Info("dependent feature is not found - don't need to change it.", "feature", feature)
				continue
			}

			reqLogger.Info("dependent feature can not be loaded.", "dependent-feature", feature)
			missingDependent = append(missingDependent, feature)

			continue
		}

		if ift.DeletionTimestamp != nil {
			continue
		}

		iftStatus := r.Client.GetInstalledFeaturePatchBase(ift)
		if instance.DeletionTimestamp == nil {
			for _, dep := range ift.Status.MissingDependencies {
				if dep.Namespace == instance.Namespace && dep.Name == instance.Name {
					r.removeMissingDependencyStatus(ift, featuresv1alpha1.InstalledFeatureRef{
						Namespace: instance.Namespace,
						Name:      instance.Name,
					}, reqLogger)
				}
			}
		} else {
			r.markDependencyAsMissing(
				ift,
				featuresv1alpha1.InstalledFeatureRef{
					Namespace: instance.Namespace,
					Name:      instance.Name,
				},
				reqLogger,
			)
		}

		err = r.Client.PatchInstalledFeatureStatus(ctx, ift, iftStatus)
		if err != nil {
			return changed, err
		}

	}

	if len(missingDependent) > 0 {
		return changed, fmt.Errorf("could not update dependent features: %v", missingDependent)
	}

	return changed, nil
}
