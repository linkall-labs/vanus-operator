/*
Copyright 2022.

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

package controllers

import (
	"context"
	"time"

	cons "github.com/linkall-labs/vanus-operator/internal/constants"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/log"

	vanusv1alpha1 "github.com/linkall-labs/vanus-operator/api/v1alpha1"
)

// StoreReconciler reconciles a Store object
type StoreReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=vanus.linkall.com,resources=stores,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=vanus.linkall.com,resources=stores/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=vanus.linkall.com,resources=stores/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Store object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.12.2/pkg/reconcile
func (r *StoreReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = log.FromContext(ctx)

	logger := log.Log.WithName("Store")
	logger.Info("Reconciling Store.")

	// Fetch the Store instance
	store := &vanusv1alpha1.Store{}
	err := r.Get(ctx, req.NamespacedName, store)
	if err != nil {
		if errors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile req.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
			// Return and don't requeue
			logger.Info("Store resource not found. Ignoring since object must be deleted.")
			return ctrl.Result{}, nil
		}
		// Error reading the object - requeue the req.
		logger.Error(err, "Failed to get Store.")
		return ctrl.Result{RequeueAfter: time.Duration(cons.RequeueIntervalInSecond) * time.Second}, err
	}

	storeStatefulSet := r.getStatefulSetForStore(store)
	// Create Store StatefulSet
	// Check if the statefulSet already exists, if not create a new one
	sts := &appsv1.StatefulSet{}
	err = r.Get(ctx, types.NamespacedName{Name: storeStatefulSet.Name, Namespace: storeStatefulSet.Namespace}, sts)
	if err != nil {
		if errors.IsNotFound(err) {
			logger.Info("Creating a new Store StatefulSet.", "StatefulSet.Namespace", storeStatefulSet.Namespace, "StatefulSet.Name", storeStatefulSet.Name)
			err = r.Create(ctx, storeStatefulSet)
			if err != nil {
				logger.Error(err, "Failed to create new Store StatefulSet", "StatefulSet.Namespace", storeStatefulSet.Namespace, "StatefulSet.Name", storeStatefulSet.Name)
				return ctrl.Result{}, err
			} else {
				logger.Info("Successfully create Store StatefulSet")
			}
		} else {
			logger.Error(err, "Failed to get Store StatefulSet.")
			return ctrl.Result{RequeueAfter: time.Duration(cons.RequeueIntervalInSecond) * time.Second}, err
		}
	}

	// TODO(jiangkai): Update Store StatefulSet

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *StoreReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&vanusv1alpha1.Store{}).
		Complete(r)
}

// returns a Store StatefulSet object
func (r *StoreReconciler) getStatefulSetForStore(store *vanusv1alpha1.Store) *appsv1.StatefulSet {
	labels := labelsForStore(store.Name)
	annotations := annotationsForStore()
	dep := &appsv1.StatefulSet{
		ObjectMeta: metav1.ObjectMeta{
			Name:      store.Name,
			Namespace: store.Namespace,
		},
		Spec: appsv1.StatefulSetSpec{
			Replicas: store.Spec.Replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: labels,
			},
			UpdateStrategy: appsv1.StatefulSetUpdateStrategy{
				Type: appsv1.RollingUpdateStatefulSetStrategyType,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels:      labels,
					Annotations: annotations,
				},
				Spec: corev1.PodSpec{
					ServiceAccountName: cons.ServiceAccountName,
					Containers: []corev1.Container{{
						Resources:       store.Spec.Resources,
						Image:           store.Spec.Image,
						Name:            cons.StoreContainerName,
						ImagePullPolicy: store.Spec.ImagePullPolicy,
						Env:             getEnvForStore(store),
						Ports:           getPortsForStore(store),
						VolumeMounts:    getVolumeMountsForStore(store),
						Command:         getCommandForStore(store),
					}},
					Volumes: getVolumesForStore(store),
				},
			},
			VolumeClaimTemplates: getVolumeClaimTemplatesForStore(store),
		},
	}
	// Set Store instance as the owner and controller
	controllerutil.SetControllerReference(store, dep, r.Scheme)

	return dep
}

func getEnvForStore(store *vanusv1alpha1.Store) []corev1.EnvVar {
	defaultEnvs := []corev1.EnvVar{{
		Name:      cons.EnvPodIP,
		ValueFrom: &corev1.EnvVarSource{FieldRef: &corev1.ObjectFieldSelector{FieldPath: "status.podIP"}},
	}, {
		Name:      cons.EnvPodName,
		ValueFrom: &corev1.EnvVarSource{FieldRef: &corev1.ObjectFieldSelector{FieldPath: "metadata.name"}},
	}, {
		Name:  cons.EnvLogLevel,
		Value: "DEBUG",
	}}
	return defaultEnvs
}

func getPortsForStore(store *vanusv1alpha1.Store) []corev1.ContainerPort {
	defaultPorts := []corev1.ContainerPort{{
		Name:          cons.ContainerPortNameGrpc,
		ContainerPort: cons.StorePortGrpc,
	}}
	return defaultPorts
}

func getVolumeMountsForStore(store *vanusv1alpha1.Store) []corev1.VolumeMount {
	defaultVolumeMounts := []corev1.VolumeMount{{
		MountPath: cons.ConfigMountPath,
		Name:      cons.StoreConfigMapName,
	}, {
		MountPath: cons.VolumeMountPath,
		Name:      cons.VolumeName,
	}}
	if len(store.Spec.VolumeClaimTemplates) != 0 && store.Spec.VolumeClaimTemplates[0].Name != "" {
		defaultVolumeMounts = append(defaultVolumeMounts, corev1.VolumeMount{
			MountPath: cons.VolumeMountPath,
			Name:      store.Spec.VolumeClaimTemplates[0].Name,
		})
	}
	return defaultVolumeMounts
}

func getCommandForStore(store *vanusv1alpha1.Store) []string {
	defaultCommand := []string{"/bin/sh", "-c", "NODE_ID=${HOSTNAME##*-} /vanus/bin/store"}
	return defaultCommand
}

func getVolumesForStore(store *vanusv1alpha1.Store) []corev1.Volume {
	defaultVolumes := []corev1.Volume{{
		Name: cons.StoreConfigMapName,
		VolumeSource: corev1.VolumeSource{
			ConfigMap: &corev1.ConfigMapVolumeSource{
				LocalObjectReference: corev1.LocalObjectReference{
					Name: cons.StoreConfigMapName,
				},
			}},
	}}
	return defaultVolumes
}

func getVolumeClaimTemplatesForStore(store *vanusv1alpha1.Store) []corev1.PersistentVolumeClaim {
	labels := labelsForStore(store.Name)
	requests := make(map[corev1.ResourceName]resource.Quantity)
	requests[corev1.ResourceStorage] = resource.MustParse(cons.VolumeStorage)
	defaultPersistentVolumeClaims := []corev1.PersistentVolumeClaim{{
		ObjectMeta: metav1.ObjectMeta{
			Labels: labels,
			Name:   cons.VolumeName,
		},
		Spec: corev1.PersistentVolumeClaimSpec{
			AccessModes: []corev1.PersistentVolumeAccessMode{corev1.ReadWriteOnce},
			Resources: corev1.ResourceRequirements{
				Requests: requests,
			},
		},
	}}
	if len(store.Spec.VolumeClaimTemplates) != 0 {
		if store.Spec.VolumeClaimTemplates[0].Name != "" {
			defaultPersistentVolumeClaims[0].Name = store.Spec.VolumeClaimTemplates[0].Name
		}
		defaultPersistentVolumeClaims[0].Spec.Resources = store.Spec.VolumeClaimTemplates[0].Spec.Resources
	}
	return defaultPersistentVolumeClaims
}

func labelsForStore(name string) map[string]string {
	return map[string]string{"app": name}
}

func annotationsForStore() map[string]string {
	return map[string]string{"vanus.dev/metrics.port": "2112"}
}
