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
	"k8s.io/apimachinery/pkg/util/intstr"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/log"

	vanusv1alpha1 "github.com/linkall-labs/vanus-operator/api/v1alpha1"
)

// ControllerReconciler reconciles a Controller object
type ControllerReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=vanus.linkall.com,resources=controllers,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=vanus.linkall.com,resources=controllers/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=vanus.linkall.com,resources=controllers/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Controller object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.12.2/pkg/reconcile
func (r *ControllerReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = log.FromContext(ctx)

	logger := log.Log.WithName("Controller")
	logger.Info("Reconciling Controller.")

	// Fetch the Controller instance
	controller := &vanusv1alpha1.Controller{}
	err := r.Get(ctx, req.NamespacedName, controller)
	if err != nil {
		if errors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile req.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
			// Return and don't requeue
			logger.Info("Controller resource not found. Ignoring since object must be deleted.")
			return ctrl.Result{}, nil
		}
		// Error reading the object - requeue the req.
		logger.Error(err, "Failed to get Controller.")
		return ctrl.Result{RequeueAfter: time.Duration(cons.RequeueIntervalInSecond) * time.Second}, err
	}

	controllerStatefulSet := r.getStatefulSetForController(controller)
	// Create Controller StatefulSet
	// Check if the statefulSet already exists, if not create a new one
	sts := &appsv1.StatefulSet{}
	err = r.Get(ctx, types.NamespacedName{Name: controllerStatefulSet.Name, Namespace: controllerStatefulSet.Namespace}, sts)
	if err != nil {
		if errors.IsNotFound(err) {
			logger.Info("Creating a new Controller StatefulSet.", "StatefulSet.Namespace", controllerStatefulSet.Namespace, "StatefulSet.Name", controllerStatefulSet.Name)
			err = r.Create(ctx, controllerStatefulSet)
			if err != nil {
				logger.Error(err, "Failed to create new Controller StatefulSet", "StatefulSet.Namespace", controllerStatefulSet.Namespace, "StatefulSet.Name", controllerStatefulSet.Name)
				return ctrl.Result{}, err
			} else {
				logger.Info("Successfully create Controller StatefulSet")
			}
		} else {
			logger.Error(err, "Failed to get Controller StatefulSet.")
			return ctrl.Result{RequeueAfter: time.Duration(cons.RequeueIntervalInSecond) * time.Second}, err
		}
	}

	// TODO(jiangkai): Update Controller StatefulSet

	controllerSvc := r.generateSvcForController(controller)
	// Create Controller Service
	// Check if the service already exists, if not create a new one
	svc := &corev1.Service{}
	err = r.Get(ctx, types.NamespacedName{Name: controllerSvc.Name, Namespace: controllerSvc.Namespace}, svc)
	if err != nil {
		if errors.IsNotFound(err) {
			logger.Info("Creating a new Controller Service.", "Service.Namespace", controllerSvc.Namespace, "Service.Name", controllerSvc.Name)
			err = r.Create(ctx, controllerSvc)
			if err != nil {
				logger.Error(err, "Failed to create new Controller Service", "Service.Namespace", controllerSvc.Namespace, "Service.Name", controllerSvc.Name)
				return ctrl.Result{}, err
			} else {
				logger.Info("Successfully create Controller Service")
			}
		} else {
			logger.Error(err, "Failed to get Controller Service.")
			return ctrl.Result{RequeueAfter: time.Duration(cons.RequeueIntervalInSecond) * time.Second}, err
		}
	}

	// TODO(jiangkai): Update Controller Service

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *ControllerReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&vanusv1alpha1.Controller{}).
		Complete(r)
}

// returns a Controller StatefulSet object
func (r *ControllerReconciler) getStatefulSetForController(controller *vanusv1alpha1.Controller) *appsv1.StatefulSet {
	labels := labelsForController(controller.Name)
	dep := &appsv1.StatefulSet{
		ObjectMeta: metav1.ObjectMeta{
			Name:      controller.Name,
			Namespace: controller.Namespace,
		},
		Spec: appsv1.StatefulSetSpec{
			Replicas: controller.Spec.Replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: labels,
			},
			UpdateStrategy: appsv1.StatefulSetUpdateStrategy{
				Type: appsv1.RollingUpdateStatefulSetStrategyType,
			},
			ServiceName: controller.Name,
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: labels,
				},
				Spec: corev1.PodSpec{
					ServiceAccountName: cons.ServiceAccountName,
					Containers: []corev1.Container{{
						Resources:       controller.Spec.Resources,
						Image:           controller.Spec.Image,
						Name:            cons.ControllerContainerName,
						ImagePullPolicy: controller.Spec.ImagePullPolicy,
						Env:             getEnvForController(controller),
						Ports:           getPortsForController(controller),
						VolumeMounts:    getVolumeMountsForController(controller),
						Command:         getCommandForController(controller),
					}},
					Volumes: getVolumesForController(controller),
				},
			},
			VolumeClaimTemplates: getVolumeClaimTemplatesForController(controller),
		},
	}
	// Set Controller instance as the owner and controller
	controllerutil.SetControllerReference(controller, dep, r.Scheme)

	return dep
}

