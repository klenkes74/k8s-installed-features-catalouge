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

//go:generate mockgen -copyright_file=../hack/boilerplate.go.txt -destination=../generated/mock_ocp_client.go -package=generated github.com/klenkes74/k8s-installed-features-catalogue/controllers OcpClient

package controllers

import (
	"context"
	"github.com/klenkes74/k8s-installed-features-catalogue/api/v1alpha1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type OcpClient interface {
	client.Client
	LoadInstalledFeature(ctx context.Context, lookup types.NamespacedName) (*v1alpha1.InstalledFeature, error)
	SaveInstalledFeature(ctx context.Context, instance *v1alpha1.InstalledFeature) error
	GetInstalledFeaturePatchBase(instance *v1alpha1.InstalledFeature) client.Patch
	PatchInstalledFeatureStatus(ctx context.Context, instance *v1alpha1.InstalledFeature, patch client.Patch) error

	LoadInstalledFeatureGroup(ctx context.Context, lookup types.NamespacedName) (*v1alpha1.InstalledFeatureGroup, error)
	SaveInstalledFeatureGroup(ctx context.Context, instance *v1alpha1.InstalledFeatureGroup) error
	GetInstalledFeatureGroupPatchBase(instance *v1alpha1.InstalledFeature) client.Patch
	PatchInstalledFeatureGroupStatus(ctx context.Context, instance *v1alpha1.InstalledFeatureGroup, patch client.Patch) error
}

var _ OcpClient = &OcpClientProd{}

type OcpClientProd struct {
	Client client.Client
}

func (o OcpClientProd) LoadInstalledFeature(ctx context.Context, lookup types.NamespacedName) (*v1alpha1.InstalledFeature, error) {
	instance := &v1alpha1.InstalledFeature{}

	err := o.Get(ctx, lookup, instance)
	if err != nil {
		return nil, err
	}

	return instance, nil
}

func (o OcpClientProd) SaveInstalledFeature(ctx context.Context, instance *v1alpha1.InstalledFeature) error {
	return o.Client.Update(ctx, instance)
}

func (o OcpClientProd) GetInstalledFeaturePatchBase(instance *v1alpha1.InstalledFeature) client.Patch {
	return client.MergeFrom(instance.DeepCopy())
}

func (o OcpClientProd) PatchInstalledFeatureStatus(ctx context.Context, instance *v1alpha1.InstalledFeature, patch client.Patch) error {
	return o.Client.Status().Patch(ctx, instance, patch)
}

func (o OcpClientProd) LoadInstalledFeatureGroup(ctx context.Context, lookup types.NamespacedName) (*v1alpha1.InstalledFeatureGroup, error) {
	instance := &v1alpha1.InstalledFeatureGroup{}

	err := o.Get(ctx, lookup, instance)
	if err != nil {
		return nil, err
	}

	return instance, nil
}

func (o OcpClientProd) SaveInstalledFeatureGroup(ctx context.Context, instance *v1alpha1.InstalledFeatureGroup) error {
	return o.Client.Update(ctx, instance)
}

func (o OcpClientProd) GetInstalledFeatureGroupPatchBase(instance *v1alpha1.InstalledFeature) client.Patch {
	return client.MergeFrom(instance.DeepCopy())
}

func (o OcpClientProd) PatchInstalledFeatureGroupStatus(ctx context.Context, instance *v1alpha1.InstalledFeatureGroup, patch client.Patch) error {
	return o.Client.Status().Patch(ctx, instance, patch)
}

func (o OcpClientProd) Get(ctx context.Context, key client.ObjectKey, obj runtime.Object) error {
	return o.Client.Get(ctx, key, obj)
}

func (o OcpClientProd) List(ctx context.Context, list runtime.Object, opts ...client.ListOption) error {
	return o.Client.List(ctx, list, opts...)
}

func (o OcpClientProd) Create(ctx context.Context, obj runtime.Object, opts ...client.CreateOption) error {
	return o.Client.Create(ctx, obj, opts...)
}

func (o OcpClientProd) Delete(ctx context.Context, obj runtime.Object, opts ...client.DeleteOption) error {
	return o.Client.Delete(ctx, obj, opts...)
}

func (o OcpClientProd) Update(ctx context.Context, obj runtime.Object, opts ...client.UpdateOption) error {
	return o.Client.Update(ctx, obj, opts...)
}

func (o OcpClientProd) Patch(ctx context.Context, obj runtime.Object, patch client.Patch, opts ...client.PatchOption) error {
	return o.Client.Patch(ctx, obj, patch, opts...)
}

func (o OcpClientProd) DeleteAllOf(ctx context.Context, obj runtime.Object, opts ...client.DeleteAllOfOption) error {
	return o.Client.DeleteAllOf(ctx, obj, opts...)
}

func (o OcpClientProd) Status() client.StatusWriter {
	return o.Client.Status()
}
