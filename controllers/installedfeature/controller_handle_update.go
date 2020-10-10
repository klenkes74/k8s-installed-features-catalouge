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
	ctrl "sigs.k8s.io/controller-runtime"
)

func (r *Reconciler) handleUpdate(ctx context.Context, instance *featuresv1alpha1.InstalledFeature, eventReason string, reqLogger logr.Logger, changed bool) (ctrl.Result, error) {
	if changed {
		err := r.Client.SaveInstalledFeature(ctx, instance)
		if err != nil {
			reqLogger.Error(err, "could not rewrite the installedfeature")
			r.Client.WarnEvent(instance, eventReason,
				NoteFeatureSaveFailed,
				instance.GetNamespace(),
				instance.GetName(),
				err.Error(),
			)

			return errorRequeue, err
		}
	}

	statusChanged := false
	status := r.Client.GetInstalledFeaturePatchBase(instance)

	if len(instance.Status.MissingDependencies) > 0 {
		instance.Status.Phase = "pending"
		instance.Status.Message = "dependencies are missing"
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

			r.Client.WarnEvent(instance, eventReason,
				NoteStatusUpdateFailed,
				instance.GetNamespace(),
				instance.GetName(),
				err.Error(),
			)

			return errorRequeue, err
		}
	}

	r.Client.InfoEvent(instance, eventReason, NoteChangedFeature, instance.GetNamespace(), instance.GetName())
	return ctrl.Result{}, nil
}