func getEnvForController(controller *vanusv1alpha1.Controller) []corev1.EnvVar {
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

func getPortsForController(controller *vanusv1alpha1.Controller) []corev1.ContainerPort {
	defaultPorts := []corev1.ContainerPort{{
		Name:          cons.ContainerPortNameGrpc,
		ContainerPort: cons.ControllerPortGrpc,
	}, {
		Name:          cons.ContainerPortNameEtcdClient,
		ContainerPort: cons.ControllerPortEtcdClient,
	}, {
		Name:          cons.ContainerPortNameEtcdPeer,
		ContainerPort: cons.ControllerPortEtcdPeer,
	}, {
		Name:          cons.ContainerPortNameMetrics,
		ContainerPort: cons.ControllerPortMetrics,
	}}
	return defaultPorts
}

func getVolumeMountsForController(controller *vanusv1alpha1.Controller) []corev1.VolumeMount {
	defaultVolumeMounts := []corev1.VolumeMount{{
		MountPath: cons.ConfigMountPath,
		Name:      cons.ControllerConfigMapName,
	}}
	if len(controller.Spec.VolumeClaimTemplates) != 0 && controller.Spec.VolumeClaimTemplates[0].Name != "" {
		defaultVolumeMounts = append(defaultVolumeMounts, corev1.VolumeMount{
			MountPath: cons.VolumeMountPath,
			Name:      controller.Spec.VolumeClaimTemplates[0].Name,
		})
	}
	return defaultVolumeMounts
}

func getCommandForController(controller *vanusv1alpha1.Controller) []string {
	defaultCommand := []string{"/bin/sh", "-c", "NODE_ID=${HOSTNAME##*-} /vanus/bin/controller"}
	return defaultCommand
}

func getVolumesForController(controller *vanusv1alpha1.Controller) []corev1.Volume {
	defaultVolumes := []corev1.Volume{{
		Name: cons.ControllerConfigMapName,
		VolumeSource: corev1.VolumeSource{
			ConfigMap: &corev1.ConfigMapVolumeSource{
				LocalObjectReference: corev1.LocalObjectReference{
					Name: cons.ControllerConfigMapName,
				},
			}},
	}}
	return defaultVolumes
}

func getVolumeClaimTemplatesForController(controller *vanusv1alpha1.Controller) []corev1.PersistentVolumeClaim {
	labels := labelsForController(controller.Name)
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
	if len(controller.Spec.VolumeClaimTemplates) != 0 {
		if controller.Spec.VolumeClaimTemplates[0].Name != "" {
			defaultPersistentVolumeClaims[0].Name = controller.Spec.VolumeClaimTemplates[0].Name
		}
		defaultPersistentVolumeClaims[0].Spec.Resources = controller.Spec.VolumeClaimTemplates[0].Spec.Resources
	}
	return defaultPersistentVolumeClaims
}

func labelsForController(name string) map[string]string {
	return map[string]string{"app": name}
}

func (r *ControllerReconciler) generateSvcForController(controller *vanusv1alpha1.Controller) *corev1.Service {
	labels := labelsForController(controller.Name)
	controllerSvc := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Namespace:  controller.Namespace,
			Name:       controller.Name,
			Labels:     labels,
			Finalizers: []string{metav1.FinalizerOrphanDependents},
		},
		Spec: corev1.ServiceSpec{
			ClusterIP: cons.HeadlessService,
			Selector:  labels,
			Ports: []corev1.ServicePort{
				{
					Name:       controller.Name,
					Port:       cons.ControllerPortGrpc,
					Protocol:   corev1.ProtocolTCP,
					TargetPort: intstr.FromInt(cons.ControllerPortGrpc),
				},
			},
		},
	}

	controllerutil.SetControllerReference(controller, controllerSvc, r.Scheme)
	return controllerSvc
}
