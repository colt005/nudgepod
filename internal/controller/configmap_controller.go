/*
Copyright 2025.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controller

import (
	"context"
	"time"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
)

// ConfigMapReconciler reconciles a ConfigMap object
type ConfigMapReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=core,resources=configmaps,verbs=get;list;watch
// +kubebuilder:rbac:groups=core,resources=configmaps/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=core,resources=configmaps/finalizers,verbs=update

// +kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch;patch;update
// +kubebuilder:rbac:groups=apps,resources=statefulsets,verbs=get;list;watch;patch;update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the ConfigMap object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.21.0/pkg/reconcile
func (r *ConfigMapReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	l := logf.FromContext(ctx)

	l.Info("Reconciling ConfigMap")

	var cm corev1.ConfigMap
	if err := r.Get(ctx, req.NamespacedName, &cm); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	// 1. List Deployments in this namespace
	var deployList appsv1.DeploymentList
	if err := r.List(ctx, &deployList, client.InNamespace(req.Namespace)); err != nil {
		return ctrl.Result{}, err
	}

	// 2. For each Deployment that references this ConfigMap, patch an annotation
	for _, d := range deployList.Items {
		if referencesConfigMap(&d.Spec.Template.Spec, cm.Name) {
			patch := client.MergeFrom(d.DeepCopy())
			if d.Spec.Template.Annotations == nil {
				d.Spec.Template.Annotations = map[string]string{}
			}
			d.Spec.Template.Annotations["configmap-reloader/updatedAt"] = time.Now().Format(time.RFC3339)
			if err := r.Patch(ctx, &d, patch); err != nil {
				l.Error(err, "failed to patch deployment", "deployment", d.Name)
				return ctrl.Result{}, err
			}
		}
	}

	// 3. Repeat for StatefulSets
	var stsList appsv1.StatefulSetList
	if err := r.List(ctx, &stsList, client.InNamespace(req.Namespace)); err != nil {
		return ctrl.Result{}, err
	}
	for _, s := range stsList.Items {
		if referencesConfigMap(&s.Spec.Template.Spec, cm.Name) {
			patch := client.MergeFrom(s.DeepCopy())
			if s.Spec.Template.Annotations == nil {
				s.Spec.Template.Annotations = map[string]string{}
			}
			s.Spec.Template.Annotations["configmap-reloader/updatedAt"] = time.Now().Format(time.RFC3339)
			if err := r.Patch(ctx, &s, patch); err != nil {
				l.Error(err, "failed to patch statefulset", "statefulset", s.Name)
				return ctrl.Result{}, err
			}
		}
	}

	return ctrl.Result{}, nil
}

// Helper to detect if a PodSpec references the given ConfigMap
func referencesConfigMap(podSpec *corev1.PodSpec, cmName string) bool {
	// Check volumes
	for _, vol := range podSpec.Volumes {
		if vol.ConfigMap != nil && vol.ConfigMap.Name == cmName {
			return true
		}
	}
	// Check env and envFrom in containers
	for _, c := range podSpec.Containers {
		for _, e := range c.Env {
			if e.ValueFrom != nil && e.ValueFrom.ConfigMapKeyRef != nil && e.ValueFrom.ConfigMapKeyRef.Name == cmName {
				return true
			}
		}
		for _, ef := range c.EnvFrom {
			if ef.ConfigMapRef != nil && ef.ConfigMapRef.Name == cmName {
				return true
			}
		}
	}
	return false
}

// SetupWithManager sets up the controller with the Manager.
func (r *ConfigMapReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&corev1.ConfigMap{}).
		Named("configmap").
		Complete(r)
}
