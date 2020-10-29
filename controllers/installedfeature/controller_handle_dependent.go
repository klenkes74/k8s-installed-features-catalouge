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

func (r *Reconciler) handleDependent(
	ctx context.Context,
	instance *featuresv1alpha1.InstalledFeature,
	reqLogger logr.Logger,
	changed bool,
) (bool, error) {
	if instance.Status.DependingFeatures == nil || len(instance.Status.DependingFeatures) == 0 {
		return changed, nil
	}

	reqLogger.Info("handling dependent features")

	patch := r.Client.GetInstalledFeaturePatchBase(instance)

	changedStatus := false
	removedFeatures := []featuresv1alpha1.InstalledFeatureRef{}
	for _, feature := range instance.Status.DependingFeatures {
		ift, err := r.Client.LoadInstalledFeature(ctx,
			types.NamespacedName{Namespace: feature.Namespace, Name: feature.Name},
		)

		if err != nil {
			if errors.IsNotFound(err) {
				removedFeatures, _ = r.addRef(removedFeatures, feature)
				changedStatus = true
				continue
			}

			reqLogger.Info("dependent feature can not be loaded.", "dependent-feature", feature)

			continue
		}

		if instance.DeletionTimestamp != nil && ift.DeletionTimestamp != nil {
			err = r.Client.ReconcileFeature(ctx, ift)
			if err != nil {
				reqLogger.Info("can not reconcile feature",
					"feature", ift,
				)

				return changed, fmt.Errorf("can not start reconcilation of depending feature: %s", err.Error())
			}
		}
	}

	if changedStatus {
		for _, feature := range removedFeatures {
			instance.Status.DependingFeatures, _ = r.removeRef(instance.Status.DependingFeatures, feature)
		}
		err := r.Client.PatchInstalledFeatureStatus(ctx, instance, patch)
		if err != nil {
			return changed, fmt.Errorf("could not update dependent features: %s", err.Error())
		}
	}

	return changed, nil
}
