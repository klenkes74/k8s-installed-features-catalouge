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
	"k8s.io/apimachinery/pkg/types"
)

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
