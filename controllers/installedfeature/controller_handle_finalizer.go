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
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	"github.com/go-logr/logr"
	featuresv1alpha1 "github.com/klenkes74/k8s-installed-features-catalogue/api/v1alpha1"
)

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
